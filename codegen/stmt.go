// codegen/stmt.go
package codegen

import (
	"fmt"

	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
)

func (cg *Codegen) genBlock(block *ast.BlockStmt) error {
	cg.pushScope()
	defer cg.popScope()
	for _, stmt := range block.List {
		if err := cg.genStmt(stmt); err != nil {
			return err
		}
		if cur := cg.Builder.CurrentBlock(); cur != nil && cur.Terminator() != nil {
			break
		}
	}
	return nil
}

func (cg *Codegen) genStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.DeclStmt:
		return cg.genDeclStmt(s)
	case *ast.AssignStmt:
		fmt.Printf("DEBUG genStmt: AssignStmt target=%T op=%s value=%T\n", s.Target, s.Op, s.Value)
		return cg.genAssignStmt(s)
	case *ast.ReturnStmt:
		return cg.genReturn(s)
	case *ast.ExprStmt:
		cg.genExpr(s.X)
	case *ast.IfStmt:
		return cg.genIf(s)
	case *ast.ForStmt:
		return cg.genFor(s)
	case *ast.ForInStmt:
		return cg.genForIn(s)
	case *ast.SwitchStmt:
		return cg.genSwitch(s)
	case *ast.BreakStmt:
		lc, ok := cg.currentLoop()
		if !ok {
			return fmt.Errorf("break outside loop")
		}
		cg.Builder.CreateBr(lc.endBlock)
	case *ast.ContinueStmt:
		lc, ok := cg.currentLoop()
		if !ok {
			return fmt.Errorf("continue outside loop")
		}
		cg.Builder.CreateBr(lc.postBlock)
	case *ast.BlockStmt:
		return cg.genBlock(s)
	case *ast.DeferStmt:
		cg.genExpr(s.Call)
	default:
		fmt.Printf("DEBUG genStmt: unknown stmt type %T\n", stmt)
	}
	return nil
}

func (cg *Codegen) genDeclStmt(s *ast.DeclStmt) error {
	switch d := s.Decl.(type) {
	case *ast.VarDecl:
		return cg.genVarDecl(d)
	case *ast.ConstDecl:
		for _, spec := range d.Specs {
			if spec.Value == nil {
				continue
			}
			val := cg.genExpr(spec.Value)
			if val == nil {
				continue
			}
			alloca := cg.Builder.CreateAlloca(val.Type(), spec.Name)
			cg.Builder.CreateStore(val, alloca)
			cg.defineVar(spec.Name, alloca)
		}
	}
	return nil
}

func (cg *Codegen) genVarDecl(d *ast.VarDecl) error {
	var alloca *ir.AllocaInst
	if d.Value != nil {
		fmt.Printf("DEBUG genVarDecl: %s has Value (type=%T)\n", d.Name, d.Value)
		val := cg.genExpr(d.Value)
		if val == nil {
			fmt.Printf("DEBUG genVarDecl: genExpr returned nil for %s\n", d.Name)
			return fmt.Errorf("genVarDecl: init expr for %q produced nil value", d.Name)
		}
		fmt.Printf("DEBUG genVarDecl: genExpr returned %T for %s\n", val, d.Name)
		alloca = cg.Builder.CreateAlloca(val.Type(), d.Name)
		cg.Builder.CreateStore(val, alloca)
	} else if d.Type != nil {
		fmt.Printf("DEBUG genVarDecl: %s has no Value, using Type\n", d.Name)
		irType := cg.TypeGen.GenType(d.Type)
		alloca = cg.Builder.CreateAlloca(irType, d.Name)
		cg.Builder.CreateStore(cg.Builder.ConstZero(irType), alloca)
	} else if d.IsNull {
		fmt.Printf("DEBUG genVarDecl: %s is null\n", d.Name)
		irType := cg.TypeGen.GenType(d.Type)
		alloca = cg.Builder.CreateAlloca(irType, d.Name)
		if pt, ok := irType.(*types.PointerType); ok {
			cg.Builder.CreateStore(cg.Builder.ConstNull(pt), alloca)
		}
	} else {
		return fmt.Errorf("genVarDecl: %q has no type and no initialiser", d.Name)
	}
	cg.defineVar(d.Name, alloca)
	return nil
}

func (cg *Codegen) genAssignStmt(s *ast.AssignStmt) error {
	if s.Op == "++" || s.Op == "--" {
		ptr := cg.genLValue(s.Target)
		if ptr == nil {
			return fmt.Errorf("genAssignStmt: increment/decrement target is not an l-value")
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
		if s.Op == "++" {
			next = cg.Builder.CreateAdd(cur, one, "")
		} else {
			next = cg.Builder.CreateSub(cur, one, "")
		}
		cg.Builder.CreateStore(next, ptr)
		return nil
	}

	rhs := cg.genExpr(s.Value)
	if rhs == nil {
		return fmt.Errorf("genAssignStmt: rhs expression produced nil value")
	}

	if s.Op != "=" {
		lhsVal := cg.genExpr(s.Target)
		if lhsVal == nil {
			return fmt.Errorf("genAssignStmt: compound assignment lhs produced nil value")
		}
		rhs = cg.genCompoundOp(s.Op, lhsVal, rhs)
	}

	ptr := cg.genLValue(s.Target)
	if ptr == nil {
		return fmt.Errorf("genAssignStmt: assignment target is not an l-value")
	}
	cg.Builder.CreateStore(rhs, ptr)
	return nil
}

// In stmt.go
func (cg *Codegen) genLValue(expr ast.Expr) ir.Value {
	switch e := expr.(type) {
	case *ast.Ident:
		return cg.lookupVar(e.Name)

	case *ast.SelectorExpr:
		structPtr := cg.genLValue(e.X)
		if structPtr == nil {
			return nil
		}
		pt, ok := structPtr.Type().(*types.PointerType)
		if !ok {
			return nil
		}
		st, ok := pt.ElementType.(*types.StructType)
		if !ok {
			return nil
		}
		idx := cg.TypeGen.FieldIndex(st.Name, e.Sel)
		if idx < 0 {
			return nil
		}
		return cg.Builder.CreateStructGEP(st, structPtr, idx, e.Sel+".ptr")

	case *ast.IndexExpr:
		basePtr := cg.genLValue(e.X)
		if basePtr == nil {
			return nil
		}
		pt, ok := basePtr.Type().(*types.PointerType)
		if !ok {
			return nil
		}
		idxVal := cg.genExpr(e.Index)
		if idxVal == nil {
			return nil
		}

		// Handle array types: [N]T requires [0, idx] indexing
		if _, ok := pt.ElementType.(*types.ArrayType); ok {
			zero := cg.Builder.ConstInt(types.I32, 0)
			return cg.Builder.CreateInBoundsGEP(pt.ElementType, basePtr, []ir.Value{zero, idxVal}, "elem.ptr")
		}

		// Handle slice/vector: extract data pointer, then index into it
		if st, ok := pt.ElementType.(*types.StructType); ok && (st.Name == "slice" || st.Name == "vector") {
			dataPtrPtr := cg.Builder.CreateStructGEP(st, basePtr, 0, "data.ptr.ptr")
			dpt := dataPtrPtr.Type().(*types.PointerType)
			dataPtr := cg.Builder.CreateLoad(dpt.ElementType, dataPtrPtr, "data.ptr")
			elemPtrType := dpt.ElementType.(*types.PointerType)
			return cg.Builder.CreateInBoundsGEP(elemPtrType.ElementType, dataPtr, []ir.Value{idxVal}, "elem.ptr")
		}

		// Handle pointer-to-pointer: dereference then index
		if _, ok := pt.ElementType.(*types.PointerType); ok {
			loadedPtr := cg.Builder.CreateLoad(pt.ElementType, basePtr, "ptr.val")
			lpt := loadedPtr.Type().(*types.PointerType)
			return cg.Builder.CreateInBoundsGEP(lpt.ElementType, loadedPtr, []ir.Value{idxVal}, "elem.ptr")
		}

		// Direct pointer indexing (shouldn't normally reach here for arrays)
		return cg.Builder.CreateInBoundsGEP(pt.ElementType, basePtr, []ir.Value{idxVal}, "elem.ptr")
	}
	return nil
}

func (cg *Codegen) genCompoundOp(op string, lhs, rhs ir.Value) ir.Value {
	switch op {
	case "+=":
		return cg.Builder.CreateAdd(lhs, rhs, "")
	case "-=":
		return cg.Builder.CreateSub(lhs, rhs, "")
	case "*=":
		return cg.Builder.CreateMul(lhs, rhs, "")
	case "/=":
		return cg.Builder.CreateSDiv(lhs, rhs, "")
	case "%=":
		return cg.Builder.CreateSRem(lhs, rhs, "")
	case "&=":
		return cg.Builder.CreateAnd(lhs, rhs, "")
	case "|=":
		return cg.Builder.CreateOr(lhs, rhs, "")
	case "^=":
		return cg.Builder.CreateXor(lhs, rhs, "")
	case "<<=":
		return cg.Builder.CreateShl(lhs, rhs, "")
	case ">>=":
		return cg.Builder.CreateAShr(lhs, rhs, "")
	}
	return rhs
}

func (cg *Codegen) genReturn(s *ast.ReturnStmt) error {
	switch len(s.Results) {
	case 0:
		cg.Builder.CreateRetVoid()
	case 1:
		val := cg.genExpr(s.Results[0])
		if val == nil {
			cg.Builder.CreateRetVoid()
		} else {
			cg.Builder.CreateRet(val)
		}
	default:
		fn := cg.Builder.CurrentFunction()
		retType := fn.FuncType.ReturnType
		var agg ir.Value = cg.Builder.ConstUndef(retType)
		for i, res := range s.Results {
			val := cg.genExpr(res)
			if val != nil {
				agg = cg.Builder.CreateInsertValue(agg, val, []int{i}, "")
			}
		}
		cg.Builder.CreateRet(agg)
	}
	return nil
}

func (cg *Codegen) genIf(s *ast.IfStmt) error {
	cond := cg.genExpr(s.Cond)
	if cond == nil {
		return fmt.Errorf("if condition produced nil")
	}
	if cond.Type().BitSize() != 1 {
		zero := cg.Builder.ConstInt(types.I32, 0)
		cond = cg.Builder.CreateICmpNE(cond, zero, "cond")
	}
	thenBlock := cg.Builder.CreateBlock("if.then")
	endBlock := cg.Builder.CreateBlock("if.end")
	var elseBlock *ir.BasicBlock
	if s.Else != nil {
		elseBlock = cg.Builder.CreateBlock("if.else")
		cg.Builder.CreateCondBr(cond, thenBlock, elseBlock)
	} else {
		cg.Builder.CreateCondBr(cond, thenBlock, endBlock)
	}
	cg.Builder.SetInsertPoint(thenBlock)
	if err := cg.genBlock(s.Body); err != nil {
		return err
	}
	if cg.Builder.CurrentBlock().Terminator() == nil {
		cg.Builder.CreateBr(endBlock)
	}
	if s.Else != nil {
		cg.Builder.SetInsertPoint(elseBlock)
		switch e := s.Else.(type) {
		case *ast.BlockStmt:
			if err := cg.genBlock(e); err != nil {
				return err
			}
		case *ast.IfStmt:
			if err := cg.genIf(e); err != nil {
				return err
			}
		}
		if cg.Builder.CurrentBlock().Terminator() == nil {
			cg.Builder.CreateBr(endBlock)
		}
	}
	cg.Builder.SetInsertPoint(endBlock)
	return nil
}

func (cg *Codegen) genFor(s *ast.ForStmt) error {
	fmt.Printf("DEBUG genFor: entering\n")
	cg.pushScope()
	defer cg.popScope()
	if s.Init != nil {
		fmt.Printf("DEBUG genFor: executing init\n")
		if err := cg.genStmt(s.Init); err != nil {
			return err
		}
		fmt.Printf("DEBUG genFor: scope after init has %d entries\n", len(cg.scopes[len(cg.scopes)-1]))
	}
	condBlock := cg.Builder.CreateBlock("for.cond")
	bodyBlock := cg.Builder.CreateBlock("for.body")
	postBlock := cg.Builder.CreateBlock("for.post")
	endBlock := cg.Builder.CreateBlock("for.end")
	cg.Builder.CreateBr(condBlock)
	cg.Builder.SetInsertPoint(condBlock)
	if s.Cond != nil {
		fmt.Printf("DEBUG genFor: evaluating condition\n")
		cond := cg.genExpr(s.Cond)
		if cond == nil {
			return fmt.Errorf("for condition produced nil")
		}
		if cond.Type().BitSize() != 1 {
			zero := cg.Builder.ConstInt(types.I32, 0)
			cond = cg.Builder.CreateICmpNE(cond, zero, "cond")
		}
		cg.Builder.CreateCondBr(cond, bodyBlock, endBlock)
	} else {
		cg.Builder.CreateBr(bodyBlock)
	}
	cg.pushLoop(postBlock, endBlock)
	cg.Builder.SetInsertPoint(bodyBlock)
	fmt.Printf("DEBUG genFor: executing body, current scopes=%d\n", len(cg.scopes))
	if err := cg.genBlock(s.Body); err != nil {
		return err
	}
	cg.popLoop()
	if cg.Builder.CurrentBlock().Terminator() == nil {
		cg.Builder.CreateBr(postBlock)
	}
	cg.Builder.SetInsertPoint(postBlock)
	if s.Post != nil {
		fmt.Printf("DEBUG genFor: executing post, scopes=%d\n", len(cg.scopes))
		if err := cg.genStmt(s.Post); err != nil {
			fmt.Printf("DEBUG genFor: post statement error: %v\n", err)
			return err
		}
	}
	if cg.Builder.CurrentBlock().Terminator() == nil {
		cg.Builder.CreateBr(condBlock)
	}
	cg.Builder.SetInsertPoint(endBlock)
	return nil
}

func (cg *Codegen) genForIn(s *ast.ForInStmt) error {
	cg.pushScope()
	defer cg.popScope()
	iterVal := cg.genExpr(s.Iter)
	if iterVal == nil {
		return fmt.Errorf("for-in: iterator expression produced nil")
	}
	idxAlloca := cg.Builder.CreateAlloca(types.I64, "for.idx")
	cg.Builder.CreateStore(cg.Builder.ConstInt(types.I64, 0), idxAlloca)
	condBlock := cg.Builder.CreateBlock("forin.cond")
	bodyBlock := cg.Builder.CreateBlock("forin.body")
	postBlock := cg.Builder.CreateBlock("forin.post")
	endBlock := cg.Builder.CreateBlock("forin.end")
	cg.Builder.CreateBr(condBlock)
	cg.Builder.SetInsertPoint(condBlock)
	idxVal := cg.Builder.CreateLoad(types.I64, idxAlloca, "idx")
	var lenVal ir.Value
	switch iterVal.Type().(type) {
	case *types.StructType:
		lenVal = cg.Builder.CreateExtractValue(iterVal, []int{1}, "len")
	default:
		lenVal = cg.Builder.ConstInt(types.I64, 0)
	}
	cond := cg.Builder.CreateICmpULT(idxVal, lenVal, "forin.cond")
	cg.Builder.CreateCondBr(cond, bodyBlock, endBlock)
	cg.pushLoop(postBlock, endBlock)
	cg.Builder.SetInsertPoint(bodyBlock)
	keyAlloca := cg.Builder.CreateAlloca(types.I64, s.Key)
	cg.Builder.CreateStore(idxVal, keyAlloca)
	cg.defineVar(s.Key, keyAlloca)
	if s.Value != "" {
		dataPtr := cg.Builder.CreateExtractValue(iterVal, []int{0}, "data.ptr")
		if pt, ok := dataPtr.Type().(*types.PointerType); ok {
			elemPtr := cg.Builder.CreateInBoundsGEP(pt.ElementType, dataPtr, []ir.Value{idxVal}, "elem.ptr")
			elemAlloca := cg.Builder.CreateAlloca(pt.ElementType, s.Value)
			elemVal := cg.Builder.CreateLoad(pt.ElementType, elemPtr, "elem")
			cg.Builder.CreateStore(elemVal, elemAlloca)
			cg.defineVar(s.Value, elemAlloca)
		}
	}
	if err := cg.genBlock(s.Body); err != nil {
		return err
	}
	cg.popLoop()
	if cg.Builder.CurrentBlock().Terminator() == nil {
		cg.Builder.CreateBr(postBlock)
	}
	cg.Builder.SetInsertPoint(postBlock)
	curIdx := cg.Builder.CreateLoad(types.I64, idxAlloca, "idx.cur")
	nextIdx := cg.Builder.CreateAdd(curIdx, cg.Builder.ConstInt(types.I64, 1), "idx.next")
	cg.Builder.CreateStore(nextIdx, idxAlloca)
	cg.Builder.CreateBr(condBlock)
	cg.Builder.SetInsertPoint(endBlock)
	return nil
}

func (cg *Codegen) genSwitch(s *ast.SwitchStmt) error {
	tag := cg.genExpr(s.Tag)
	if tag == nil {
		return fmt.Errorf("switch tag produced nil")
	}
	endBlock := cg.Builder.CreateBlock("switch.end")
	var defaultBlock *ir.BasicBlock
	if s.Default != nil {
		defaultBlock = cg.Builder.CreateBlock("switch.default")
	} else {
		defaultBlock = endBlock
	}
	sw := cg.Builder.CreateSwitch(tag, defaultBlock, len(s.Cases))
	for _, c := range s.Cases {
		caseBlock := cg.Builder.CreateBlock("switch.case")
		for _, val := range c.Values {
			constVal := cg.genExpr(val)
			if ci, ok := constVal.(*ir.ConstantInt); ok {
				cg.Builder.AddCase(sw, ci, caseBlock)
			}
		}
		cg.pushLoop(endBlock, endBlock)
		cg.Builder.SetInsertPoint(caseBlock)
		for _, st := range c.Body {
			if err := cg.genStmt(st); err != nil {
				return err
			}
			if cg.Builder.CurrentBlock().Terminator() != nil {
				break
			}
		}
		cg.popLoop()
		if cg.Builder.CurrentBlock().Terminator() == nil {
			cg.Builder.CreateBr(endBlock)
		}
	}
	if s.Default != nil {
		cg.Builder.SetInsertPoint(defaultBlock)
		for _, st := range s.Default {
			if err := cg.genStmt(st); err != nil {
				return err
			}
			if cg.Builder.CurrentBlock().Terminator() != nil {
				break
			}
		}
		if cg.Builder.CurrentBlock().Terminator() == nil {
			cg.Builder.CreateBr(endBlock)
		}
	}
	cg.Builder.SetInsertPoint(endBlock)
	return nil
}

func fieldIndex(st *types.StructType, fieldName string) int {
	return -1
}