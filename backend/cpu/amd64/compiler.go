package amd64

import (
	"bytes"
	"fmt"
	"math"

	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

type Artifact struct {
	Text    []byte
	Data    []byte
	Relocs  []RelocationRecord
	Symbols []SymbolDef
}

type SymbolDef struct {
	Name   string
	Offset uint64
	Size   uint64
	IsFunc bool
}

type compiler struct {
	asm          *Assembler
	runtime      *Runtime 
	data         *bytes.Buffer
	stackMap     map[ir.Value]int // Value -> RBP offset (negative)
	frameSize    int
	blockOffsets map[*ir.BasicBlock]int
	jumpsToFix   []jumpFixup

	currentFunc *ir.Function
}

type jumpFixup struct {
	asmOffset int
	target    *ir.BasicBlock
}

func Compile(m *ir.Module) (*Artifact, error) {
	c := &compiler{
		asm:      NewAssembler(),
		data:     new(bytes.Buffer),
		stackMap: make(map[ir.Value]int),
	}

	c.runtime = NewRuntime(c.asm)

	var syms []SymbolDef

	// 1. Compile Globals
	for _, g := range m.Globals {
		for c.data.Len()%8 != 0 {
			c.data.WriteByte(0)
		}

		offset := c.data.Len()
		if err := c.emitGlobal(g); err != nil {
			return nil, err
		}

		syms = append(syms, SymbolDef{
			Name: g.Name(), Offset: uint64(offset), Size: uint64(c.data.Len() - offset), IsFunc: false,
		})
	}

	// 2. Compile Functions
	for _, fn := range m.Functions {
		if len(fn.Blocks) == 0 {
			continue
		}

		for c.asm.Len()%16 != 0 {
			c.asm.emitByte(0x90)
		}

		start := c.asm.Len()
		if err := c.compileFunction(fn); err != nil {
			return nil, fmt.Errorf("in function %s: %w", fn.Name(), err)
		}
		end := c.asm.Len()

		syms = append(syms, SymbolDef{
			Name: fn.Name(), Offset: uint64(start), Size: uint64(end - start), IsFunc: true,
		})
	}

	return &Artifact{
		Text:    c.asm.Bytes(),
		Data:    c.data.Bytes(),
		Relocs:  c.asm.Relocs,
		Symbols: syms,
	}, nil
}

func (c *compiler) compileFunction(fn *ir.Function) error {
	c.currentFunc = fn
	c.stackMap = make(map[ir.Value]int)
	c.blockOffsets = make(map[*ir.BasicBlock]int)
	c.jumpsToFix = nil

    if fn.Name() == "main" {
        c.runtime.EmitInitialization()
    }

	offset := 0

	// Arguments
	for _, arg := range fn.Arguments {
		size := SizeOf(arg.Type())
		offset += size
		if offset%8 != 0 { offset += 8 - (offset % 8) }
		c.stackMap[arg] = -offset
	}

	// Instructions
	for _, block := range fn.Blocks {
		for _, inst := range block.Instructions {
			if alloca, ok := inst.(*ir.AllocaInst); ok {
				allocSize := SizeOf(alloca.AllocatedType)
				if count, ok := alloca.NumElements.(*ir.ConstantInt); ok {
					allocSize *= int(count.Value)
				}
				if offset%16 != 0 { offset += 16 - (offset % 16) }
				offset += allocSize
				c.stackMap[ir.Value(alloca)] = -offset
				continue
			}

			if inst.Type() != nil && inst.Type().BitSize() > 0 {
				size := SizeOf(inst.Type())
				if size == 0 { size = 1 }
				offset += size
				if offset%8 != 0 { offset += 8 - (offset % 8) }
				c.stackMap[inst] = -offset
			}
		}
	}

	if offset%16 != 0 { offset += 16 - (offset % 16) }
	c.frameSize = offset

	// Prologue
	c.asm.Push(RBP)
	c.asm.Mov(RegOp(RBP), RegOp(RSP), 64)
	if c.frameSize > 0 {
		c.asm.Sub(RegOp(RSP), ImmOp(c.frameSize))
	}

	// Save Register Arguments
	regs := []Register{RDI, RSI, RDX, RCX, R8, R9}
	xmmRegs := []Register{0, 1, 2, 3, 4, 5, 6, 7}
	gprIdx := 0
	xmmIdx := 0

	for _, arg := range fn.Arguments {
		if types.IsFloat(arg.Type()) {
			if xmmIdx < len(xmmRegs) {
				slot := c.getStackSlot(arg)
				if arg.Type().BitSize() == 64 {
					c.asm.Movsd(slot, xmmRegs[xmmIdx])
				} else {
					c.asm.Movss(slot, xmmRegs[xmmIdx])
				}
				xmmIdx++
			}
		} else {
			if gprIdx < len(regs) {
				slot := c.getStackSlot(arg)
				c.asm.Mov(slot, RegOp(regs[gprIdx]), 64)
				gprIdx++
			}
		}
	}

	// Compile Body
	for _, block := range fn.Blocks {
		c.blockOffsets[block] = c.asm.Len()
		for _, inst := range block.Instructions {
			if err := c.compileInst(inst); err != nil {
				return err
			}
		}
	}

	// Fixup Jumps
	for _, fix := range c.jumpsToFix {
		targetOff, ok := c.blockOffsets[fix.target]
		if !ok {
			return fmt.Errorf("jump target block not found")
		}
		rel := int32(targetOff - (fix.asmOffset + 4))
		c.asm.PatchInt32(fix.asmOffset, rel)
	}

	return nil
}

func (c *compiler) compileInst(inst ir.Instruction) error {
	if inst == nil { return fmt.Errorf("nil instruction encountered") }

	// DEBUG LOGGING
	// fmt.Printf("[CodeGen] Compiling Op %d: %s\n", inst.Opcode(), inst.String())

	if inst.Opcode() == ir.OpAsyncTaskCreate {
		return c.compileAsyncTaskCreate(inst.(*ir.AsyncTaskCreateInst))
	}
	if inst.Opcode() == ir.OpAsyncTaskAwait {
		return c.compileAsyncTaskAwait(inst.(*ir.AsyncTaskAwaitInst))
	}
	if inst.Opcode() == ir.OpProcessCreate {
		return c.compileProcessCreate(inst.(*ir.ProcessCreateInst))
	}

	// Safety helper
	requireOps := func(n int) error {
		if len(inst.Operands()) < n {
			return fmt.Errorf("instruction %s (Op %d) requires %d operands, got %d", 
				inst.String(), inst.Opcode(), n, len(inst.Operands()))
		}
		return nil
	}

	switch inst.Opcode() {
	case ir.OpAdd:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Add(RegOp(RAX), RegOp(RCX))
		c.store(RAX, inst)

	case ir.OpSub:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Sub(RegOp(RAX), RegOp(RCX))
		c.store(RAX, inst)

	case ir.OpMul:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Imul(RAX, RCX)
		c.store(RAX, inst)

	case ir.OpSDiv:
		if err := requireOps(2); err != nil { return err }
		op0 := inst.Operands()[0]
		op1 := inst.Operands()[1]
		c.load(RAX, op0)
		if op0.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
		c.load(RCX, op1)
		if op1.Type().BitSize() == 32 { c.asm.Movsxd(RCX, RCX) }
		c.asm.Cqo()
		c.asm.Div(RCX, true)
		c.store(RAX, inst)

	case ir.OpUDiv:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Xor(RegOp(RDX), RegOp(RDX))
		c.asm.Div(RCX, false)
		c.store(RAX, inst)

	case ir.OpSRem:
		if err := requireOps(2); err != nil { return err }
		op0 := inst.Operands()[0]
		op1 := inst.Operands()[1]
		c.load(RAX, op0)
		if op0.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
		c.load(RCX, op1)
		if op1.Type().BitSize() == 32 { c.asm.Movsxd(RCX, RCX) }
		c.asm.Cqo()
		c.asm.Div(RCX, true)
		c.store(RDX, inst)

	case ir.OpURem:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Xor(RegOp(RDX), RegOp(RDX))
		c.asm.Div(RCX, false)
		c.store(RDX, inst)

	case ir.OpAnd:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.And(RegOp(RAX), RegOp(RCX))
		c.store(RAX, inst)

	case ir.OpOr:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Or(RegOp(RAX), RegOp(RCX))
		c.store(RAX, inst)

	case ir.OpXor:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Xor(RegOp(RAX), RegOp(RCX))
		c.store(RAX, inst)

	case ir.OpShl:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Shl(RAX, RCX)
		c.store(RAX, inst)

	case ir.OpLShr:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])
		c.asm.Shr(RAX, RCX)
		c.store(RAX, inst)

	case ir.OpAShr:
		if err := requireOps(2); err != nil { return err }
		op0 := inst.Operands()[0]
		c.load(RAX, op0)
		if op0.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
		c.load(RCX, inst.Operands()[1])
		c.asm.Sar(RAX, RCX)
		c.store(RAX, inst)

	case ir.OpTrunc:
		if err := requireOps(1); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.store(RAX, inst)

	case ir.OpBitcast:
		if err := requireOps(1); err != nil { return err }
		c.moveValue(RBP, c.stackMap[inst], inst.Operands()[0])

	case ir.OpZExt:
		if err := requireOps(1); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.store(RAX, inst)

	case ir.OpSExt:
		if err := requireOps(1); err != nil { return err }
		src := inst.Operands()[0]
		srcSize := SizeOf(src.Type())
		c.load(RAX, src)
		if srcSize == 4 {
			c.asm.Movsxd(RAX, RAX)
		} else if srcSize == 1 {
			c.asm.Movsx(RAX, RegOp(RAX), 8)
		}
		c.store(RAX, inst)

	case ir.OpFPToSI:
		if err := requireOps(1); err != nil { return err }
		src := inst.Operands()[0]
		c.load(RAX, src)
		c.asm.Push(RAX)
		c.asm.Cvttss2si(RAX, NewMem(RSP, 0))
		c.asm.Pop(RCX)
		c.store(RAX, inst)

	case ir.OpSIToFP:
		if err := requireOps(1); err != nil { return err }
		src := inst.Operands()[0]
		c.load(RAX, src)
		c.asm.Push(RAX)
		c.asm.Cvtsi2ss(RAX, NewMem(RSP, 0))
		c.asm.Pop(RCX)
		destSlot := c.getStackSlot(inst)
		c.asm.Movss(destSlot, RAX)

	case ir.OpAlloca:
		return nil

	case ir.OpLoad:
		if err := requireOps(1); err != nil { return err }
		ptr := inst.Operands()[0]
		c.load(RCX, ptr) 
		c.moveFromMem(RBP, c.stackMap[inst], RCX, 0, SizeOf(inst.Type()))

	case ir.OpStore:
		if err := requireOps(2); err != nil { return err }
		val := inst.Operands()[0]
		ptr := inst.Operands()[1]
		c.load(RCX, ptr)
		c.moveValue(RCX, 0, val)

	case ir.OpGetElementPtr:
		gep := inst.(*ir.GetElementPtrInst)
		if err := requireOps(1); err != nil { return err }
		base := gep.Operands()[0]
		c.load(RAX, base)

		indices := gep.Operands()[1:]
		if len(indices) == 0 {
			c.store(RAX, inst)
			return nil
		}

		firstIdx := indices[0]
		baseType := gep.SourceElementType
		baseSize := SizeOf(baseType)

		if cIdx, ok := firstIdx.(*ir.ConstantInt); ok {
			if cIdx.Value != 0 {
				c.asm.Add(RegOp(RAX), ImmOp(int64(int(cIdx.Value)*baseSize)))
			}
		} else {
			c.load(RCX, firstIdx)
			c.asm.ImulImm(RCX, int32(baseSize))
			c.asm.Add(RegOp(RAX), RegOp(RCX))
		}

		currentType := baseType
		for _, idxVal := range indices[1:] {
			if st, ok := currentType.(*types.StructType); ok {
				if cIdx, ok := idxVal.(*ir.ConstantInt); ok {
					idx := int(cIdx.Value)
					offset := GetStructFieldOffset(st, idx)
					if offset != 0 {
						c.asm.Add(RegOp(RAX), ImmOp(int64(offset)))
					}
					if idx >= 0 && idx < len(st.Fields) {
						currentType = st.Fields[idx]
					} else {
						return fmt.Errorf("field index %d out of bounds for struct %s", idx, st.Name)
					}
				} else {
					return fmt.Errorf("non-constant struct index in GEP")
				}
			} else if at, ok := currentType.(*types.ArrayType); ok {
				elemSize := SizeOf(at.ElementType)
				if cIdx, ok := idxVal.(*ir.ConstantInt); ok {
					c.asm.Add(RegOp(RAX), ImmOp(int64(int(cIdx.Value)*elemSize)))
				} else {
					c.load(RCX, idxVal)
					c.asm.ImulImm(RCX, int32(elemSize))
					c.asm.Add(RegOp(RAX), RegOp(RCX))
				}
				currentType = at.ElementType
			}
		}
		c.store(RAX, inst)

	case ir.OpInsertValue:
		if err := requireOps(2); err != nil { return err }
		iv := inst.(*ir.InsertValueInst)
		agg := iv.Operands()[0]
		val := iv.Operands()[1]
		destOff := c.stackMap[inst]

		c.moveValue(RBP, destOff, agg)

		currentType := agg.Type()
		offset := 0
		for _, idx := range iv.Indices {
			if st, ok := currentType.(*types.StructType); ok {
				offset += GetStructFieldOffset(st, idx)
				currentType = st.Fields[idx]
			} else if at, ok := currentType.(*types.ArrayType); ok {
				offset += idx * SizeOf(at.ElementType)
				currentType = at.ElementType
			}
		}

		c.moveValue(RBP, destOff+offset, val)

	case ir.OpExtractValue:
		if err := requireOps(1); err != nil { return err }
		ev := inst.(*ir.ExtractValueInst)
		agg := ev.Operands()[0]

		currentType := agg.Type()
		offset := 0
		for _, idx := range ev.Indices {
			if st, ok := currentType.(*types.StructType); ok {
				offset += GetStructFieldOffset(st, idx)
				currentType = st.Fields[idx]
			} else if at, ok := currentType.(*types.ArrayType); ok {
				offset += idx * SizeOf(at.ElementType)
				currentType = at.ElementType
			}
		}

		if srcSlot, ok := c.stackMap[agg]; ok {
			c.moveFromMem(RBP, c.stackMap[inst], RBP, srcSlot+offset, SizeOf(inst.Type()))
		}

	case ir.OpSizeOf:
		val := SizeOf(inst.(*ir.SizeOfInst).QueryType)
		c.asm.Mov(RegOp(RAX), ImmOp(int64(val)), 64)
		c.store(RAX, inst)

	case ir.OpAlignOf:
		val := AlignOf(inst.(*ir.AlignOfInst).QueryType)
		c.asm.Mov(RegOp(RAX), ImmOp(int64(val)), 64)
		c.store(RAX, inst)

	case ir.OpStrLen:
		if err := requireOps(1); err != nil { return err }
		c.load(RDI, inst.Operands()[0])
		c.asm.CallRelative("strlen")
		c.store(RAX, inst)

	case ir.OpSyscall:
		ops := inst.Operands()
		if len(ops) == 0 { return fmt.Errorf("syscall requires at least 1 operand") }
		regs := []Register{RDI, RSI, RDX, R10, R8, R9}
		c.load(RAX, ops[0])
		for i, arg := range ops[1:] {
			if i < len(regs) {
				c.load(regs[i], arg)
			}
		}
		c.asm.Syscall()
		c.store(RAX, inst)

	case ir.OpCall:
		call := inst.(*ir.CallInst)
		gprRegs := []Register{RDI, RSI, RDX, RCX, R8, R9}
		xmmRegs := []Register{0, 1, 2, 3, 4, 5, 6, 7}
		gprIdx := 0
		xmmIdx := 0

		fixedParams := 0
		calleeVariadic := false
		if call.Callee != nil {
			fixedParams = len(call.Callee.FuncType.ParamTypes)
			calleeVariadic = call.Callee.FuncType.Variadic
		} else if call.CalleeVal != nil {
			if ptr, ok := call.CalleeVal.Type().(*types.PointerType); ok {
				if ft, ok := ptr.ElementType.(*types.FunctionType); ok {
					fixedParams = len(ft.ParamTypes)
					calleeVariadic = ft.Variadic
				}
			} else if ft, ok := call.CalleeVal.Type().(*types.FunctionType); ok {
				fixedParams = len(ft.ParamTypes)
				calleeVariadic = ft.Variadic
			}
		}

		for i, arg := range call.Operands() {
			if arg == nil { continue }
			isVariadic := i >= fixedParams && calleeVariadic
			if types.IsFloat(arg.Type()) {
				if xmmIdx < len(xmmRegs) {
					c.load(RAX, arg)
					if arg.Type().BitSize() == 64 {
						c.asm.Movq(xmmRegs[xmmIdx], RAX)
					} else {
						c.asm.Movd(xmmRegs[xmmIdx], RAX)
						if arg.Type().BitSize() == 32 && isVariadic {
							c.asm.Cvtss2sd(xmmRegs[xmmIdx], xmmRegs[xmmIdx])
						}
					}
					xmmIdx++
				}
			} else {
				if gprIdx < len(gprRegs) {
					c.load(gprRegs[gprIdx], arg)
					gprIdx++
				}
			}
		}

		c.asm.Mov(RegOp(RAX), ImmOp(int64(xmmIdx)), 8)

		if call.Callee != nil {
			c.asm.CallRelative(call.Callee.Name())
		} else if call.CalleeName != "" {
			c.asm.CallRelative(call.CalleeName)
		} else if call.CalleeVal != nil {
			c.load(RAX, call.CalleeVal)
			c.asm.CallReg(RAX)
		} else {
			return fmt.Errorf("invalid call instruction: no callee")
		}

		// Handle Return Values
		if call.Type() != nil && call.Type().Kind() != types.VoidKind {
			size := SizeOf(call.Type())
			if types.IsFloat(call.Type()) {
				if call.Type().BitSize() == 32 {
					c.asm.MovdXmmToGpr(RAX, 0)
				} else {
					c.asm.MovqXmmToGpr(RAX, 0)
				}
				c.store(RAX, inst)
			} else {
				if size <= 8 {
					c.store(RAX, inst)
				} else if size <= 16 {
					// Handle 16-byte structs returned in RAX:RDX
					if off, ok := c.stackMap[inst]; ok {
						c.asm.Mov(NewMem(RBP, off), RegOp(RAX), 64)
						c.asm.Mov(NewMem(RBP, off+8), RegOp(RDX), 64)
					}
				}
			}
		}

	case ir.OpRet:
		if len(inst.Operands()) > 0 {
			val := inst.Operands()[0]
			size := SizeOf(val.Type())
			if size <= 8 {
				c.load(RAX, val)
			} else if size <= 16 {
				// Return 16-byte struct in RAX:RDX
				if off, ok := c.stackMap[val]; ok {
					c.asm.Mov(RegOp(RAX), NewMem(RBP, off), 64)
					c.asm.Mov(RegOp(RDX), NewMem(RBP, off+8), 64)
				}
			}
		}
		c.asm.Mov(RegOp(RSP), RegOp(RBP), 64)
		c.asm.Pop(RBP)
		c.asm.Ret()

	case ir.OpBr:
		br := inst.(*ir.BrInst)
		c.handlePhi(inst.Parent(), br.Target)
		off := c.asm.JmpRel(0)
		c.jumpsToFix = append(c.jumpsToFix, jumpFixup{asmOffset: off, target: br.Target})

	case ir.OpCondBr:
		cbr := inst.(*ir.CondBrInst)
		if cbr.Condition == nil { return fmt.Errorf("CondBr missing condition") }
		
		c.load(RAX, cbr.Condition)
		c.asm.Test(RAX, RAX)
		offFalse := c.asm.JccRel(CondEq, 0)
		c.jumpsToFix = append(c.jumpsToFix, jumpFixup{asmOffset: offFalse, target: cbr.FalseBlock})
		c.handlePhi(inst.Parent(), cbr.TrueBlock)
		offTrue := c.asm.JmpRel(0)
		c.jumpsToFix = append(c.jumpsToFix, jumpFixup{asmOffset: offTrue, target: cbr.TrueBlock})

	case ir.OpICmp:
		if err := requireOps(2); err != nil { return err }
		c.load(RAX, inst.Operands()[0])
		c.load(RCX, inst.Operands()[1])

		icmp := inst.(*ir.ICmpInst)
		isSigned := false
		switch icmp.Predicate {
		case ir.ICmpSLT, ir.ICmpSLE, ir.ICmpSGT, ir.ICmpSGE:
			isSigned = true
		}
		if isSigned && inst.Operands()[0].Type().BitSize() == 32 {
			c.asm.Movsxd(RAX, RAX)
		}
		if isSigned && inst.Operands()[1].Type().BitSize() == 32 {
			c.asm.Movsxd(RCX, RCX)
		}

		c.asm.Cmp(RegOp(RAX), RegOp(RCX))

		var cc CondCode
		switch icmp.Predicate {
		case ir.ICmpEQ:
			cc = CondEq
		case ir.ICmpNE:
			cc = CondNe
		case ir.ICmpSLT:
			cc = CondLt
		case ir.ICmpSLE:
			cc = CondLe
		case ir.ICmpSGT:
			cc = CondGt
		case ir.ICmpSGE:
			cc = CondGe
		case ir.ICmpULT:
			cc = CondBlo
		case ir.ICmpULE:
			cc = CondBle
		case ir.ICmpUGT:
			cc = CondA
		case ir.ICmpUGE:
			cc = CondAe
		}
		c.asm.Setcc(cc, RAX)
		c.asm.MovZX(RAX, RegOp(RAX), 8)
		c.store(RAX, inst)

	case ir.OpPhi:
		return nil

	case ir.OpSelect:
		if err := requireOps(3); err != nil { return err }
		cond := inst.Operands()[0]
		trueVal := inst.Operands()[1]
		falseVal := inst.Operands()[2]

		c.load(RAX, cond)
		c.asm.Test(RAX, RAX)
		
		offFalse := c.asm.JccRel(CondEq, 0)
		c.load(RAX, trueVal)
		offDone := c.asm.JmpRel(0)
		
		c.asm.PatchInt32(offFalse, int32(c.asm.Len() - (offFalse + 4)))
		c.load(RAX, falseVal)
		
		c.asm.PatchInt32(offDone, int32(c.asm.Len() - (offDone + 4)))
		c.store(RAX, inst)

	case ir.OpUnreachable:
		c.asm.emitByte(0x0F)
		c.asm.emitByte(0x0B)

	default:
		return fmt.Errorf("unknown opcode: %s", inst.Opcode())
	}
	return nil
}

func (c *compiler) compileAsyncTaskCreate(inst *ir.AsyncTaskCreateInst) error {
	if inst.Callee == nil {
		return fmt.Errorf("async_task_create requires direct function call")
	}
    loader := func(reg Register, val ir.Value) {
        c.load(reg, val)
    }
	c.runtime.EmitAsyncTaskCreate(inst.Callee, inst.Operands(), 64*1024, loader)
	c.store(RAX, inst)
	return nil
}

func (c *compiler) compileAsyncTaskAwait(inst *ir.AsyncTaskAwaitInst) error {
	if len(inst.Operands()) == 0 { return fmt.Errorf("await missing operand") }
	handleOp := inst.Operands()[0]
	c.load(RCX, handleOp)
	c.runtime.EmitAsyncTaskAwait(RCX)
	c.store(RAX, inst)
	return nil
}

func (c *compiler) getStackSlot(v ir.Value) MemOp {
	off, ok := c.stackMap[v]
	if !ok {
		// Not panic, return error? No, logic error.
		panic(fmt.Sprintf("Value %v (%s) not allocated in stack map", v, v.Name()))
	}
	return NewMem(RBP, off)
}

func (c *compiler) load(dst Register, src ir.Value) {
	switch v := src.(type) {
	case *ir.ConstantInt:
		c.asm.Mov(RegOp(dst), ImmOp(v.Value), 64)
	case *ir.ConstantFloat:
		bits := math.Float64bits(v.Value)
		if v.Type().BitSize() == 32 {
			bits = uint64(math.Float32bits(float32(v.Value)))
		}
		c.asm.Mov(RegOp(dst), ImmOp(int64(bits)), 64)
	case *ir.ConstantNull:
		c.asm.Xor(RegOp(dst), RegOp(dst))
	case *ir.ConstantZero:
		c.asm.Xor(RegOp(dst), RegOp(dst))
	case *ir.Global:
		c.asm.LeaRel(dst, v.Name())
	case *ir.Function:
		c.asm.LeaRel(dst, v.Name())
	case *ir.Argument:
		// FIX: Arguments are already in their stack slots (saved by prologue)
		// Just load them directly from the stack
		slot := c.getStackSlot(v)
		typ := v.Type()
		size := SizeOf(typ)
		if size == 8 {
			c.asm.Mov(RegOp(dst), slot, 64)
		} else if size == 4 {
			c.asm.Mov(RegOp(dst), slot, 32)
		} else if size == 2 {
			c.asm.Mov(RegOp(dst), slot, 16)
		} else if size == 1 {
			isSigned := false
			if intTy, ok := typ.(*types.IntType); ok && intTy.Signed {
				isSigned = true
			}
			if isSigned {
				c.asm.Movsx(dst, slot, 8)
			} else {
				c.asm.MovZX(dst, slot, 8)
			}
		} else {
			c.asm.Mov(RegOp(dst), slot, 64)
		}
	case *ir.AllocaInst:
		off := c.stackMap[v]
		c.asm.Lea(dst, NewMem(RBP, off))
	case *ir.LoadInst:
		// LoadInst that hasn't been compiled yet - compile it inline
		if _, hasSlot := c.stackMap[v]; !hasSlot {
			ptr := v.Operands()[0]
			c.load(dst, ptr)
			c.asm.Mov(RegOp(dst), NewMem(dst, 0), 64)
			return
		}
		// LoadInst that was compiled - load from stack
		slot := c.getStackSlot(v)
		typ := v.Type()
		size := SizeOf(typ)
		if size == 8 {
			c.asm.Mov(RegOp(dst), slot, 64)
		} else if size == 4 {
			c.asm.Mov(RegOp(dst), slot, 32)
		} else if size == 2 {
			c.asm.Mov(RegOp(dst), slot, 16)
		} else if size == 1 {
			isSigned := false
			if intTy, ok := typ.(*types.IntType); ok && intTy.Signed {
				isSigned = true
			}
			if isSigned {
				c.asm.Movsx(dst, slot, 8)
			} else {
				c.asm.MovZX(dst, slot, 8)
			}
		} else {
			c.asm.Mov(RegOp(dst), slot, 64)
		}
	case *ir.ICmpInst:
		// ICmpInst that hasn't been compiled yet - compile it inline
		if _, hasSlot := c.stackMap[v]; !hasSlot {
			lhs := v.Operands()[0]
			rhs := v.Operands()[1]
			
			c.load(RAX, lhs)
			
			// Handle signed comparisons
			isSigned := false
			switch v.Predicate {
			case ir.ICmpSLT, ir.ICmpSLE, ir.ICmpSGT, ir.ICmpSGE:
				isSigned = true
			}
			if isSigned && lhs.Type().BitSize() == 32 {
				c.asm.Movsxd(RAX, RAX)
			}
			
			c.load(RCX, rhs)
			if isSigned && rhs.Type().BitSize() == 32 {
				c.asm.Movsxd(RCX, RCX)
			}
			
			c.asm.Cmp(RegOp(RAX), RegOp(RCX))
			
			var cc CondCode
			switch v.Predicate {
			case ir.ICmpEQ:
				cc = CondEq
			case ir.ICmpNE:
				cc = CondNe
			case ir.ICmpSLT:
				cc = CondLt
			case ir.ICmpSLE:
				cc = CondLe
			case ir.ICmpSGT:
				cc = CondGt
			case ir.ICmpSGE:
				cc = CondGe
			case ir.ICmpULT:
				cc = CondBlo
			case ir.ICmpULE:
				cc = CondBle
			case ir.ICmpUGT:
				cc = CondA
			case ir.ICmpUGE:
				cc = CondAe
			}
			c.asm.Setcc(cc, dst)
			c.asm.MovZX(dst, RegOp(dst), 8)
			return
		}
		// ICmpInst that was compiled - load from stack
		slot := c.getStackSlot(v)
		c.asm.Mov(RegOp(dst), slot, 8)
		c.asm.MovZX(dst, RegOp(dst), 8)
	case *ir.BinaryInst:
		// BinaryInst that hasn't been compiled yet - compile it inline
		if _, hasSlot := c.stackMap[v]; !hasSlot {
			lhs := v.Operands()[0]
			rhs := v.Operands()[1]
			
			c.load(RAX, lhs)
			c.load(RCX, rhs)
			
			switch v.Opcode() {
			case ir.OpAdd:
				c.asm.Add(RegOp(RAX), RegOp(RCX))
			case ir.OpSub:
				c.asm.Sub(RegOp(RAX), RegOp(RCX))
			case ir.OpMul:
				c.asm.Imul(RAX, RCX)
			case ir.OpSDiv:
				if lhs.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
				if rhs.Type().BitSize() == 32 { c.asm.Movsxd(RCX, RCX) }
				c.asm.Cqo()
				c.asm.Div(RCX, true)
			case ir.OpUDiv:
				c.asm.Xor(RegOp(RDX), RegOp(RDX))
				c.asm.Div(RCX, false)
			case ir.OpSRem:
				if lhs.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
				if rhs.Type().BitSize() == 32 { c.asm.Movsxd(RCX, RCX) }
				c.asm.Cqo()
				c.asm.Div(RCX, true)
				c.asm.Mov(RegOp(RAX), RegOp(RDX), 64)
			case ir.OpURem:
				c.asm.Xor(RegOp(RDX), RegOp(RDX))
				c.asm.Div(RCX, false)
				c.asm.Mov(RegOp(RAX), RegOp(RDX), 64)
			case ir.OpAnd:
				c.asm.And(RegOp(RAX), RegOp(RCX))
			case ir.OpOr:
				c.asm.Or(RegOp(RAX), RegOp(RCX))
			case ir.OpXor:
				c.asm.Xor(RegOp(RAX), RegOp(RCX))
			case ir.OpShl:
				c.asm.Shl(RAX, RCX)
			case ir.OpLShr:
				c.asm.Shr(RAX, RCX)
			case ir.OpAShr:
				if lhs.Type().BitSize() == 32 { c.asm.Movsxd(RAX, RAX) }
				c.asm.Sar(RAX, RCX)
			}
			
			if dst != RAX {
				c.asm.Mov(RegOp(dst), RegOp(RAX), 64)
			}
			return
		}
		// BinaryInst that was compiled - load from stack
		slot := c.getStackSlot(v)
		c.asm.Mov(RegOp(dst), slot, 64)

	default:
		slot := c.getStackSlot(v)
		typ := v.Type()
		size := SizeOf(typ)

		if size == 8 {
			c.asm.Mov(RegOp(dst), slot, 64)
		} else if size == 4 {
			c.asm.Mov(RegOp(dst), slot, 32)
		} else if size == 2 {
			c.asm.Mov(RegOp(dst), slot, 16)
		} else if size == 1 {
			isSigned := false
			if intTy, ok := typ.(*types.IntType); ok && intTy.Signed {
				isSigned = true
			}
			if isSigned {
				c.asm.Movsx(dst, slot, 8)
			} else {
				c.asm.MovZX(dst, slot, 8)
			}
		} else {
			c.asm.Mov(RegOp(dst), slot, 64)
		}
	}
}

func (c *compiler) store(src Register, dst ir.Value) {
	if _, ok := c.stackMap[dst]; !ok {
		return
	}

	slot := c.getStackSlot(dst)
	size := SizeOf(dst.Type())
	if size == 8 {
		c.asm.Mov(slot, RegOp(src), 64)
	} else if size == 4 {
		c.asm.Mov(slot, RegOp(src), 32)
	} else if size == 1 {
		c.asm.Mov(slot, RegOp(src), 8)
	} else {
		c.asm.Mov(slot, RegOp(src), 64)
	}
}

func (c *compiler) moveValue(dstBase Register, dstDisp int, src ir.Value) {
	if alloca, ok := src.(*ir.AllocaInst); ok {
		off := c.stackMap[alloca]
		c.asm.Lea(RAX, NewMem(RBP, off))
		c.asm.Mov(NewMem(dstBase, dstDisp), RegOp(RAX), 64)
		return
	}

	size := SizeOf(src.Type())

	if cInt, ok := src.(*ir.ConstantInt); ok {
		if size <= 4 || (cInt.Value >= -2147483648 && cInt.Value <= 2147483647) {
			c.asm.Mov(NewMem(dstBase, dstDisp), ImmOp(cInt.Value), size*8)
		} else {
			c.asm.Mov(RegOp(RAX), ImmOp(cInt.Value), 64)
			c.asm.Mov(NewMem(dstBase, dstDisp), RegOp(RAX), 64)
		}
		return
	}

	if cFloat, ok := src.(*ir.ConstantFloat); ok {
		bits := math.Float64bits(cFloat.Value)
		if size == 4 {
			bits = uint64(math.Float32bits(float32(cFloat.Value)))
		}

		if size <= 4 || (int64(bits) >= -2147483648 && int64(bits) <= 2147483647) {
			c.asm.Mov(NewMem(dstBase, dstDisp), ImmOp(int64(bits)), size*8)
		} else {
			c.asm.Mov(RegOp(RAX), ImmOp(int64(bits)), 64)
			c.asm.Mov(NewMem(dstBase, dstDisp), RegOp(RAX), 64)
		}
		return
	}

	if _, ok := src.(*ir.ConstantZero); ok {
		// Optimize: Write 64-bit chunks where possible
		offset := 0
		for offset+8 <= size {
			c.asm.Mov(NewMem(dstBase, dstDisp+offset), ImmOp(0), 64)
			offset += 8
		}
		// Write remaining bytes
		for offset < size {
			c.asm.Mov(NewMem(dstBase, dstDisp+offset), ImmOp(0), 8)
			offset++
		}
		return
	}

	if cArr, ok := src.(*ir.ConstantArray); ok {
		elemSize := SizeOf(cArr.Type().(*types.ArrayType).ElementType)
		for i, elem := range cArr.Elements {
			c.moveValue(dstBase, dstDisp+(i*elemSize), elem)
		}
		return
	}

	if cStruct, ok := src.(*ir.ConstantStruct); ok {
		st := cStruct.Type().(*types.StructType)
		for i, field := range cStruct.Fields {
			offset := GetStructFieldOffset(st, i)
			c.moveValue(dstBase, dstDisp+offset, field)
		}
		return
	}

	if srcSlot, ok := c.stackMap[src]; ok {
		c.moveFromMem(dstBase, dstDisp, RBP, srcSlot, size)
		return
	}

	if size <= 8 {
		c.load(RAX, src)
		c.asm.Mov(NewMem(dstBase, dstDisp), RegOp(RAX), size*8)
	} else {
		if g, ok := src.(*ir.Global); ok {
			c.asm.LeaRel(RCX, g.Name())
			c.moveFromMem(dstBase, dstDisp, RCX, 0, size)
		} else {
			panic(fmt.Sprintf("Unsupported large move from %T", src))
		}
	}
}

func (c *compiler) moveFromMem(dstBase Register, dstDisp int, srcBase Register, srcDisp int, size int) {
	offset := 0
	for offset+8 <= size {
		c.asm.Mov(RegOp(RAX), NewMem(srcBase, srcDisp+offset), 64)
		c.asm.Mov(NewMem(dstBase, dstDisp+offset), RegOp(RAX), 64)
		offset += 8
	}
	if offset+4 <= size {
		c.asm.Mov(RegOp(RAX), NewMem(srcBase, srcDisp+offset), 32)
		c.asm.Mov(NewMem(dstBase, dstDisp+offset), RegOp(RAX), 32)
		offset += 4
	}
	for offset < size {
		c.asm.MovZX(RAX, NewMem(srcBase, srcDisp+offset), 8)
		c.asm.Mov(NewMem(dstBase, dstDisp+offset), RegOp(RAX), 8)
		offset++
	}
}

func (c *compiler) handlePhi(from, to *ir.BasicBlock) {
	for _, inst := range to.Instructions {
		if phi, ok := inst.(*ir.PhiInst); ok {
			for _, incoming := range phi.Incoming {
				if incoming.Block == from {
					c.moveValue(RBP, c.stackMap[phi], incoming.Value)
					break
				}
			}
		}
	}
}

func (c *compiler) emitGlobal(g *ir.Global) error {
	if g.Initializer != nil {
		return c.emitConstant(g.Initializer)
	}
	size := SizeOf(g.Type())
	c.data.Write(make([]byte, size))
	return nil
}

func (c *compiler) emitConstant(k ir.Constant) error {
	switch v := k.(type) {
	case *ir.ConstantInt:
		size := SizeOf(v.Type())
		val := uint64(v.Value)
		for i := 0; i < size; i++ {
			c.data.WriteByte(byte(val))
			val >>= 8
		}
	case *ir.ConstantFloat:
		size := SizeOf(v.Type())
		var val uint64
		if size == 4 {
			val = uint64(math.Float32bits(float32(v.Value)))
		} else {
			val = math.Float64bits(v.Value)
		}
		for i := 0; i < size; i++ {
			c.data.WriteByte(byte(val))
			val >>= 8
		}
	case *ir.ConstantArray:
		for _, elem := range v.Elements {
			if err := c.emitConstant(elem); err != nil { return err }
		}
	case *ir.ConstantStruct:
		st, ok := v.Type().(*types.StructType)
		if !ok { return fmt.Errorf("ConstantStruct has non-struct type") }

		currentOffset := 0
		for i, field := range v.Fields {
			targetOffset := GetStructFieldOffset(st, i)
			if targetOffset > currentOffset {
				padding := targetOffset - currentOffset
				c.data.Write(make([]byte, padding))
				currentOffset += padding
			}
			if err := c.emitConstant(field); err != nil { return err }
			currentOffset += SizeOf(field.Type())
		}
		totalSize := SizeOf(v.Type())
		if currentOffset < totalSize {
			c.data.Write(make([]byte, totalSize-currentOffset))
		}

	case *ir.ConstantZero:
		size := SizeOf(v.Type())
		c.data.Write(make([]byte, size))
	case *ir.ConstantNull:
		c.data.Write(make([]byte, 8))
	default:
		return fmt.Errorf("unsupported constant type: %T", k)
	}
	return nil
}

func (c *compiler) compileProcessCreate(inst *ir.ProcessCreateInst) error {
	preservedRegs := []Register{R12, R13, R14, R15, RBX}
	if len(inst.Operands()) > len(preservedRegs) {
		return fmt.Errorf("process_create currently supports max 5 arguments")
	}
	for i, arg := range inst.Operands() {
		c.load(preservedRegs[i], arg)
	}
	c.runtime.EmitProcessCreate(inst.Callee, len(inst.Operands()))
	c.store(RAX, inst)
	return nil
}