package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Entry       string
	BaseAddr    uint64
	Interpreter string
}

type Linker struct {
	Config      Config
	Objects     []*InputObject
	SharedLibs  []*SharedObject
	GlobalTable map[string]*ResolvedSymbol

	DynStrTab  []byte
	DynSyms    []Elf64Sym
	RelaDyn    []Elf64Rela
	GotEntries []string

	InterpSect  []byte
	DynSect     []byte
	DynSymSect  []byte
	DynStrSect  []byte
	RelaDynSect []byte
	TextSection []byte
	DataSection []byte

	ShStrTab []byte

	TextAddr  uint64
	DataAddr  uint64
	BssAddr   uint64
	EntryAddr uint64
	BssSize   uint64
	
	SectionOffsets map[*InputSection]uint64
}

type ResolvedSymbol struct {
	Name    string
	Value   uint64
	Section string
	Defined bool
}

func NewLinker(cfg Config) *Linker {
	if cfg.Entry == "" { cfg.Entry = "_start" }
	if cfg.BaseAddr == 0 { cfg.BaseAddr = 0x400000 }
	if cfg.Interpreter == "" { cfg.Interpreter = "/lib64/ld-linux-x86-64.so.2" }

	return &Linker{
		Config:         cfg,
		GlobalTable:    make(map[string]*ResolvedSymbol),
		DynSyms:        []Elf64Sym{{}},
		DynStrTab:      []byte{0},
		ShStrTab:       []byte{0},
		SectionOffsets: make(map[*InputSection]uint64),
	}
}

func (l *Linker) AddObject(name string, data []byte) error {
	obj, err := LoadObject(name, data)
	if err != nil { return err }
	l.Objects = append(l.Objects, obj)
	return nil
}

func (l *Linker) AddArchive(path string) error {
	objs, err := LoadArchive(path)
	if err != nil { return err }
	l.Objects = append(l.Objects, objs...)
	return nil
}

func (l *Linker) AddSharedLib(path string, data []byte) error {
	so, err := LoadSharedObject(path, data)
	if err != nil { return err }
	l.SharedLibs = append(l.SharedLibs, so)
	return nil
}

func (l *Linker) Link(outPath string) error {
	if err := l.scanSymbols(); err != nil { return err }
	l.layout()
	if err := l.applyRelocations(); err != nil { return err }

	return l.write(outPath)
}

func (l *Linker) scanSymbols() error {
	l.GlobalTable[l.Config.Entry] = &ResolvedSymbol{Name: l.Config.Entry, Defined: false}

	for _, obj := range l.Objects {
		for _, sym := range obj.Symbols {
			if sym.Bind == STB_GLOBAL {
				if sym.Section != nil {
					if existing, ok := l.GlobalTable[sym.Name]; ok && existing.Defined {
						return fmt.Errorf("duplicate symbol: %s", sym.Name)
					}
					l.GlobalTable[sym.Name] = &ResolvedSymbol{Name: sym.Name, Defined: true, Section: "text"}
				} else {
					if _, ok := l.GlobalTable[sym.Name]; !ok {
						l.GlobalTable[sym.Name] = &ResolvedSymbol{Name: sym.Name, Defined: false}
					}
				}
			}
		}
	}

	for _, obj := range l.Objects {
		for _, sec := range obj.Sections {
			for _, reloc := range sec.Relocs {
				if reloc.Sym.Section == nil && reloc.Sym.Bind == STB_GLOBAL {
					if _, ok := l.GlobalTable[reloc.Sym.Name]; !ok {
						l.GlobalTable[reloc.Sym.Name] = &ResolvedSymbol{
							Name:    reloc.Sym.Name,
							Defined: false,
						}
					}
				}
			}
		}
	}

	if _, hasMain := l.GlobalTable["main"]; hasMain {
		if _, ok := l.GlobalTable["__libc_start_main"]; !ok {
			l.GlobalTable["__libc_start_main"] = &ResolvedSymbol{Name: "__libc_start_main", Defined: false}
		}
	}

	for _, so := range l.SharedLibs {
		for _, symName := range so.Symbols {
			cleanName := symName
			if idx := strings.Index(cleanName, "@"); idx != -1 {
				cleanName = cleanName[:idx]
			}
			if target, ok := l.GlobalTable[cleanName]; ok && !target.Defined {
				target.Defined = true
				target.Section = "dynamic"
				l.addDynamicSymbol(cleanName)
				l.GotEntries = append(l.GotEntries, cleanName)
			}
		}
	}

	if !l.GlobalTable[l.Config.Entry].Defined {
		if l.GlobalTable["main"].Defined {
			l.GlobalTable[l.Config.Entry].Defined = true
			l.GlobalTable[l.Config.Entry].Section = "stub"
		} else {
			return fmt.Errorf("entry point %s not found", l.Config.Entry)
		}
	}
	
	return nil
}

func (l *Linker) addDynamicSymbol(name string) {
	nameIdx := uint32(len(l.DynStrTab))
	l.DynStrTab = append(l.DynStrTab, []byte(name)...)
	l.DynStrTab = append(l.DynStrTab, 0)
	l.DynSyms = append(l.DynSyms, Elf64Sym{Name: nameIdx, Info: 0x12, Shndx: 0, Value: 0, Size: 0})
}

func (l *Linker) addShStr(s string) uint32 {
	idx := uint32(len(l.ShStrTab))
	l.ShStrTab = append(l.ShStrTab, []byte(s)...)
	l.ShStrTab = append(l.ShStrTab, 0)
	return idx
}

func align(val, align uint64) uint64 {
	if val%align != 0 { return val + (align - (val % align)) }
	return val
}

func (l *Linker) layout() {

	headerSize := uint64(64 + 56*5)
	l.TextAddr = align(l.Config.BaseAddr+headerSize, 4096)
	l.InterpSect = append([]byte(l.Config.Interpreter), 0)

	pltSize := uint64(len(l.GotEntries) * 16)

	var objText []byte
	
	needsStub := l.GlobalTable[l.Config.Entry].Section == "stub"
	
	for _, obj := range l.Objects {
		for _, sec := range obj.Sections {
			// Include Executable (Code) OR Read-Only Data (Alloc + !Write + Progbits)
			isCode := sec.Flags&SHF_EXECINSTR != 0
			isRoData := (sec.Flags&SHF_ALLOC != 0) && (sec.Flags&SHF_WRITE == 0) && (sec.Type == SHT_PROGBITS)

			if isCode || isRoData {
				alignVal := 16
				if !isCode { alignVal = 4 }
				
				pad := (alignVal - (len(objText) % alignVal)) % alignVal
				for i := 0; i < int(pad); i++ { 
					if isCode {
						objText = append(objText, 0x90) // NOP
					} else {
						objText = append(objText, 0)    // Zero
					}
				}
				l.SectionOffsets[sec] = uint64(len(objText))
				objText = append(objText, sec.Data...)
			}
		}
	}

	baseTextOffset := uint64(len(l.InterpSect)) + pltSize
	
	stubSize := uint64(0)
	if needsStub {
		stubSize = 32
		l.GlobalTable[l.Config.Entry].Value = l.TextAddr + baseTextOffset
		objText = append(make([]byte, stubSize), objText...)
		for sec := range l.SectionOffsets {
			l.SectionOffsets[sec] += stubSize
		}
	}
	
	for sec, relOff := range l.SectionOffsets {
		finalOffset := baseTextOffset + relOff
		sec.OutputOffset = finalOffset
		sec.VirtualAddress = l.TextAddr + finalOffset
		l.SectionOffsets[sec] = finalOffset
	}

	totalTextSize := baseTextOffset + uint64(len(objText))
	l.DataAddr = align(l.TextAddr+totalTextSize, 4096)

	libNameOffsets := make(map[string]uint32)
	for _, lib := range l.SharedLibs {
		libName := lib.Name
		if idx := strings.LastIndex(libName, "/"); idx != -1 {
			libName = libName[idx+1:]
		}
		
		if libName == "" {
			continue
		}
		
		nameOffset := uint32(len(l.DynStrTab))
		l.DynStrTab = append(l.DynStrTab, []byte(libName)...)
		l.DynStrTab = append(l.DynStrTab, 0)
		libNameOffsets[lib.Name] = nameOffset
	}

	symBuf := new(bytes.Buffer)
	for _, ds := range l.DynSyms {
		binary.Write(symBuf, Le, ds.Name); symBuf.WriteByte(ds.Info); symBuf.WriteByte(ds.Other)
		binary.Write(symBuf, Le, ds.Shndx); binary.Write(symBuf, Le, ds.Value); binary.Write(symBuf, Le, ds.Size)
	}
	l.DynSymSect = symBuf.Bytes()
	relaSize := uint64(len(l.GotEntries) * 24)

	dynBuf := new(bytes.Buffer)
	writeDyn := func(tag int64, val uint64) {
		binary.Write(dynBuf, Le, tag); binary.Write(dynBuf, Le, val)
	}
	
	for range l.SharedLibs {
		writeDyn(DT_NEEDED, 0)
	}
	writeDyn(DT_PLTGOT, 0)
	writeDyn(DT_BIND_NOW, 1)
	writeDyn(DT_STRTAB, 0)
	writeDyn(DT_SYMTAB, 0)
	writeDyn(DT_STRSZ, uint64(len(l.DynStrTab)))
	writeDyn(DT_SYMENT, 24)
	writeDyn(DT_PLTREL, 7)
	writeDyn(DT_JMPREL, 0)
	writeDyn(DT_PLTRELSZ, relaSize)
	writeDyn(DT_RELAENT, 24)
	writeDyn(DT_NULL, 0)
	
	dynSectSize := uint64(dynBuf.Len())

	offset := dynSectSize
	offset = align(offset, 8); offsetSym := offset; offset += uint64(len(l.DynSymSect))
	offset = align(offset, 8); offsetStr := offset; offset += uint64(len(l.DynStrTab))
	offset = align(offset, 8); offsetRela := offset; offset += relaSize
	offset = align(offset, 8); offsetGot := offset; sizeGot := uint64(len(l.GotEntries) * 8); offset += sizeGot

	symAddr := l.DataAddr + offsetSym
	strAddr := l.DataAddr + offsetStr
	relaAddr := l.DataAddr + offsetRela
	gotAddr := l.DataAddr + offsetGot

	dynBuf.Reset()
	for _, lib := range l.SharedLibs {
		offset, ok := libNameOffsets[lib.Name]
		if !ok || offset == 0 {
			continue
		}
		writeDyn(DT_NEEDED, uint64(offset))
	}
	writeDyn(DT_PLTGOT, gotAddr)
	writeDyn(DT_BIND_NOW, 1)
	writeDyn(DT_STRTAB, strAddr)
	writeDyn(DT_SYMTAB, symAddr)
	writeDyn(DT_STRSZ, uint64(len(l.DynStrTab)))
	writeDyn(DT_SYMENT, 24)
	writeDyn(DT_PLTREL, 7)
	writeDyn(DT_JMPREL, relaAddr)
	writeDyn(DT_PLTRELSZ, relaSize)
	writeDyn(DT_RELAENT, 24)
	writeDyn(DT_NULL, 0)

	l.DynSect = dynBuf.Bytes()

	relaBuf := new(bytes.Buffer)
	for i := range l.GotEntries {
		rOff := gotAddr + uint64(i*8)
		rInfo := uint64((i+1)<<32) | uint64(R_X86_64_JMP_SLOT)
		binary.Write(relaBuf, Le, rOff); binary.Write(relaBuf, Le, rInfo); binary.Write(relaBuf, Le, int64(0))
	}
	l.RelaDynSect = relaBuf.Bytes()

	buf := new(bytes.Buffer)
	buf.Write(l.DynSect)
	for uint64(buf.Len()) < offsetSym { buf.WriteByte(0) }
	buf.Write(l.DynSymSect)
	for uint64(buf.Len()) < offsetStr { buf.WriteByte(0) }
	buf.Write(l.DynStrTab)
	for uint64(buf.Len()) < offsetRela { buf.WriteByte(0) }
	buf.Write(l.RelaDynSect)
	for uint64(buf.Len()) < offsetGot { buf.WriteByte(0) }
	buf.Write(make([]byte, sizeGot))
	l.DataSection = buf.Bytes()

	pltBuf := new(bytes.Buffer)
	for i, symName := range l.GotEntries {
		targetGot := gotAddr + uint64(i*8)
		pltBuf.WriteByte(0xFF); pltBuf.WriteByte(0x25)
		currentPC := l.TextAddr + uint64(len(l.InterpSect)) + uint64(pltBuf.Len()) + 4
		rel := int32(targetGot - currentPC)
		binary.Write(pltBuf, Le, rel)
		for k:=0; k<10; k++ { pltBuf.WriteByte(0x90) }
		
		l.GlobalTable[symName].Value = l.TextAddr + uint64(len(l.InterpSect)) + uint64(i*16)
		l.GlobalTable[symName].Section = "plt"
	}

	l.TextSection = append(l.InterpSect, pltBuf.Bytes()...)
	l.TextSection = append(l.TextSection, objText...)

	currentDataLen := uint64(len(l.DataSection))
	for _, obj := range l.Objects {
		for _, sec := range obj.Sections {
			if sec.Flags&SHF_WRITE != 0 && sec.Type == SHT_PROGBITS {
				pad := (8 - ((currentDataLen) % 8)) % 8
				l.DataSection = append(l.DataSection, make([]byte, int(pad))...)
				currentDataLen += uint64(pad)
				l.SectionOffsets[sec] = currentDataLen
				sec.OutputOffset = currentDataLen
				sec.VirtualAddress = l.DataAddr + sec.OutputOffset
				l.DataSection = append(l.DataSection, sec.Data...)
				currentDataLen += uint64(len(sec.Data))
			}
		}
	}
	l.BssAddr = l.DataAddr + currentDataLen
	l.BssSize = 0 
	
	for _, obj := range l.Objects {
		for _, sym := range obj.Symbols {
			if sym.Section != nil {
				if gSym, ok := l.GlobalTable[sym.Name]; ok {
					gSym.Value = sym.Section.VirtualAddress + sym.Value
				}
			}
		}
	}
	
	l.EntryAddr = l.GlobalTable[l.Config.Entry].Value
}

func (l *Linker) applyRelocations() error {
	if l.GlobalTable[l.Config.Entry].Section == "stub" {
		stub := []byte{
			0x31, 0xed, 0x49, 0x89, 0xd1, 0x5e, 0x48, 0x89, 0xe2, 0x48, 0x83, 0xe4, 0xf0, 
			0x48, 0x8d, 0x3d, 0, 0, 0, 0,
			0x48, 0x31, 0xc9, 0x45, 0x31, 0xc0, 
			0xe8, 0, 0, 0, 0,
			0xf4,
		}
		stubOffset := l.GlobalTable[l.Config.Entry].Value - l.TextAddr
		leaPC := l.GlobalTable[l.Config.Entry].Value + 20
		
		mainSym, ok := l.GlobalTable["main"]
		if !ok || !mainSym.Defined {
			return fmt.Errorf("main symbol not found or undefined")
		}
		mainAddr := mainSym.Value
		
		fmt.Printf("DEBUG: Stub at 0x%x, main at 0x%x, lea PC at 0x%x\n", 
			l.GlobalTable[l.Config.Entry].Value, mainAddr, leaPC)
		
		binary.LittleEndian.PutUint32(stub[16:], uint32(int32(mainAddr - leaPC)))
		callPC := l.GlobalTable[l.Config.Entry].Value + 31
		libcStartAddr := l.GlobalTable["__libc_start_main"].Value
		
		fmt.Printf("DEBUG: call PC at 0x%x, __libc_start_main at 0x%x\n", 
			callPC, libcStartAddr)
		
		binary.LittleEndian.PutUint32(stub[27:], uint32(int32(libcStartAddr - callPC)))
		copy(l.TextSection[stubOffset:], stub)
	}

	for _, obj := range l.Objects {
		for _, sec := range obj.Sections {
			for _, r := range sec.Relocs {
				var symVal uint64
				if r.Sym.Bind == STB_LOCAL {
					symVal = r.Sym.Section.VirtualAddress + r.Sym.Value
				} else {
					global, ok := l.GlobalTable[r.Sym.Name]
					if !ok || !global.Defined { return fmt.Errorf("undefined symbol: %s", r.Sym.Name) }
					symVal = global.Value
				}
				
				P := sec.VirtualAddress + r.Offset
				
				var buf []byte
				var bufOff uint64
				
				if off, ok := l.SectionOffsets[sec]; ok {
					isCode := sec.Flags&SHF_EXECINSTR != 0
					isRoData := (sec.Flags&SHF_ALLOC != 0) && (sec.Flags&SHF_WRITE == 0) && (sec.Type == SHT_PROGBITS)

					if isCode || isRoData {
						buf = l.TextSection
					} else {
						buf = l.DataSection
					}
					bufOff = off + r.Offset
				} else {
					continue
				}

				if bufOff >= uint64(len(buf)) { panic("Buffer overflow in relocation") }
				
				if r.Sym.Name == "printf" || strings.HasPrefix(r.Sym.Name, ".str") {
					fmt.Printf("DEBUG: Reloc %s: Type=%d Off=%x Addend=%d SymVal=%x P=%x Val=%d\n", 
						r.Sym.Name, r.Type, r.Offset, r.Addend, symVal, P, int32(int64(symVal)+r.Addend-int64(P)))
				}

				switch r.Type {
				case R_X86_64_64:
					binary.LittleEndian.PutUint64(buf[bufOff:], uint64(int64(symVal)+r.Addend))
				case R_X86_64_32:
					binary.LittleEndian.PutUint32(buf[bufOff:], uint32(int64(symVal)+r.Addend))
				case R_X86_64_PC32, R_X86_64_PLT32:
					val := int32(int64(symVal) + r.Addend - int64(P))
					binary.LittleEndian.PutUint32(buf[bufOff:], uint32(val))
				}
			}
		}
	}
	return nil
}

func (l *Linker) write(path string) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()

	hdrSize := uint64(64 + 56*5)
	textOff := align(hdrSize, 4096)
	textSize := uint64(len(l.TextSection))
	dataOff := align(textOff+textSize, 4096)
	dataSize := uint64(len(l.DataSection))

	l.addShStr("")
	idxInterp := l.addShStr(".interp")
	idxText := l.addShStr(".text")
	idxDyn := l.addShStr(".dynamic")
	idxData := l.addShStr(".data")
	idxBss := l.addShStr(".bss")
	idxShstr := l.addShStr(".shstrtab")
	
	shStrOff := dataOff + dataSize
	shTableOff := shStrOff + uint64(len(l.ShStrTab))

	interpSize := uint64(len(l.InterpSect))
	pltSize := uint64(len(l.GotEntries) * 16)
	
	codeStartInText := interpSize + pltSize

	ehdr := Header{
		Type: ET_EXEC, Machine: EM_X86_64, Version: EV_CURRENT, 
		Entry: l.EntryAddr,
		Phoff: 64, Ehsize: 64, Phentsize: 56, Phnum: 5, 
		Shoff: shTableOff, Shentsize: 64, Shnum: 7, Shstrndx: 6,
	}
	ehdr.Ident[0]=0x7F; ehdr.Ident[1]='E'; ehdr.Ident[2]='L'; ehdr.Ident[3]='F'
	ehdr.Ident[4]=ELFCLASS64; ehdr.Ident[5]=ELFDATA2LSB; ehdr.Ident[6]=EV_CURRENT; ehdr.Ident[7]=3
	binary.Write(f, Le, ehdr)

	binary.Write(f, Le, ProgHeader{Type: PT_PHDR, Flags: PF_R, Off: 64, Vaddr: l.Config.BaseAddr + 64, Paddr: l.Config.BaseAddr + 64, Filesz: 56 * 5, Memsz: 56 * 5, Align: 8})
	binary.Write(f, Le, ProgHeader{Type: PT_INTERP, Flags: PF_R, Off: textOff, Vaddr: l.TextAddr, Paddr: l.TextAddr, Filesz: interpSize, Memsz: interpSize, Align: 1})
	binary.Write(f, Le, ProgHeader{Type: PT_LOAD, Flags: PF_R | PF_X, Off: 0, Vaddr: l.Config.BaseAddr, Paddr: l.Config.BaseAddr, Filesz: textOff + textSize, Memsz: textOff + textSize, Align: 4096})
	binary.Write(f, Le, ProgHeader{Type: PT_LOAD, Flags: PF_R | PF_W, Off: dataOff, Vaddr: l.DataAddr, Paddr: l.DataAddr, Filesz: dataSize, Memsz: dataSize + l.BssSize, Align: 4096})
	binary.Write(f, Le, ProgHeader{Type: PT_DYNAMIC, Flags: PF_R | PF_W, Off: dataOff, Vaddr: l.DataAddr, Paddr: l.DataAddr, Filesz: uint64(len(l.DynSect)), Memsz: uint64(len(l.DynSect)), Align: 8})

	if textOff > 64+56*5 { 
		f.Write(make([]byte, textOff-(64+56*5))) 
	}
	f.Write(l.TextSection)
	
	if dataOff > textOff+textSize { 
		f.Write(make([]byte, dataOff-(textOff+textSize))) 
	}
	f.Write(l.DataSection)
	f.Write(l.ShStrTab)

	binary.Write(f, Le, SectionHeader{})
	
	binary.Write(f, Le, SectionHeader{
		Name: idxInterp, Type: SHT_PROGBITS, Flags: SHF_ALLOC, 
		Addr: l.TextAddr, 
		Offset: textOff, 
		Size: interpSize, 
		Addralign: 1,
	})
	
	binary.Write(f, Le, SectionHeader{
		Name: idxText, Type: SHT_PROGBITS, Flags: SHF_ALLOC | SHF_EXECINSTR, 
		Addr: l.TextAddr + codeStartInText, 
		Offset: textOff + codeStartInText, 
		Size: textSize - codeStartInText, 
		Addralign: 16,
	})
	
	dynStructsSize := uint64(len(l.DynSect))
	dynStructsSize = align(dynStructsSize, 8)
	dynStructsSize += uint64(len(l.DynSymSect))
	dynStructsSize = align(dynStructsSize, 8)
	dynStructsSize += uint64(len(l.DynStrTab))
	dynStructsSize = align(dynStructsSize, 8)
	dynStructsSize += uint64(len(l.RelaDynSect))
	dynStructsSize = align(dynStructsSize, 8)
	dynStructsSize += uint64(len(l.GotEntries) * 8)
	dynStructsSize = align(dynStructsSize, 8)

	binary.Write(f, Le, SectionHeader{Name: idxDyn, Type: SHT_DYNAMIC, Flags: SHF_ALLOC | SHF_WRITE, Addr: l.DataAddr, Offset: dataOff, Size: uint64(len(l.DynSect)), Link: 0, Addralign: 8})
	
	dataContentSize := dataSize
	if dataSize > dynStructsSize {
		dataContentSize = dataSize - dynStructsSize
	} else {
		dataContentSize = 0
	}

	binary.Write(f, Le, SectionHeader{
		Name: idxData, Type: SHT_PROGBITS, Flags: SHF_ALLOC | SHF_WRITE, 
		Addr: l.DataAddr + dynStructsSize, 
		Offset: dataOff + dynStructsSize, 
		Size: dataContentSize, 
		Addralign: 4096,
	})

	binary.Write(f, Le, SectionHeader{Name: idxBss, Type: SHT_NOBITS, Flags: SHF_ALLOC | SHF_WRITE, Addr: l.BssAddr, Offset: dataOff + dataSize, Size: l.BssSize, Addralign: 8})
	binary.Write(f, Le, SectionHeader{Name: idxShstr, Type: SHT_STRTAB, Flags: 0, Addr: 0, Offset: shStrOff, Size: uint64(len(l.ShStrTab)), Addralign: 1})

	f.Chmod(0755)
	return nil
}