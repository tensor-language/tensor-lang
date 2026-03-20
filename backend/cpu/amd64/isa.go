package amd64

import "fmt"

// Register represents a machine register
type Register int

const (
	// General Purpose Registers
	RAX Register = 0
	RCX Register = 1
	RDX Register = 2
	RBX Register = 3
	RSP Register = 4
	RBP Register = 5
	RSI Register = 6
	RDI Register = 7
	R8  Register = 8
	R9  Register = 9
	R10 Register = 10
	R11 Register = 11
	R12 Register = 12
	R13 Register = 13
	R14 Register = 14
	R15 Register = 15

	// Pseudo-register to indicate "no register"
	NoReg Register = -1
)

// String returns the register name (for debugging)
func (r Register) String() string {
	names := []string{
		"RAX", "RCX", "RDX", "RBX", "RSP", "RBP", "RSI", "RDI",
		"R8", "R9", "R10", "R11", "R12", "R13", "R14", "R15",
	}
	if r >= 0 && int(r) < len(names) {
		return names[r]
	}
	return fmt.Sprintf("Reg(%d)", r)
}

// Operand is a marker interface for assembly operands
type Operand interface {
	isOperand()
}

// RegOp wraps a Register to satisfy Operand
type RegOp Register

func (r RegOp) isOperand() {}

// ImmOp represents an immediate value
type ImmOp int64

func (i ImmOp) isOperand() {}

// MemOp represents a memory reference: [Base + Index*Scale + Disp]
type MemOp struct {
	Base  Register
	Index Register
	Scale int // 1, 2, 4, 8
	Disp  int32
}

func (m MemOp) isOperand() {}

// NewMem creates a simple [Base + Disp] memory operand
func NewMem(base Register, disp int) MemOp {
	return MemOp{
		Base:  base,
		Index: NoReg,
		Scale: 1,
		Disp:  int32(disp),
	}
}

// Condition Codes for CMOVcc / Jcc / SETcc
type CondCode byte

const (
	CondEq  CondCode = 0x84 // Equal (ZF=1)
	CondNe  CondCode = 0x85 // Not Equal (ZF=0)
	CondLt  CondCode = 0x8C // Less (Signed)
	CondLe  CondCode = 0x8E // Less or Equal (Signed)
	CondGt  CondCode = 0x8F // Greater (Signed)
	CondGe  CondCode = 0x8D // Greater or Equal (Signed)
	CondBlo CondCode = 0x82 // Below (Unsigned)
	CondBle CondCode = 0x86 // Below or Equal (Unsigned)
	CondA   CondCode = 0x87 // Above (Unsigned)
	CondAe  CondCode = 0x83 // Above or Equal (Unsigned)
)