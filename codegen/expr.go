// codegen/expr.go
package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

func (cg *Codegen) genExpr(expr ast.Expr) ir.Value {
	if expr == nil {
		return nil
	}
	switch e := expr.(type) {
	case *ast.BasicLit:
		return cg.genBasicLit(e)
	case *ast.Ident:
		switch e.Name {
		case "true":
			return cg.Builder.True()
		case "false":
			return cg.Builder.False()
		case "null":
			return cg.Builder.ConstNull(types.NewPointer(types.Void))
		}
		alloca := cg.lookupVar(e.Name)
		if alloca != nil {
			pt := alloca.Type().(*types.PointerType)
			return cg.Builder.CreateLoad(pt.ElementType, alloca, e.Name)
		}
		if fn := cg.Module.GetFunction(e.Name); fn != nil {
			return fn
		}
		if g := cg.Module.GetGlobal(e.Name); g != nil {
			pt := g.Type().(*types.PointerType)
			return cg.Builder.CreateLoad(pt.ElementType, g, e.Name)
		}
		panic(fmt.Sprintf("codegen: undefined identifier %q", e.Name))
	case *ast.SelectorExpr:
		return cg.genSelector(e)
	case *ast.IndexExpr:
		return cg.genIndex(e)
	case *ast.UnaryExpr:
		return cg.genUnary(e)
	case *ast.BinaryExpr:
		return cg.genBinary(e)
	case *ast.CallExpr:
		return cg.genCall(e)
	case *ast.CastExpr:
		return cg.genCast(e)
	case *ast.CompositeLit:
		return cg.genCompositeLit(e)
	case *ast.TupleLit:
		return cg.genTupleLit(e)
	case *ast.NewExpr:
		return cg.genNew(e)
	case *ast.NewArrayExpr:
		return cg.genNewArray(e)
	case *ast.DeleteExpr:
		return cg.genDelete(e)
	case *ast.RangeExpr:
		return cg.genExpr(e.Low)
	case *ast.AwaitExpr:
		handle := cg.genExpr(e.X)
		inst := cg.Builder.CreateAwaitTask(handle, types.Void, "await")
		return inst
	}
	return nil
}

func (cg *Codegen) genBasicLit(e *ast.BasicLit) ir.Value {
	switch e.Kind {
	case "INT":
		v, err := strconv.ParseInt(e.Value, 0, 64)
		if err != nil {
			v = 0
		}
		return cg.Builder.ConstInt(types.I32, v)
	case "HEX":
		raw := strings.TrimPrefix(e.Value, "0x")
		raw = strings.TrimPrefix(raw, "0X")
		v, err := strconv.ParseInt(raw, 16, 64)
		if err != nil {
			v = 0
		}
		return cg.Builder.ConstInt(types.I64, v)
	case "FLOAT":
		v, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			v = 0
		}
		return cg.Builder.ConstFloat(types.F64, v)
	case "STRING":
		s := e.Value
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		s = processEscapes(s)
		return cg.createGlobalString(s)
	case "CHAR":
		s := e.Value
		if len(s) >= 2 && s[0] == '\'' {
			s = s[1 : len(s)-1]
		}
		s = processEscapes(s)
		r := rune(0)
		if len(s) > 0 {
			r = rune(s[0])
		}
		return cg.Builder.ConstInt(types.I32, int64(r))
	case "BOOL":
		if e.Value == "true" {
			return cg.Builder.True()
		}
		return cg.Builder.False()
	case "NULL":
		return cg.Builder.ConstNull(types.NewPointer(types.Void))
	}
	return nil
}

func processEscapes(s string) string {
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\r`, "\r")
	s = strings.ReplaceAll(s, `\\`, "\\")
	s = strings.ReplaceAll(s, `\"`, "\"")
	s = strings.ReplaceAll(s, `\'`, "'")
	s = strings.ReplaceAll(s, `\0`, "\x00")
	return s
}

func (cg *Codegen) genSelector(e *ast.SelectorExpr) ir.Value {
	basePtr := cg.genLValue(e.X)
	if basePtr != nil {
		pt, ok := basePtr.Type().(*types.PointerType)
		if ok {
			if st, ok := pt.ElementType.(*types.StructType); ok {
				idx := cg.TypeGen.FieldIndex(st.Name, e.Sel)
				if idx >= 0 {
					fieldPtr := cg.Builder.CreateStructGEP(st, basePtr, idx, e.Sel+".ptr")
					fpt := fieldPtr.Type().(*types.PointerType)
					return cg.Builder.CreateLoad(fpt.ElementType, fieldPtr, e.Sel)
				}
			}
		}
	}
	qualName := fmt.Sprintf("%s.%s", cg.exprName(e.X), e.Sel)
	if fn := cg.Module.GetFunction(qualName); fn != nil {
		return fn
	}
	return nil
}

func (cg *Codegen) exprName(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		return cg.exprName(v.X) + "." + v.Sel
	}
	return ""
}

func (cg *Codegen) genIndex(e *ast.IndexExpr) ir.Value {
	ptr := cg.genLValue(e)
	if ptr == nil {
		return nil
	}
	pt, ok := ptr.Type().(*types.PointerType)
	if !ok {
		return nil
	}
	return cg.Builder.CreateLoad(pt.ElementType, ptr, "elem")
}

func (cg *Codegen) genUnary(e *ast.UnaryExpr) ir.Value {
	x := cg.genExpr(e.X)
	if x == nil {
		return nil
	}
	switch e.Op {
	case "-":
		zero := cg.Builder.ConstInt(types.I32, 0)
		if _, ok := x.Type().(*types.FloatType); ok {
			return cg.Builder.CreateFSub(cg.Builder.ConstFloat(x.Type().(*types.FloatType), 0), x, "neg")
		}
		return cg.Builder.CreateSub(zero, x, "neg")
	case "!":
		zero := cg.Builder.ConstInt(types.I1, 0)
		return cg.Builder.CreateICmpEQ(x, zero, "not")
	case "~":
		allOnes := cg.Builder.ConstInt(types.I32, -1)
		return cg.Builder.CreateXor(x, allOnes, "bitnot")
	case "&":
		return cg.genLValue(e.X)
	case "++", "--":
		// Postfix increment/decrement: load, modify, store, return original value
		ptr := cg.genLValue(e.X)
		if ptr == nil {
			return nil
		}
		pt := ptr.Type().(*types.PointerType)
		cur := cg.Builder.CreateLoad(pt.ElementType, ptr, "")
		one := cg.Builder.ConstInt(types.I32, 1)
		if it, ok := pt.ElementType.(*types.IntType); ok && it.BitWidth != 32 {
			one = &ir.ConstantInt{}
			one.SetType(pt.ElementType)
			one.Value = 1
		}
		var next ir.Value
		if e.Op == "++" {
			next = cg.Builder.CreateAdd(cur, one, "")
		} else {
			next = cg.Builder.CreateSub(cur, one, "")
		}
		cg.Builder.CreateStore(next, ptr)
		return cur // Return original value for postfix operators
	}
	return x
}

func (cg *Codegen) genBinary(e *ast.BinaryExpr) ir.Value {
	if e.Op == "&&" {
		return cg.genLogicalAnd(e)
	}
	if e.Op == "||" {
		return cg.genLogicalOr(e)
	}
	lhs := cg.genExpr(e.Left)
	rhs := cg.genExpr(e.Right)
	if lhs == nil || rhs == nil {
		return nil
	}
	isFloat := lhs.Type().Kind() == types.FloatKind
	switch e.Op {
	case "+":
		if isFloat {
			return cg.Builder.CreateFAdd(lhs, rhs, "")
		}
		return cg.Builder.CreateAdd(lhs, rhs, "")
	case "-":
		if isFloat {
			return cg.Builder.CreateFSub(lhs, rhs, "")
		}
		return cg.Builder.CreateSub(lhs, rhs, "")
	case "*":
		if isFloat {
			return cg.Builder.CreateFMul(lhs, rhs, "")
		}
		return cg.Builder.CreateMul(lhs, rhs, "")
	case "/":
		if isFloat {
			return cg.Builder.CreateFDiv(lhs, rhs, "")
		}
		if it, ok := lhs.Type().(*types.IntType); ok && !it.Signed {
			return cg.Builder.CreateUDiv(lhs, rhs, "")
		}
		return cg.Builder.CreateSDiv(lhs, rhs, "")
	case "%":
		if it, ok := lhs.Type().(*types.IntType); ok && !it.Signed {
			return cg.Builder.CreateURem(lhs, rhs, "")
		}
		return cg.Builder.CreateSRem(lhs, rhs, "")
	case "&":
		return cg.Builder.CreateAnd(lhs, rhs, "")
	case "|":
		return cg.Builder.CreateOr(lhs, rhs, "")
	case "^":
		return cg.Builder.CreateXor(lhs, rhs, "")
	case "<<":
		return cg.Builder.CreateShl(lhs, rhs, "")
	case ">>":
		if it, ok := lhs.Type().(*types.IntType); ok && !it.Signed {
			return cg.Builder.CreateLShr(lhs, rhs, "")
		}
		return cg.Builder.CreateAShr(lhs, rhs, "")
	case "==":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpOEQ, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpEQ(lhs, rhs, "")
	case "!=":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpONE, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpNE(lhs, rhs, "")
	case "<":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpOLT, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpSLT(lhs, rhs, "")
	case "<=":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpOLE, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpSLE(lhs, rhs, "")
	case ">":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpOGT, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpSGT(lhs, rhs, "")
	case ">=":
		if isFloat {
			return cg.Builder.CreateFCmp(ir.FCmpOGE, lhs, rhs, "")
		}
		return cg.Builder.CreateICmpSGE(lhs, rhs, "")
	}
	return nil
}

func (cg *Codegen) genLogicalAnd(e *ast.BinaryExpr) ir.Value {
	lhs := cg.genExpr(e.Left)
	if lhs.Type().BitSize() != 1 {
		lhs = cg.Builder.CreateICmpNE(lhs, cg.Builder.ConstInt(types.I32, 0), "")
	}
	rhsBlock := cg.Builder.CreateBlock("and.rhs")
	mergeBlock := cg.Builder.CreateBlock("and.merge")
	shortBlock := cg.Builder.CreateBlock("and.short")
	cg.Builder.CreateCondBr(lhs, rhsBlock, shortBlock)
	cg.Builder.SetInsertPoint(shortBlock)
	cg.Builder.CreateBr(mergeBlock)
	cg.Builder.SetInsertPoint(rhsBlock)
	rhs := cg.genExpr(e.Right)
	if rhs.Type().BitSize() != 1 {
		rhs = cg.Builder.CreateICmpNE(rhs, cg.Builder.ConstInt(types.I32, 0), "")
	}
	rhsDone := cg.Builder.CurrentBlock()
	cg.Builder.CreateBr(mergeBlock)
	cg.Builder.SetInsertPoint(mergeBlock)
	phi := cg.Builder.CreatePhi(types.I1, "and")
	phi.AddIncoming(cg.Builder.False(), shortBlock)
	phi.AddIncoming(rhs, rhsDone)
	return phi
}

func (cg *Codegen) genLogicalOr(e *ast.BinaryExpr) ir.Value {
	lhs := cg.genExpr(e.Left)
	if lhs.Type().BitSize() != 1 {
		lhs = cg.Builder.CreateICmpNE(lhs, cg.Builder.ConstInt(types.I32, 0), "")
	}
	rhsBlock := cg.Builder.CreateBlock("or.rhs")
	mergeBlock := cg.Builder.CreateBlock("or.merge")
	shortBlock := cg.Builder.CreateBlock("or.short")
	cg.Builder.CreateCondBr(lhs, shortBlock, rhsBlock)
	cg.Builder.SetInsertPoint(shortBlock)
	cg.Builder.CreateBr(mergeBlock)
	cg.Builder.SetInsertPoint(rhsBlock)
	rhs := cg.genExpr(e.Right)
	if rhs.Type().BitSize() != 1 {
		rhs = cg.Builder.CreateICmpNE(rhs, cg.Builder.ConstInt(types.I32, 0), "")
	}
	rhsDone := cg.Builder.CurrentBlock()
	cg.Builder.CreateBr(mergeBlock)
	cg.Builder.SetInsertPoint(mergeBlock)
	phi := cg.Builder.CreatePhi(types.I1, "or")
	phi.AddIncoming(cg.Builder.True(), shortBlock)
	phi.AddIncoming(rhs, rhsDone)
	return phi
}

func (cg *Codegen) genCall(e *ast.CallExpr) ir.Value {
	args := make([]ir.Value, 0, len(e.Args))
	for _, arg := range e.Args {
		v := cg.genExpr(arg)
		if v != nil {
			args = append(args, v)
		}
	}
	fnName := cg.exprName(e.Fun)
	switch fnName {
	case "decref":
		if len(args) > 0 {
			return cg.Builder.CreateCallByName("arc_decref", types.Void, args, "")
		}
		return nil
	case "syscall_spawn":
		if len(args) >= 1 {
			if fn, ok := args[0].(*ir.Function); ok {
				return cg.Builder.CreateAsyncTask(fn, args[1:], "handle")
			}
		}
		return nil
	case "thread_join":
		if len(args) == 1 {
			return cg.Builder.CreateAwaitTask(args[0], types.Void, "join")
		}
		return nil
	case "thread_exit":
		cg.Builder.CreateRetVoid()
		return nil
	case "sizeof":
		if len(args) >= 1 {
			return cg.Builder.CreateSizeOf(args[0].Type(), "")
		}
		return nil
	case "alignof":
		if len(args) >= 1 {
			return cg.Builder.CreateAlignOf(args[0].Type(), "")
		}
		return nil
	}
	if irFn := cg.Module.GetFunction(fnName); irFn != nil {
		call := cg.Builder.CreateCall(irFn, args, "")
		if irFn.FuncType.ReturnType.Kind() == types.VoidKind {
			return nil
		}
		return call
	}
	calleeVal := cg.genExpr(e.Fun)
	if calleeVal != nil {
		call := cg.Builder.CreateIndirectCall(calleeVal, args, "")
		if call.Type().Kind() == types.VoidKind {
			return nil
		}
		return call
	}
	panic(fmt.Sprintf("codegen: unresolved callee %q", fnName))
}

func (cg *Codegen) genCast(e *ast.CastExpr) ir.Value {
	src := cg.genExpr(e.X)
	dest := cg.TypeGen.GenType(e.Type)
	if src == nil || dest == nil {
		return nil
	}
	return cg.emitCast(src, dest)
}

func (cg *Codegen) emitCast(src ir.Value, dest types.Type) ir.Value {
	srcType := src.Type()
	if srcType.Equal(dest) {
		return src
	}
	srcBits := srcType.BitSize()
	dstBits := dest.BitSize()
	switch {
	case srcType.Kind() == types.PointerKind && dest.Kind() == types.PointerKind:
		return cg.Builder.CreateBitCast(src, dest, "")
	case srcType.Kind() == types.PointerKind && dest.Kind() == types.IntegerKind:
		return cg.Builder.CreatePtrToInt(src, dest, "")
	case srcType.Kind() == types.IntegerKind && dest.Kind() == types.PointerKind:
		return cg.Builder.CreateIntToPtr(src, dest, "")
	case srcType.Kind() == types.IntegerKind && dest.Kind() == types.IntegerKind:
		if dstBits < srcBits {
			return cg.Builder.CreateTrunc(src, dest, "")
		}
		if it, ok := srcType.(*types.IntType); ok && it.Signed {
			return cg.Builder.CreateSExt(src, dest, "")
		}
		return cg.Builder.CreateZExt(src, dest, "")
	case srcType.Kind() == types.FloatKind && dest.Kind() == types.FloatKind:
		if dstBits < srcBits {
			return cg.Builder.CreateFPTrunc(src, dest, "")
		}
		return cg.Builder.CreateFPExt(src, dest, "")
	case srcType.Kind() == types.IntegerKind && dest.Kind() == types.FloatKind:
		if it, ok := srcType.(*types.IntType); ok && it.Signed {
			return cg.Builder.CreateSIToFP(src, dest, "")
		}
		return cg.Builder.CreateUIToFP(src, dest, "")
	case srcType.Kind() == types.FloatKind && dest.Kind() == types.IntegerKind:
		if it, ok := dest.(*types.IntType); ok && it.Signed {
			return cg.Builder.CreateFPToSI(src, dest, "")
		}
		return cg.Builder.CreateFPToUI(src, dest, "")
	default:
		return cg.Builder.CreateBitCast(src, dest, "")
	}
}

func (cg *Codegen) genCompositeLit(e *ast.CompositeLit) ir.Value {
	if e.Type == nil {
		return nil
	}
	irType := cg.TypeGen.GenType(e.Type)

	if st, ok := irType.(*types.StructType); ok {
		alloca := cg.Builder.CreateAlloca(st, "composite")
		cg.Builder.CreateStore(cg.Builder.ConstZero(st), alloca)
		for _, field := range e.Fields {
			switch f := field.(type) {
			case *ast.KeyValueExpr:
				keyName := ""
				if id, ok := f.Key.(*ast.Ident); ok {
					keyName = id.Name
				}
				idx := cg.TypeGen.FieldIndex(st.Name, keyName)
				if idx < 0 {
					continue
				}
				val := cg.genExpr(f.Value)
				if val == nil {
					continue
				}
				fieldPtr := cg.Builder.CreateStructGEP(st, alloca, idx, keyName+".ptr")
				cg.Builder.CreateStore(val, fieldPtr)
			}
		}
		return cg.Builder.CreateLoad(st, alloca, "composite.val")
	}

	if at, ok := irType.(*types.ArrayType); ok {
		alloca := cg.Builder.CreateAlloca(at, "array.lit")
		cg.Builder.CreateStore(cg.Builder.ConstZero(at), alloca)
		for i, field := range e.Fields {
			var valExpr ast.Expr = field
			if kv, ok := field.(*ast.KeyValueExpr); ok {
				valExpr = kv.Value
			}
			val := cg.genExpr(valExpr)
			if val == nil {
				continue
			}
			val = cg.emitCast(val, at.ElementType)
			zero := cg.Builder.ConstInt(types.I32, 0)
			idx := cg.Builder.ConstInt(types.I32, int64(i))
			elemPtr := cg.Builder.CreateInBoundsGEP(at, alloca, []ir.Value{zero, idx}, fmt.Sprintf("elem.%d.ptr", i))
			cg.Builder.CreateStore(val, elemPtr)
		}
		return cg.Builder.CreateLoad(at, alloca, "array.val")
	}
	return nil
}

func (cg *Codegen) genTupleLit(e *ast.TupleLit) ir.Value {
	fieldTypes := make([]types.Type, len(e.Elems))
	vals := make([]ir.Value, len(e.Elems))
	for i, el := range e.Elems {
		v := cg.genExpr(el)
		if v == nil {
			return nil
		}
		vals[i] = v
		fieldTypes[i] = v.Type()
	}
	st := types.NewStruct("", fieldTypes, false)
	var agg ir.Value = cg.Builder.ConstUndef(st)
	for i, v := range vals {
		agg = cg.Builder.CreateInsertValue(agg, v, []int{i}, "")
	}
	return agg
}

func (cg *Codegen) genNew(e *ast.NewExpr) ir.Value {
	irType := cg.TypeGen.GenType(e.Type)
	sz := cg.Builder.CreateSizeOf(irType, "sz")
	rawPtr := cg.Builder.CreateCallByName("arc_alloc", types.NewPointer(types.Void), []ir.Value{sz}, "raw")
	ptr := cg.emitCast(rawPtr, types.NewPointer(irType))
	if e.Init != nil {
		if st, ok := irType.(*types.StructType); ok {
			for _, field := range e.Init.Fields {
				if kv, ok := field.(*ast.KeyValueExpr); ok {
					keyName := ""
					if id, ok := kv.Key.(*ast.Ident); ok {
						keyName = id.Name
					}
					idx := cg.TypeGen.FieldIndex(st.Name, keyName)
					if idx < 0 {
						continue
					}
					val := cg.genExpr(kv.Value)
					if val != nil {
						fieldPtr := cg.Builder.CreateStructGEP(st, ptr, idx, keyName+".ptr")
						cg.Builder.CreateStore(val, fieldPtr)
					}
				}
			}
		}
	}
	return ptr
}

func (cg *Codegen) genNewArray(e *ast.NewArrayExpr) ir.Value {
	elemType := cg.TypeGen.GenType(e.Elem)
	count := cg.genExpr(e.Len)
	if count == nil {
		count = cg.Builder.ConstInt(types.I64, 0)
	}
	if count.Type().BitSize() < 64 {
		count = cg.emitCast(count, types.I64)
	}
	elemSz := cg.Builder.CreateSizeOf(elemType, "elem.sz")
	totalSz := cg.Builder.CreateMul(count, elemSz, "total.sz")
	rawPtr := cg.Builder.CreateCallByName("arc_alloc", types.NewPointer(types.Void), []ir.Value{totalSz}, "raw")
	return cg.emitCast(rawPtr, types.NewPointer(elemType))
}

func (cg *Codegen) genDelete(e *ast.DeleteExpr) ir.Value {
	ptr := cg.genExpr(e.X)
	if ptr != nil {
		voidPtr := cg.emitCast(ptr, types.NewPointer(types.Void))
		cg.Builder.CreateCallByName("arc_free", types.Void, []ir.Value{voidPtr}, "")
	}
	return nil
}