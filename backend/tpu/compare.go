// compare.go - Comparison operation emission
package tpu

import (
	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitICmp(i *ir.ICmpInst) error {
	dst := g.assignName(i)
	lhs := g.getOperand(i.Operands()[0])
	rhs := g.getOperand(i.Operands()[1])
	dir := g.icmpToHloDirection(i.Predicate)

	g.emit("%s = pred[] compare(%s, %s), direction=%s", dst, lhs, rhs, dir)
	return nil
}

func (g *Generator) emitFCmp(i *ir.FCmpInst) error {
	dst := g.assignName(i)
	lhs := g.getOperand(i.Operands()[0])
	rhs := g.getOperand(i.Operands()[1])
	dir := g.fcmpToHloDirection(i.Predicate)

	g.emit("%s = pred[] compare(%s, %s), direction=%s", dst, lhs, rhs, dir)
	return nil
}

func (g *Generator) icmpToHloDirection(pred ir.ICmpPredicate) string {
	switch pred {
	case ir.ICmpEQ:
		return "EQ"
	case ir.ICmpNE:
		return "NE"
	case ir.ICmpUGT, ir.ICmpSGT:
		return "GT"
	case ir.ICmpUGE, ir.ICmpSGE:
		return "GE"
	case ir.ICmpULT, ir.ICmpSLT:
		return "LT"
	case ir.ICmpULE, ir.ICmpSLE:
		return "LE"
	default:
		return "EQ"
	}
}

func (g *Generator) fcmpToHloDirection(pred ir.FCmpPredicate) string {
	switch pred {
	case ir.FCmpOEQ, ir.FCmpUEQ:
		return "EQ"
	case ir.FCmpONE, ir.FCmpUNE:
		return "NE"
	case ir.FCmpOGT, ir.FCmpUGT:
		return "GT"
	case ir.FCmpOGE, ir.FCmpUGE:
		return "GE"
	case ir.FCmpOLT, ir.FCmpULT:
		return "LT"
	case ir.FCmpOLE, ir.FCmpULE:
		return "LE"
	default:
		return "EQ"
	}
}