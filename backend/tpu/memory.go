// memory.go - Memory operation emission (alloca, load, store, GEP)
package tpu

import (
	"strings"

	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitAlloca(i *ir.AllocaInst) error {
	dst := g.assignName(i)
	shape := g.inferShape(i.AllocatedType)

	// In HLO, create a broadcast of zeros for allocation
	zeroConst := g.nextName("zero")
	g.emit("%s = %s constant(0)", zeroConst, shape.ElemType+"[]")
	g.emit("%s = %s broadcast(%s), dimensions={}", dst, shape.String(), zeroConst)
	g.shapeMap[i] = shape
	return nil
}

func (g *Generator) emitLoad(i *ir.LoadInst) error {
	dst := g.assignName(i)
	ptr := i.Operands()[0]
	typ := g.toHloType(i.Type())

	// Check if this is an indexed load (from GEP)
	if gep, ok := ptr.(*ir.GetElementPtrInst); ok {
		base := g.getOperand(gep.Operands()[0])
		indices := g.getGEPIndices(gep)

		if len(indices) == 1 {
			g.emit("%s = %s dynamic-slice(%s, %s), dynamic_slice_sizes={1}", dst, typ, base, indices[0])
		} else if len(indices) == 2 {
			g.emit("%s = %s dynamic-slice(%s, %s, %s), dynamic_slice_sizes={1,1}",
				dst, typ, base, indices[0], indices[1])
		} else {
			g.emit("%s = %s copy(%s)", dst, typ, base)
		}
	} else {
		ptrVal := g.getOperand(ptr)
		g.emit("%s = %s copy(%s)", dst, typ, ptrVal)
	}
	return nil
}

func (g *Generator) emitStore(i *ir.StoreInst) error {
	val := g.getOperand(i.Operands()[0])
	ptr := i.Operands()[1]

	if gep, ok := ptr.(*ir.GetElementPtrInst); ok {
		base := g.getOperand(gep.Operands()[0])
		indices := g.getGEPIndices(gep)
		resultName := g.nextName("updated")

		if len(indices) == 1 {
			g.emit("%s = %s dynamic-update-slice(%s, %s, %s)",
				resultName, g.toHloType(gep.Operands()[0].Type()), base, val, indices[0])
		} else {
			g.emit("%s = %s dynamic-update-slice(%s, %s, %s)",
				resultName, g.toHloType(gep.Operands()[0].Type()), base, val, strings.Join(indices, ", "))
		}
		g.memMap[gep.Operands()[0]] = resultName
	} else {
		ptrVal := g.getOperand(ptr)
		g.emit("// store %s -> %s", val, ptrVal)
	}
	return nil
}

func (g *Generator) emitGEP(i *ir.GetElementPtrInst) error {
	if len(i.Operands()) > 1 {
		base := g.getOperand(i.Operands()[0])
		indices := g.getGEPIndices(i)

		if len(indices) == 1 {
			g.valMap[i] = indices[0]
		} else {
			dst := g.assignName(i)
			g.emit("// GEP: %s[%s]", base, strings.Join(indices, ", "))
			g.emit("%s = s32[] constant(0)", dst)
		}
	}
	return nil
}

func (g *Generator) getGEPIndices(gep *ir.GetElementPtrInst) []string {
	var indices []string
	for _, op := range gep.Operands()[1:] {
		indices = append(indices, g.getOperand(op))
	}
	return indices
}

func (g *Generator) emitMemSet(i *ir.MemSetInst) error {
	dest := g.getOperand(i.Operands()[0])
	val := g.getOperand(i.Operands()[1])

	resultName := g.nextName("memset")
	g.emit("%s = f32[] broadcast(%s), dimensions={}", resultName, val)
	g.memMap[i.Operands()[0]] = resultName
	_ = dest
	return nil
}

func (g *Generator) emitMemCpy(i *ir.MemCpyInst) error {
	dest := g.getOperand(i.Operands()[0])
	src := g.getOperand(i.Operands()[1])

	resultName := g.nextName("memcpy")
	g.emit("%s = f32[] copy(%s)", resultName, src)
	g.memMap[i.Operands()[0]] = resultName
	_ = dest
	return nil
}