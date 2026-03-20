package elf

import "encoding/binary"

// ELF Constants
const (
	EI_NIDENT = 16

	// Magic
	ELFMAG0 = 0x7f
	ELFMAG1 = 'E'
	ELFMAG2 = 'L'
	ELFMAG3 = 'F'

	// Class / Data
	ELFCLASS64  = 2
	ELFDATA2LSB = 1
	EV_CURRENT  = 1

	// OS ABI
	ELFOSABI_NONE  = 0
	ELFOSABI_LINUX = 3

	// Object Types
	ET_REL  = 1 // Relocatable file (.o)
	ET_EXEC = 2 // Executable file
	ET_DYN  = 3 // Shared object (.so)

	// Machine
	EM_X86_64 = 62

	// Section Types
	SHT_NULL     = 0
	SHT_PROGBITS = 1
	SHT_SYMTAB   = 2
	SHT_STRTAB   = 3
	SHT_RELA     = 4
	SHT_DYNAMIC  = 6 // Added: Dynamic linking information
	SHT_NOBITS   = 8

	// Section Flags
	SHF_WRITE     = 0x1
	SHF_ALLOC     = 0x2
	SHF_EXECINSTR = 0x4

	// Segment Types (Phdr)
	PT_NULL    = 0
	PT_LOAD    = 1
	PT_DYNAMIC = 2
	PT_INTERP  = 3
	PT_PHDR    = 6

	// Segment Flags
	PF_X = 0x1 // Execute
	PF_W = 0x2 // Write
	PF_R = 0x4 // Read

	// Symbol Bindings
	STB_LOCAL  = 0
	STB_GLOBAL = 1
	STB_WEAK   = 2

	// Symbol Types
	STT_NOTYPE = 0
	STT_OBJECT = 1
	STT_FUNC   = 2

	// Relocations (AMD64)
	R_X86_64_NONE     = 0
	R_X86_64_64       = 1
	R_X86_64_PC32     = 2
	R_X86_64_PLT32    = 4
	R_X86_64_GLOB_DAT = 6
	R_X86_64_JMP_SLOT = 7
	R_X86_64_RELATIVE = 8
	R_X86_64_32       = 10
	R_X86_64_32S      = 11

	// Dynamic Array Tags (d_tag)
	DT_NULL     = 0
	DT_NEEDED   = 1
	DT_STRTAB   = 5
	DT_SYMTAB   = 6
	DT_RELA     = 7
	DT_RELASZ   = 8
	DT_RELAENT  = 9
	DT_STRSZ    = 10
	DT_SYMENT   = 11
	DT_INIT     = 12
	DT_FINI     = 13
	DT_SONAME   = 14
	DT_RPATH    = 15
	DT_REL      = 17
	DT_DEBUG    = 21
	DT_TEXTREL  = 22
	DT_JMPREL   = 23
	DT_BIND_NOW = 24

	DT_PLTRELSZ = 2
	DT_PLTREL   = 20

	DT_PLTGOT   = 3
)

// Internal ELF Headers used for writing

// Header represents the ELF File Header (Elf64_Ehdr)
type Header struct {
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

// ProgHeader represents a Program Header (Elf64_Phdr)
type ProgHeader struct {
	Type   uint32
	Flags  uint32
	Off    uint64
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

// SectionHeader represents a Section Header (Elf64_Shdr)
type SectionHeader struct {
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

// Elf64Sym represents the symbol structure in the binary (24 bytes)
type Elf64Sym struct {
	Name  uint32
	Info  uint8
	Other uint8
	Shndx uint16
	Value uint64
	Size  uint64
}

// Elf64Rela represents a relocation entry (24 bytes)
type Elf64Rela struct {
	Offset uint64
	Info   uint64
	Addend int64
}

// ByteOrder helper
var Le = binary.LittleEndian