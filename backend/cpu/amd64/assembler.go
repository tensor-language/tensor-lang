package amd64

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// RelocationType defines x86-64 relocation types
type RelocationType int

const (
	RelocPC32  RelocationType = 2 // R_X86_64_PC32
	RelocPLT32 RelocationType = 4 // R_X86_64_PLT32
)

// RelocationRecord records where patching is needed
type RelocationRecord struct {
	Offset int
	Symbol string
	Type   RelocationType
	Addend int64
}

// Assembler handles low-level instruction encoding
type Assembler struct {
	buf    *bytes.Buffer
	Relocs []RelocationRecord
}

func NewAssembler() *Assembler {
	return &Assembler{
		buf: new(bytes.Buffer),
	}
}

func (a *Assembler) Bytes() []byte {
	return a.buf.Bytes()
}

func (a *Assembler) Len() int {
	return a.buf.Len()
}

// --- Labels & Patching ---

// NewLabel returns the current position as a label target
func (a *Assembler) NewLabel() int {
	return a.Len()
}

// Label marks the current position (no-op in this simple assembler, used for readability)
func (a *Assembler) Label(l int) {
	// In a single-pass assembler, we don't need to do anything here 
	// because NewLabel() already captured the offset.
}

// PatchJump calculates the relative offset for a previously emitted jump
// jumpDataStart is the offset where the 32-bit immediate begins.
func (a *Assembler) PatchJump(jumpDataStart int) {
	target := a.Len()
	// The jump offset is relative to the *next* instruction.
	// 32-bit jump immediate is always 4 bytes.
	// So NextIP = jumpDataStart + 4
	rel := int32(target - (jumpDataStart + 4))
	a.PatchInt32(jumpDataStart, rel)
}

// PatchInt32 writes a value at a specific offset
func (a *Assembler) PatchInt32(offset int, value int32) {
	old := a.buf.Bytes()
	binary.LittleEndian.PutUint32(old[offset:], uint32(value))
}

// --- Encoding Primitives ---

func (a *Assembler) emitByte(b byte) {
	a.buf.WriteByte(b)
}

func (a *Assembler) emitInt32(v int32) {
	binary.Write(a.buf, binary.LittleEndian, v)
}

func (a *Assembler) emitInt64(v int64) {
	binary.Write(a.buf, binary.LittleEndian, v)
}

func (a *Assembler) encodeRex(w bool, r, x, b Register) {
	var rex byte = 0x40
	needsRex := false

	if w {
		rex |= 0x08
		needsRex = true
	}
	if r >= 8 {
		rex |= 0x04
		needsRex = true
	}
	if x >= 8 {
		rex |= 0x02
		needsRex = true
	}
	if b >= 8 {
		rex |= 0x01
		needsRex = true
	}

	if needsRex {
		a.emitByte(rex)
	}
}

func (a *Assembler) encodeModRM(reg Register, rm Operand) {
	regCode := byte(reg) & 7

	switch op := rm.(type) {
	case RegOp:
		rmCode := byte(op) & 7
		a.emitByte(0xC0 | (regCode << 3) | rmCode)

	case MemOp:
		baseCode := byte(op.Base) & 7
		
		if op.Index == NoReg {
			if op.Base == RSP || op.Base == R12 {
				a.emitByte(0x80 | (regCode << 3) | 0x04)
				a.emitByte(0x24) 
			} else {
				a.emitByte(0x80 | (regCode << 3) | baseCode)
			}
			a.emitInt32(op.Disp)
		} else {
			panic("Complex SIB addressing not implemented in this simplified assembler")
		}
	
	default:
		panic(fmt.Sprintf("unsupported operand for ModRM: %T", rm))
	}
}

// --- Instructions ---

func (a *Assembler) Mov(dst, src Operand, size int) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(size == 64, Register(s), NoReg, Register(d))
			if size == 8 {
				a.emitByte(0x88)
			} else if size == 16 {
				a.emitByte(0x66)
				a.emitByte(0x89)
			} else {
				a.emitByte(0x89)
			}
			a.encodeModRM(Register(s), d) 
			return
		}
	}

	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(MemOp); ok {
			a.encodeRex(size == 64, Register(d), NoReg, s.Base)
			if size == 8 {
				a.emitByte(0x0F); a.emitByte(0xB6) 
			} else if size == 16 {
				a.emitByte(0x0F); a.emitByte(0xB7)
			} else {
				a.emitByte(0x8B) 
			}
			a.encodeModRM(Register(d), s)
			return
		}
	}

	if d, ok := dst.(MemOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(size == 64, Register(s), NoReg, d.Base)
			if size == 8 {
				a.emitByte(0x88)
			} else if size == 16 {
				a.emitByte(0x66); a.emitByte(0x89)
			} else {
				a.emitByte(0x89)
			}
			a.encodeModRM(Register(s), d)
			return
		}
	}

	if d, ok := dst.(RegOp); ok {
		if imm, ok := src.(ImmOp); ok {
			if imm == 0 && size == 64 {
				a.Xor(d, d)
				return
			}
			
			reg := Register(d)
			a.encodeRex(true, NoReg, NoReg, reg)
			a.emitByte(0xB8 | (byte(reg) & 7))
			a.emitInt64(int64(imm))
			return
		}
	}

	if d, ok := dst.(MemOp); ok {
		if imm, ok := src.(ImmOp); ok {
			if size == 8 {
				a.encodeRex(false, 0, NoReg, d.Base)
				a.emitByte(0xC6)
				a.encodeModRM(0, d)
				a.emitByte(byte(imm))
				return
			}
			
			if size == 16 {
				a.emitByte(0x66)
				a.encodeRex(false, 0, NoReg, d.Base)
				a.emitByte(0xC7)
				a.encodeModRM(0, d)
				v := uint16(imm)
				a.emitByte(byte(v))
				a.emitByte(byte(v >> 8))
				return
			}
			
			if size == 32 {
				a.encodeRex(false, 0, NoReg, d.Base)
				a.emitByte(0xC7)
				a.encodeModRM(0, d)
				a.emitInt32(int32(imm))
				return
			}
			
			if size == 64 {
				if int64(int32(imm)) != int64(imm) {
					panic(fmt.Sprintf("MOV Mem64, Imm requires 32-bit signed immediate, got %d", imm))
				}
				a.encodeRex(true, 0, NoReg, d.Base)
				a.emitByte(0xC7)
				a.encodeModRM(0, d)
				a.emitInt32(int32(imm))
				return
			}
		}
	}

	panic(fmt.Sprintf("Unsupported MOV combination: %T -> %T", src, dst))
}

func (a *Assembler) Add(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x01)
			a.encodeModRM(Register(s), d)
		} else if imm, ok := src.(ImmOp); ok {
			a.encodeRex(true, 0, NoReg, Register(d))
			if imm >= -128 && imm <= 127 {
				a.emitByte(0x83)
				a.encodeModRM(0, d)
				a.emitByte(byte(imm))
			} else {
				a.emitByte(0x81)
				a.encodeModRM(0, d)
				a.emitInt32(int32(imm))
			}
		}
	}
}

func (a *Assembler) Sub(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x29)
			a.encodeModRM(Register(s), d)
		} else if imm, ok := src.(ImmOp); ok {
			a.encodeRex(true, 0, NoReg, Register(d))
			if imm >= -128 && imm <= 127 {
				a.emitByte(0x83)
				a.encodeModRM(5, d) // /5
				a.emitByte(byte(imm))
			} else {
				a.emitByte(0x81)
				a.encodeModRM(5, d)
				a.emitInt32(int32(imm))
			}
		}
	}
}

func (a *Assembler) Imul(dst, src Register) {
	a.encodeRex(true, dst, NoReg, src)
	a.emitByte(0x0F)
	a.emitByte(0xAF)
	a.encodeModRM(dst, RegOp(src))
}

func (a *Assembler) Div(src Register, signed bool) {
	a.encodeRex(true, 0, NoReg, src)
	a.emitByte(0xF7)
	subOp := 6 // DIV
	if signed {
		subOp = 7 // IDIV
	}
	a.encodeModRM(Register(subOp), RegOp(src))
}

func (a *Assembler) Xor(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x31)
			a.encodeModRM(Register(s), d)
		}
	}
}

func (a *Assembler) Push(src Register) {
	if src >= 8 {
		a.emitByte(0x41) // REX.B
	}
	a.emitByte(0x50 | (byte(src) & 7))
}

func (a *Assembler) Pop(dst Register) {
	if dst >= 8 {
		a.emitByte(0x41) // REX.B
	}
	a.emitByte(0x58 | (byte(dst) & 7))
}

func (a *Assembler) Ret() {
	a.emitByte(0xC3)
}

func (a *Assembler) Syscall() {
	a.emitByte(0x0F)
	a.emitByte(0x05)
}

func (a *Assembler) Lea(dst Register, mem MemOp) {
	a.encodeRex(true, dst, NoReg, mem.Base)
	a.emitByte(0x8D)
	a.encodeModRM(dst, mem)
}

func (a *Assembler) CallRelative(symbol string) {
	a.emitByte(0xE8)
	a.Relocs = append(a.Relocs, RelocationRecord{
		Offset: a.Len(),
		Symbol: symbol,
		Type:   RelocPLT32,
		Addend: -4,
	})
	a.emitInt32(0) 
}

// JmpRelSymbol emits a JMP (E9) to a named symbol
func (a *Assembler) JmpRelSymbol(symbol string) {
	a.emitByte(0xE9)
	a.Relocs = append(a.Relocs, RelocationRecord{
		Offset: a.Len(),
		Symbol: symbol,
		Type:   RelocPLT32,
		Addend: -4,
	})
	a.emitInt32(0) 
}

func (a *Assembler) CallReg(reg Register) {
	a.encodeRex(false, 0, NoReg, reg)
	a.emitByte(0xFF)
	a.encodeModRM(2, RegOp(reg)) // /2
}

// JmpRel emits a E9 jump. 
// If val is a Label (previous offset), we calculate relative jump backwards.
// If val is 0, we emit 0 and return offset for patching (forward jump).
func (a *Assembler) JmpRel(val int) int {
	a.emitByte(0xE9)
	pos := a.Len()
	
	offset := int32(0)
	
	// Heuristic: If val is a valid backward position, it's a loop
	if val > 0 && val < pos {
		// Target = val
		// NextIP = pos + 4
		offset = int32(val - (pos + 4))
	}
	
	a.emitInt32(offset)
	return pos
}

func (a *Assembler) JccRel(cond CondCode, offset int32) int {
	a.emitByte(0x0F)
	a.emitByte(byte(cond))
	pos := a.Len()
	a.emitInt32(offset)
	return pos
}

func (a *Assembler) Cqo() {
	a.emitByte(0x48)
	a.emitByte(0x99)
}

func (a *Assembler) And(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x21)
			a.encodeModRM(Register(s), d)
		} else if imm, ok := src.(ImmOp); ok {
			// AND r/m64, imm8/32
			a.encodeRex(true, 0, NoReg, Register(d))
			if imm >= -128 && imm <= 127 {
				a.emitByte(0x83)
				a.encodeModRM(4, d) // /4
				a.emitByte(byte(imm))
			} else {
				a.emitByte(0x81)
				a.encodeModRM(4, d) // /4
				a.emitInt32(int32(imm))
			}
		}
	}
}

func (a *Assembler) Or(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x09)
			a.encodeModRM(Register(s), d)
		}
	}
}

func (a *Assembler) Shl(dst, src Register) {
	a.encodeRex(true, 0, NoReg, dst)
	a.emitByte(0xD3)
	a.encodeModRM(4, RegOp(dst))
}

func (a *Assembler) Shr(dst, src Register) {
	a.encodeRex(true, 0, NoReg, dst)
	a.emitByte(0xD3)
	a.encodeModRM(5, RegOp(dst))
}

func (a *Assembler) Sar(dst, src Register) {
	a.encodeRex(true, 0, NoReg, dst)
	a.emitByte(0xD3)
	a.encodeModRM(7, RegOp(dst))
}

func (a *Assembler) Cmp(dst, src Operand) {
	if d, ok := dst.(RegOp); ok {
		if s, ok := src.(RegOp); ok {
			a.encodeRex(true, Register(s), NoReg, Register(d))
			a.emitByte(0x39)
			a.encodeModRM(Register(s), d)
		}
	}
}

func (a *Assembler) Test(dst, src Register) {
	a.encodeRex(true, src, NoReg, dst)
	a.emitByte(0x85)
	a.encodeModRM(src, RegOp(dst))
}

func (a *Assembler) MovZX(dst Register, src Operand, srcSize int) {
	if srcSize == 8 {
		if m, ok := src.(MemOp); ok {
			a.encodeRex(true, dst, NoReg, m.Base)
			a.emitByte(0x0F); a.emitByte(0xB6)
			a.encodeModRM(dst, src)
		} else if s, ok := src.(RegOp); ok {
			a.encodeRex(true, dst, NoReg, Register(s))
			a.emitByte(0x0F); a.emitByte(0xB6)
			a.encodeModRM(dst, src)
		}
	}
}

func (a *Assembler) Movsxd(dst, src Register) {
	a.encodeRex(true, dst, NoReg, src)
	a.emitByte(0x63)
	a.encodeModRM(dst, RegOp(src))
}

func (a *Assembler) Movsx(dst Register, src Operand, srcSize int) {
	if srcSize == 8 {
		if m, ok := src.(MemOp); ok {
			a.encodeRex(true, dst, NoReg, m.Base)
			a.emitByte(0x0F); a.emitByte(0xBE)
			a.encodeModRM(dst, src)
		} else if s, ok := src.(RegOp); ok {
			a.encodeRex(true, dst, NoReg, Register(s))
			a.emitByte(0x0F); a.emitByte(0xBE)
			a.encodeModRM(dst, src)
		}
	}
}

func (a *Assembler) Setcc(cc CondCode, dst Register) {
	a.emitByte(0x0F)
	a.emitByte(byte(cc) + 0x10) 
	a.encodeModRM(0, RegOp(dst))
}

func (a *Assembler) ImulImm(dst Register, imm int32) {
	a.encodeRex(true, dst, NoReg, dst)
	a.emitByte(0x69)
	a.encodeModRM(dst, RegOp(dst))
	a.emitInt32(imm)
}

func (a *Assembler) LeaRel(dst Register, symbol string) {
	a.encodeRex(true, dst, NoReg, 0)
	a.emitByte(0x8D)
	reg := byte(dst) & 7
	a.emitByte(0x05 | (reg << 3))
	
	a.Relocs = append(a.Relocs, RelocationRecord{
		Offset: a.Len(),
		Symbol: symbol,
		Type:   RelocPC32,
		Addend: -4,
	})
	a.emitInt32(0)
}

func (a *Assembler) Cvttss2si(dst Register, src MemOp) {
	a.emitByte(0xF3)
	a.encodeRex(false, dst, NoReg, src.Base) 
	a.emitByte(0x0F)
	a.emitByte(0x2C)
	a.encodeModRM(dst, src)
}

func (a *Assembler) Cvtsi2ss(dst Register, src MemOp) {
	a.emitByte(0xF3)
	a.encodeRex(false, dst, NoReg, src.Base)
	a.emitByte(0x0F)
	a.emitByte(0x2A)
	a.encodeModRM(dst, src)
}

func (a *Assembler) Movss(dst MemOp, src Register) {
	// F3 0F 11 /r (MOVSS m32, xmm1)
	a.emitByte(0xF3)
	a.encodeRex(false, src, NoReg, dst.Base)
	a.emitByte(0x0F)
	a.emitByte(0x11)
	a.encodeModRM(src, dst)
}

func (a *Assembler) Movsd(dst MemOp, src Register) {
	// F2 0F 11 /r (MOVSD m64, xmm1)
	a.emitByte(0xF2)
	a.encodeRex(false, src, NoReg, dst.Base)
	a.emitByte(0x0F)
	a.emitByte(0x11)
	a.encodeModRM(src, dst)
}

// MovssLoad loads a 32-bit float from memory to XMM
func (a *Assembler) MovssLoad(dst Register, src MemOp) {
	// F3 0F 10 /r (MOVSS xmm1, m32)
	a.emitByte(0xF3)
	a.encodeRex(false, dst, NoReg, src.Base)
	a.emitByte(0x0F)
	a.emitByte(0x10)
	a.encodeModRM(dst, src)
}

// MovsdLoad loads a 64-bit float from memory to XMM
func (a *Assembler) MovsdLoad(dst Register, src MemOp) {
	// F2 0F 10 /r (MOVSD xmm1, m64)
	a.emitByte(0xF2)
	a.encodeRex(false, dst, NoReg, src.Base)
	a.emitByte(0x0F)
	a.emitByte(0x10)
	a.encodeModRM(dst, src)
}

func (a *Assembler) Movd(dst Register, src Register) {
	a.emitByte(0x66)
	a.encodeRex(false, dst, NoReg, src)
	a.emitByte(0x0F)
	a.emitByte(0x6E)
	a.encodeModRM(dst, RegOp(src))
}

func (a *Assembler) Movq(dst Register, src Register) {
	a.emitByte(0x66)
	a.encodeRex(true, dst, NoReg, src)
	a.emitByte(0x0F)
	a.emitByte(0x6E)
	a.encodeModRM(dst, RegOp(src))
}

func (a *Assembler) Cvtss2sd(dst Register, src Register) {
	a.emitByte(0xF3)
	a.encodeRex(false, dst, NoReg, src)
	a.emitByte(0x0F)
	a.emitByte(0x5A)
	a.encodeModRM(dst, RegOp(src))
}

func (a *Assembler) MovdXmmToGpr(dst Register, src Register) {
	// MOVD r32, xmm (Transfer xmm to r32)
	a.emitByte(0x66)
	a.encodeRex(false, src, NoReg, dst) // Src is XMM (reg), Dst is GPR (rm)
	a.emitByte(0x0F)
	a.emitByte(0x7E)
	a.encodeModRM(src, RegOp(dst))
}

func (a *Assembler) MovqXmmToGpr(dst Register, src Register) {
	// MOVQ r64, xmm (Transfer xmm to r64)
	a.emitByte(0x66)
	a.encodeRex(true, src, NoReg, dst)
	a.emitByte(0x0F)
	a.emitByte(0x7E)
	a.encodeModRM(src, RegOp(dst))
}

// Xchg performs an atomic exchange between a register and memory/register
// LOCK prefix is implicit when the operand is memory
func (a *Assembler) Xchg(dst MemOp, src Register, size int) {
	if size == 32 {
		a.encodeRex(false, src, NoReg, dst.Base)
		a.emitByte(0x87)
		a.encodeModRM(src, dst)
	} else if size == 64 {
		a.encodeRex(true, src, NoReg, dst.Base)
		a.emitByte(0x87)
		a.encodeModRM(src, dst)
	} else {
		panic("Xchg only supports 32 or 64 bit size")
	}
}