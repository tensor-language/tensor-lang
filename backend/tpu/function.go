// function.go - Function emission
package tpu

import (
	"fmt"
	"strings"

	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

func (g *Generator) emitFunction(fn *ir.Function) error {
	g.currentFunc = fn
	g.valMap = make(map[ir.Value]string)
	g.memMap = make(map[ir.Value]string)
	g.shapeMap = make(map[ir.Value]Shape)
	g.blockVisited = make(map[*ir.BasicBlock]bool)
	g.nextID = 0

	// Build parameter signature
	var paramSigs []string
	for i, arg := range fn.Arguments {
		hloType := g.toHloType(arg.Type())
		paramName := fmt.Sprintf("p%d.%s", i, sanitizeName(arg.Name()))
		paramSigs = append(paramSigs, fmt.Sprintf("%s: %s", paramName, hloType))
	}

	retType := g.toHloType(fn.FuncType.ReturnType)

	g.printf("ENTRY %s(%s) -> %s {\n", sanitizeName(fn.Name()), strings.Join(paramSigs, ", "), retType)
	g.indentLevel++

	// Emit parameters
	for i, arg := range fn.Arguments {
		hloType := g.toHloType(arg.Type())
		name := g.nextName("param")
		g.valMap[arg] = name
		g.memMap[arg] = name
		g.emit("%s = %s parameter(%d)", name, hloType, i)

		if ptr, ok := arg.Type().(*types.PointerType); ok {
			g.shapeMap[arg] = g.inferShape(ptr.ElementType)
		}
	}

	// Analyze control flow for loops
	loopHeaders := g.detectLoops(fn)

	// Emit blocks
	for _, block := range fn.Blocks {
		if g.blockVisited[block] {
			continue
		}

		if info, isLoop := loopHeaders[block]; isLoop {
			if err := g.emitWhileLoop(block, info); err != nil {
				return err
			}
		} else {
			if err := g.emitBlock(block); err != nil {
				return err
			}
		}
	}

	g.indentLevel--
	g.printf("}\n\n")
	return nil
}

func (g *Generator) emitBlock(block *ir.BasicBlock) error {
	if g.blockVisited[block] {
		return nil
	}
	g.blockVisited[block] = true

	g.emit("// Block: %s", block.Name())

	for _, inst := range block.Instructions {
		if err := g.emitInstruction(inst); err != nil {
			g.errors = append(g.errors, err)
		}
	}
	return nil
}