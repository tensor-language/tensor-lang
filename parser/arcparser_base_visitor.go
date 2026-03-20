// Code generated from ArcParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ArcParser

import "github.com/antlr4-go/antlr/v4"

type BaseArcParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseArcParserVisitor) VisitCompilationUnit(ctx *CompilationUnitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitNamespaceDecl(ctx *NamespaceDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSemi(ctx *SemiContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitImportAlias(ctx *ImportAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitConstDecl(ctx *ConstDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitConstSpec(ctx *ConstSpecContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTopLevelVarDecl(ctx *TopLevelVarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTopLevelLetDecl(ctx *TopLevelLetDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFuncModifier(ctx *FuncModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitDeinitDecl(ctx *DeinitDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitParamList(ctx *ParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSelfParam(ctx *SelfParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitParamType(ctx *ParamTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitReturnType(ctx *ReturnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTupleType(ctx *TupleTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitGenericArgs(ctx *GenericArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitInterfaceDecl(ctx *InterfaceDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitInterfaceField(ctx *InterfaceFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitEnumMember(ctx *EnumMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAttribute(ctx *AttributeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTypeRef(ctx *TypeRefContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBaseType(ctx *BaseTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTypeList(ctx *TypeListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternDecl(ctx *ExternDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternMember(ctx *ExternMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternFuncDecl(ctx *ExternFuncDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitCallingConvention(ctx *CallingConventionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternSymbol(ctx *ExternSymbolContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternParamList(ctx *ExternParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternParam(ctx *ExternParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternReturnType(ctx *ExternReturnTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternType(ctx *ExternTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternNamespace(ctx *ExternNamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternClass(ctx *ExternClassContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternClassMember(ctx *ExternClassMemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternVirtualMethod(ctx *ExternVirtualMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternStaticMethod(ctx *ExternStaticMethodContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternConstructor(ctx *ExternConstructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternDestructor(ctx *ExternDestructorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternMethodParamList(ctx *ExternMethodParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternMethodParam(ctx *ExternMethodParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternTypeAlias(ctx *ExternTypeAliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExternFunctionPtrType(ctx *ExternFunctionPtrTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBlock(ctx *BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLetStatement(ctx *LetStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitVarStatement(ctx *VarStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitForStatement(ctx *ForStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitForHeader(ctx *ForHeaderContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitForInit(ctx *ForInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitForPost(ctx *ForPostContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitForIterator(ctx *ForIteratorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSwitchDefault(ctx *SwitchDefaultContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAssignmentTarget(ctx *AssignmentTargetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAssignOp(ctx *AssignOpContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitMulExpr(ctx *MulExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBitwiseAndExpr(ctx *BitwiseAndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBitwiseOrExpr(ctx *BitwiseOrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitPostDecrement(ctx *PostDecrementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitUnaryMinus(ctx *UnaryMinusContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAddExpr(ctx *AddExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitRelationalExpr(ctx *RelationalExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitRangeExpr(ctx *RangeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitIndexExpr(ctx *IndexExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAwaitExpr(ctx *AwaitExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitEqualityExpr(ctx *EqualityExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitMemberAccess(ctx *MemberAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitAddressOf(ctx *AddressOfContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitSliceExpr(ctx *SliceExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitCallExpr(ctx *CallExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitPostIncrement(ctx *PostIncrementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBitwiseXorExpr(ctx *BitwiseXorExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBitwiseNot(ctx *BitwiseNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitShiftExpr(ctx *ShiftExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitIntLiteral(ctx *IntLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitHexLiteral(ctx *HexLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitCharLiteral(ctx *CharLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTrueLiteral(ctx *TrueLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFalseLiteral(ctx *FalseLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitNullLiteral(ctx *NullLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTypedInitExpr(ctx *TypedInitExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitBareInitExpr(ctx *BareInitExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitVectorLiteral(ctx *VectorLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitQualifiedExpr(ctx *QualifiedExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitIdentExpr(ctx *IdentExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitPrimitiveTypeExpr(ctx *PrimitiveTypeExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitParenExpr(ctx *ParenExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitTupleLiteral(ctx *TupleLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitNewExpr(ctx *NewExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitNewArrayExpr(ctx *NewArrayExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitDeleteExpr(ctx *DeleteExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLambdaExpr(ctx *LambdaExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitProcessExpr(ctx *ProcessExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitInitializerBlock(ctx *InitializerBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitFieldInit(ctx *FieldInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitMapEntry(ctx *MapEntryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitArgumentList(ctx *ArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitArgument(ctx *ArgumentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLambdaParamList(ctx *LambdaParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitLambdaParam(ctx *LambdaParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseArcParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{} {
	return v.VisitChildren(ctx)
}
