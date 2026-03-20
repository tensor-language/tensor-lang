// controlflow.go - Control flow instruction emission
package tpu

import (
	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitPhi(i *ir.PhiInst) error {
	// Phis are typically handled in loop/branch context
	// For linear code, we track the value
	if len(i.Incoming) > 0 {
		g.valMap[i] = g.getOperand(i.Incoming[0].Value)
	}
	return nil
}

func (g *Generator) emitSelect(i *ir.SelectInst) error {
	dst := g.assignName(i)
	cond := g.getOperand(i.Operands()[0])
	tVal := g.getOperand(i.Operands()[1])
	fVal := g.getOperand(i.Operands()[2])
	typ := g.toHloType(i.Type())

	g.emit("%s = %s select(%s, %s, %s)", dst, typ, cond, tVal, fVal)
	return nil
}

func (g *Generator) emitCondBr(i *ir.CondBrInst) error {
	// Conditional branches may create conditional computation
	cond := g.getOperand(i.Condition)
	g.emit("// conditional branch on %s -> %s / %s",
		cond, i.TrueBlock.Name(), i.FalseBlock.Name())

	// For simple if-then-else, we might emit a conditional
	// More complex control flow needs restructuring
	return nil
}

func (g *Generator) emitRet(i *ir.RetInst) error {
	if len(i.Operands()) > 0 && i.Operands()[0] != nil {
		val := g.getOperand(i.Operands()[0])
		typ := g.toHloType(i.Operands()[0].Type())
		rootName := g.nextName("return")
		g.emit("ROOT %s = %s copy(%s)", rootName, typ, val)
	} else {
		g.emit("ROOT %s = () tuple()", g.nextName("void_return"))
	}
	return nil
}

func (g *Generator) emitSwitch(i *ir.SwitchInst) error {
	// Switch can be lowered to nested selects or a lookup table
	cond := g.getOperand(i.Condition)
	g.emit("// switch on %s with %d cases", cond, len(i.Cases))

	// For now, emit as a series of comparisons and selects
	// A more efficient approach would use gather or case instruction
	return nil
}