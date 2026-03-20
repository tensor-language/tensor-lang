package builder

import (
	"fmt"

	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

// Builder constructs IR instructions
type Builder struct {
	module       *ir.Module
	currentFunc  *ir.Function
	currentBlock *ir.BasicBlock
	insertPoint  int // -1 means append to end
	nameCounter  int
}

// New creates a new IR builder
func New() *Builder {
	return &Builder{
		insertPoint: -1,
	}
}

// NewWithModule creates a builder with an existing module
func NewWithModule(m *ir.Module) *Builder {
	return &Builder{
		module:      m,
		insertPoint: -1,
	}
}

// Module returns the current module
func (b *Builder) Module() *ir.Module {
	return b.module
}

// CreateModule creates a new module
func (b *Builder) CreateModule(name string) *ir.Module {
	b.module = ir.NewModule(name)
	return b.module
}

// CurrentFunction returns the function being built
func (b *Builder) CurrentFunction() *ir.Function {
	return b.currentFunc
}

// CurrentBlock returns the current basic block
func (b *Builder) CurrentBlock() *ir.BasicBlock {
	return b.currentBlock
}

// SetInsertPoint sets where instructions will be inserted
func (b *Builder) SetInsertPoint(block *ir.BasicBlock) {
	b.currentBlock = block
	if block != nil {
		b.currentFunc = block.Parent
	}
	b.insertPoint = -1
}

// SetInsertPointBefore sets insertion point before an instruction
func (b *Builder) SetInsertPointBefore(inst ir.Instruction) {
	b.currentBlock = inst.Parent()
	if b.currentBlock != nil {
		for i, in := range b.currentBlock.Instructions {
			if in == inst {
				b.insertPoint = i
				return
			}
		}
	}
	b.insertPoint = -1 // Fallback
}

// GetInsertBlock returns the current insertion block
func (b *Builder) GetInsertBlock() *ir.BasicBlock {
	return b.currentBlock
}

// generateName creates a unique name for unnamed values
func (b *Builder) generateName() string {
	name := fmt.Sprintf("%d", b.nameCounter)
	b.nameCounter++
	return name
}

// insert adds an instruction at the current insertion point
func (b *Builder) insert(inst ir.Instruction) {
	if b.currentBlock == nil {
		// Silently fail or panic? Panic helps debug.
		panic("IR Builder: Attempted to insert instruction with no active BasicBlock set.")
	}
	
	if b.insertPoint < 0 || b.insertPoint >= len(b.currentBlock.Instructions) {
		// Append mode
		b.currentBlock.AddInstruction(inst)
	} else {
		// Insert mode
		insts := b.currentBlock.Instructions
		newInsts := make([]ir.Instruction, len(insts)+1)
		
		copy(newInsts, insts[:b.insertPoint])
		newInsts[b.insertPoint] = inst
		copy(newInsts[b.insertPoint+1:], insts[b.insertPoint:])
		
		b.currentBlock.Instructions = newInsts
		inst.SetParent(b.currentBlock)
		b.insertPoint++
	}
}

// ============================================================================
// Module-level operations
// ============================================================================

func (b *Builder) CreateFunction(name string, retType types.Type, params []types.Type, variadic bool) *ir.Function {
	fnType := types.NewFunction(retType, params, variadic)
	fn := ir.NewFunction(name, fnType)
	if b.module != nil {
		b.module.AddFunction(fn)
	}
	b.currentFunc = fn
	return fn
}

func (b *Builder) DeclareFunction(name string, retType types.Type, params []types.Type, variadic bool) *ir.Function {
	fnType := types.NewFunction(retType, params, variadic)
	fn := ir.NewFunction(name, fnType)
	fn.Linkage = ir.ExternalLinkage
	if b.module != nil {
		b.module.AddFunction(fn)
	}
	return fn
}

func (b *Builder) SetCallConv(fn *ir.Function, cc ir.CallingConvention) {
	fn.CallConv = cc
}

func (b *Builder) DefineStruct(st *types.StructType) {
	if b.module != nil && st.Name != "" {
		b.module.Types[st.Name] = st
	}
}

func (b *Builder) CreateGlobalVariable(name string, typ types.Type, initializer ir.Constant) *ir.Global {
	g := &ir.Global{
		Initializer: initializer,
		Linkage:     ir.ExternalLinkage,
	}
	g.SetType(types.NewPointer(typ))
	g.SetName(name)
	if b.module != nil {
		b.module.AddGlobal(g)
	}
	return g
}

func (b *Builder) CreateGlobalConstant(name string, initializer ir.Constant) *ir.Global {
	g := &ir.Global{
		Initializer: initializer,
		IsConstant:  true,
		Linkage:     ir.ExternalLinkage,
	}
	g.SetName(name)
	g.SetType(types.NewPointer(initializer.Type()))
	if b.module != nil {
		b.module.AddGlobal(g)
	}
	return g
}

// ============================================================================
// Basic Block operations
// ============================================================================

func (b *Builder) CreateBlock(baseName string) *ir.BasicBlock {
	uniqueName := fmt.Sprintf("%s.%d", baseName, b.nameCounter)
	b.nameCounter++
	
	block := ir.NewBasicBlock(uniqueName)
	if b.currentFunc != nil {
		b.currentFunc.AddBlock(block)
	}
	return block
}

func (b *Builder) CreateBlockInFunction(baseName string, fn *ir.Function) *ir.BasicBlock {
	uniqueName := fmt.Sprintf("%s.%d", baseName, b.nameCounter)
	b.nameCounter++
	
	block := ir.NewBasicBlock(uniqueName)
	fn.AddBlock(block)
	return block
}

// ============================================================================
// Terminator instructions
// ============================================================================

func (b *Builder) CreateRet(v ir.Value) *ir.RetInst {
	inst := &ir.RetInst{}
	inst.Self = inst
	inst.Op = ir.OpRet
	if v != nil {
		inst.SetOperand(0, v)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateRetVoid() *ir.RetInst {
	inst := &ir.RetInst{}
	inst.Self = inst
	inst.Op = ir.OpRet
	b.insert(inst)
	return inst
}

func (b *Builder) CreateBr(target *ir.BasicBlock) *ir.BrInst {
	inst := &ir.BrInst{Target: target}
	inst.Self = inst
	inst.Op = ir.OpBr
	b.insert(inst)
	b.currentBlock.Successors = append(b.currentBlock.Successors, target)
	target.Predecessors = append(target.Predecessors, b.currentBlock)
	return inst
}

func (b *Builder) CreateCondBr(cond ir.Value, trueBlock, falseBlock *ir.BasicBlock) *ir.CondBrInst {
	inst := &ir.CondBrInst{
		Condition:  cond,
		TrueBlock:  trueBlock,
		FalseBlock: falseBlock,
	}
	inst.Self = inst
	inst.Op = ir.OpCondBr
	
	// FIX: Register user to prevent DCE deletion
	if cond != nil {
		if tracker, ok := cond.(ir.TrackableValue); ok {
			tracker.AddUser(inst)
		}
	}

	b.insert(inst)
	b.currentBlock.Successors = append(b.currentBlock.Successors, trueBlock, falseBlock)
	trueBlock.Predecessors = append(trueBlock.Predecessors, b.currentBlock)
	falseBlock.Predecessors = append(falseBlock.Predecessors, b.currentBlock)
	return inst
}

func (b *Builder) CreateSwitch(cond ir.Value, defaultBlock *ir.BasicBlock, numCases int) *ir.SwitchInst {
	inst := &ir.SwitchInst{
		Condition:    cond,
		DefaultBlock: defaultBlock,
		Cases:        make([]ir.SwitchCase, 0, numCases),
	}
	inst.Self = inst
	inst.Op = ir.OpSwitch
	
	// FIX: Register user
	if cond != nil {
		if tracker, ok := cond.(ir.TrackableValue); ok {
			tracker.AddUser(inst)
		}
	}

	b.insert(inst)
	b.currentBlock.Successors = append(b.currentBlock.Successors, defaultBlock)
	defaultBlock.Predecessors = append(defaultBlock.Predecessors, b.currentBlock)
	return inst
}

func (b *Builder) AddCase(sw *ir.SwitchInst, val *ir.ConstantInt, block *ir.BasicBlock) {
	sw.Cases = append(sw.Cases, ir.SwitchCase{Value: val, Block: block})
	parent := sw.Parent()
	if parent != nil {
		parent.Successors = append(parent.Successors, block)
		block.Predecessors = append(block.Predecessors, parent)
	}
}

func (b *Builder) CreateUnreachable() *ir.UnreachableInst {
	inst := &ir.UnreachableInst{}
	inst.Self = inst
	inst.Op = ir.OpUnreachable
	b.insert(inst)
	return inst
}

// ============================================================================
// Binary operations
// ============================================================================

func (b *Builder) createBinaryOp(op ir.Opcode, lhs, rhs ir.Value, name string) *ir.BinaryInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.BinaryInst{}
	inst.Self = inst
	inst.Op = op
	inst.SetName(name)
	inst.SetOperand(0, lhs)
	inst.SetOperand(1, rhs)
	inst.SetType(lhs.Type())
	b.insert(inst)
	return inst
}

func (b *Builder) CreateAdd(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpAdd, lhs, rhs, name)
}

func (b *Builder) CreateNSWAdd(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpAdd, lhs, rhs, name)
	inst.NoSignedWrap = true
	return inst
}

func (b *Builder) CreateNUWAdd(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpAdd, lhs, rhs, name)
	inst.NoUnsignedWrap = true
	return inst
}

func (b *Builder) CreateSub(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpSub, lhs, rhs, name)
}

func (b *Builder) CreateNSWSub(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpSub, lhs, rhs, name)
	inst.NoSignedWrap = true
	return inst
}

func (b *Builder) CreateMul(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpMul, lhs, rhs, name)
}

func (b *Builder) CreateNSWMul(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpMul, lhs, rhs, name)
	inst.NoSignedWrap = true
	return inst
}

func (b *Builder) CreateUDiv(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpUDiv, lhs, rhs, name)
}

func (b *Builder) CreateExactUDiv(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpUDiv, lhs, rhs, name)
	inst.Exact = true
	return inst
}

func (b *Builder) CreateSDiv(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpSDiv, lhs, rhs, name)
}

func (b *Builder) CreateExactSDiv(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	inst := b.createBinaryOp(ir.OpSDiv, lhs, rhs, name)
	inst.Exact = true
	return inst
}

func (b *Builder) CreateURem(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpURem, lhs, rhs, name)
}

func (b *Builder) CreateSRem(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpSRem, lhs, rhs, name)
}

func (b *Builder) CreateFAdd(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpFAdd, lhs, rhs, name)
}

func (b *Builder) CreateFSub(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpFSub, lhs, rhs, name)
}

func (b *Builder) CreateFMul(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpFMul, lhs, rhs, name)
}

func (b *Builder) CreateFDiv(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpFDiv, lhs, rhs, name)
}

func (b *Builder) CreateFRem(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpFRem, lhs, rhs, name)
}

// ============================================================================
// Bitwise operations
// ============================================================================

func (b *Builder) CreateShl(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpShl, lhs, rhs, name)
}

func (b *Builder) CreateLShr(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpLShr, lhs, rhs, name)
}

func (b *Builder) CreateAShr(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpAShr, lhs, rhs, name)
}

func (b *Builder) CreateAnd(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpAnd, lhs, rhs, name)
}

func (b *Builder) CreateOr(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpOr, lhs, rhs, name)
}

func (b *Builder) CreateXor(lhs, rhs ir.Value, name string) *ir.BinaryInst {
	return b.createBinaryOp(ir.OpXor, lhs, rhs, name)
}

// ============================================================================
// Memory operations
// ============================================================================

func (b *Builder) CreateAlloca(typ types.Type, name string) *ir.AllocaInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.AllocaInst{
		AllocatedType: typ,
	}
	inst.Self = inst
	inst.Op = ir.OpAlloca
	inst.SetType(types.NewPointer(typ))
	inst.SetName(name)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateAllocaWithCount(typ types.Type, count ir.Value, name string) *ir.AllocaInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.AllocaInst{
		AllocatedType: typ,
		NumElements:   count,
	}
	inst.Self = inst
	inst.Op = ir.OpAlloca
	inst.SetType(types.NewPointer(typ))
	inst.SetName(name)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateLoad(typ types.Type, ptr ir.Value, name string) *ir.LoadInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.LoadInst{}
	inst.Self = inst
	inst.Op = ir.OpLoad
	inst.SetName(name)
	inst.SetType(typ)
	inst.SetOperand(0, ptr)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateVolatileLoad(typ types.Type, ptr ir.Value, name string) *ir.LoadInst {
	inst := b.CreateLoad(typ, ptr, name)
	inst.Volatile = true
	return inst
}

func (b *Builder) CreateAlignedLoad(typ types.Type, ptr ir.Value, align int, name string) *ir.LoadInst {
	inst := b.CreateLoad(typ, ptr, name)
	inst.Alignment = align
	return inst
}

func (b *Builder) CreateStore(val ir.Value, ptr ir.Value) *ir.StoreInst {
	inst := &ir.StoreInst{}
	inst.Self = inst
	inst.Op = ir.OpStore
	inst.SetOperand(0, val)
	inst.SetOperand(1, ptr)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateVolatileStore(val ir.Value, ptr ir.Value) *ir.StoreInst {
	inst := b.CreateStore(val, ptr)
	inst.Volatile = true
	return inst
}

func (b *Builder) CreateAlignedStore(val ir.Value, ptr ir.Value, align int) *ir.StoreInst {
	inst := b.CreateStore(val, ptr)
	inst.Alignment = align
	return inst
}

func (b *Builder) CreateGEP(pointeeType types.Type, ptr ir.Value, indices []ir.Value, name string) *ir.GetElementPtrInst {
	if name == "" {
		name = b.generateName()
	}
	operands := make([]ir.Value, 1+len(indices))
	operands[0] = ptr
	copy(operands[1:], indices)

	inst := &ir.GetElementPtrInst{
		SourceElementType: pointeeType,
	}
	inst.Self = inst
	inst.Op = ir.OpGetElementPtr
	inst.SetName(name)

	resultType := pointeeType

	if len(indices) > 1 {
		for _, idxVal := range indices[1:] {
			if st, ok := resultType.(*types.StructType); ok {
				if cIdx, ok := idxVal.(*ir.ConstantInt); ok {
					if int(cIdx.Value) < len(st.Fields) {
						resultType = st.Fields[cIdx.Value]
					}
				}
			} else if at, ok := resultType.(*types.ArrayType); ok {
				resultType = at.ElementType
			}
		}
	}

	inst.SetType(types.NewPointer(resultType))

	for i, op := range operands {
		inst.SetOperand(i, op)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateInBoundsGEP(pointeeType types.Type, ptr ir.Value, indices []ir.Value, name string) *ir.GetElementPtrInst {
	inst := b.CreateGEP(pointeeType, ptr, indices, name)
	inst.InBounds = true
	return inst
}

func (b *Builder) CreateStructGEP(structType types.Type, ptr ir.Value, idx int, name string) *ir.GetElementPtrInst {
	zero := b.ConstInt(types.I32, 0)
	idxVal := b.ConstInt(types.I32, int64(idx))
	return b.CreateGEP(structType, ptr, []ir.Value{zero, idxVal}, name)
}

// ============================================================================
// Cast operations
// ============================================================================

func (b *Builder) createCast(op ir.Opcode, v ir.Value, destTy types.Type, name string) *ir.CastInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.CastInst{
		DestType: destTy,
	}
	inst.Self = inst
	inst.Op = op
	inst.SetName(name)
	inst.SetType(destTy)
	inst.SetOperand(0, v)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateTrunc(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpTrunc, v, destTy, name)
}

func (b *Builder) CreateZExt(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpZExt, v, destTy, name)
}

func (b *Builder) CreateSExt(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpSExt, v, destTy, name)
}

func (b *Builder) CreateFPTrunc(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpFPTrunc, v, destTy, name)
}

func (b *Builder) CreateFPExt(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpFPExt, v, destTy, name)
}

func (b *Builder) CreateFPToUI(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpFPToUI, v, destTy, name)
}

func (b *Builder) CreateFPToSI(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpFPToSI, v, destTy, name)
}

func (b *Builder) CreateUIToFP(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpUIToFP, v, destTy, name)
}

func (b *Builder) CreateSIToFP(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpSIToFP, v, destTy, name)
}

func (b *Builder) CreatePtrToInt(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpPtrToInt, v, destTy, name)
}

func (b *Builder) CreateIntToPtr(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpIntToPtr, v, destTy, name)
}

func (b *Builder) CreateBitCast(v ir.Value, destTy types.Type, name string) *ir.CastInst {
	return b.createCast(ir.OpBitcast, v, destTy, name)
}

// ============================================================================
// Comparison operations
// ============================================================================

func (b *Builder) CreateICmp(pred ir.ICmpPredicate, lhs, rhs ir.Value, name string) *ir.ICmpInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.ICmpInst{
		Predicate: pred,
	}
	inst.Self = inst
	inst.Op = ir.OpICmp
	inst.SetName(name)
	inst.SetType(types.I1)
	inst.SetOperand(0, lhs)
	inst.SetOperand(1, rhs)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateFCmp(pred ir.FCmpPredicate, lhs, rhs ir.Value, name string) *ir.FCmpInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.FCmpInst{
		Predicate: pred,
	}
	inst.Self = inst
	inst.Op = ir.OpFCmp
	inst.SetName(name)
	inst.SetType(types.I1)
	inst.SetOperand(0, lhs)
	inst.SetOperand(1, rhs)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateICmpEQ(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpEQ, lhs, rhs, name)
}

func (b *Builder) CreateICmpNE(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpNE, lhs, rhs, name)
}

func (b *Builder) CreateICmpSLT(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpSLT, lhs, rhs, name)
}

func (b *Builder) CreateICmpSLE(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpSLE, lhs, rhs, name)
}

func (b *Builder) CreateICmpSGT(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpSGT, lhs, rhs, name)
}

func (b *Builder) CreateICmpSGE(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpSGE, lhs, rhs, name)
}

func (b *Builder) CreateICmpULT(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpULT, lhs, rhs, name)
}

func (b *Builder) CreateICmpULE(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpULE, lhs, rhs, name)
}

func (b *Builder) CreateICmpUGT(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpUGT, lhs, rhs, name)
}

func (b *Builder) CreateICmpUGE(lhs, rhs ir.Value, name string) *ir.ICmpInst {
	return b.CreateICmp(ir.ICmpUGE, lhs, rhs, name)
}

// ============================================================================
// Other operations
// ============================================================================

func (b *Builder) CreatePhi(typ types.Type, name string) *ir.PhiInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.PhiInst{}
	inst.Self = inst
	inst.Op = ir.OpPhi
	inst.SetName(name)
	inst.SetType(typ)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateSelect(cond ir.Value, trueVal, falseVal ir.Value, name string) *ir.SelectInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.SelectInst{}
	inst.Self = inst
	inst.Op = ir.OpSelect
	inst.SetName(name)
	inst.SetType(trueVal.Type())
	inst.SetOperand(0, cond)
	inst.SetOperand(1, trueVal)
	inst.SetOperand(2, falseVal)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateCall(fn *ir.Function, args []ir.Value, name string) *ir.CallInst {
	if name == "" && fn.FuncType.ReturnType.Kind() != types.VoidKind {
		name = b.generateName()
	}
	inst := &ir.CallInst{
		Callee:    fn,
		CalleeVal: fn,
		CallConv:  fn.CallConv,
	}
	inst.Self = inst
	inst.Op = ir.OpCall
	inst.SetName(name)
	inst.SetType(fn.FuncType.ReturnType)
	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateIndirectCall(callee ir.Value, args []ir.Value, name string) *ir.CallInst {
	var retType types.Type = types.Void
	if ptr, ok := callee.Type().(*types.PointerType); ok {
		if fn, ok := ptr.ElementType.(*types.FunctionType); ok {
			retType = fn.ReturnType
		}
	}

	if name == "" && retType.Kind() != types.VoidKind {
		name = b.generateName()
	}

	inst := &ir.CallInst{
		CalleeVal: callee,
		CallConv:  ir.CC_C,
	}
	inst.Self = inst
	inst.Op = ir.OpCall
	inst.SetName(name)
	inst.SetType(retType)
	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateCallByName(name string, retType types.Type, args []ir.Value, resultName string) *ir.CallInst {
	if resultName == "" && retType.Kind() != types.VoidKind {
		resultName = b.generateName()
	}
	inst := &ir.CallInst{
		CalleeName: name,
	}
	inst.Self = inst
	inst.Op = ir.OpCall
	inst.SetName(resultName)
	inst.SetType(retType)
	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateSyscall(args []ir.Value) *ir.SyscallInst {
	if b.currentBlock == nil {
		panic("no insertion block set")
	}
	name := b.generateName()
	inst := &ir.SyscallInst{}
	inst.Self = inst
	inst.Op = ir.OpSyscall
	inst.SetName(name)
	inst.SetType(types.I64)

	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateExtractValue(agg ir.Value, indices []int, name string) *ir.ExtractValueInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.ExtractValueInst{
		Indices: indices,
	}
	inst.Self = inst
	inst.Op = ir.OpExtractValue
	inst.SetName(name)

	typ := agg.Type()
	for _, idx := range indices {
		if st, ok := typ.(*types.StructType); ok {
			if idx >= 0 && idx < len(st.Fields) {
				typ = st.Fields[idx]
			}
		} else if at, ok := typ.(*types.ArrayType); ok {
			typ = at.ElementType
		}
	}
	inst.SetType(typ)

	inst.SetOperand(0, agg)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateInsertValue(agg ir.Value, val ir.Value, indices []int, name string) *ir.InsertValueInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.InsertValueInst{
		Indices: indices,
	}
	inst.Self = inst
	inst.Op = ir.OpInsertValue
	inst.SetName(name)
	inst.SetType(agg.Type())
	inst.SetOperand(0, agg)
	inst.SetOperand(1, val)
	b.insert(inst)
	return inst
}

// ============================================================================
// Constant creation
// ============================================================================

func (b *Builder) ConstInt(typ *types.IntType, val int64) *ir.ConstantInt {
	c := &ir.ConstantInt{
		Value: val,
	}
	c.SetType(typ)
	return c
}

func (b *Builder) ConstFloat(typ *types.FloatType, val float64) *ir.ConstantFloat {
	c := &ir.ConstantFloat{
		Value: val,
	}
	c.SetType(typ)
	return c
}

func (b *Builder) ConstNull(ptrType *types.PointerType) *ir.ConstantNull {
	c := &ir.ConstantNull{}
	c.SetType(ptrType)
	return c
}

func (b *Builder) ConstUndef(typ types.Type) *ir.ConstantUndef {
	c := &ir.ConstantUndef{}
	c.SetType(typ)
	return c
}

func (b *Builder) ConstZero(typ types.Type) *ir.ConstantZero {
	c := &ir.ConstantZero{}
	c.SetType(typ)
	return c
}

func (b *Builder) True() *ir.ConstantInt {
	return b.ConstInt(types.I1, 1)
}

func (b *Builder) False() *ir.ConstantInt {
	return b.ConstInt(types.I1, 0)
}

// ============================================================================
// Intrinsic operations
// ============================================================================

func (b *Builder) CreateVaStart(vaList ir.Value) *ir.VaStartInst {
	inst := &ir.VaStartInst{}
	inst.Self = inst
	inst.Op = ir.OpVaStart
	inst.SetOperand(0, vaList)
	inst.SetType(types.Void)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateVaArg(vaList ir.Value, argType types.Type, name string) *ir.VaArgInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.VaArgInst{
		ArgType: argType,
	}
	inst.Self = inst
	inst.Op = ir.OpVaArg
	inst.SetName(name)
	inst.SetType(argType)
	inst.SetOperand(0, vaList)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateVaEnd(vaList ir.Value) *ir.VaEndInst {
	inst := &ir.VaEndInst{}
	inst.Self = inst
	inst.Op = ir.OpVaEnd
	inst.SetOperand(0, vaList)
	inst.SetType(types.Void)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateSizeOf(typ types.Type, name string) *ir.SizeOfInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.SizeOfInst{
		QueryType: typ,
	}
	inst.Self = inst
	inst.Op = ir.OpSizeOf
	inst.SetName(name)
	inst.SetType(types.U64)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateAlignOf(typ types.Type, name string) *ir.AlignOfInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.AlignOfInst{
		QueryType: typ,
	}
	inst.Self = inst
	inst.Op = ir.OpAlignOf
	inst.SetName(name)
	inst.SetType(types.U64)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateMemSet(dest ir.Value, val ir.Value, count ir.Value) *ir.MemSetInst {
	inst := &ir.MemSetInst{}
	inst.Self = inst
	inst.Op = ir.OpMemSet
	inst.SetType(types.Void)
	inst.SetOperand(0, dest)
	inst.SetOperand(1, val)
	inst.SetOperand(2, count)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateMemCpy(dest ir.Value, src ir.Value, count ir.Value) *ir.MemCpyInst {
	inst := &ir.MemCpyInst{}
	inst.Self = inst
	inst.Op = ir.OpMemCpy
	inst.SetType(types.Void)
	inst.SetOperand(0, dest)
	inst.SetOperand(1, src)
	inst.SetOperand(2, count)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateMemMove(dest ir.Value, src ir.Value, count ir.Value) *ir.MemMoveInst {
	inst := &ir.MemMoveInst{}
	inst.Self = inst
	inst.Op = ir.OpMemMove
	inst.SetType(types.Void)
	inst.SetOperand(0, dest)
	inst.SetOperand(1, src)
	inst.SetOperand(2, count)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateStrLen(str ir.Value, name string) *ir.StrLenInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.StrLenInst{}
	inst.Self = inst
	inst.Op = ir.OpStrLen
	inst.SetName(name)
	inst.SetType(types.U64)
	inst.SetOperand(0, str)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateMemChr(ptr ir.Value, val ir.Value, count ir.Value, name string) *ir.MemChrInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.MemChrInst{}
	inst.Self = inst
	inst.Op = ir.OpMemChr
	inst.SetName(name)
	inst.SetType(types.NewPointer(types.Void))
	inst.SetOperand(0, ptr)
	inst.SetOperand(1, val)
	inst.SetOperand(2, count)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateMemCmp(ptr1 ir.Value, ptr2 ir.Value, count ir.Value, name string) *ir.MemCmpInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.MemCmpInst{}
	inst.Self = inst
	inst.Op = ir.OpMemCmp
	inst.SetName(name)
	inst.SetType(types.I32)
	inst.SetOperand(0, ptr1)
	inst.SetOperand(1, ptr2)
	inst.SetOperand(2, count)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateRaise(message ir.Value) *ir.RaiseInst {
	inst := &ir.RaiseInst{}
	inst.Self = inst
	inst.Op = ir.OpRaise
	inst.SetType(types.Void)
	inst.SetOperand(0, message)
	b.insert(inst)
	return inst
}

// ============================================================================
// Async Task operations
// ============================================================================

func (b *Builder) CreateAsyncTask(fn *ir.Function, args []ir.Value, name string) *ir.AsyncTaskCreateInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.AsyncTaskCreateInst{
		Callee: fn,
	}
	inst.Self = inst
	inst.Op = ir.OpAsyncTaskCreate
	inst.SetName(name)
	inst.SetType(types.NewPointer(types.I8))
	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}

func (b *Builder) CreateAwaitTask(handle ir.Value, resultType types.Type, name string) *ir.AsyncTaskAwaitInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.AsyncTaskAwaitInst{}
	inst.Self = inst
	inst.Op = ir.OpAsyncTaskAwait
	inst.SetName(name)
	inst.SetType(resultType)
	inst.SetOperand(0, handle)
	b.insert(inst)
	return inst
}

func (b *Builder) CreateProcess(fn *ir.Function, args []ir.Value, name string) *ir.ProcessCreateInst {
	if name == "" {
		name = b.generateName()
	}
	inst := &ir.ProcessCreateInst{
		Callee: fn,
	}
	inst.Self = inst
	inst.Op = ir.OpProcessCreate
	inst.SetName(name)
	inst.SetType(types.I32)
	for i, arg := range args {
		inst.SetOperand(i, arg)
	}
	b.insert(inst)
	return inst
}