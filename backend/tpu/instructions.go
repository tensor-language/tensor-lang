// instructions.go - Main instruction dispatch and emission
package tpu

import (
	"fmt"

	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitInstruction(inst ir.Instruction) error {
	switch i := inst.(type) {
	// ========== Binary Operations ==========
	case *ir.BinaryInst:
		return g.emitBinary(i)

	// ========== Memory Operations ==========
	case *ir.AllocaInst:
		return g.emitAlloca(i)
	case *ir.LoadInst:
		return g.emitLoad(i)
	case *ir.StoreInst:
		return g.emitStore(i)
	case *ir.GetElementPtrInst:
		return g.emitGEP(i)

	// ========== Cast Operations ==========
	case *ir.CastInst:
		return g.emitCast(i)

	// ========== Comparison Operations ==========
	case *ir.ICmpInst:
		return g.emitICmp(i)
	case *ir.FCmpInst:
		return g.emitFCmp(i)

	// ========== Control Flow ==========
	case *ir.PhiInst:
		return g.emitPhi(i)
	case *ir.SelectInst:
		return g.emitSelect(i)
	case *ir.CondBrInst:
		return g.emitCondBr(i)
	case *ir.BrInst:
		// Unconditional branches handled by block structure
		return nil
	case *ir.RetInst:
		return g.emitRet(i)

	// ========== Call Operations ==========
	case *ir.CallInst:
		return g.emitCall(i)

	// ========== Aggregate Operations ==========
	case *ir.ExtractValueInst:
		return g.emitExtractValue(i)
	case *ir.InsertValueInst:
		return g.emitInsertValue(i)

	// ========== Intrinsics ==========
	case *ir.MemSetInst:
		return g.emitMemSet(i)
	case *ir.MemCpyInst:
		return g.emitMemCpy(i)

	case *ir.SwitchInst:
		return g.emitSwitch(i)
	case *ir.UnreachableInst:
		g.emit("// unreachable")
		return nil

	default:
		g.emit("// Unsupported: %T", inst)
		return nil
	}
}

func (g *Generator) emitInstructionToSub(inst ir.Instruction, localValMap map[ir.Value]string, stateTypes []string) {
	switch i := inst.(type) {
	case *ir.BinaryInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		lhs := g.getOperandLocal(i.Operands()[0], localValMap)
		rhs := g.getOperandLocal(i.Operands()[1], localValMap)
		op := g.binaryOpToHlo(i.Opcode())
		typ := g.toHloType(i.Type())
		g.subPrintf("  %s = %s %s(%s, %s)\n", dst, typ, op, lhs, rhs)

	case *ir.LoadInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		ptr := g.getOperandLocal(i.Operands()[0], localValMap)
		typ := g.toHloType(i.Type())
		g.subPrintf("  %s = %s dynamic-slice(%s, %%idx), dynamic_slice_sizes={1}\n", dst, typ, ptr)

	case *ir.StoreInst:
		val := g.getOperandLocal(i.Operands()[0], localValMap)
		ptr := g.getOperandLocal(i.Operands()[1], localValMap)
		g.subPrintf("  // store %s -> %s\n", val, ptr)

	case *ir.GetElementPtrInst:
		if len(i.Operands()) > 1 {
			idx := g.getOperandLocal(i.Operands()[len(i.Operands())-1], localValMap)
			localValMap[inst] = idx
		}

	case *ir.CastInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		src := g.getOperandLocal(i.Operands()[0], localValMap)
		typ := g.toHloType(i.Type())
		g.subPrintf("  %s = %s convert(%s)\n", dst, typ, src)

	case *ir.SelectInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		cond := g.getOperandLocal(i.Operands()[0], localValMap)
		tVal := g.getOperandLocal(i.Operands()[1], localValMap)
		fVal := g.getOperandLocal(i.Operands()[2], localValMap)
		typ := g.toHloType(i.Type())
		g.subPrintf("  %s = %s select(%s, %s, %s)\n", dst, typ, cond, tVal, fVal)

	case *ir.ICmpInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		lhs := g.getOperandLocal(i.Operands()[0], localValMap)
		rhs := g.getOperandLocal(i.Operands()[1], localValMap)
		dir := g.icmpToHloDirection(i.Predicate)
		g.subPrintf("  %s = pred[] compare(%s, %s), direction=%s\n", dst, lhs, rhs, dir)

	case *ir.FCmpInst:
		dst := fmt.Sprintf("%%v%d", g.nextID)
		g.nextID++
		localValMap[inst] = dst
		lhs := g.getOperandLocal(i.Operands()[0], localValMap)
		rhs := g.getOperandLocal(i.Operands()[1], localValMap)
		dir := g.fcmpToHloDirection(i.Predicate)
		g.subPrintf("  %s = pred[] compare(%s, %s), direction=%s\n", dst, lhs, rhs, dir)

	case *ir.PhiInst:
		// Phis handled separately in loop state
	}
}