// binary.go - Binary operation emission
package tpu

import (
	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitBinary(i *ir.BinaryInst) error {
	dst := g.assignName(i)
	lhs := g.getOperand(i.Operands()[0])
	rhs := g.getOperand(i.Operands()[1])
	typ := g.toHloType(i.Type())
	op := g.binaryOpToHlo(i.Opcode())

	// Handle overflow flags for integer ops
	if i.NoSignedWrap || i.NoUnsignedWrap {
		g.emit("// nsw/nuw flags (informational)")
	}

	g.emit("%s = %s %s(%s, %s)", dst, typ, op, lhs, rhs)
	return nil
}

func (g *Generator) binaryOpToHlo(op ir.Opcode) string {
	switch op {
	case ir.OpAdd, ir.OpFAdd:
		return "add"
	case ir.OpSub, ir.OpFSub:
		return "subtract"
	case ir.OpMul, ir.OpFMul:
		return "multiply"
	case ir.OpUDiv, ir.OpSDiv, ir.OpFDiv:
		return "divide"
	case ir.OpURem, ir.OpSRem, ir.OpFRem:
		return "remainder"
	case ir.OpShl:
		return "shift-left"
	case ir.OpLShr:
		return "shift-right-logical"
	case ir.OpAShr:
		return "shift-right-arithmetic"
	case ir.OpAnd:
		return "and"
	case ir.OpOr:
		return "or"
	case ir.OpXor:
		return "xor"
	default:
		return "unknown-binary-op"
	}
}