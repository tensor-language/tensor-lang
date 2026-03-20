package elf

import (
	"bytes"
	stdelf "debug/elf" // Kept only for LoadSharedObject convenience
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// InputObject represents a parsed .o file
type InputObject struct {
	Name     string
	Sections []*InputSection
	Symbols  []*InputSymbol
}

type InputSection struct {
	Name   string
	Type   uint32
	Flags  uint64
	Data   []byte
	Relocs []InputReloc

	// Output mapping
	VirtualAddress uint64
	OutputOffset   uint64
}

type InputReloc struct {
	Offset uint64
	Type   uint32
	Addend int64
	Sym    *InputSymbol
}

type InputSymbol struct {
	Name    string
	Type    uint8
	Bind    uint8
	Section *InputSection // Nil if Undefined
	Value   uint64
	Size    uint64
}

// SharedObject represents a dynamic library (e.g., libc.so)
type SharedObject struct {
	Name    string
	Symbols []string 
}

// Internal ELF structures for manual parsing
type elfHeader struct {
	Ident     [16]byte
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     uint64
	Phoff     uint64
	Shoff     uint64
	Flags     uint32
	Ehsize    uint16
	Phentsize uint16
	Phnum     uint16
	Shentsize uint16
	Shnum     uint16
	Shstrndx  uint16
}

type sectionHeader struct {
	Name      uint32
	Type      uint32
	Flags     uint64
	Addr      uint64
	Offset    uint64
	Size      uint64
	Link      uint32
	Info      uint32
	Addralign uint64
	Entsize   uint64
}

type elfSymbol struct {
	Name  uint32
	Info  uint8
	Other uint8
	Shndx uint16
	Value uint64
	Size  uint64
}

type elfRela struct {
	Offset uint64
	Info   uint64
	Addend int64
}

// LoadObject parses an ELF object manually to ensure 1:1 index mapping
func LoadObject(name string, data []byte) (*InputObject, error) {
	r := bytes.NewReader(data)
	var hdr elfHeader
	if err := binary.Read(r, Le, &hdr); err != nil {
		return nil, fmt.Errorf("failed to read elf header: %w", err)
	}

	if hdr.Ident[0] != 0x7f || hdr.Ident[1] != 'E' || hdr.Ident[2] != 'L' || hdr.Ident[3] != 'F' {
		return nil, fmt.Errorf("not a valid ELF file")
	}

	// Read Section Headers
	shdrs := make([]sectionHeader, hdr.Shnum)
	if _, err := r.Seek(int64(hdr.Shoff), io.SeekStart); err != nil {
		return nil, err
	}
	if err := binary.Read(r, Le, &shdrs); err != nil {
		return nil, err
	}

	// Read ShStrTab
	if hdr.Shstrndx >= hdr.Shnum {
		return nil, fmt.Errorf("invalid shstrndx")
	}
	shstrSec := shdrs[hdr.Shstrndx]
	shstrData := make([]byte, shstrSec.Size)
	r.Seek(int64(shstrSec.Offset), io.SeekStart)
	io.ReadFull(r, shstrData)

	getName := func(idx uint32) string {
		if idx >= uint32(len(shstrData)) { return "" }
		end := idx
		for end < uint32(len(shstrData)) && shstrData[end] != 0 {
			end++
		}
		return string(shstrData[idx:end])
	}

	obj := &InputObject{Name: name}
	sectionsByIndex := make(map[int]*InputSection)

	// 1. Load Sections
	for i, sh := range shdrs {
		if sh.Type == SHT_PROGBITS || sh.Type == SHT_NOBITS {
			secName := getName(sh.Name)
			var secData []byte
			if sh.Type == SHT_PROGBITS {
				secData = make([]byte, sh.Size)
				r.Seek(int64(sh.Offset), io.SeekStart)
				io.ReadFull(r, secData)
			}

			isec := &InputSection{
				Name:  secName,
				Type:  sh.Type,
				Flags: sh.Flags,
				Data:  secData,
			}
			obj.Sections = append(obj.Sections, isec)
			sectionsByIndex[i] = isec
		}
	}

	// 2. Load Symbols
	var symTabHdr *sectionHeader
	var strTabHdr *sectionHeader

	for i := range shdrs {
		if shdrs[i].Type == SHT_SYMTAB {
			symTabHdr = &shdrs[i]
			if shdrs[i].Link < uint32(len(shdrs)) {
				strTabHdr = &shdrs[shdrs[i].Link]
			}
			break
		}
	}

	if symTabHdr != nil && strTabHdr != nil {
		strTabData := make([]byte, strTabHdr.Size)
		r.Seek(int64(strTabHdr.Offset), io.SeekStart)
		io.ReadFull(r, strTabData)

		getSymName := func(idx uint32) string {
			if idx >= uint32(len(strTabData)) { return "" }
			end := idx
			for end < uint32(len(strTabData)) && strTabData[end] != 0 {
				end++
			}
			return string(strTabData[idx:end])
		}

		numSyms := symTabHdr.Size / symTabHdr.Entsize
		r.Seek(int64(symTabHdr.Offset), io.SeekStart)

		for k := uint64(0); k < numSyms; k++ {
			var sym elfSymbol
			binary.Read(r, Le, &sym)

			isym := &InputSymbol{
				Name:  getSymName(sym.Name),
				Type:  sym.Info & 0xf,
				Bind:  sym.Info >> 4,
				Value: sym.Value,
				Size:  sym.Size,
			}

			if sym.Shndx < 0xFF00 {
				if sec, ok := sectionsByIndex[int(sym.Shndx)]; ok {
					isym.Section = sec
				}
			}
			obj.Symbols = append(obj.Symbols, isym)
		}
	}

	// 3. Load Relocations
	for _, sh := range shdrs {
		if sh.Type == SHT_RELA {
			targetSec := sectionsByIndex[int(sh.Info)]
			if targetSec == nil { continue }

			numRelocs := sh.Size / sh.Entsize
			r.Seek(int64(sh.Offset), io.SeekStart)

			for k := uint64(0); k < numRelocs; k++ {
				var rel elfRela
				binary.Read(r, Le, &rel)

				symIdx := int(rel.Info >> 32)
				rType := uint32(rel.Info & 0xffffffff)

				if symIdx < len(obj.Symbols) {
					targetSec.Relocs = append(targetSec.Relocs, InputReloc{
						Offset: rel.Offset,
						Type:   rType,
						Addend: rel.Addend,
						Sym:    obj.Symbols[symIdx],
					})
				}
			}
		}
	}

	return obj, nil
}

// LoadArchive iterates a .a file and returns all contained ELF objects.
func LoadArchive(path string) ([]*InputObject, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	magic := make([]byte, 8)
	if _, err := f.Read(magic); err != nil || string(magic) != "!<arch>\n" {
		return nil, fmt.Errorf("not a valid archive: %s", path)
	}

	var objects []*InputObject

	for {
		header := make([]byte, 60)
		if _, err := io.ReadFull(f, header); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		name := strings.TrimSpace(string(header[0:16]))
		sizeStr := strings.TrimSpace(string(header[48:58]))
		size, _ := strconv.ParseInt(sizeStr, 10, 64)

		content := make([]byte, size)
		if _, err := io.ReadFull(f, content); err != nil {
			return nil, err
		}

		if size%2 != 0 {
			f.Seek(1, io.SeekCurrent)
		}

		if name == "/" || name == "//" || name == "/SYM64/" {
			continue
		}

		// Heuristic: Check for ELF magic inside archive member
		if len(content) > 4 && string(content[1:4]) == "ELF" {
			obj, err := LoadObject(name, content)
			if err == nil {
				objects = append(objects, obj)
			}
		}
	}
	return objects, nil
}

func LoadSharedObject(name string, data []byte) (*SharedObject, error) {
	f, err := stdelf.NewFile(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse so %s: %w", name, err)
	}

	so := &SharedObject{Name: name}

	syms, err := f.DynamicSymbols()
	if err != nil {
		fmt.Printf("WARNING: No dynamic symbols in %s: %v\n", name, err)
		return so, nil
	}

	for _, s := range syms {
		if s.Section != stdelf.SHN_UNDEF && (stdelf.ST_BIND(s.Info) == stdelf.STB_GLOBAL || stdelf.ST_BIND(s.Info) == stdelf.STB_WEAK) {
			so.Symbols = append(so.Symbols, s.Name)
		}
	}
	
	fmt.Printf("DEBUG: Loaded shared object '%s' with %d exported symbols\n", name, len(so.Symbols))
	return so, nil
}