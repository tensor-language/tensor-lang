package backend

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/backend/cpu/amd64"
	"github.com/arc-language/arc-lang/backend/format/elf"
	"github.com/arc-language/arc-lang/backend/gpu/nvidia"
	"github.com/arc-language/arc-lang/backend/gpu/amd"
	"github.com/arc-language/arc-lang/backend/tpu"
)

// Generate handles the backend code generation for the module.
// It automatically detects GPU/TPU kernels and invokes the appropriate backend.
func Generate(module *ir.Module) error {
	// 1. Scan for Hardware Specific Functions
	hasGPUFunctions  := false
	hasROCmFunctions := false
	hasTPUFunctions  := false

	for _, fn := range module.Functions {
		switch fn.CallConv {
		case ir.CC_PTX:
			hasGPUFunctions = true
		case ir.CC_ROCM:
			hasROCmFunctions = true
		case ir.CC_TPU:
			hasTPUFunctions = true
		}
	}

	// 2. Generate NVIDIA PTX if needed
	if hasGPUFunctions {
		log.Println("[backend] NVIDIA GPU functions detected (CC_PTX). Generating PTX Assembly...")

		ptxCode, err := nvidia.Generate(module)
		if err != nil {
			return fmt.Errorf("NVIDIA backend failed: %w", err)
		}

		fmt.Println("\n// ==========================================")
		fmt.Println("// GENERATED NVIDIA PTX ASSEMBLY")
		fmt.Println("// ==========================================")
		fmt.Println(ptxCode)
		fmt.Println("// ==========================================\n")
	}

	// 3. Generate AMD ROCm (GCN) if needed
	if hasROCmFunctions {
		log.Println("[backend] AMD GPU functions detected (CC_ROCM). Generating GCN Assembly...")

		gcnCode, err := amd.Generate(module)
		if err != nil {
			return fmt.Errorf("AMD backend failed: %w", err)
		}

		fmt.Println("\n// ==========================================")
		fmt.Println("// GENERATED AMD GCN ASSEMBLY")
		fmt.Println("// ==========================================")
		fmt.Println(gcnCode)
		fmt.Println("// ==========================================\n")
	}

	// 4. Generate Google TPU HLO if needed
	if hasTPUFunctions {
		log.Println("[backend] TPU functions detected (CC_TPU). Generating Google HLO...")

		hloCode, err := tpu.Generate(module)
		if err != nil {
			return fmt.Errorf("TPU backend failed: %w", err)
		}

		fmt.Println("\n// ==========================================")
		fmt.Println("// GENERATED TPU HLO IR")
		fmt.Println("// ==========================================")
		fmt.Println(hloCode)
		fmt.Println("// ==========================================\n")
	}

	// 5. Always proceed with standard CPU code generation.
	log.Println("[backend] Generating x86-64 CPU code...")

	objBytes, err := GenerateObject(module)
	if err != nil {
		return fmt.Errorf("CPU backend failed: %w", err)
	}

	log.Printf("[backend] Generated %d bytes of x86-64 object code.", len(objBytes))
	return nil
}

// GenerateObject creates a relocatable ELF object file (.o)
func GenerateObject(m *ir.Module) ([]byte, error) {
	// 1. Compile IR to machine code
	artifact, err := amd64.Compile(m)
	if err != nil {
		return nil, err
	}

	// 2. Wrap in ELF container
	f := elf.NewFile()
	f.Type = elf.ET_REL

	textSec := f.AddSection(".text", elf.SHT_PROGBITS, elf.SHF_ALLOC|elf.SHF_EXECINSTR, artifact.Text)
	textSec.Addralign = 16

	dataSec := f.AddSection(".data", elf.SHT_PROGBITS, elf.SHF_ALLOC|elf.SHF_WRITE, artifact.Data)
	dataSec.Addralign = 8

	// 3. Add symbols
	symMap := make(map[string]*elf.Symbol)

	for _, symDef := range artifact.Symbols {
		var sec *elf.Section
		if symDef.IsFunc {
			sec = textSec
		} else {
			sec = dataSec
		}

		info := elf.MakeSymbolInfo(elf.STB_GLOBAL, elf.STT_OBJECT)
		if symDef.IsFunc {
			info = elf.MakeSymbolInfo(elf.STB_GLOBAL, elf.STT_FUNC)
		}

		es := f.AddSymbol(symDef.Name, info, sec, symDef.Offset, symDef.Size)
		symMap[symDef.Name] = es
	}

	// 4. Process relocations
	for _, reloc := range artifact.Relocs {
		sym, ok := symMap[reloc.Symbol]
		if !ok {
			sym = f.AddSymbol(reloc.Symbol, elf.MakeSymbolInfo(elf.STB_GLOBAL, elf.STT_NOTYPE), nil, 0, 0)
			symMap[reloc.Symbol] = sym
		}
		f.AddRelocation(textSec, sym, uint64(reloc.Offset), uint32(reloc.Type), reloc.Addend)
	}

	// 5. Write to buffer
	buf := new(bytes.Buffer)
	if err := f.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GenerateExecutable creates a static ELF executable directly (no libc).
func GenerateExecutable(m *ir.Module) ([]byte, error) {
	// 1. Compile to machine code
	artifact, err := amd64.Compile(m)
	if err != nil {
		return nil, err
	}

	// 2. Setup memory layout
	const (
		BaseAddr = 0x400000
		PageSize = 0x1000
	)

	// Entry stub: _start calls main then exit(0)
	entryStub := []byte{
		0xE8, 0x00, 0x00, 0x00, 0x00, // call main (placeholder)
		0x48, 0x31, 0xFF,             // xor rdi, rdi
		0x48, 0xC7, 0xC0, 0x3C, 0x00, 0x00, 0x00, // mov rax, 60
		0x0F, 0x05,                   // syscall
	}

	stubSize  := len(entryStub)
	finalText := append(entryStub, artifact.Text...)

	textVAddr := uint64(BaseAddr + PageSize)
	dataVAddr := textVAddr + uint64(len(finalText))
	if dataVAddr%PageSize != 0 {
		dataVAddr += PageSize - (dataVAddr % PageSize)
	}

	// 3. Resolve main for the stub
	var mainOffset uint64
	foundMain := false
	for _, sym := range artifact.Symbols {
		if sym.Name == "main" {
			mainOffset = sym.Offset
			foundMain = true
			break
		}
	}
	if !foundMain {
		return nil, fmt.Errorf("entry point 'main' not found")
	}

	rel := int32((stubSize + int(mainOffset)) - 5)
	binary.LittleEndian.PutUint32(finalText[1:], uint32(rel))

	// 4. Resolve internal relocations
	for _, reloc := range artifact.Relocs {
		var symAddr uint64
		found := false

		for _, sym := range artifact.Symbols {
			if sym.Name == reloc.Symbol {
				if sym.IsFunc {
					symAddr = textVAddr + uint64(stubSize) + sym.Offset
				} else {
					symAddr = dataVAddr + sym.Offset
				}
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("undefined symbol '%s' (static linking does not support external libs)", reloc.Symbol)
		}

		pc          := textVAddr + uint64(stubSize) + uint64(reloc.Offset)
		val         := int32(int64(symAddr) - int64(pc) + reloc.Addend)
		patchOffset := stubSize + int(reloc.Offset)
		binary.LittleEndian.PutUint32(finalText[patchOffset:], uint32(val))
	}

	// 5. Create ELF container
	f := elf.NewFile()
	f.Type  = elf.ET_EXEC
	f.Entry = textVAddr

	textSize := uint64(len(finalText))
	f.AddProgramHeader(elf.PT_LOAD, elf.PF_R|elf.PF_X, 0x1000, textVAddr, textSize, textSize, PageSize)

	dataSize := uint64(len(artifact.Data))
	if dataSize > 0 {
		dataFileOff := uint64(0x1000) + textSize
		if dataFileOff%PageSize != 0 {
			dataFileOff += PageSize - (dataFileOff % PageSize)
		}
		f.AddProgramHeader(elf.PT_LOAD, elf.PF_R|elf.PF_W, dataFileOff, dataVAddr, dataSize, dataSize, PageSize)
	}

	t := f.AddSection(".text", elf.SHT_PROGBITS, elf.SHF_ALLOC|elf.SHF_EXECINSTR, finalText)
	t.Addr      = textVAddr
	t.Addralign = PageSize

	d := f.AddSection(".data", elf.SHT_PROGBITS, elf.SHF_ALLOC|elf.SHF_WRITE, artifact.Data)
	d.Addr      = dataVAddr
	d.Addralign = PageSize

	buf := new(bytes.Buffer)
	if err := f.WriteTo(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}