// helpers.go - Helper functions for emission
package tpu

import (
	"fmt"
	"strings"

	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emit(format string, args ...interface{}) {
	indent := strings.Repeat("  ", g.indentLevel)
	g.printf("%s%s\n", indent, fmt.Sprintf(format, args...))
}

func (g *Generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) subPrintf(format string, args ...interface{}) {
	fmt.Fprintf(g.subComputations, format, args...)
}

func (g *Generator) nextName(prefix string) string {
	name := fmt.Sprintf("%%%s.%d", prefix, g.nextID)
	g.nextID++
	return name
}

func (g *Generator) assignName(v ir.Value) string {
	name := g.nextName("v")
	g.valMap[v] = name
	return name
}

func (g *Generator) getOperand(v ir.Value) string {
	if v == nil {
		return "/*null*/"
	}

	switch c := v.(type) {
	case *ir.ConstantInt:
		return fmt.Sprintf("s32[] constant(%d)", c.Value)
	case *ir.ConstantFloat:
		return fmt.Sprintf("f32[] constant(%f)", c.Value)
	case *ir.ConstantNull:
		return "/*null*/ s32[] constant(0)"
	case *ir.ConstantZero:
		return g.toHloType(c.Type()) + " constant(0)"
	case *ir.ConstantUndef:
		return "/*undef*/ " + g.toHloType(c.Type()) + " constant(0)"
	case *ir.Global:
		if name, ok := g.valMap[v]; ok {
			return name
		}
		return fmt.Sprintf("@%s", c.Name())
	}

	if name, ok := g.valMap[v]; ok {
		return name
	}
	if name, ok := g.memMap[v]; ok {
		return name
	}

	return fmt.Sprintf("/*unknown:%s*/", v.Name())
}

func (g *Generator) getOperandLocal(v ir.Value, localMap map[ir.Value]string) string {
	if name, ok := localMap[v]; ok {
		return name
	}
	return g.getOperand(v)
}

func (g *Generator) getOperandForSub(v ir.Value) string {
	if c, ok := v.(*ir.ConstantInt); ok {
		return fmt.Sprintf("s32[] constant(%d)", c.Value)
	}
	if c, ok := v.(*ir.ConstantFloat); ok {
		return fmt.Sprintf("f32[] constant(%f)", c.Value)
	}
	return g.getOperand(v)
}

func (g *Generator) constantValue(c ir.Constant) string {
	switch v := c.(type) {
	case *ir.ConstantInt:
		return fmt.Sprintf("%d", v.Value)
	case *ir.ConstantFloat:
		return fmt.Sprintf("%f", v.Value)
	default:
		return "0"
	}
}

func sanitizeName(name string) string {
	// HLO names can't have certain characters
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "<", "_")
	name = strings.ReplaceAll(name, ">", "_")
	return name
}