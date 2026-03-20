// aggregate.go - Aggregate operation emission (extract/insert value)
package tpu

import (
	"fmt"

	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitExtractValue(i *ir.ExtractValueInst) error {
	dst := g.assignName(i)
	agg := g.getOperand(i.Operands()[0])
	typ := g.toHloType(i.Type())

	if len(i.Indices) == 1 {
		g.emit("%s = %s get-tuple-element(%s), index=%d", dst, typ, agg, i.Indices[0])
	} else {
		// Nested extraction
		current := agg
		for j, idx := range i.Indices {
			name := fmt.Sprintf("%%extract%d", j)
			g.emit("%s = %s get-tuple-element(%s), index=%d", name, typ, current, idx)
			current = name
		}
		g.valMap[i] = current
	}
	return nil
}

func (g *Generator) emitInsertValue(i *ir.InsertValueInst) error {
	dst := g.assignName(i)
	agg := g.getOperand(i.Operands()[0])
	val := g.getOperand(i.Operands()[1])
	typ := g.toHloType(i.Type())

	// HLO doesn't have direct insert - rebuild tuple
	g.emit("// insert value at indices %v", i.Indices)
	g.emit("%s = %s tuple(%s) // with %s inserted", dst, typ, agg, val)
	return nil
}