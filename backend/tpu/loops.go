// loops.go - Loop detection and while loop emission
package tpu

import (
	"fmt"
	"strings"

	"github.com/arc-language/arc-lang/builder/ir"
)

// LoopInfo contains analyzed loop structure
type LoopInfo struct {
	Header    *ir.BasicBlock
	Body      []*ir.BasicBlock
	Exit      *ir.BasicBlock
	Phis      []*ir.PhiInst
	Condition ir.Instruction
	Induction *ir.PhiInst
}

func (g *Generator) detectLoops(fn *ir.Function) map[*ir.BasicBlock]*LoopInfo {
	loops := make(map[*ir.BasicBlock]*LoopInfo)

	for _, block := range fn.Blocks {
		// Check for back edges
		for _, succ := range block.Successors {
			for _, pred := range succ.Predecessors {
				if pred == block && g.dominates(succ, block, fn) {
					// Found a loop header
					info := &LoopInfo{Header: succ}

					// Collect phis
					for _, inst := range succ.Instructions {
						if phi, ok := inst.(*ir.PhiInst); ok {
							info.Phis = append(info.Phis, phi)
							// Heuristic: first phi is often induction variable
							if info.Induction == nil {
								info.Induction = phi
							}
						}
					}

					// Find condition and exit
					if term := succ.Terminator(); term != nil {
						if cond, ok := term.(*ir.CondBrInst); ok {
							info.Condition = cond
							// Determine exit block
							for _, s := range succ.Successors {
								if !g.isInLoop(s, succ, block) {
									info.Exit = s
								}
							}
						}
					}

					// Collect body blocks
					info.Body = g.collectLoopBody(succ, block)

					loops[succ] = info
				}
			}
		}
	}
	return loops
}

func (g *Generator) dominates(a, b *ir.BasicBlock, fn *ir.Function) bool {
	// Simplified: check if a appears before b in block order
	aIdx, bIdx := -1, -1
	for i, blk := range fn.Blocks {
		if blk == a {
			aIdx = i
		}
		if blk == b {
			bIdx = i
		}
	}
	return aIdx <= bIdx
}

func (g *Generator) isInLoop(block, header, latch *ir.BasicBlock) bool {
	if block == header {
		return true
	}
	visited := make(map[*ir.BasicBlock]bool)
	return g.reachesWithout(block, latch, header, visited)
}

func (g *Generator) reachesWithout(from, to, without *ir.BasicBlock, visited map[*ir.BasicBlock]bool) bool {
	if from == to {
		return true
	}
	if from == without || visited[from] {
		return false
	}
	visited[from] = true
	for _, succ := range from.Successors {
		if g.reachesWithout(succ, to, without, visited) {
			return true
		}
	}
	return false
}

func (g *Generator) collectLoopBody(header, latch *ir.BasicBlock) []*ir.BasicBlock {
	var body []*ir.BasicBlock
	visited := make(map[*ir.BasicBlock]bool)

	var collect func(*ir.BasicBlock)
	collect = func(b *ir.BasicBlock) {
		if visited[b] || b == header {
			return
		}
		visited[b] = true
		body = append(body, b)
		for _, pred := range b.Predecessors {
			collect(pred)
		}
	}

	collect(latch)
	return body
}

func (g *Generator) emitWhileLoop(header *ir.BasicBlock, info *LoopInfo) error {
	g.blockVisited[header] = true
	for _, b := range info.Body {
		g.blockVisited[b] = true
	}

	g.emit("// --- While Loop: %s ---", header.Name())

	// Build state tuple type
	var stateTypes []string
	var initialValues []string

	for _, phi := range info.Phis {
		stateTypes = append(stateTypes, g.toHloType(phi.Type()))

		// Find initial value (from outside loop)
		for _, inc := range phi.Incoming {
			isFromOutside := true
			for _, bodyBlock := range info.Body {
				if inc.Block == bodyBlock || inc.Block == header {
					isFromOutside = false
					break
				}
			}
			if isFromOutside {
				initialValues = append(initialValues, g.getOperand(inc.Value))
				break
			}
		}
	}

	tupleType := fmt.Sprintf("(%s)", strings.Join(stateTypes, ", "))

	// Create initial tuple
	initTuple := g.nextName("loop_init")
	g.emit("%s = %s tuple(%s)", initTuple, tupleType, strings.Join(initialValues, ", "))

	// Generate condition and body computation names
	condName := fmt.Sprintf("%s_condition", sanitizeName(header.Name()))
	bodyName := fmt.Sprintf("%s_body", sanitizeName(header.Name()))

	// Emit while instruction
	loopResult := g.nextName("loop_result")
	g.emit("%s = %s while(%s), condition=%%%s, body=%%%s",
		loopResult, tupleType, initTuple, condName, bodyName)

	// Extract results
	for i, phi := range info.Phis {
		resultName := g.nextName("loop_val")
		g.emit("%s = %s get-tuple-element(%s), index=%d", resultName, stateTypes[i], loopResult, i)
		g.valMap[phi] = resultName
	}

	// Generate condition sub-computation
	g.emitConditionComputation(condName, tupleType, stateTypes, info)

	// Generate body sub-computation
	g.emitBodyComputation(bodyName, tupleType, stateTypes, info)

	return nil
}

func (g *Generator) emitConditionComputation(name, tupleType string, stateTypes []string, info *LoopInfo) {
	g.subPrintf("\n%s {\n", name)
	g.subPrintf("  %%param = %s parameter(0)\n", tupleType)

	// Extract induction variable (usually first phi)
	if len(stateTypes) > 0 {
		g.subPrintf("  %%idx = %s get-tuple-element(%%param), index=0\n", stateTypes[0])
	}

	// Find the comparison in header
	var cmpOp string = "LT"
	var limitVal string = "s32[] constant(1024)" // Default

	for _, inst := range info.Header.Instructions {
		if icmp, ok := inst.(*ir.ICmpInst); ok {
			cmpOp = g.icmpToHloDirection(icmp.Predicate)
			if len(icmp.Operands()) > 1 {
				limitVal = g.getOperandForSub(icmp.Operands()[1])
			}
			break
		}
	}

	g.subPrintf("  %%limit = %s\n", limitVal)
	g.subPrintf("  ROOT %%cond = pred[] compare(%%idx, %%limit), direction=%s\n", cmpOp)
	g.subPrintf("}\n")
}

func (g *Generator) emitBodyComputation(name, tupleType string, stateTypes []string, info *LoopInfo) {
	g.subPrintf("\n%s {\n", name)
	g.subPrintf("  %%param = %s parameter(0)\n", tupleType)

	// Create local value map for body
	localValMap := make(map[ir.Value]string)

	// Extract state elements
	for i, phi := range info.Phis {
		elemName := fmt.Sprintf("%%state%d", i)
		g.subPrintf("  %s = %s get-tuple-element(%%param), index=%d\n", elemName, stateTypes[i], i)
		localValMap[phi] = elemName
	}

	// Emit body instructions
	for _, block := range info.Body {
		for _, inst := range block.Instructions {
			if inst.IsTerminator() {
				continue
			}
			g.emitInstructionToSub(inst, localValMap, stateTypes)
		}
	}

	// Build next state tuple
	var nextValues []string
	for _, phi := range info.Phis {
		// Find value from back edge
		var nextVal string
		for _, inc := range phi.Incoming {
			for _, bodyBlock := range info.Body {
				if inc.Block == bodyBlock || inc.Block == info.Header {
					if name, ok := localValMap[inc.Value]; ok {
						nextVal = name
					} else {
						nextVal = g.getOperandLocal(inc.Value, localValMap)
					}
					break
				}
			}
		}
		if nextVal == "" {
			nextVal = localValMap[phi]
		}
		nextValues = append(nextValues, nextVal)
	}

	g.subPrintf("  ROOT %%next = %s tuple(%s)\n", tupleType, strings.Join(nextValues, ", "))
	g.subPrintf("}\n")
}