// generator.go - Core generator struct and main entry point
package tpu

import (
	"bytes"
	"fmt"

	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

// Generator handles the emission of Google TPU HLO IR
type Generator struct {
	buf             *bytes.Buffer
	subComputations *bytes.Buffer

	// SSA value mapping
	valMap map[ir.Value]string
	memMap map[ir.Value]string

	// Shape inference
	shapeMap map[ir.Value]Shape

	nextID       int
	indentLevel  int
	currentFunc  *ir.Function
	loopDepth    int
	blockVisited map[*ir.BasicBlock]bool

	errors []error
}

// Generate compiles an IR module into HLO text.
func Generate(m *ir.Module) (string, error) {
	g := &Generator{
		buf:             new(bytes.Buffer),
		subComputations: new(bytes.Buffer),
		valMap:          make(map[ir.Value]string),
		memMap:          make(map[ir.Value]string),
		shapeMap:        make(map[ir.Value]Shape),
		blockVisited:    make(map[*ir.BasicBlock]bool),
	}

	g.emitHeader(m)

	// Emit global constants
	for _, global := range m.Globals {
		if err := g.emitGlobal(global); err != nil {
			g.errors = append(g.errors, err)
		}
	}

	// Emit TPU functions
	for _, fn := range m.Functions {
		if fn.CallConv == ir.CC_TPU {
			if err := g.emitFunction(fn); err != nil {
				g.errors = append(g.errors, err)
			}
		}
	}

	// Append sub-computations (while loop bodies, conditions, etc.)
	g.buf.Write(g.subComputations.Bytes())

	if len(g.errors) > 0 {
		return g.buf.String(), fmt.Errorf("HLO generation had %d errors: %v", len(g.errors), g.errors[0])
	}

	return g.buf.String(), nil
}

func (g *Generator) emitHeader(m *ir.Module) {
	g.printf("HloModule %s, is_scheduled=false\n\n", sanitizeName(m.Name))
}

func (g *Generator) emitGlobal(global *ir.Global) error {
	name := g.nextName("global")
	g.valMap[global] = name

	if global.Initializer != nil {
		switch c := global.Initializer.(type) {
		case *ir.ConstantArray:
			elemType := g.toHloElemType(c.Type())
			if arr, ok := c.Type().(*types.ArrayType); ok {
				g.printf("%s = %s[%d] constant({", name, elemType, arr.Length)
				for i, elem := range c.Elements {
					if i > 0 {
						g.printf(", ")
					}
					g.printf("%s", g.constantValue(elem))
				}
				g.printf("})\n\n")
			}
		case *ir.ConstantFloat:
			g.printf("%s = f32[] constant(%f)\n", name, c.Value)
		case *ir.ConstantInt:
			g.printf("%s = s32[] constant(%d)\n", name, c.Value)
		}
	}
	return nil
}