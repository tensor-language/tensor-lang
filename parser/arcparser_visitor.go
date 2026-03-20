// Code generated from ArcParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // ArcParser

import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by ArcParser.
type ArcParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ArcParser#compilationUnit.
	VisitCompilationUnit(ctx *CompilationUnitContext) interface{}

	// Visit a parse tree produced by ArcParser#namespaceDecl.
	VisitNamespaceDecl(ctx *NamespaceDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#topLevelDecl.
	VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#semi.
	VisitSemi(ctx *SemiContext) interface{}

	// Visit a parse tree produced by ArcParser#importDecl.
	VisitImportDecl(ctx *ImportDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#importSpec.
	VisitImportSpec(ctx *ImportSpecContext) interface{}

	// Visit a parse tree produced by ArcParser#importAlias.
	VisitImportAlias(ctx *ImportAliasContext) interface{}

	// Visit a parse tree produced by ArcParser#constDecl.
	VisitConstDecl(ctx *ConstDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#constSpec.
	VisitConstSpec(ctx *ConstSpecContext) interface{}

	// Visit a parse tree produced by ArcParser#topLevelVarDecl.
	VisitTopLevelVarDecl(ctx *TopLevelVarDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#topLevelLetDecl.
	VisitTopLevelLetDecl(ctx *TopLevelLetDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#funcDecl.
	VisitFuncDecl(ctx *FuncDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#funcModifier.
	VisitFuncModifier(ctx *FuncModifierContext) interface{}

	// Visit a parse tree produced by ArcParser#deinitDecl.
	VisitDeinitDecl(ctx *DeinitDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#paramList.
	VisitParamList(ctx *ParamListContext) interface{}

	// Visit a parse tree produced by ArcParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by ArcParser#selfParam.
	VisitSelfParam(ctx *SelfParamContext) interface{}

	// Visit a parse tree produced by ArcParser#paramType.
	VisitParamType(ctx *ParamTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#returnType.
	VisitReturnType(ctx *ReturnTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#tupleType.
	VisitTupleType(ctx *TupleTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#genericParams.
	VisitGenericParams(ctx *GenericParamsContext) interface{}

	// Visit a parse tree produced by ArcParser#genericArgs.
	VisitGenericArgs(ctx *GenericArgsContext) interface{}

	// Visit a parse tree produced by ArcParser#interfaceDecl.
	VisitInterfaceDecl(ctx *InterfaceDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#interfaceField.
	VisitInterfaceField(ctx *InterfaceFieldContext) interface{}

	// Visit a parse tree produced by ArcParser#enumDecl.
	VisitEnumDecl(ctx *EnumDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#enumMember.
	VisitEnumMember(ctx *EnumMemberContext) interface{}

	// Visit a parse tree produced by ArcParser#typeAliasDecl.
	VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#attribute.
	VisitAttribute(ctx *AttributeContext) interface{}

	// Visit a parse tree produced by ArcParser#typeRef.
	VisitTypeRef(ctx *TypeRefContext) interface{}

	// Visit a parse tree produced by ArcParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#baseType.
	VisitBaseType(ctx *BaseTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#primitiveType.
	VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by ArcParser#externDecl.
	VisitExternDecl(ctx *ExternDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#externMember.
	VisitExternMember(ctx *ExternMemberContext) interface{}

	// Visit a parse tree produced by ArcParser#externFuncDecl.
	VisitExternFuncDecl(ctx *ExternFuncDeclContext) interface{}

	// Visit a parse tree produced by ArcParser#callingConvention.
	VisitCallingConvention(ctx *CallingConventionContext) interface{}

	// Visit a parse tree produced by ArcParser#externSymbol.
	VisitExternSymbol(ctx *ExternSymbolContext) interface{}

	// Visit a parse tree produced by ArcParser#externParamList.
	VisitExternParamList(ctx *ExternParamListContext) interface{}

	// Visit a parse tree produced by ArcParser#externParam.
	VisitExternParam(ctx *ExternParamContext) interface{}

	// Visit a parse tree produced by ArcParser#externReturnType.
	VisitExternReturnType(ctx *ExternReturnTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#externType.
	VisitExternType(ctx *ExternTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#externNamespace.
	VisitExternNamespace(ctx *ExternNamespaceContext) interface{}

	// Visit a parse tree produced by ArcParser#externClass.
	VisitExternClass(ctx *ExternClassContext) interface{}

	// Visit a parse tree produced by ArcParser#externClassMember.
	VisitExternClassMember(ctx *ExternClassMemberContext) interface{}

	// Visit a parse tree produced by ArcParser#externVirtualMethod.
	VisitExternVirtualMethod(ctx *ExternVirtualMethodContext) interface{}

	// Visit a parse tree produced by ArcParser#externStaticMethod.
	VisitExternStaticMethod(ctx *ExternStaticMethodContext) interface{}

	// Visit a parse tree produced by ArcParser#externConstructor.
	VisitExternConstructor(ctx *ExternConstructorContext) interface{}

	// Visit a parse tree produced by ArcParser#externDestructor.
	VisitExternDestructor(ctx *ExternDestructorContext) interface{}

	// Visit a parse tree produced by ArcParser#externMethodParamList.
	VisitExternMethodParamList(ctx *ExternMethodParamListContext) interface{}

	// Visit a parse tree produced by ArcParser#externMethodParam.
	VisitExternMethodParam(ctx *ExternMethodParamContext) interface{}

	// Visit a parse tree produced by ArcParser#externTypeAlias.
	VisitExternTypeAlias(ctx *ExternTypeAliasContext) interface{}

	// Visit a parse tree produced by ArcParser#externFunctionPtrType.
	VisitExternFunctionPtrType(ctx *ExternFunctionPtrTypeContext) interface{}

	// Visit a parse tree produced by ArcParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by ArcParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by ArcParser#letStatement.
	VisitLetStatement(ctx *LetStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#varStatement.
	VisitVarStatement(ctx *VarStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#returnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#breakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#continueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#deferStatement.
	VisitDeferStatement(ctx *DeferStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#forStatement.
	VisitForStatement(ctx *ForStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#forHeader.
	VisitForHeader(ctx *ForHeaderContext) interface{}

	// Visit a parse tree produced by ArcParser#forInit.
	VisitForInit(ctx *ForInitContext) interface{}

	// Visit a parse tree produced by ArcParser#forPost.
	VisitForPost(ctx *ForPostContext) interface{}

	// Visit a parse tree produced by ArcParser#forIterator.
	VisitForIterator(ctx *ForIteratorContext) interface{}

	// Visit a parse tree produced by ArcParser#switchStatement.
	VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#switchCase.
	VisitSwitchCase(ctx *SwitchCaseContext) interface{}

	// Visit a parse tree produced by ArcParser#switchDefault.
	VisitSwitchDefault(ctx *SwitchDefaultContext) interface{}

	// Visit a parse tree produced by ArcParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by ArcParser#assignmentStatement.
	VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#assignmentTarget.
	VisitAssignmentTarget(ctx *AssignmentTargetContext) interface{}

	// Visit a parse tree produced by ArcParser#assignOp.
	VisitAssignOp(ctx *AssignOpContext) interface{}

	// Visit a parse tree produced by ArcParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by ArcParser#MulExpr.
	VisitMulExpr(ctx *MulExprContext) interface{}

	// Visit a parse tree produced by ArcParser#BitwiseAndExpr.
	VisitBitwiseAndExpr(ctx *BitwiseAndExprContext) interface{}

	// Visit a parse tree produced by ArcParser#BitwiseOrExpr.
	VisitBitwiseOrExpr(ctx *BitwiseOrExprContext) interface{}

	// Visit a parse tree produced by ArcParser#PostDecrement.
	VisitPostDecrement(ctx *PostDecrementContext) interface{}

	// Visit a parse tree produced by ArcParser#UnaryMinus.
	VisitUnaryMinus(ctx *UnaryMinusContext) interface{}

	// Visit a parse tree produced by ArcParser#AddExpr.
	VisitAddExpr(ctx *AddExprContext) interface{}

	// Visit a parse tree produced by ArcParser#RelationalExpr.
	VisitRelationalExpr(ctx *RelationalExprContext) interface{}

	// Visit a parse tree produced by ArcParser#RangeExpr.
	VisitRangeExpr(ctx *RangeExprContext) interface{}

	// Visit a parse tree produced by ArcParser#LogicalAndExpr.
	VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{}

	// Visit a parse tree produced by ArcParser#IndexExpr.
	VisitIndexExpr(ctx *IndexExprContext) interface{}

	// Visit a parse tree produced by ArcParser#LogicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by ArcParser#LogicalOrExpr.
	VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{}

	// Visit a parse tree produced by ArcParser#AwaitExpr.
	VisitAwaitExpr(ctx *AwaitExprContext) interface{}

	// Visit a parse tree produced by ArcParser#EqualityExpr.
	VisitEqualityExpr(ctx *EqualityExprContext) interface{}

	// Visit a parse tree produced by ArcParser#MemberAccess.
	VisitMemberAccess(ctx *MemberAccessContext) interface{}

	// Visit a parse tree produced by ArcParser#AddressOf.
	VisitAddressOf(ctx *AddressOfContext) interface{}

	// Visit a parse tree produced by ArcParser#PrimaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by ArcParser#SliceExpr.
	VisitSliceExpr(ctx *SliceExprContext) interface{}

	// Visit a parse tree produced by ArcParser#CallExpr.
	VisitCallExpr(ctx *CallExprContext) interface{}

	// Visit a parse tree produced by ArcParser#PostIncrement.
	VisitPostIncrement(ctx *PostIncrementContext) interface{}

	// Visit a parse tree produced by ArcParser#BitwiseXorExpr.
	VisitBitwiseXorExpr(ctx *BitwiseXorExprContext) interface{}

	// Visit a parse tree produced by ArcParser#BitwiseNot.
	VisitBitwiseNot(ctx *BitwiseNotContext) interface{}

	// Visit a parse tree produced by ArcParser#ShiftExpr.
	VisitShiftExpr(ctx *ShiftExprContext) interface{}

	// Visit a parse tree produced by ArcParser#IntLiteral.
	VisitIntLiteral(ctx *IntLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#HexLiteral.
	VisitHexLiteral(ctx *HexLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#FloatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#StringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#CharLiteral.
	VisitCharLiteral(ctx *CharLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#TrueLiteral.
	VisitTrueLiteral(ctx *TrueLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#FalseLiteral.
	VisitFalseLiteral(ctx *FalseLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#NullLiteral.
	VisitNullLiteral(ctx *NullLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#TypedInitExpr.
	VisitTypedInitExpr(ctx *TypedInitExprContext) interface{}

	// Visit a parse tree produced by ArcParser#BareInitExpr.
	VisitBareInitExpr(ctx *BareInitExprContext) interface{}

	// Visit a parse tree produced by ArcParser#VectorLiteral.
	VisitVectorLiteral(ctx *VectorLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#MapLiteral.
	VisitMapLiteral(ctx *MapLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#QualifiedExpr.
	VisitQualifiedExpr(ctx *QualifiedExprContext) interface{}

	// Visit a parse tree produced by ArcParser#IdentExpr.
	VisitIdentExpr(ctx *IdentExprContext) interface{}

	// Visit a parse tree produced by ArcParser#PrimitiveTypeExpr.
	VisitPrimitiveTypeExpr(ctx *PrimitiveTypeExprContext) interface{}

	// Visit a parse tree produced by ArcParser#ParenExpr.
	VisitParenExpr(ctx *ParenExprContext) interface{}

	// Visit a parse tree produced by ArcParser#TupleLiteral.
	VisitTupleLiteral(ctx *TupleLiteralContext) interface{}

	// Visit a parse tree produced by ArcParser#NewExpr.
	VisitNewExpr(ctx *NewExprContext) interface{}

	// Visit a parse tree produced by ArcParser#NewArrayExpr.
	VisitNewArrayExpr(ctx *NewArrayExprContext) interface{}

	// Visit a parse tree produced by ArcParser#DeleteExpr.
	VisitDeleteExpr(ctx *DeleteExprContext) interface{}

	// Visit a parse tree produced by ArcParser#LambdaExpr.
	VisitLambdaExpr(ctx *LambdaExprContext) interface{}

	// Visit a parse tree produced by ArcParser#ProcessExpr.
	VisitProcessExpr(ctx *ProcessExprContext) interface{}

	// Visit a parse tree produced by ArcParser#initializerBlock.
	VisitInitializerBlock(ctx *InitializerBlockContext) interface{}

	// Visit a parse tree produced by ArcParser#fieldInit.
	VisitFieldInit(ctx *FieldInitContext) interface{}

	// Visit a parse tree produced by ArcParser#mapEntry.
	VisitMapEntry(ctx *MapEntryContext) interface{}

	// Visit a parse tree produced by ArcParser#argumentList.
	VisitArgumentList(ctx *ArgumentListContext) interface{}

	// Visit a parse tree produced by ArcParser#argument.
	VisitArgument(ctx *ArgumentContext) interface{}

	// Visit a parse tree produced by ArcParser#lambdaParamList.
	VisitLambdaParamList(ctx *LambdaParamListContext) interface{}

	// Visit a parse tree produced by ArcParser#lambdaParam.
	VisitLambdaParam(ctx *LambdaParamContext) interface{}

	// Visit a parse tree produced by ArcParser#qualifiedName.
	VisitQualifiedName(ctx *QualifiedNameContext) interface{}
}
