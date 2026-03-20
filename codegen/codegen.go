// codegen/codegen.go
package codegen

import (
	"fmt"

	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/builder/builder"
	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

// Codegen translates a lowered Arc AST into IR.
//
// Architecture:
//   - Two passes over the file: forward-declare all functions (so mutual
//     recursion works), then emit bodies.
//   - Variables are stack-allocated with alloca + store/load (mem2reg handles
//     promotion to SSA registers in a later pass).
//   - Scopes are a slice of maps; lookupVar searches innermost-first.
//   - globalStrIdx is the only global mutable counter; all other counters live
//     inside the builder's nameCounter.
type Codegen struct {
	Builder      *builder.Builder
	Module       *ir.Module
	TypeGen      *TypeGenerator
	scopes       []map[string]ir.Value // name -> alloca pointer
	globalStrIdx int                   // counter for .str.N globals

	// loopCtx is a stack of loop control blocks pushed by genFor/genForIn.
	// break/continue pop from this.
	loopCtx []loopContext
}

type loopContext struct {
	postBlock *ir.BasicBlock // continue target
	endBlock  *ir.BasicBlock // break target
}

func New(moduleName string) *Codegen {
	b := builder.New()
	mod := b.CreateModule(moduleName)
	return &Codegen{
		Builder: b,
		Module:  mod,
		TypeGen: NewTypeGenerator(),
	}
}

// Generate is the entry point. It populates the IR module and returns it.
func (cg *Codegen) Generate(file *ast.File) (*ir.Module, error) {
	// Pass 1: forward-declare every function so call sites can resolve them
	// regardless of declaration order.
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			cg.declareFunction(d)
		case *ast.InterfaceDecl:
			cg.declareStruct(d)
		case *ast.ExternDecl:
			cg.declareExtern(d)
		}
	}

	// Pass 2: emit function bodies.
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Body == nil {
				continue // declaration-only (extern forwarded already)
			}
			if err := cg.genFuncBody(fn); err != nil {
				return nil, fmt.Errorf("codegen %s: %w", fn.Name, err)
			}
		}
	}

	return cg.Module, nil
}

// ─── Scope management ────────────────────────────────────────────────────────

func (cg *Codegen) pushScope() {
	cg.scopes = append(cg.scopes, make(map[string]ir.Value))
}

func (cg *Codegen) popScope() {
	if len(cg.scopes) == 0 {
		panic("codegen: popScope on empty scope stack")
	}
	cg.scopes = cg.scopes[:len(cg.scopes)-1]
}

// defineVar registers name → alloca in the innermost scope.
func (cg *Codegen) defineVar(name string, alloca ir.Value) {
	if len(cg.scopes) == 0 {
		panic("codegen: defineVar with no active scope")
	}
	cg.scopes[len(cg.scopes)-1][name] = alloca
}

// lookupVar walks scopes innermost-first and returns the alloca pointer, or nil.
func (cg *Codegen) lookupVar(name string) ir.Value {
	for i := len(cg.scopes) - 1; i >= 0; i-- {
		if v, ok := cg.scopes[i][name]; ok {
			return v
		}
	}
	return nil
}

// ─── Loop context ─────────────────────────────────────────────────────────────

func (cg *Codegen) pushLoop(post, end *ir.BasicBlock) {
	cg.loopCtx = append(cg.loopCtx, loopContext{postBlock: post, endBlock: end})
}

func (cg *Codegen) popLoop() {
	cg.loopCtx = cg.loopCtx[:len(cg.loopCtx)-1]
}

func (cg *Codegen) currentLoop() (loopContext, bool) {
	if len(cg.loopCtx) == 0 {
		return loopContext{}, false
	}
	return cg.loopCtx[len(cg.loopCtx)-1], true
}

// ─── Global strings ──────────────────────────────────────────────────────────

// createGlobalString creates a null-terminated [N x i8] global and returns a
// GEP pointer to its first byte (*i8). The string value should already have
// surrounding quotes stripped.
func (cg *Codegen) createGlobalString(s string) ir.Value {
	// Build element constants.
	chars := make([]ir.Constant, len(s)+1)
	for i := 0; i < len(s); i++ {
		chars[i] = cg.Builder.ConstInt(types.I8, int64(s[i]))
	}
	chars[len(s)] = cg.Builder.ConstInt(types.I8, 0) // NUL terminator

	arrType := types.NewArray(types.I8, int64(len(chars)))
	arrConst := &ir.ConstantArray{Elements: chars}
	arrConst.SetType(arrType)

	name := fmt.Sprintf(".str.%d", cg.globalStrIdx)
	cg.globalStrIdx++

	glob := cg.Builder.CreateGlobalConstant(name, arrConst)

	zero := cg.Builder.ConstInt(types.I32, 0)
	return cg.Builder.CreateInBoundsGEP(arrType, glob, []ir.Value{zero, zero}, name+".ptr")
}

// ─── Top-level declarations ───────────────────────────────────────────────────

// declareFunction emits a function declaration (no body) for use as a forward
// reference. If the function was already declared (e.g. by declareExtern), the
// existing declaration is reused.
func (cg *Codegen) declareFunction(fn *ast.FuncDecl) *ir.Function {
	if existing := cg.Module.GetFunction(fn.Name); existing != nil {
		return existing
	}

	paramTypes := make([]types.Type, len(fn.Params))
	for i, p := range fn.Params {
		paramTypes[i] = cg.TypeGen.GenType(p.Type)
	}

	retType := cg.TypeGen.GenType(fn.ReturnType)
	irFn := cg.Builder.DeclareFunction(fn.Name, retType, paramTypes, fn.IsVariadic)

	// Name the arguments so IR output is readable.
	for i, arg := range irFn.Arguments {
		if i < len(fn.Params) {
			arg.SetName(fn.Params[i].Name)
		}
	}

	// Apply calling convention for GPU kernels.
	if fn.IsGpu {
		cg.Builder.SetCallConv(irFn, ir.CC_PTX)
	}

	return irFn
}

// declareStruct registers a named struct type in the module's type table.
func (cg *Codegen) declareStruct(d *ast.InterfaceDecl) {
	fieldTypes := make([]types.Type, len(d.Fields))
	for i, f := range d.Fields {
		fieldTypes[i] = cg.TypeGen.GenType(f.Type)
	}
	st := types.NewStruct(d.Name, fieldTypes, false)
	cg.Builder.DefineStruct(st)
}

// declareExtern registers extern function declarations into the module.
func (cg *Codegen) declareExtern(d *ast.ExternDecl) {
	for _, m := range d.Members {
		ef, ok := m.(*ast.ExternFunc)
		if !ok {
			continue
		}
		if cg.Module.GetFunction(ef.Name) != nil {
			continue
		}
		paramTypes := make([]types.Type, len(ef.Params))
		for i, p := range ef.Params {
			paramTypes[i] = cg.TypeGen.GenExternType(p)
		}
		var retType types.Type = types.Void
		if ef.Return != nil {
			retType = cg.TypeGen.GenExternType(ef.Return)
		}
		irFn := cg.Builder.DeclareFunction(ef.Name, retType, paramTypes, ef.IsVariadic)

		// Map calling convention.
		switch ef.Convention {
		case "stdcall":
			cg.Builder.SetCallConv(irFn, ir.CC_StdCall)
		case "fastcall":
			cg.Builder.SetCallConv(irFn, ir.CC_FastCall)
		case "thiscall":
			cg.Builder.SetCallConv(irFn, ir.CC_ThisCall)
		case "vectorcall":
			cg.Builder.SetCallConv(irFn, ir.CC_VectorCall)
		}
	}
}

// ─── Function body emission ───────────────────────────────────────────────────

func (cg *Codegen) genFuncBody(fn *ast.FuncDecl) error {
	irFn := cg.Module.GetFunction(fn.Name)
	if irFn == nil {
		// May have been skipped above; declare now.
		irFn = cg.declareFunction(fn)
	}

	// The function must not already have a body.
	if len(irFn.Blocks) > 0 {
		return nil
	}

	entry := cg.Builder.CreateBlockInFunction("entry", irFn)
	cg.Builder.SetInsertPoint(entry)

	cg.pushScope()
	defer cg.popScope()

	// Promote parameters to alloca slots so they behave like mutable locals.
	for i, param := range fn.Params {
		argVal := irFn.Arguments[i]
		argType := argVal.Type()
		alloca := cg.Builder.CreateAlloca(argType, param.Name+".addr")
		cg.Builder.CreateStore(argVal, alloca)
		cg.defineVar(param.Name, alloca)
	}

	if err := cg.genBlock(fn.Body); err != nil {
		return err
	}

	// Seal the current block if it has no terminator.
	if cur := cg.Builder.CurrentBlock(); cur != nil && cur.Terminator() == nil {
		if irFn.FuncType.ReturnType.Kind() == types.VoidKind {
			cg.Builder.CreateRetVoid()
		} else {
			// Emit undef to keep IR valid; a later verifier pass should flag this.
			cg.Builder.CreateRet(cg.Builder.ConstUndef(irFn.FuncType.ReturnType))
		}
	}

	return nil
}