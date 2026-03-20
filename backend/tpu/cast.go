// cast.go - Cast operation emission
package tpu

import (
	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitCast(i *ir.CastInst) error {
	dst := g.assignName(i)
	src := g.getOperand(i.Operands()[0])
	destType := g.toHloType(i.DestType)

	switch i.Opcode() {
	case ir.OpTrunc, ir.OpZExt, ir.OpSExt:
		g.emit("%s = %s convert(%s)", dst, destType, src)
	case ir.OpFPTrunc, ir.OpFPExt:
		g.emit("%s = %s convert(%s)", dst, destType, src)
	case ir.OpFPToUI, ir.OpFPToSI:
		g.emit("%s = %s convert(%s)", dst, destType, src)
	case ir.OpUIToFP, ir.OpSIToFP:
		g.emit("%s = %s convert(%s)", dst, destType, src)
	case ir.OpBitcast:
		g.emit("%s = %s bitcast-convert(%s)", dst, destType, src)
	case ir.OpPtrToInt, ir.OpIntToPtr:
		// TPU doesn't really have pointers - this is mostly for compatibility
		g.emit("%s = %s copy(%s)", dst, destType, src)
	default:
		g.emit("%s = %s convert(%s)", dst, destType, src)
	}
	return nil
}