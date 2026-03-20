package elf

import (
	"bytes"
	"encoding/binary"
	"io"
)

// --- ELF64 Constants ---
const (
	EI_NIDENT = 16
	EI_MAG0   = 0
	ELFMAG0   = 0x7f
	ELFMAG1   = 'E'
	ELFMAG2   = 'L'
	ELFMAG3   = 'F'
	EI_CLASS  = 4
	EI_DATA   = 5
	EI_VERSION = 6

	ELFCLASS64  = 2
	ELFDATA2LSB = 1
	EV_CURRENT  = 1

	// Object Types
	ET_NONE = 0
	ET_REL  = 1
	ET_EXEC = 2

	// Machine
	EM_X86_64 = 62

	// Section Types
	SHT_NULL     = 0
	SHT_PROGBITS = 1
	SHT_SYMTAB   = 2
	SHT_STRTAB   = 3
	SHT_RELA     = 4
	SHT_NOBITS   = 8

	// Section Flags
	SHF_WRITE     = 0x1
	SHF_ALLOC     = 0x2
	SHF_EXECINSTR = 0x4

	// Segment Types
	PT_NULL = 0
	PT_LOAD = 1

	// Segment Flags
	PF_X = 0x1
	PF_W = 0x2
	PF_R = 0x4

	// Symbol Binding
	STB_LOCAL  = 0
	STB_GLOBAL = 1

	// Symbol Types
	STT_NOTYPE = 0
	STT_OBJECT = 1
	STT_FUNC   = 2
	STT_SECTION = 3

	// Special Indices
	SHN_UNDEF = 0
)

// --- Structures ---

type File struct {
	Type           uint16
	Sections       []*Section
	Symbols        []*Symbol
	ProgramHeaders []*ProgramHeader
	
	// Internal tables
	StrTab   *StringTable
	ShStrTab *StringTable
	
	Entry uint64
}

type Section struct {
	Name      string
	Type      uint32
	Flags     uint64
	Addr      uint64
	Addralign uint64
	Entsize   uint64
	Link      uint32
	Info      uint32
	Content   []byte

	// Relocations for this section (if Type == SHT_PROGBITS)
	Relocations []Relocation

	// Internal
	Index   uint16
	nameIdx uint32
	offset  uint64
	size    uint64
}

type Symbol struct {
	Name    string
	Info    byte
	Section *Section
	Value   uint64
	Size    uint64
	
	// Internal
	nameIdx uint32
	symIdx  int
}

type Relocation struct {
	Offset uint64
	Symbol *Symbol
	Type   uint32
	Addend int64
}

type ProgramHeader struct {
	Type   uint32
	Flags  uint32
	Off    uint64
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

// --- String Table Helper ---

type StringTable struct {
	Data []byte
	strs map[string]uint32
}

func NewStringTable() *StringTable {
	return &StringTable{
		Data: []byte{0}, 
		strs: make(map[string]uint32),
	}
}

func (st *StringTable) Add(s string) uint32 {
	if s == "" { return 0 }
	if idx, ok := st.strs[s]; ok { return idx }
	idx := uint32(len(st.Data))
	st.Data = append(st.Data, []byte(s)...)
	st.Data = append(st.Data, 0)
	st.strs[s] = idx
	return idx
}

// --- Constructor ---

func NewFile() *File {
	f := &File{
		Type:     ET_REL,
		StrTab:   NewStringTable(),
		ShStrTab: NewStringTable(),
	}
	// Null Section [0]
	f.Sections = append(f.Sections, &Section{Name: "", Type: SHT_NULL})
	return f
}

// --- API ---

func (f *File) AddSection(name string, typ uint32, flags uint64, content []byte) *Section {
	s := &Section{
		Name:    name,
		Type:    typ,
		Flags:   flags,
		Content: content,
		Index:   uint16(len(f.Sections)),
	}
	f.Sections = append(f.Sections, s)
	return s
}

func (f *File) AddSymbol(name string, info byte, section *Section, value, size uint64) *Symbol {
	s := &Symbol{
		Name:    name,
		Info:    info,
		Section: section,
		Value:   value,
		Size:    size,
	}
	f.Symbols = append(f.Symbols, s)
	return s
}

func (f *File) AddProgramHeader(typ, flags uint32, off, vaddr, filesz, memsz, align uint64) {
	f.ProgramHeaders = append(f.ProgramHeaders, &ProgramHeader{
		Type: typ, Flags: flags, Off: off, Vaddr: vaddr, Paddr: vaddr, Filesz: filesz, Memsz: memsz, Align: align,
	})
}

// AddRelocation records that 'sec' needs a patch at 'offset' using 'sym'
func (f *File) AddRelocation(sec *Section, sym *Symbol, offset uint64, typ uint32, addend int64) {
	sec.Relocations = append(sec.Relocations, Relocation{
		Offset: offset,
		Symbol: sym,
		Type:   typ,
		Addend: addend,
	})
}

func MakeSymbolInfo(binding, typ byte) byte {
	return (binding << 4) | (typ & 0xf)
}

// --- Writer ---

func (f *File) WriteTo(w io.Writer) error {
	// 1. Prepare Tables
	shstrtab := f.AddSection(".shstrtab", SHT_STRTAB, 0, nil)
	strtab   := f.AddSection(".strtab",   SHT_STRTAB, 0, nil)
	symtab   := f.AddSection(".symtab",   SHT_SYMTAB, 0, nil)
	symtab.Entsize = 24
	symtab.Addralign = 8

	// 2. Create .rela sections for sections that have relocations
	// We must do this before finalizing indices
	for _, sec := range f.Sections {
		if len(sec.Relocations) > 0 {
			rela := f.AddSection(".rela"+sec.Name, SHT_RELA, 0, nil)
			rela.Link = uint32(symtab.Index) // Points to symbol table
			rela.Info = uint32(sec.Index)    // Points to modified section
			rela.Entsize = 24
			rela.Addralign = 8
			
			// Store relocations temporarily to serialize later
			// (We need finalized symbol indices first)
		}
	}

	// 3. Finalize Symbols & Strings
	// Note: Local symbols must come before Global symbols in the table.
	// For simplicity in this non-optimizing writer, we write them as provided
	// but we must calculate the "First Global" index.
	// A proper implementation sorts them. We will just index them linearly.
	
	// Add symbol strings
	for _, sym := range f.Symbols {
		sym.nameIdx = f.StrTab.Add(sym.Name)
	}

	// Prepare symbol buffer
	// Add Null Symbol first
	allSyms := append([]*Symbol{{}}, f.Symbols...)
	symBuf := new(bytes.Buffer)
	
	firstGlobal := 0
	for i, sym := range allSyms {
		sym.symIdx = i
		// Naive check for first global (if not sorted, this might be inaccurate per spec, 
		// but works for simple object files)
		if (sym.Info>>4) == STB_GLOBAL && firstGlobal == 0 {
			firstGlobal = i
		}
		
		secIdx := uint16(SHN_UNDEF)
		if sym.Section != nil {
			secIdx = sym.Section.Index
		}
		
		// Write Symbol
		binary.Write(symBuf, binary.LittleEndian, sym.nameIdx)
		symBuf.WriteByte(sym.Info)
		symBuf.WriteByte(0) // Other
		binary.Write(symBuf, binary.LittleEndian, secIdx)
		binary.Write(symBuf, binary.LittleEndian, sym.Value)
		binary.Write(symBuf, binary.LittleEndian, sym.Size)
	}
	symtab.Content = symBuf.Bytes()
	symtab.Info = uint32(firstGlobal)
	symtab.Link = uint32(strtab.Index)

	// 4. Serialize Relocations
	for _, sec := range f.Sections {
		// Find the corresponding rela section (brute force for simplicity)
		if len(sec.Relocations) > 0 {
			relaName := ".rela" + sec.Name
			var relaSec *Section
			for _, s := range f.Sections {
				if s.Name == relaName { relaSec = s; break }
			}
			
			if relaSec != nil {
				rbuf := new(bytes.Buffer)
				for _, r := range sec.Relocations {
					// Elf64_Rela
					// r_offset (8)
					// r_info (8) = (symIdx << 32) | type
					// r_addend (8)
					
					symIdx := 0
					if r.Symbol != nil {
						symIdx = r.Symbol.symIdx
					}
					
					info := (uint64(symIdx) << 32) | uint64(r.Type)
					
					binary.Write(rbuf, binary.LittleEndian, r.Offset)
					binary.Write(rbuf, binary.LittleEndian, info)
					binary.Write(rbuf, binary.LittleEndian, r.Addend)
				}
				relaSec.Content = rbuf.Bytes()
			}
		}
	}

	// 5. Finalize Section Names & Strings
	for _, sec := range f.Sections {
		sec.nameIdx = f.ShStrTab.Add(sec.Name)
		if sec.size == 0 {
			sec.size = uint64(len(sec.Content))
		}
	}
	strtab.Content = f.StrTab.Data
	shstrtab.Content = f.ShStrTab.Data

	// 6. Calculate Offsets
	// Header (64) + Phdrs
	offset := uint64(64) + uint64(len(f.ProgramHeaders)*56)
	
	for _, sec := range f.Sections {
		// Align
		if sec.Addralign > 0 && offset % sec.Addralign != 0 {
			offset += sec.Addralign - (offset % sec.Addralign)
		}
		
		// If explicit offset (from static linking)
		if sec.offset == 0 && sec.Type != SHT_NULL && sec.Type != SHT_NOBITS {
			sec.offset = offset
			offset += uint64(len(sec.Content))
		} else if sec.Type == SHT_NOBITS {
			// NOBITS sections have virtual address but take no file space
			// We track the virtual offset if needed, but file offset stays
			sec.offset = offset // usually conceptuall same
		}
	}
	shOff := offset // Section Header Table starts here

	// 7. Write Data
	
	// Header
	hdr := elfHeader{
		Type: f.Type, Machine: EM_X86_64, Version: EV_CURRENT, Entry: f.Entry,
		Phoff: 64, Shoff: shOff, Ehsize: 64, Phentsize: 56, Phnum: uint16(len(f.ProgramHeaders)),
		Shentsize: 64, Shnum: uint16(len(f.Sections)), Shstrndx: uint16(shstrtab.Index),
	}
	hdr.Ident[0]=0x7F; hdr.Ident[1]='E'; hdr.Ident[2]='L'; hdr.Ident[3]='F'
	hdr.Ident[4]=ELFCLASS64; hdr.Ident[5]=ELFDATA2LSB; hdr.Ident[6]=EV_CURRENT
	
	binary.Write(w, binary.LittleEndian, hdr)
	
	// Phdrs
	for _, ph := range f.ProgramHeaders {
		binary.Write(w, binary.LittleEndian, elfProgramHeader{
			Type: ph.Type, Flags: ph.Flags, Off: ph.Off, Vaddr: ph.Vaddr, Paddr: ph.Paddr,
			Filesz: ph.Filesz, Memsz: ph.Memsz, Align: ph.Align,
		})
	}

	// Section Content
	curOff := uint64(64) + uint64(len(f.ProgramHeaders)*56)
	for _, sec := range f.Sections {
		if sec.Type == SHT_NULL || sec.Type == SHT_NOBITS { continue }
		
		// Pad
		if sec.offset > curOff {
			w.Write(make([]byte, sec.offset - curOff))
			curOff = sec.offset
		}
		w.Write(sec.Content)
		curOff += uint64(len(sec.Content))
	}
	
	// Section Headers
	if shOff > curOff {
		w.Write(make([]byte, shOff - curOff))
	}
	
	for _, sec := range f.Sections {
		binary.Write(w, binary.LittleEndian, elfSectionHeader{
			Name: sec.nameIdx, Type: sec.Type, Flags: sec.Flags, Addr: sec.Addr,
			Offset: sec.offset, Size: uint64(len(sec.Content)), Link: sec.Link,
			Info: sec.Info, Addralign: sec.Addralign, Entsize: sec.Entsize,
		})
	}

	return nil
}

// Internal Structs for serialization
type elfHeader struct {
	Ident [16]byte; Type, Machine uint16; Version uint32; Entry, Phoff, Shoff uint64; Flags uint32
	Ehsize, Phentsize, Phnum, Shentsize, Shnum, Shstrndx uint16
}
type elfProgramHeader struct {
	Type, Flags uint32; Off, Vaddr, Paddr, Filesz, Memsz, Align uint64
}
type elfSectionHeader struct {
	Name, Type uint32; Flags, Addr, Offset, Size uint64; Link, Info uint32; Addralign, Entsize uint64
}