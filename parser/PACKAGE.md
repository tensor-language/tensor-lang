

CONSTANTS

const (
	ArcLexerNAMESPACE     = 1
	ArcLexerIMPORT        = 2
	ArcLexerFUNC          = 3
	ArcLexerASYNC         = 4
	ArcLexerGPU           = 5
	ArcLexerINTERFACE     = 6
	ArcLexerENUM          = 7
	ArcLexerCONST         = 8
	ArcLexerLET           = 9
	ArcLexerVAR           = 10
	ArcLexerNEW           = 11
	ArcLexerDELETE        = 12
	ArcLexerDEFER         = 13
	ArcLexerDEINIT        = 14
	ArcLexerRETURN        = 15
	ArcLexerIF            = 16
	ArcLexerELSE          = 17
	ArcLexerFOR           = 18
	ArcLexerIN            = 19
	ArcLexerSWITCH        = 20
	ArcLexerCASE          = 21
	ArcLexerDEFAULT       = 22
	ArcLexerBREAK         = 23
	ArcLexerCONTINUE      = 24
	ArcLexerPROCESS       = 25
	ArcLexerAWAIT         = 26
	ArcLexerEXTERN        = 27
	ArcLexerTYPE          = 28
	ArcLexerOPAQUE        = 29
	ArcLexerSELF          = 30
	ArcLexerMUT           = 31
	ArcLexerVOID          = 32
	ArcLexerNULL          = 33
	ArcLexerTRUE          = 34
	ArcLexerFALSE         = 35
	ArcLexerCLASS         = 36
	ArcLexerVIRTUAL       = 37
	ArcLexerSTATIC        = 38
	ArcLexerABSTRACT      = 39
	ArcLexerCDECL         = 40
	ArcLexerSTDCALL       = 41
	ArcLexerTHISCALL      = 42
	ArcLexerVECTORCALL    = 43
	ArcLexerFASTCALL      = 44
	ArcLexerINT8          = 45
	ArcLexerINT16         = 46
	ArcLexerINT32         = 47
	ArcLexerINT64         = 48
	ArcLexerUINT8         = 49
	ArcLexerUINT16        = 50
	ArcLexerUINT32        = 51
	ArcLexerUINT64        = 52
	ArcLexerUSIZE         = 53
	ArcLexerISIZE         = 54
	ArcLexerFLOAT32       = 55
	ArcLexerFLOAT64       = 56
	ArcLexerBYTE          = 57
	ArcLexerBOOL          = 58
	ArcLexerCHAR          = 59
	ArcLexerSTRING        = 60
	ArcLexerVECTOR        = 61
	ArcLexerMAP           = 62
	ArcLexerARROW         = 63
	ArcLexerELLIPSIS      = 64
	ArcLexerRANGE         = 65
	ArcLexerLSHIFT        = 66
	ArcLexerRSHIFT        = 67
	ArcLexerLE            = 68
	ArcLexerGE            = 69
	ArcLexerEQ            = 70
	ArcLexerNEQ           = 71
	ArcLexerAND           = 72
	ArcLexerOR            = 73
	ArcLexerINC           = 74
	ArcLexerDEC           = 75
	ArcLexerADD_ASSIGN    = 76
	ArcLexerSUB_ASSIGN    = 77
	ArcLexerMUL_ASSIGN    = 78
	ArcLexerDIV_ASSIGN    = 79
	ArcLexerMOD_ASSIGN    = 80
	ArcLexerAND_ASSIGN    = 81
	ArcLexerOR_ASSIGN     = 82
	ArcLexerXOR_ASSIGN    = 83
	ArcLexerSHL_ASSIGN    = 84
	ArcLexerSHR_ASSIGN    = 85
	ArcLexerLPAREN        = 86
	ArcLexerRPAREN        = 87
	ArcLexerLBRACKET      = 88
	ArcLexerRBRACKET      = 89
	ArcLexerLBRACE        = 90
	ArcLexerRBRACE        = 91
	ArcLexerDOT           = 92
	ArcLexerCOMMA         = 93
	ArcLexerCOLON         = 94
	ArcLexerSEMI          = 95
	ArcLexerASSIGN        = 96
	ArcLexerPLUS          = 97
	ArcLexerMINUS         = 98
	ArcLexerSTAR          = 99
	ArcLexerSLASH         = 100
	ArcLexerPERCENT       = 101
	ArcLexerAMP           = 102
	ArcLexerPIPE          = 103
	ArcLexerCARET         = 104
	ArcLexerTILDE         = 105
	ArcLexerBANG          = 106
	ArcLexerLT            = 107
	ArcLexerGT            = 108
	ArcLexerAT            = 109
	ArcLexerUNDERSCORE    = 110
	ArcLexerHEX_LIT       = 111
	ArcLexerFLOAT_LIT     = 112
	ArcLexerINT_LIT       = 113
	ArcLexerCHAR_LIT      = 114
	ArcLexerSTRING_LIT    = 115
	ArcLexerIDENTIFIER    = 116
	ArcLexerNL            = 117
	ArcLexerWS            = 118
	ArcLexerLINE_COMMENT  = 119
	ArcLexerBLOCK_COMMENT = 120
)
    ArcLexer tokens.

const (
	ArcParserEOF           = antlr.TokenEOF
	ArcParserNAMESPACE     = 1
	ArcParserIMPORT        = 2
	ArcParserFUNC          = 3
	ArcParserASYNC         = 4
	ArcParserGPU           = 5
	ArcParserINTERFACE     = 6
	ArcParserENUM          = 7
	ArcParserCONST         = 8
	ArcParserLET           = 9
	ArcParserVAR           = 10
	ArcParserNEW           = 11
	ArcParserDELETE        = 12
	ArcParserDEFER         = 13
	ArcParserDEINIT        = 14
	ArcParserRETURN        = 15
	ArcParserIF            = 16
	ArcParserELSE          = 17
	ArcParserFOR           = 18
	ArcParserIN            = 19
	ArcParserSWITCH        = 20
	ArcParserCASE          = 21
	ArcParserDEFAULT       = 22
	ArcParserBREAK         = 23
	ArcParserCONTINUE      = 24
	ArcParserPROCESS       = 25
	ArcParserAWAIT         = 26
	ArcParserEXTERN        = 27
	ArcParserTYPE          = 28
	ArcParserOPAQUE        = 29
	ArcParserSELF          = 30
	ArcParserMUT           = 31
	ArcParserVOID          = 32
	ArcParserNULL          = 33
	ArcParserTRUE          = 34
	ArcParserFALSE         = 35
	ArcParserCLASS         = 36
	ArcParserVIRTUAL       = 37
	ArcParserSTATIC        = 38
	ArcParserABSTRACT      = 39
	ArcParserCDECL         = 40
	ArcParserSTDCALL       = 41
	ArcParserTHISCALL      = 42
	ArcParserVECTORCALL    = 43
	ArcParserFASTCALL      = 44
	ArcParserINT8          = 45
	ArcParserINT16         = 46
	ArcParserINT32         = 47
	ArcParserINT64         = 48
	ArcParserUINT8         = 49
	ArcParserUINT16        = 50
	ArcParserUINT32        = 51
	ArcParserUINT64        = 52
	ArcParserUSIZE         = 53
	ArcParserISIZE         = 54
	ArcParserFLOAT32       = 55
	ArcParserFLOAT64       = 56
	ArcParserBYTE          = 57
	ArcParserBOOL          = 58
	ArcParserCHAR          = 59
	ArcParserSTRING        = 60
	ArcParserVECTOR        = 61
	ArcParserMAP           = 62
	ArcParserARROW         = 63
	ArcParserELLIPSIS      = 64
	ArcParserRANGE         = 65
	ArcParserLSHIFT        = 66
	ArcParserRSHIFT        = 67
	ArcParserLE            = 68
	ArcParserGE            = 69
	ArcParserEQ            = 70
	ArcParserNEQ           = 71
	ArcParserAND           = 72
	ArcParserOR            = 73
	ArcParserINC           = 74
	ArcParserDEC           = 75
	ArcParserADD_ASSIGN    = 76
	ArcParserSUB_ASSIGN    = 77
	ArcParserMUL_ASSIGN    = 78
	ArcParserDIV_ASSIGN    = 79
	ArcParserMOD_ASSIGN    = 80
	ArcParserAND_ASSIGN    = 81
	ArcParserOR_ASSIGN     = 82
	ArcParserXOR_ASSIGN    = 83
	ArcParserSHL_ASSIGN    = 84
	ArcParserSHR_ASSIGN    = 85
	ArcParserLPAREN        = 86
	ArcParserRPAREN        = 87
	ArcParserLBRACKET      = 88
	ArcParserRBRACKET      = 89
	ArcParserLBRACE        = 90
	ArcParserRBRACE        = 91
	ArcParserDOT           = 92
	ArcParserCOMMA         = 93
	ArcParserCOLON         = 94
	ArcParserSEMI          = 95
	ArcParserASSIGN        = 96
	ArcParserPLUS          = 97
	ArcParserMINUS         = 98
	ArcParserSTAR          = 99
	ArcParserSLASH         = 100
	ArcParserPERCENT       = 101
	ArcParserAMP           = 102
	ArcParserPIPE          = 103
	ArcParserCARET         = 104
	ArcParserTILDE         = 105
	ArcParserBANG          = 106
	ArcParserLT            = 107
	ArcParserGT            = 108
	ArcParserAT            = 109
	ArcParserUNDERSCORE    = 110
	ArcParserHEX_LIT       = 111
	ArcParserFLOAT_LIT     = 112
	ArcParserINT_LIT       = 113
	ArcParserCHAR_LIT      = 114
	ArcParserSTRING_LIT    = 115
	ArcParserIDENTIFIER    = 116
	ArcParserNL            = 117
	ArcParserWS            = 118
	ArcParserLINE_COMMENT  = 119
	ArcParserBLOCK_COMMENT = 120
)
    ArcParser tokens.

const (
	ArcParserRULE_compilationUnit       = 0
	ArcParserRULE_namespaceDecl         = 1
	ArcParserRULE_topLevelDecl          = 2
	ArcParserRULE_semi                  = 3
	ArcParserRULE_importDecl            = 4
	ArcParserRULE_importSpec            = 5
	ArcParserRULE_importAlias           = 6
	ArcParserRULE_constDecl             = 7
	ArcParserRULE_constSpec             = 8
	ArcParserRULE_topLevelVarDecl       = 9
	ArcParserRULE_topLevelLetDecl       = 10
	ArcParserRULE_funcDecl              = 11
	ArcParserRULE_funcModifier          = 12
	ArcParserRULE_deinitDecl            = 13
	ArcParserRULE_paramList             = 14
	ArcParserRULE_param                 = 15
	ArcParserRULE_selfParam             = 16
	ArcParserRULE_paramType             = 17
	ArcParserRULE_returnType            = 18
	ArcParserRULE_tupleType             = 19
	ArcParserRULE_genericParams         = 20
	ArcParserRULE_genericArgs           = 21
	ArcParserRULE_interfaceDecl         = 22
	ArcParserRULE_interfaceField        = 23
	ArcParserRULE_enumDecl              = 24
	ArcParserRULE_enumMember            = 25
	ArcParserRULE_typeAliasDecl         = 26
	ArcParserRULE_attribute             = 27
	ArcParserRULE_typeRef               = 28
	ArcParserRULE_functionType          = 29
	ArcParserRULE_baseType              = 30
	ArcParserRULE_primitiveType         = 31
	ArcParserRULE_typeList              = 32
	ArcParserRULE_externDecl            = 33
	ArcParserRULE_externMember          = 34
	ArcParserRULE_externFuncDecl        = 35
	ArcParserRULE_callingConvention     = 36
	ArcParserRULE_externSymbol          = 37
	ArcParserRULE_externParamList       = 38
	ArcParserRULE_externParam           = 39
	ArcParserRULE_externReturnType      = 40
	ArcParserRULE_externType            = 41
	ArcParserRULE_externNamespace       = 42
	ArcParserRULE_externClass           = 43
	ArcParserRULE_externClassMember     = 44
	ArcParserRULE_externVirtualMethod   = 45
	ArcParserRULE_externStaticMethod    = 46
	ArcParserRULE_externConstructor     = 47
	ArcParserRULE_externDestructor      = 48
	ArcParserRULE_externMethodParamList = 49
	ArcParserRULE_externMethodParam     = 50
	ArcParserRULE_externTypeAlias       = 51
	ArcParserRULE_externFunctionPtrType = 52
	ArcParserRULE_block                 = 53
	ArcParserRULE_statement             = 54
	ArcParserRULE_letStatement          = 55
	ArcParserRULE_varStatement          = 56
	ArcParserRULE_returnStatement       = 57
	ArcParserRULE_breakStatement        = 58
	ArcParserRULE_continueStatement     = 59
	ArcParserRULE_deferStatement        = 60
	ArcParserRULE_ifStatement           = 61
	ArcParserRULE_forStatement          = 62
	ArcParserRULE_forHeader             = 63
	ArcParserRULE_forInit               = 64
	ArcParserRULE_forPost               = 65
	ArcParserRULE_forIterator           = 66
	ArcParserRULE_switchStatement       = 67
	ArcParserRULE_switchCase            = 68
	ArcParserRULE_switchDefault         = 69
	ArcParserRULE_expressionList        = 70
	ArcParserRULE_assignmentStatement   = 71
	ArcParserRULE_assignmentTarget      = 72
	ArcParserRULE_assignOp              = 73
	ArcParserRULE_expressionStatement   = 74
	ArcParserRULE_expression            = 75
	ArcParserRULE_primary               = 76
	ArcParserRULE_initializerBlock      = 77
	ArcParserRULE_fieldInit             = 78
	ArcParserRULE_mapEntry              = 79
	ArcParserRULE_argumentList          = 80
	ArcParserRULE_argument              = 81
	ArcParserRULE_lambdaParamList       = 82
	ArcParserRULE_lambdaParam           = 83
	ArcParserRULE_qualifiedName         = 84
)
    ArcParser rules.


VARIABLES

var ArcLexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	ChannelNames           []string
	ModeNames              []string
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}
var ArcParserParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

FUNCTIONS

func ArcLexerInit()
    ArcLexerInit initializes any static state used to implement ArcLexer. By
    default the static state used to implement the lexer is lazily initialized
    during the first call to NewArcLexer(). You can call this function if you
    wish to initialize the static state ahead of time.

func ArcParserInit()
    ArcParserInit initializes any static state used to implement ArcParser. By
    default the static state used to implement the parser is lazily initialized
    during the first call to NewArcParser(). You can call this function if you
    wish to initialize the static state ahead of time.

func InitEmptyArgumentContext(p *ArgumentContext)
func InitEmptyArgumentListContext(p *ArgumentListContext)
func InitEmptyAssignOpContext(p *AssignOpContext)
func InitEmptyAssignmentStatementContext(p *AssignmentStatementContext)
func InitEmptyAssignmentTargetContext(p *AssignmentTargetContext)
func InitEmptyAttributeContext(p *AttributeContext)
func InitEmptyBaseTypeContext(p *BaseTypeContext)
func InitEmptyBlockContext(p *BlockContext)
func InitEmptyBreakStatementContext(p *BreakStatementContext)
func InitEmptyCallingConventionContext(p *CallingConventionContext)
func InitEmptyCompilationUnitContext(p *CompilationUnitContext)
func InitEmptyConstDeclContext(p *ConstDeclContext)
func InitEmptyConstSpecContext(p *ConstSpecContext)
func InitEmptyContinueStatementContext(p *ContinueStatementContext)
func InitEmptyDeferStatementContext(p *DeferStatementContext)
func InitEmptyDeinitDeclContext(p *DeinitDeclContext)
func InitEmptyEnumDeclContext(p *EnumDeclContext)
func InitEmptyEnumMemberContext(p *EnumMemberContext)
func InitEmptyExpressionContext(p *ExpressionContext)
func InitEmptyExpressionListContext(p *ExpressionListContext)
func InitEmptyExpressionStatementContext(p *ExpressionStatementContext)
func InitEmptyExternClassContext(p *ExternClassContext)
func InitEmptyExternClassMemberContext(p *ExternClassMemberContext)
func InitEmptyExternConstructorContext(p *ExternConstructorContext)
func InitEmptyExternDeclContext(p *ExternDeclContext)
func InitEmptyExternDestructorContext(p *ExternDestructorContext)
func InitEmptyExternFuncDeclContext(p *ExternFuncDeclContext)
func InitEmptyExternFunctionPtrTypeContext(p *ExternFunctionPtrTypeContext)
func InitEmptyExternMemberContext(p *ExternMemberContext)
func InitEmptyExternMethodParamContext(p *ExternMethodParamContext)
func InitEmptyExternMethodParamListContext(p *ExternMethodParamListContext)
func InitEmptyExternNamespaceContext(p *ExternNamespaceContext)
func InitEmptyExternParamContext(p *ExternParamContext)
func InitEmptyExternParamListContext(p *ExternParamListContext)
func InitEmptyExternReturnTypeContext(p *ExternReturnTypeContext)
func InitEmptyExternStaticMethodContext(p *ExternStaticMethodContext)
func InitEmptyExternSymbolContext(p *ExternSymbolContext)
func InitEmptyExternTypeAliasContext(p *ExternTypeAliasContext)
func InitEmptyExternTypeContext(p *ExternTypeContext)
func InitEmptyExternVirtualMethodContext(p *ExternVirtualMethodContext)
func InitEmptyFieldInitContext(p *FieldInitContext)
func InitEmptyForHeaderContext(p *ForHeaderContext)
func InitEmptyForInitContext(p *ForInitContext)
func InitEmptyForIteratorContext(p *ForIteratorContext)
func InitEmptyForPostContext(p *ForPostContext)
func InitEmptyForStatementContext(p *ForStatementContext)
func InitEmptyFuncDeclContext(p *FuncDeclContext)
func InitEmptyFuncModifierContext(p *FuncModifierContext)
func InitEmptyFunctionTypeContext(p *FunctionTypeContext)
func InitEmptyGenericArgsContext(p *GenericArgsContext)
func InitEmptyGenericParamsContext(p *GenericParamsContext)
func InitEmptyIfStatementContext(p *IfStatementContext)
func InitEmptyImportAliasContext(p *ImportAliasContext)
func InitEmptyImportDeclContext(p *ImportDeclContext)
func InitEmptyImportSpecContext(p *ImportSpecContext)
func InitEmptyInitializerBlockContext(p *InitializerBlockContext)
func InitEmptyInterfaceDeclContext(p *InterfaceDeclContext)
func InitEmptyInterfaceFieldContext(p *InterfaceFieldContext)
func InitEmptyLambdaParamContext(p *LambdaParamContext)
func InitEmptyLambdaParamListContext(p *LambdaParamListContext)
func InitEmptyLetStatementContext(p *LetStatementContext)
func InitEmptyMapEntryContext(p *MapEntryContext)
func InitEmptyNamespaceDeclContext(p *NamespaceDeclContext)
func InitEmptyParamContext(p *ParamContext)
func InitEmptyParamListContext(p *ParamListContext)
func InitEmptyParamTypeContext(p *ParamTypeContext)
func InitEmptyPrimaryContext(p *PrimaryContext)
func InitEmptyPrimitiveTypeContext(p *PrimitiveTypeContext)
func InitEmptyQualifiedNameContext(p *QualifiedNameContext)
func InitEmptyReturnStatementContext(p *ReturnStatementContext)
func InitEmptyReturnTypeContext(p *ReturnTypeContext)
func InitEmptySelfParamContext(p *SelfParamContext)
func InitEmptySemiContext(p *SemiContext)
func InitEmptyStatementContext(p *StatementContext)
func InitEmptySwitchCaseContext(p *SwitchCaseContext)
func InitEmptySwitchDefaultContext(p *SwitchDefaultContext)
func InitEmptySwitchStatementContext(p *SwitchStatementContext)
func InitEmptyTopLevelDeclContext(p *TopLevelDeclContext)
func InitEmptyTopLevelLetDeclContext(p *TopLevelLetDeclContext)
func InitEmptyTopLevelVarDeclContext(p *TopLevelVarDeclContext)
func InitEmptyTupleTypeContext(p *TupleTypeContext)
func InitEmptyTypeAliasDeclContext(p *TypeAliasDeclContext)
func InitEmptyTypeListContext(p *TypeListContext)
func InitEmptyTypeRefContext(p *TypeRefContext)
func InitEmptyVarStatementContext(p *VarStatementContext)

TYPES

type AddExprContext struct {
	ExpressionContext
	// Has unexported fields.
}

func NewAddExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AddExprContext

func (s *AddExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AddExprContext) AllExpression() []IExpressionContext

func (s *AddExprContext) Expression(i int) IExpressionContext

func (s *AddExprContext) GetOp() antlr.Token

func (s *AddExprContext) GetRuleContext() antlr.RuleContext

func (s *AddExprContext) MINUS() antlr.TerminalNode

func (s *AddExprContext) PLUS() antlr.TerminalNode

func (s *AddExprContext) SetOp(v antlr.Token)

type AddressOfContext struct {
	ExpressionContext
}

func NewAddressOfContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AddressOfContext

func (s *AddressOfContext) AMP() antlr.TerminalNode

func (s *AddressOfContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AddressOfContext) Expression() IExpressionContext

func (s *AddressOfContext) GetRuleContext() antlr.RuleContext

type ArcLexer struct {
	*antlr.BaseLexer

	// Has unexported fields.
}

func NewArcLexer(input antlr.CharStream) *ArcLexer
    NewArcLexer produces a new lexer instance for the optional input
    antlr.CharStream.

type ArcParser struct {
	*antlr.BaseParser
}

func NewArcParser(input antlr.TokenStream) *ArcParser
    NewArcParser produces a new parser instance for the optional input
    antlr.TokenStream.

func (p *ArcParser) Argument() (localctx IArgumentContext)

func (p *ArcParser) ArgumentList() (localctx IArgumentListContext)

func (p *ArcParser) AssignOp() (localctx IAssignOpContext)

func (p *ArcParser) AssignmentStatement() (localctx IAssignmentStatementContext)

func (p *ArcParser) AssignmentTarget() (localctx IAssignmentTargetContext)

func (p *ArcParser) Attribute() (localctx IAttributeContext)

func (p *ArcParser) BaseType() (localctx IBaseTypeContext)

func (p *ArcParser) Block() (localctx IBlockContext)

func (p *ArcParser) BreakStatement() (localctx IBreakStatementContext)

func (p *ArcParser) CallingConvention() (localctx ICallingConventionContext)

func (p *ArcParser) CompilationUnit() (localctx ICompilationUnitContext)

func (p *ArcParser) ConstDecl() (localctx IConstDeclContext)

func (p *ArcParser) ConstSpec() (localctx IConstSpecContext)

func (p *ArcParser) ContinueStatement() (localctx IContinueStatementContext)

func (p *ArcParser) DeferStatement() (localctx IDeferStatementContext)

func (p *ArcParser) DeinitDecl() (localctx IDeinitDeclContext)

func (p *ArcParser) EnumDecl() (localctx IEnumDeclContext)

func (p *ArcParser) EnumMember() (localctx IEnumMemberContext)

func (p *ArcParser) Expression() (localctx IExpressionContext)

func (p *ArcParser) ExpressionList() (localctx IExpressionListContext)

func (p *ArcParser) ExpressionStatement() (localctx IExpressionStatementContext)

func (p *ArcParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool

func (p *ArcParser) ExternClass() (localctx IExternClassContext)

func (p *ArcParser) ExternClassMember() (localctx IExternClassMemberContext)

func (p *ArcParser) ExternConstructor() (localctx IExternConstructorContext)

func (p *ArcParser) ExternDecl() (localctx IExternDeclContext)

func (p *ArcParser) ExternDestructor() (localctx IExternDestructorContext)

func (p *ArcParser) ExternFuncDecl() (localctx IExternFuncDeclContext)

func (p *ArcParser) ExternFunctionPtrType() (localctx IExternFunctionPtrTypeContext)

func (p *ArcParser) ExternMember() (localctx IExternMemberContext)

func (p *ArcParser) ExternMethodParam() (localctx IExternMethodParamContext)

func (p *ArcParser) ExternMethodParamList() (localctx IExternMethodParamListContext)

func (p *ArcParser) ExternNamespace() (localctx IExternNamespaceContext)

func (p *ArcParser) ExternParam() (localctx IExternParamContext)

func (p *ArcParser) ExternParamList() (localctx IExternParamListContext)

func (p *ArcParser) ExternReturnType() (localctx IExternReturnTypeContext)

func (p *ArcParser) ExternStaticMethod() (localctx IExternStaticMethodContext)

func (p *ArcParser) ExternSymbol() (localctx IExternSymbolContext)

func (p *ArcParser) ExternType() (localctx IExternTypeContext)

func (p *ArcParser) ExternTypeAlias() (localctx IExternTypeAliasContext)

func (p *ArcParser) ExternVirtualMethod() (localctx IExternVirtualMethodContext)

func (p *ArcParser) FieldInit() (localctx IFieldInitContext)

func (p *ArcParser) ForHeader() (localctx IForHeaderContext)

func (p *ArcParser) ForInit() (localctx IForInitContext)

func (p *ArcParser) ForIterator() (localctx IForIteratorContext)

func (p *ArcParser) ForPost() (localctx IForPostContext)

func (p *ArcParser) ForStatement() (localctx IForStatementContext)

func (p *ArcParser) FuncDecl() (localctx IFuncDeclContext)

func (p *ArcParser) FuncModifier() (localctx IFuncModifierContext)

func (p *ArcParser) FunctionType() (localctx IFunctionTypeContext)

func (p *ArcParser) GenericArgs() (localctx IGenericArgsContext)

func (p *ArcParser) GenericParams() (localctx IGenericParamsContext)

func (p *ArcParser) IfStatement() (localctx IIfStatementContext)

func (p *ArcParser) ImportAlias() (localctx IImportAliasContext)

func (p *ArcParser) ImportDecl() (localctx IImportDeclContext)

func (p *ArcParser) ImportSpec() (localctx IImportSpecContext)

func (p *ArcParser) InitializerBlock() (localctx IInitializerBlockContext)

func (p *ArcParser) InterfaceDecl() (localctx IInterfaceDeclContext)

func (p *ArcParser) InterfaceField() (localctx IInterfaceFieldContext)

func (p *ArcParser) LambdaParam() (localctx ILambdaParamContext)

func (p *ArcParser) LambdaParamList() (localctx ILambdaParamListContext)

func (p *ArcParser) LetStatement() (localctx ILetStatementContext)

func (p *ArcParser) MapEntry() (localctx IMapEntryContext)

func (p *ArcParser) NamespaceDecl() (localctx INamespaceDeclContext)

func (p *ArcParser) Param() (localctx IParamContext)

func (p *ArcParser) ParamList() (localctx IParamListContext)

func (p *ArcParser) ParamType() (localctx IParamTypeContext)

func (p *ArcParser) Primary() (localctx IPrimaryContext)

func (p *ArcParser) PrimitiveType() (localctx IPrimitiveTypeContext)

func (p *ArcParser) QualifiedName() (localctx IQualifiedNameContext)

func (p *ArcParser) ReturnStatement() (localctx IReturnStatementContext)

func (p *ArcParser) ReturnType() (localctx IReturnTypeContext)

func (p *ArcParser) SelfParam() (localctx ISelfParamContext)

func (p *ArcParser) Semi() (localctx ISemiContext)

func (p *ArcParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool

func (p *ArcParser) Statement() (localctx IStatementContext)

func (p *ArcParser) SwitchCase() (localctx ISwitchCaseContext)

func (p *ArcParser) SwitchDefault() (localctx ISwitchDefaultContext)

func (p *ArcParser) SwitchStatement() (localctx ISwitchStatementContext)

func (p *ArcParser) TopLevelDecl() (localctx ITopLevelDeclContext)

func (p *ArcParser) TopLevelLetDecl() (localctx ITopLevelLetDeclContext)

func (p *ArcParser) TopLevelVarDecl() (localctx ITopLevelVarDeclContext)

func (p *ArcParser) TupleType() (localctx ITupleTypeContext)

func (p *ArcParser) TypeAliasDecl() (localctx ITypeAliasDeclContext)

func (p *ArcParser) TypeList() (localctx ITypeListContext)

func (p *ArcParser) TypeRef() (localctx ITypeRefContext)

func (p *ArcParser) VarStatement() (localctx IVarStatementContext)

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
    A complete Visitor for a parse tree produced by ArcParser.

type ArgumentContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentContext

func NewEmptyArgumentContext() *ArgumentContext

func (s *ArgumentContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgumentContext) Expression() IExpressionContext

func (s *ArgumentContext) GetParser() antlr.Parser

func (s *ArgumentContext) GetRuleContext() antlr.RuleContext

func (*ArgumentContext) IsArgumentContext()

func (s *ArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ArgumentListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewArgumentListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentListContext

func NewEmptyArgumentListContext() *ArgumentListContext

func (s *ArgumentListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ArgumentListContext) AllArgument() []IArgumentContext

func (s *ArgumentListContext) AllCOMMA() []antlr.TerminalNode

func (s *ArgumentListContext) Argument(i int) IArgumentContext

func (s *ArgumentListContext) COMMA(i int) antlr.TerminalNode

func (s *ArgumentListContext) GetParser() antlr.Parser

func (s *ArgumentListContext) GetRuleContext() antlr.RuleContext

func (*ArgumentListContext) IsArgumentListContext()

func (s *ArgumentListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AssignOpContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignOpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignOpContext

func NewEmptyAssignOpContext() *AssignOpContext

func (s *AssignOpContext) ADD_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) AND_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignOpContext) DIV_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) GetParser() antlr.Parser

func (s *AssignOpContext) GetRuleContext() antlr.RuleContext

func (*AssignOpContext) IsAssignOpContext()

func (s *AssignOpContext) MOD_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) MUL_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) OR_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) SHL_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) SHR_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) SUB_ASSIGN() antlr.TerminalNode

func (s *AssignOpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *AssignOpContext) XOR_ASSIGN() antlr.TerminalNode

type AssignmentStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignmentStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentStatementContext

func NewEmptyAssignmentStatementContext() *AssignmentStatementContext

func (s *AssignmentStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignmentStatementContext) AssignOp() IAssignOpContext

func (s *AssignmentStatementContext) AssignmentTarget() IAssignmentTargetContext

func (s *AssignmentStatementContext) DEC() antlr.TerminalNode

func (s *AssignmentStatementContext) Expression() IExpressionContext

func (s *AssignmentStatementContext) GetParser() antlr.Parser

func (s *AssignmentStatementContext) GetRuleContext() antlr.RuleContext

func (s *AssignmentStatementContext) INC() antlr.TerminalNode

func (*AssignmentStatementContext) IsAssignmentStatementContext()

func (s *AssignmentStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AssignmentTargetContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAssignmentTargetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentTargetContext

func NewEmptyAssignmentTargetContext() *AssignmentTargetContext

func (s *AssignmentTargetContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AssignmentTargetContext) AllExpression() []IExpressionContext

func (s *AssignmentTargetContext) DOT() antlr.TerminalNode

func (s *AssignmentTargetContext) Expression(i int) IExpressionContext

func (s *AssignmentTargetContext) GetParser() antlr.Parser

func (s *AssignmentTargetContext) GetRuleContext() antlr.RuleContext

func (s *AssignmentTargetContext) IDENTIFIER() antlr.TerminalNode

func (*AssignmentTargetContext) IsAssignmentTargetContext()

func (s *AssignmentTargetContext) LBRACKET() antlr.TerminalNode

func (s *AssignmentTargetContext) RBRACKET() antlr.TerminalNode

func (s *AssignmentTargetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AttributeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewAttributeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeContext

func NewEmptyAttributeContext() *AttributeContext

func (s *AttributeContext) AT() antlr.TerminalNode

func (s *AttributeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AttributeContext) Expression() IExpressionContext

func (s *AttributeContext) GetParser() antlr.Parser

func (s *AttributeContext) GetRuleContext() antlr.RuleContext

func (s *AttributeContext) IDENTIFIER() antlr.TerminalNode

func (*AttributeContext) IsAttributeContext()

func (s *AttributeContext) LPAREN() antlr.TerminalNode

func (s *AttributeContext) RPAREN() antlr.TerminalNode

func (s *AttributeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type AwaitExprContext struct {
	ExpressionContext
}

func NewAwaitExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AwaitExprContext

func (s *AwaitExprContext) AWAIT() antlr.TerminalNode

func (s *AwaitExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *AwaitExprContext) Expression() IExpressionContext

func (s *AwaitExprContext) GetRuleContext() antlr.RuleContext

type BareInitExprContext struct {
	PrimaryContext
}

func NewBareInitExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BareInitExprContext

func (s *BareInitExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BareInitExprContext) GetRuleContext() antlr.RuleContext

func (s *BareInitExprContext) InitializerBlock() IInitializerBlockContext

type BaseArcParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseArcParserVisitor) VisitAddExpr(ctx *AddExprContext) interface{}

func (v *BaseArcParserVisitor) VisitAddressOf(ctx *AddressOfContext) interface{}

func (v *BaseArcParserVisitor) VisitArgument(ctx *ArgumentContext) interface{}

func (v *BaseArcParserVisitor) VisitArgumentList(ctx *ArgumentListContext) interface{}

func (v *BaseArcParserVisitor) VisitAssignOp(ctx *AssignOpContext) interface{}

func (v *BaseArcParserVisitor) VisitAssignmentStatement(ctx *AssignmentStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitAssignmentTarget(ctx *AssignmentTargetContext) interface{}

func (v *BaseArcParserVisitor) VisitAttribute(ctx *AttributeContext) interface{}

func (v *BaseArcParserVisitor) VisitAwaitExpr(ctx *AwaitExprContext) interface{}

func (v *BaseArcParserVisitor) VisitBareInitExpr(ctx *BareInitExprContext) interface{}

func (v *BaseArcParserVisitor) VisitBaseType(ctx *BaseTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitBitwiseAndExpr(ctx *BitwiseAndExprContext) interface{}

func (v *BaseArcParserVisitor) VisitBitwiseNot(ctx *BitwiseNotContext) interface{}

func (v *BaseArcParserVisitor) VisitBitwiseOrExpr(ctx *BitwiseOrExprContext) interface{}

func (v *BaseArcParserVisitor) VisitBitwiseXorExpr(ctx *BitwiseXorExprContext) interface{}

func (v *BaseArcParserVisitor) VisitBlock(ctx *BlockContext) interface{}

func (v *BaseArcParserVisitor) VisitBreakStatement(ctx *BreakStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitCallExpr(ctx *CallExprContext) interface{}

func (v *BaseArcParserVisitor) VisitCallingConvention(ctx *CallingConventionContext) interface{}

func (v *BaseArcParserVisitor) VisitCharLiteral(ctx *CharLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitCompilationUnit(ctx *CompilationUnitContext) interface{}

func (v *BaseArcParserVisitor) VisitConstDecl(ctx *ConstDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitConstSpec(ctx *ConstSpecContext) interface{}

func (v *BaseArcParserVisitor) VisitContinueStatement(ctx *ContinueStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitDeferStatement(ctx *DeferStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitDeinitDecl(ctx *DeinitDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitDeleteExpr(ctx *DeleteExprContext) interface{}

func (v *BaseArcParserVisitor) VisitEnumDecl(ctx *EnumDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitEnumMember(ctx *EnumMemberContext) interface{}

func (v *BaseArcParserVisitor) VisitEqualityExpr(ctx *EqualityExprContext) interface{}

func (v *BaseArcParserVisitor) VisitExpressionList(ctx *ExpressionListContext) interface{}

func (v *BaseArcParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitExternClass(ctx *ExternClassContext) interface{}

func (v *BaseArcParserVisitor) VisitExternClassMember(ctx *ExternClassMemberContext) interface{}

func (v *BaseArcParserVisitor) VisitExternConstructor(ctx *ExternConstructorContext) interface{}

func (v *BaseArcParserVisitor) VisitExternDecl(ctx *ExternDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitExternDestructor(ctx *ExternDestructorContext) interface{}

func (v *BaseArcParserVisitor) VisitExternFuncDecl(ctx *ExternFuncDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitExternFunctionPtrType(ctx *ExternFunctionPtrTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitExternMember(ctx *ExternMemberContext) interface{}

func (v *BaseArcParserVisitor) VisitExternMethodParam(ctx *ExternMethodParamContext) interface{}

func (v *BaseArcParserVisitor) VisitExternMethodParamList(ctx *ExternMethodParamListContext) interface{}

func (v *BaseArcParserVisitor) VisitExternNamespace(ctx *ExternNamespaceContext) interface{}

func (v *BaseArcParserVisitor) VisitExternParam(ctx *ExternParamContext) interface{}

func (v *BaseArcParserVisitor) VisitExternParamList(ctx *ExternParamListContext) interface{}

func (v *BaseArcParserVisitor) VisitExternReturnType(ctx *ExternReturnTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitExternStaticMethod(ctx *ExternStaticMethodContext) interface{}

func (v *BaseArcParserVisitor) VisitExternSymbol(ctx *ExternSymbolContext) interface{}

func (v *BaseArcParserVisitor) VisitExternType(ctx *ExternTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitExternTypeAlias(ctx *ExternTypeAliasContext) interface{}

func (v *BaseArcParserVisitor) VisitExternVirtualMethod(ctx *ExternVirtualMethodContext) interface{}

func (v *BaseArcParserVisitor) VisitFalseLiteral(ctx *FalseLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitFieldInit(ctx *FieldInitContext) interface{}

func (v *BaseArcParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitForHeader(ctx *ForHeaderContext) interface{}

func (v *BaseArcParserVisitor) VisitForInit(ctx *ForInitContext) interface{}

func (v *BaseArcParserVisitor) VisitForIterator(ctx *ForIteratorContext) interface{}

func (v *BaseArcParserVisitor) VisitForPost(ctx *ForPostContext) interface{}

func (v *BaseArcParserVisitor) VisitForStatement(ctx *ForStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitFuncDecl(ctx *FuncDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitFuncModifier(ctx *FuncModifierContext) interface{}

func (v *BaseArcParserVisitor) VisitFunctionType(ctx *FunctionTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitGenericArgs(ctx *GenericArgsContext) interface{}

func (v *BaseArcParserVisitor) VisitGenericParams(ctx *GenericParamsContext) interface{}

func (v *BaseArcParserVisitor) VisitHexLiteral(ctx *HexLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitIdentExpr(ctx *IdentExprContext) interface{}

func (v *BaseArcParserVisitor) VisitIfStatement(ctx *IfStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitImportAlias(ctx *ImportAliasContext) interface{}

func (v *BaseArcParserVisitor) VisitImportDecl(ctx *ImportDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitImportSpec(ctx *ImportSpecContext) interface{}

func (v *BaseArcParserVisitor) VisitIndexExpr(ctx *IndexExprContext) interface{}

func (v *BaseArcParserVisitor) VisitInitializerBlock(ctx *InitializerBlockContext) interface{}

func (v *BaseArcParserVisitor) VisitIntLiteral(ctx *IntLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitInterfaceDecl(ctx *InterfaceDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitInterfaceField(ctx *InterfaceFieldContext) interface{}

func (v *BaseArcParserVisitor) VisitLambdaExpr(ctx *LambdaExprContext) interface{}

func (v *BaseArcParserVisitor) VisitLambdaParam(ctx *LambdaParamContext) interface{}

func (v *BaseArcParserVisitor) VisitLambdaParamList(ctx *LambdaParamListContext) interface{}

func (v *BaseArcParserVisitor) VisitLetStatement(ctx *LetStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitLogicalAndExpr(ctx *LogicalAndExprContext) interface{}

func (v *BaseArcParserVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{}

func (v *BaseArcParserVisitor) VisitLogicalOrExpr(ctx *LogicalOrExprContext) interface{}

func (v *BaseArcParserVisitor) VisitMapEntry(ctx *MapEntryContext) interface{}

func (v *BaseArcParserVisitor) VisitMapLiteral(ctx *MapLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitMemberAccess(ctx *MemberAccessContext) interface{}

func (v *BaseArcParserVisitor) VisitMulExpr(ctx *MulExprContext) interface{}

func (v *BaseArcParserVisitor) VisitNamespaceDecl(ctx *NamespaceDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitNewArrayExpr(ctx *NewArrayExprContext) interface{}

func (v *BaseArcParserVisitor) VisitNewExpr(ctx *NewExprContext) interface{}

func (v *BaseArcParserVisitor) VisitNullLiteral(ctx *NullLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitParam(ctx *ParamContext) interface{}

func (v *BaseArcParserVisitor) VisitParamList(ctx *ParamListContext) interface{}

func (v *BaseArcParserVisitor) VisitParamType(ctx *ParamTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitParenExpr(ctx *ParenExprContext) interface{}

func (v *BaseArcParserVisitor) VisitPostDecrement(ctx *PostDecrementContext) interface{}

func (v *BaseArcParserVisitor) VisitPostIncrement(ctx *PostIncrementContext) interface{}

func (v *BaseArcParserVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

func (v *BaseArcParserVisitor) VisitPrimitiveType(ctx *PrimitiveTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitPrimitiveTypeExpr(ctx *PrimitiveTypeExprContext) interface{}

func (v *BaseArcParserVisitor) VisitProcessExpr(ctx *ProcessExprContext) interface{}

func (v *BaseArcParserVisitor) VisitQualifiedExpr(ctx *QualifiedExprContext) interface{}

func (v *BaseArcParserVisitor) VisitQualifiedName(ctx *QualifiedNameContext) interface{}

func (v *BaseArcParserVisitor) VisitRangeExpr(ctx *RangeExprContext) interface{}

func (v *BaseArcParserVisitor) VisitRelationalExpr(ctx *RelationalExprContext) interface{}

func (v *BaseArcParserVisitor) VisitReturnStatement(ctx *ReturnStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitReturnType(ctx *ReturnTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitSelfParam(ctx *SelfParamContext) interface{}

func (v *BaseArcParserVisitor) VisitSemi(ctx *SemiContext) interface{}

func (v *BaseArcParserVisitor) VisitShiftExpr(ctx *ShiftExprContext) interface{}

func (v *BaseArcParserVisitor) VisitSliceExpr(ctx *SliceExprContext) interface{}

func (v *BaseArcParserVisitor) VisitStatement(ctx *StatementContext) interface{}

func (v *BaseArcParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitSwitchCase(ctx *SwitchCaseContext) interface{}

func (v *BaseArcParserVisitor) VisitSwitchDefault(ctx *SwitchDefaultContext) interface{}

func (v *BaseArcParserVisitor) VisitSwitchStatement(ctx *SwitchStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitTopLevelDecl(ctx *TopLevelDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitTopLevelLetDecl(ctx *TopLevelLetDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitTopLevelVarDecl(ctx *TopLevelVarDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitTrueLiteral(ctx *TrueLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitTupleLiteral(ctx *TupleLiteralContext) interface{}

func (v *BaseArcParserVisitor) VisitTupleType(ctx *TupleTypeContext) interface{}

func (v *BaseArcParserVisitor) VisitTypeAliasDecl(ctx *TypeAliasDeclContext) interface{}

func (v *BaseArcParserVisitor) VisitTypeList(ctx *TypeListContext) interface{}

func (v *BaseArcParserVisitor) VisitTypeRef(ctx *TypeRefContext) interface{}

func (v *BaseArcParserVisitor) VisitTypedInitExpr(ctx *TypedInitExprContext) interface{}

func (v *BaseArcParserVisitor) VisitUnaryMinus(ctx *UnaryMinusContext) interface{}

func (v *BaseArcParserVisitor) VisitVarStatement(ctx *VarStatementContext) interface{}

func (v *BaseArcParserVisitor) VisitVectorLiteral(ctx *VectorLiteralContext) interface{}

type BaseTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBaseTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BaseTypeContext

func NewEmptyBaseTypeContext() *BaseTypeContext

func (s *BaseTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BaseTypeContext) AllTypeRef() []ITypeRefContext

func (s *BaseTypeContext) BOOL() antlr.TerminalNode

func (s *BaseTypeContext) BYTE() antlr.TerminalNode

func (s *BaseTypeContext) CHAR() antlr.TerminalNode

func (s *BaseTypeContext) Expression() IExpressionContext

func (s *BaseTypeContext) GenericArgs() IGenericArgsContext

func (s *BaseTypeContext) GetParser() antlr.Parser

func (s *BaseTypeContext) GetRuleContext() antlr.RuleContext

func (s *BaseTypeContext) IDENTIFIER() antlr.TerminalNode

func (*BaseTypeContext) IsBaseTypeContext()

func (s *BaseTypeContext) LBRACKET() antlr.TerminalNode

func (s *BaseTypeContext) MAP() antlr.TerminalNode

func (s *BaseTypeContext) PrimitiveType() IPrimitiveTypeContext

func (s *BaseTypeContext) QualifiedName() IQualifiedNameContext

func (s *BaseTypeContext) RBRACKET() antlr.TerminalNode

func (s *BaseTypeContext) STRING() antlr.TerminalNode

func (s *BaseTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *BaseTypeContext) TypeRef(i int) ITypeRefContext

func (s *BaseTypeContext) VECTOR() antlr.TerminalNode

func (s *BaseTypeContext) VOID() antlr.TerminalNode

type BitwiseAndExprContext struct {
	ExpressionContext
}

func NewBitwiseAndExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitwiseAndExprContext

func (s *BitwiseAndExprContext) AMP() antlr.TerminalNode

func (s *BitwiseAndExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BitwiseAndExprContext) AllExpression() []IExpressionContext

func (s *BitwiseAndExprContext) Expression(i int) IExpressionContext

func (s *BitwiseAndExprContext) GetRuleContext() antlr.RuleContext

type BitwiseNotContext struct {
	ExpressionContext
}

func NewBitwiseNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitwiseNotContext

func (s *BitwiseNotContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BitwiseNotContext) Expression() IExpressionContext

func (s *BitwiseNotContext) GetRuleContext() antlr.RuleContext

func (s *BitwiseNotContext) TILDE() antlr.TerminalNode

type BitwiseOrExprContext struct {
	ExpressionContext
}

func NewBitwiseOrExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitwiseOrExprContext

func (s *BitwiseOrExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BitwiseOrExprContext) AllExpression() []IExpressionContext

func (s *BitwiseOrExprContext) Expression(i int) IExpressionContext

func (s *BitwiseOrExprContext) GetRuleContext() antlr.RuleContext

func (s *BitwiseOrExprContext) PIPE() antlr.TerminalNode

type BitwiseXorExprContext struct {
	ExpressionContext
}

func NewBitwiseXorExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BitwiseXorExprContext

func (s *BitwiseXorExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BitwiseXorExprContext) AllExpression() []IExpressionContext

func (s *BitwiseXorExprContext) CARET() antlr.TerminalNode

func (s *BitwiseXorExprContext) Expression(i int) IExpressionContext

func (s *BitwiseXorExprContext) GetRuleContext() antlr.RuleContext

type BlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BlockContext

func NewEmptyBlockContext() *BlockContext

func (s *BlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BlockContext) AllStatement() []IStatementContext

func (s *BlockContext) GetParser() antlr.Parser

func (s *BlockContext) GetRuleContext() antlr.RuleContext

func (*BlockContext) IsBlockContext()

func (s *BlockContext) LBRACE() antlr.TerminalNode

func (s *BlockContext) RBRACE() antlr.TerminalNode

func (s *BlockContext) Statement(i int) IStatementContext

func (s *BlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type BreakStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewBreakStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BreakStatementContext

func NewEmptyBreakStatementContext() *BreakStatementContext

func (s *BreakStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *BreakStatementContext) BREAK() antlr.TerminalNode

func (s *BreakStatementContext) GetParser() antlr.Parser

func (s *BreakStatementContext) GetRuleContext() antlr.RuleContext

func (*BreakStatementContext) IsBreakStatementContext()

func (s *BreakStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type CallExprContext struct {
	ExpressionContext
}

func NewCallExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CallExprContext

func (s *CallExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CallExprContext) ArgumentList() IArgumentListContext

func (s *CallExprContext) Expression() IExpressionContext

func (s *CallExprContext) GetRuleContext() antlr.RuleContext

func (s *CallExprContext) LPAREN() antlr.TerminalNode

func (s *CallExprContext) RPAREN() antlr.TerminalNode

type CallingConventionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCallingConventionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CallingConventionContext

func NewEmptyCallingConventionContext() *CallingConventionContext

func (s *CallingConventionContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CallingConventionContext) CDECL() antlr.TerminalNode

func (s *CallingConventionContext) FASTCALL() antlr.TerminalNode

func (s *CallingConventionContext) GetParser() antlr.Parser

func (s *CallingConventionContext) GetRuleContext() antlr.RuleContext

func (*CallingConventionContext) IsCallingConventionContext()

func (s *CallingConventionContext) STDCALL() antlr.TerminalNode

func (s *CallingConventionContext) THISCALL() antlr.TerminalNode

func (s *CallingConventionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *CallingConventionContext) VECTORCALL() antlr.TerminalNode

type CharLiteralContext struct {
	PrimaryContext
}

func NewCharLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CharLiteralContext

func (s *CharLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CharLiteralContext) CHAR_LIT() antlr.TerminalNode

func (s *CharLiteralContext) GetRuleContext() antlr.RuleContext

type CompilationUnitContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewCompilationUnitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompilationUnitContext

func NewEmptyCompilationUnitContext() *CompilationUnitContext

func (s *CompilationUnitContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *CompilationUnitContext) AllTopLevelDecl() []ITopLevelDeclContext

func (s *CompilationUnitContext) EOF() antlr.TerminalNode

func (s *CompilationUnitContext) GetParser() antlr.Parser

func (s *CompilationUnitContext) GetRuleContext() antlr.RuleContext

func (*CompilationUnitContext) IsCompilationUnitContext()

func (s *CompilationUnitContext) NamespaceDecl() INamespaceDeclContext

func (s *CompilationUnitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *CompilationUnitContext) TopLevelDecl(i int) ITopLevelDeclContext

type ConstDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConstDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstDeclContext

func NewEmptyConstDeclContext() *ConstDeclContext

func (s *ConstDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConstDeclContext) AllConstSpec() []IConstSpecContext

func (s *ConstDeclContext) CONST() antlr.TerminalNode

func (s *ConstDeclContext) ConstSpec(i int) IConstSpecContext

func (s *ConstDeclContext) GetParser() antlr.Parser

func (s *ConstDeclContext) GetRuleContext() antlr.RuleContext

func (*ConstDeclContext) IsConstDeclContext()

func (s *ConstDeclContext) LPAREN() antlr.TerminalNode

func (s *ConstDeclContext) RPAREN() antlr.TerminalNode

func (s *ConstDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ConstSpecContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewConstSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstSpecContext

func NewEmptyConstSpecContext() *ConstSpecContext

func (s *ConstSpecContext) ASSIGN() antlr.TerminalNode

func (s *ConstSpecContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ConstSpecContext) COLON() antlr.TerminalNode

func (s *ConstSpecContext) Expression() IExpressionContext

func (s *ConstSpecContext) GetParser() antlr.Parser

func (s *ConstSpecContext) GetRuleContext() antlr.RuleContext

func (s *ConstSpecContext) IDENTIFIER() antlr.TerminalNode

func (*ConstSpecContext) IsConstSpecContext()

func (s *ConstSpecContext) Semi() ISemiContext

func (s *ConstSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ConstSpecContext) TypeRef() ITypeRefContext

type ContinueStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewContinueStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ContinueStatementContext

func NewEmptyContinueStatementContext() *ContinueStatementContext

func (s *ContinueStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ContinueStatementContext) CONTINUE() antlr.TerminalNode

func (s *ContinueStatementContext) GetParser() antlr.Parser

func (s *ContinueStatementContext) GetRuleContext() antlr.RuleContext

func (*ContinueStatementContext) IsContinueStatementContext()

func (s *ContinueStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeferStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeferStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeferStatementContext

func NewEmptyDeferStatementContext() *DeferStatementContext

func (s *DeferStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeferStatementContext) DEFER() antlr.TerminalNode

func (s *DeferStatementContext) Expression() IExpressionContext

func (s *DeferStatementContext) GetParser() antlr.Parser

func (s *DeferStatementContext) GetRuleContext() antlr.RuleContext

func (*DeferStatementContext) IsDeferStatementContext()

func (s *DeferStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeinitDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewDeinitDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeinitDeclContext

func NewEmptyDeinitDeclContext() *DeinitDeclContext

func (s *DeinitDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeinitDeclContext) Block() IBlockContext

func (s *DeinitDeclContext) DEINIT() antlr.TerminalNode

func (s *DeinitDeclContext) GetParser() antlr.Parser

func (s *DeinitDeclContext) GetRuleContext() antlr.RuleContext

func (*DeinitDeclContext) IsDeinitDeclContext()

func (s *DeinitDeclContext) LPAREN() antlr.TerminalNode

func (s *DeinitDeclContext) RPAREN() antlr.TerminalNode

func (s *DeinitDeclContext) SelfParam() ISelfParamContext

func (s *DeinitDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type DeleteExprContext struct {
	PrimaryContext
}

func NewDeleteExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DeleteExprContext

func (s *DeleteExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *DeleteExprContext) DELETE() antlr.TerminalNode

func (s *DeleteExprContext) Expression() IExpressionContext

func (s *DeleteExprContext) GetRuleContext() antlr.RuleContext

func (s *DeleteExprContext) LPAREN() antlr.TerminalNode

func (s *DeleteExprContext) RPAREN() antlr.TerminalNode

type EnumDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumDeclContext() *EnumDeclContext

func NewEnumDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumDeclContext

func (s *EnumDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumDeclContext) AllEnumMember() []IEnumMemberContext

func (s *EnumDeclContext) COLON() antlr.TerminalNode

func (s *EnumDeclContext) ENUM() antlr.TerminalNode

func (s *EnumDeclContext) EnumMember(i int) IEnumMemberContext

func (s *EnumDeclContext) GetParser() antlr.Parser

func (s *EnumDeclContext) GetRuleContext() antlr.RuleContext

func (s *EnumDeclContext) IDENTIFIER() antlr.TerminalNode

func (*EnumDeclContext) IsEnumDeclContext()

func (s *EnumDeclContext) LBRACE() antlr.TerminalNode

func (s *EnumDeclContext) PrimitiveType() IPrimitiveTypeContext

func (s *EnumDeclContext) RBRACE() antlr.TerminalNode

func (s *EnumDeclContext) Semi() ISemiContext

func (s *EnumDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EnumMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyEnumMemberContext() *EnumMemberContext

func NewEnumMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EnumMemberContext

func (s *EnumMemberContext) ASSIGN() antlr.TerminalNode

func (s *EnumMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EnumMemberContext) Expression() IExpressionContext

func (s *EnumMemberContext) GetParser() antlr.Parser

func (s *EnumMemberContext) GetRuleContext() antlr.RuleContext

func (s *EnumMemberContext) IDENTIFIER() antlr.TerminalNode

func (*EnumMemberContext) IsEnumMemberContext()

func (s *EnumMemberContext) Semi() ISemiContext

func (s *EnumMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type EqualityExprContext struct {
	ExpressionContext
	// Has unexported fields.
}

func NewEqualityExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityExprContext

func (s *EqualityExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *EqualityExprContext) AllExpression() []IExpressionContext

func (s *EqualityExprContext) EQ() antlr.TerminalNode

func (s *EqualityExprContext) Expression(i int) IExpressionContext

func (s *EqualityExprContext) GetOp() antlr.Token

func (s *EqualityExprContext) GetRuleContext() antlr.RuleContext

func (s *EqualityExprContext) NEQ() antlr.TerminalNode

func (s *EqualityExprContext) SetOp(v antlr.Token)

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExpressionContext() *ExpressionContext

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext

func (s *ExpressionContext) CopyAll(ctx *ExpressionContext)

func (s *ExpressionContext) GetParser() antlr.Parser

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext

func (*ExpressionContext) IsExpressionContext()

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExpressionListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExpressionListContext() *ExpressionListContext

func NewExpressionListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionListContext

func (s *ExpressionListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExpressionListContext) AllCOMMA() []antlr.TerminalNode

func (s *ExpressionListContext) AllExpression() []IExpressionContext

func (s *ExpressionListContext) COMMA(i int) antlr.TerminalNode

func (s *ExpressionListContext) Expression(i int) IExpressionContext

func (s *ExpressionListContext) GetParser() antlr.Parser

func (s *ExpressionListContext) GetRuleContext() antlr.RuleContext

func (*ExpressionListContext) IsExpressionListContext()

func (s *ExpressionListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExpressionStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExpressionStatementContext() *ExpressionStatementContext

func NewExpressionStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionStatementContext

func (s *ExpressionStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExpressionStatementContext) Expression() IExpressionContext

func (s *ExpressionStatementContext) GetParser() antlr.Parser

func (s *ExpressionStatementContext) GetRuleContext() antlr.RuleContext

func (*ExpressionStatementContext) IsExpressionStatementContext()

func (s *ExpressionStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternClassContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternClassContext() *ExternClassContext

func NewExternClassContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternClassContext

func (s *ExternClassContext) ABSTRACT() antlr.TerminalNode

func (s *ExternClassContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternClassContext) AllExternClassMember() []IExternClassMemberContext

func (s *ExternClassContext) CLASS() antlr.TerminalNode

func (s *ExternClassContext) ExternClassMember(i int) IExternClassMemberContext

func (s *ExternClassContext) ExternSymbol() IExternSymbolContext

func (s *ExternClassContext) GetParser() antlr.Parser

func (s *ExternClassContext) GetRuleContext() antlr.RuleContext

func (s *ExternClassContext) IDENTIFIER() antlr.TerminalNode

func (*ExternClassContext) IsExternClassContext()

func (s *ExternClassContext) LBRACE() antlr.TerminalNode

func (s *ExternClassContext) RBRACE() antlr.TerminalNode

func (s *ExternClassContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternClassMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternClassMemberContext() *ExternClassMemberContext

func NewExternClassMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternClassMemberContext

func (s *ExternClassMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternClassMemberContext) ExternConstructor() IExternConstructorContext

func (s *ExternClassMemberContext) ExternDestructor() IExternDestructorContext

func (s *ExternClassMemberContext) ExternStaticMethod() IExternStaticMethodContext

func (s *ExternClassMemberContext) ExternVirtualMethod() IExternVirtualMethodContext

func (s *ExternClassMemberContext) GetParser() antlr.Parser

func (s *ExternClassMemberContext) GetRuleContext() antlr.RuleContext

func (*ExternClassMemberContext) IsExternClassMemberContext()

func (s *ExternClassMemberContext) Semi() ISemiContext

func (s *ExternClassMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternConstructorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternConstructorContext() *ExternConstructorContext

func NewExternConstructorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternConstructorContext

func (s *ExternConstructorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternConstructorContext) ExternParamList() IExternParamListContext

func (s *ExternConstructorContext) ExternType() IExternTypeContext

func (s *ExternConstructorContext) GetParser() antlr.Parser

func (s *ExternConstructorContext) GetRuleContext() antlr.RuleContext

func (*ExternConstructorContext) IsExternConstructorContext()

func (s *ExternConstructorContext) LPAREN() antlr.TerminalNode

func (s *ExternConstructorContext) NEW() antlr.TerminalNode

func (s *ExternConstructorContext) RPAREN() antlr.TerminalNode

func (s *ExternConstructorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternDeclContext() *ExternDeclContext

func NewExternDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternDeclContext

func (s *ExternDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternDeclContext) AllExternMember() []IExternMemberContext

func (s *ExternDeclContext) EXTERN() antlr.TerminalNode

func (s *ExternDeclContext) ExternMember(i int) IExternMemberContext

func (s *ExternDeclContext) GetParser() antlr.Parser

func (s *ExternDeclContext) GetRuleContext() antlr.RuleContext

func (s *ExternDeclContext) IDENTIFIER() antlr.TerminalNode

func (*ExternDeclContext) IsExternDeclContext()

func (s *ExternDeclContext) LBRACE() antlr.TerminalNode

func (s *ExternDeclContext) RBRACE() antlr.TerminalNode

func (s *ExternDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternDestructorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternDestructorContext() *ExternDestructorContext

func NewExternDestructorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternDestructorContext

func (s *ExternDestructorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternDestructorContext) DELETE() antlr.TerminalNode

func (s *ExternDestructorContext) ExternMethodParam() IExternMethodParamContext

func (s *ExternDestructorContext) GetParser() antlr.Parser

func (s *ExternDestructorContext) GetRuleContext() antlr.RuleContext

func (*ExternDestructorContext) IsExternDestructorContext()

func (s *ExternDestructorContext) LPAREN() antlr.TerminalNode

func (s *ExternDestructorContext) RPAREN() antlr.TerminalNode

func (s *ExternDestructorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExternDestructorContext) VOID() antlr.TerminalNode

type ExternFuncDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternFuncDeclContext() *ExternFuncDeclContext

func NewExternFuncDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternFuncDeclContext

func (s *ExternFuncDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternFuncDeclContext) CallingConvention() ICallingConventionContext

func (s *ExternFuncDeclContext) ExternParamList() IExternParamListContext

func (s *ExternFuncDeclContext) ExternReturnType() IExternReturnTypeContext

func (s *ExternFuncDeclContext) ExternSymbol() IExternSymbolContext

func (s *ExternFuncDeclContext) FUNC() antlr.TerminalNode

func (s *ExternFuncDeclContext) GetParser() antlr.Parser

func (s *ExternFuncDeclContext) GetRuleContext() antlr.RuleContext

func (s *ExternFuncDeclContext) IDENTIFIER() antlr.TerminalNode

func (*ExternFuncDeclContext) IsExternFuncDeclContext()

func (s *ExternFuncDeclContext) LPAREN() antlr.TerminalNode

func (s *ExternFuncDeclContext) RPAREN() antlr.TerminalNode

func (s *ExternFuncDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternFunctionPtrTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternFunctionPtrTypeContext() *ExternFunctionPtrTypeContext

func NewExternFunctionPtrTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternFunctionPtrTypeContext

func (s *ExternFunctionPtrTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternFunctionPtrTypeContext) ExternParamList() IExternParamListContext

func (s *ExternFunctionPtrTypeContext) ExternReturnType() IExternReturnTypeContext

func (s *ExternFunctionPtrTypeContext) FUNC() antlr.TerminalNode

func (s *ExternFunctionPtrTypeContext) GetParser() antlr.Parser

func (s *ExternFunctionPtrTypeContext) GetRuleContext() antlr.RuleContext

func (*ExternFunctionPtrTypeContext) IsExternFunctionPtrTypeContext()

func (s *ExternFunctionPtrTypeContext) LPAREN() antlr.TerminalNode

func (s *ExternFunctionPtrTypeContext) RPAREN() antlr.TerminalNode

func (s *ExternFunctionPtrTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternMemberContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternMemberContext() *ExternMemberContext

func NewExternMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternMemberContext

func (s *ExternMemberContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternMemberContext) ExternClass() IExternClassContext

func (s *ExternMemberContext) ExternFuncDecl() IExternFuncDeclContext

func (s *ExternMemberContext) ExternNamespace() IExternNamespaceContext

func (s *ExternMemberContext) ExternTypeAlias() IExternTypeAliasContext

func (s *ExternMemberContext) GetParser() antlr.Parser

func (s *ExternMemberContext) GetRuleContext() antlr.RuleContext

func (*ExternMemberContext) IsExternMemberContext()

func (s *ExternMemberContext) Semi() ISemiContext

func (s *ExternMemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternMethodParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternMethodParamContext() *ExternMethodParamContext

func NewExternMethodParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternMethodParamContext

func (s *ExternMethodParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternMethodParamContext) ExternType() IExternTypeContext

func (s *ExternMethodParamContext) GetParser() antlr.Parser

func (s *ExternMethodParamContext) GetRuleContext() antlr.RuleContext

func (*ExternMethodParamContext) IsExternMethodParamContext()

func (s *ExternMethodParamContext) SELF() antlr.TerminalNode

func (s *ExternMethodParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternMethodParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternMethodParamListContext() *ExternMethodParamListContext

func NewExternMethodParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternMethodParamListContext

func (s *ExternMethodParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternMethodParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *ExternMethodParamListContext) AllExternParam() []IExternParamContext

func (s *ExternMethodParamListContext) COMMA(i int) antlr.TerminalNode

func (s *ExternMethodParamListContext) ELLIPSIS() antlr.TerminalNode

func (s *ExternMethodParamListContext) ExternMethodParam() IExternMethodParamContext

func (s *ExternMethodParamListContext) ExternParam(i int) IExternParamContext

func (s *ExternMethodParamListContext) GetParser() antlr.Parser

func (s *ExternMethodParamListContext) GetRuleContext() antlr.RuleContext

func (*ExternMethodParamListContext) IsExternMethodParamListContext()

func (s *ExternMethodParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternNamespaceContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternNamespaceContext() *ExternNamespaceContext

func NewExternNamespaceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternNamespaceContext

func (s *ExternNamespaceContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternNamespaceContext) AllDOT() []antlr.TerminalNode

func (s *ExternNamespaceContext) AllExternMember() []IExternMemberContext

func (s *ExternNamespaceContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *ExternNamespaceContext) DOT(i int) antlr.TerminalNode

func (s *ExternNamespaceContext) ExternMember(i int) IExternMemberContext

func (s *ExternNamespaceContext) GetParser() antlr.Parser

func (s *ExternNamespaceContext) GetRuleContext() antlr.RuleContext

func (s *ExternNamespaceContext) IDENTIFIER(i int) antlr.TerminalNode

func (*ExternNamespaceContext) IsExternNamespaceContext()

func (s *ExternNamespaceContext) LBRACE() antlr.TerminalNode

func (s *ExternNamespaceContext) NAMESPACE() antlr.TerminalNode

func (s *ExternNamespaceContext) RBRACE() antlr.TerminalNode

func (s *ExternNamespaceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternParamContext() *ExternParamContext

func NewExternParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternParamContext

func (s *ExternParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternParamContext) ExternType() IExternTypeContext

func (s *ExternParamContext) GetParser() antlr.Parser

func (s *ExternParamContext) GetRuleContext() antlr.RuleContext

func (*ExternParamContext) IsExternParamContext()

func (s *ExternParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternParamListContext() *ExternParamListContext

func NewExternParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternParamListContext

func (s *ExternParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *ExternParamListContext) AllExternParam() []IExternParamContext

func (s *ExternParamListContext) COMMA(i int) antlr.TerminalNode

func (s *ExternParamListContext) ELLIPSIS() antlr.TerminalNode

func (s *ExternParamListContext) ExternParam(i int) IExternParamContext

func (s *ExternParamListContext) GetParser() antlr.Parser

func (s *ExternParamListContext) GetRuleContext() antlr.RuleContext

func (*ExternParamListContext) IsExternParamListContext()

func (s *ExternParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternReturnTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternReturnTypeContext() *ExternReturnTypeContext

func NewExternReturnTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternReturnTypeContext

func (s *ExternReturnTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternReturnTypeContext) CONST() antlr.TerminalNode

func (s *ExternReturnTypeContext) ExternType() IExternTypeContext

func (s *ExternReturnTypeContext) GetParser() antlr.Parser

func (s *ExternReturnTypeContext) GetRuleContext() antlr.RuleContext

func (*ExternReturnTypeContext) IsExternReturnTypeContext()

func (s *ExternReturnTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternStaticMethodContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternStaticMethodContext() *ExternStaticMethodContext

func NewExternStaticMethodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternStaticMethodContext

func (s *ExternStaticMethodContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternStaticMethodContext) ExternParamList() IExternParamListContext

func (s *ExternStaticMethodContext) ExternReturnType() IExternReturnTypeContext

func (s *ExternStaticMethodContext) ExternSymbol() IExternSymbolContext

func (s *ExternStaticMethodContext) FUNC() antlr.TerminalNode

func (s *ExternStaticMethodContext) GetParser() antlr.Parser

func (s *ExternStaticMethodContext) GetRuleContext() antlr.RuleContext

func (s *ExternStaticMethodContext) IDENTIFIER() antlr.TerminalNode

func (*ExternStaticMethodContext) IsExternStaticMethodContext()

func (s *ExternStaticMethodContext) LPAREN() antlr.TerminalNode

func (s *ExternStaticMethodContext) RPAREN() antlr.TerminalNode

func (s *ExternStaticMethodContext) STATIC() antlr.TerminalNode

func (s *ExternStaticMethodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternSymbolContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternSymbolContext() *ExternSymbolContext

func NewExternSymbolContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternSymbolContext

func (s *ExternSymbolContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternSymbolContext) GetParser() antlr.Parser

func (s *ExternSymbolContext) GetRuleContext() antlr.RuleContext

func (*ExternSymbolContext) IsExternSymbolContext()

func (s *ExternSymbolContext) STRING_LIT() antlr.TerminalNode

func (s *ExternSymbolContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternTypeAliasContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternTypeAliasContext() *ExternTypeAliasContext

func NewExternTypeAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternTypeAliasContext

func (s *ExternTypeAliasContext) ASSIGN() antlr.TerminalNode

func (s *ExternTypeAliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternTypeAliasContext) ExternFunctionPtrType() IExternFunctionPtrTypeContext

func (s *ExternTypeAliasContext) GetParser() antlr.Parser

func (s *ExternTypeAliasContext) GetRuleContext() antlr.RuleContext

func (s *ExternTypeAliasContext) IDENTIFIER() antlr.TerminalNode

func (*ExternTypeAliasContext) IsExternTypeAliasContext()

func (s *ExternTypeAliasContext) TYPE() antlr.TerminalNode

func (s *ExternTypeAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ExternTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternTypeContext() *ExternTypeContext

func NewExternTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternTypeContext

func (s *ExternTypeContext) AMP() antlr.TerminalNode

func (s *ExternTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternTypeContext) AllSTAR() []antlr.TerminalNode

func (s *ExternTypeContext) BOOL() antlr.TerminalNode

func (s *ExternTypeContext) BYTE() antlr.TerminalNode

func (s *ExternTypeContext) CHAR() antlr.TerminalNode

func (s *ExternTypeContext) CONST() antlr.TerminalNode

func (s *ExternTypeContext) Expression() IExpressionContext

func (s *ExternTypeContext) ExternType() IExternTypeContext

func (s *ExternTypeContext) GetParser() antlr.Parser

func (s *ExternTypeContext) GetRuleContext() antlr.RuleContext

func (s *ExternTypeContext) IDENTIFIER() antlr.TerminalNode

func (s *ExternTypeContext) ISIZE() antlr.TerminalNode

func (*ExternTypeContext) IsExternTypeContext()

func (s *ExternTypeContext) LBRACKET() antlr.TerminalNode

func (s *ExternTypeContext) PrimitiveType() IPrimitiveTypeContext

func (s *ExternTypeContext) QualifiedName() IQualifiedNameContext

func (s *ExternTypeContext) RBRACKET() antlr.TerminalNode

func (s *ExternTypeContext) STAR(i int) antlr.TerminalNode

func (s *ExternTypeContext) STRING() antlr.TerminalNode

func (s *ExternTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExternTypeContext) USIZE() antlr.TerminalNode

func (s *ExternTypeContext) VOID() antlr.TerminalNode

type ExternVirtualMethodContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyExternVirtualMethodContext() *ExternVirtualMethodContext

func NewExternVirtualMethodContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternVirtualMethodContext

func (s *ExternVirtualMethodContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ExternVirtualMethodContext) CallingConvention() ICallingConventionContext

func (s *ExternVirtualMethodContext) ExternMethodParamList() IExternMethodParamListContext

func (s *ExternVirtualMethodContext) ExternReturnType() IExternReturnTypeContext

func (s *ExternVirtualMethodContext) FUNC() antlr.TerminalNode

func (s *ExternVirtualMethodContext) GetParser() antlr.Parser

func (s *ExternVirtualMethodContext) GetRuleContext() antlr.RuleContext

func (s *ExternVirtualMethodContext) IDENTIFIER() antlr.TerminalNode

func (*ExternVirtualMethodContext) IsExternVirtualMethodContext()

func (s *ExternVirtualMethodContext) LPAREN() antlr.TerminalNode

func (s *ExternVirtualMethodContext) RPAREN() antlr.TerminalNode

func (s *ExternVirtualMethodContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ExternVirtualMethodContext) VIRTUAL() antlr.TerminalNode

type FalseLiteralContext struct {
	PrimaryContext
}

func NewFalseLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FalseLiteralContext

func (s *FalseLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FalseLiteralContext) FALSE() antlr.TerminalNode

func (s *FalseLiteralContext) GetRuleContext() antlr.RuleContext

type FieldInitContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFieldInitContext() *FieldInitContext

func NewFieldInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldInitContext

func (s *FieldInitContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FieldInitContext) COLON() antlr.TerminalNode

func (s *FieldInitContext) Expression() IExpressionContext

func (s *FieldInitContext) GetParser() antlr.Parser

func (s *FieldInitContext) GetRuleContext() antlr.RuleContext

func (s *FieldInitContext) IDENTIFIER() antlr.TerminalNode

func (*FieldInitContext) IsFieldInitContext()

func (s *FieldInitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FloatLiteralContext struct {
	PrimaryContext
}

func NewFloatLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FloatLiteralContext

func (s *FloatLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FloatLiteralContext) FLOAT_LIT() antlr.TerminalNode

func (s *FloatLiteralContext) GetRuleContext() antlr.RuleContext

type ForHeaderContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForHeaderContext() *ForHeaderContext

func NewForHeaderContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForHeaderContext

func (s *ForHeaderContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForHeaderContext) AllSEMI() []antlr.TerminalNode

func (s *ForHeaderContext) Expression() IExpressionContext

func (s *ForHeaderContext) ForInit() IForInitContext

func (s *ForHeaderContext) ForIterator() IForIteratorContext

func (s *ForHeaderContext) ForPost() IForPostContext

func (s *ForHeaderContext) GetParser() antlr.Parser

func (s *ForHeaderContext) GetRuleContext() antlr.RuleContext

func (*ForHeaderContext) IsForHeaderContext()

func (s *ForHeaderContext) SEMI(i int) antlr.TerminalNode

func (s *ForHeaderContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ForInitContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForInitContext() *ForInitContext

func NewForInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForInitContext

func (s *ForInitContext) ASSIGN() antlr.TerminalNode

func (s *ForInitContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForInitContext) COLON() antlr.TerminalNode

func (s *ForInitContext) Expression() IExpressionContext

func (s *ForInitContext) GetParser() antlr.Parser

func (s *ForInitContext) GetRuleContext() antlr.RuleContext

func (s *ForInitContext) IDENTIFIER() antlr.TerminalNode

func (*ForInitContext) IsForInitContext()

func (s *ForInitContext) LET() antlr.TerminalNode

func (s *ForInitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ForInitContext) TypeRef() ITypeRefContext

type ForIteratorContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForIteratorContext() *ForIteratorContext

func NewForIteratorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForIteratorContext

func (s *ForIteratorContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForIteratorContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *ForIteratorContext) COMMA() antlr.TerminalNode

func (s *ForIteratorContext) Expression() IExpressionContext

func (s *ForIteratorContext) GetParser() antlr.Parser

func (s *ForIteratorContext) GetRuleContext() antlr.RuleContext

func (s *ForIteratorContext) IDENTIFIER(i int) antlr.TerminalNode

func (s *ForIteratorContext) IN() antlr.TerminalNode

func (*ForIteratorContext) IsForIteratorContext()

func (s *ForIteratorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ForPostContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForPostContext() *ForPostContext

func NewForPostContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForPostContext

func (s *ForPostContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForPostContext) AssignOp() IAssignOpContext

func (s *ForPostContext) AssignmentTarget() IAssignmentTargetContext

func (s *ForPostContext) DEC() antlr.TerminalNode

func (s *ForPostContext) Expression() IExpressionContext

func (s *ForPostContext) GetParser() antlr.Parser

func (s *ForPostContext) GetRuleContext() antlr.RuleContext

func (s *ForPostContext) INC() antlr.TerminalNode

func (*ForPostContext) IsForPostContext()

func (s *ForPostContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ForStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyForStatementContext() *ForStatementContext

func NewForStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ForStatementContext

func (s *ForStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ForStatementContext) Block() IBlockContext

func (s *ForStatementContext) FOR() antlr.TerminalNode

func (s *ForStatementContext) ForHeader() IForHeaderContext

func (s *ForStatementContext) GetParser() antlr.Parser

func (s *ForStatementContext) GetRuleContext() antlr.RuleContext

func (*ForStatementContext) IsForStatementContext()

func (s *ForStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncDeclContext() *FuncDeclContext

func NewFuncDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncDeclContext

func (s *FuncDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncDeclContext) AllFuncModifier() []IFuncModifierContext

func (s *FuncDeclContext) Block() IBlockContext

func (s *FuncDeclContext) FUNC() antlr.TerminalNode

func (s *FuncDeclContext) FuncModifier(i int) IFuncModifierContext

func (s *FuncDeclContext) GenericParams() IGenericParamsContext

func (s *FuncDeclContext) GetParser() antlr.Parser

func (s *FuncDeclContext) GetRuleContext() antlr.RuleContext

func (s *FuncDeclContext) IDENTIFIER() antlr.TerminalNode

func (*FuncDeclContext) IsFuncDeclContext()

func (s *FuncDeclContext) LPAREN() antlr.TerminalNode

func (s *FuncDeclContext) ParamList() IParamListContext

func (s *FuncDeclContext) RPAREN() antlr.TerminalNode

func (s *FuncDeclContext) ReturnType() IReturnTypeContext

func (s *FuncDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FuncModifierContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFuncModifierContext() *FuncModifierContext

func NewFuncModifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FuncModifierContext

func (s *FuncModifierContext) ASYNC() antlr.TerminalNode

func (s *FuncModifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FuncModifierContext) GPU() antlr.TerminalNode

func (s *FuncModifierContext) GetParser() antlr.Parser

func (s *FuncModifierContext) GetRuleContext() antlr.RuleContext

func (*FuncModifierContext) IsFuncModifierContext()

func (s *FuncModifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type FunctionTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyFunctionTypeContext() *FunctionTypeContext

func NewFunctionTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionTypeContext

func (s *FunctionTypeContext) ASYNC() antlr.TerminalNode

func (s *FunctionTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *FunctionTypeContext) FUNC() antlr.TerminalNode

func (s *FunctionTypeContext) GetParser() antlr.Parser

func (s *FunctionTypeContext) GetRuleContext() antlr.RuleContext

func (*FunctionTypeContext) IsFunctionTypeContext()

func (s *FunctionTypeContext) LPAREN() antlr.TerminalNode

func (s *FunctionTypeContext) RPAREN() antlr.TerminalNode

func (s *FunctionTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *FunctionTypeContext) TypeList() ITypeListContext

func (s *FunctionTypeContext) TypeRef() ITypeRefContext

type GenericArgsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericArgsContext() *GenericArgsContext

func NewGenericArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericArgsContext

func (s *GenericArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericArgsContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericArgsContext) AllTypeRef() []ITypeRefContext

func (s *GenericArgsContext) COMMA(i int) antlr.TerminalNode

func (s *GenericArgsContext) GetParser() antlr.Parser

func (s *GenericArgsContext) GetRuleContext() antlr.RuleContext

func (*GenericArgsContext) IsGenericArgsContext()

func (s *GenericArgsContext) LBRACKET() antlr.TerminalNode

func (s *GenericArgsContext) RBRACKET() antlr.TerminalNode

func (s *GenericArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *GenericArgsContext) TypeRef(i int) ITypeRefContext

type GenericParamsContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyGenericParamsContext() *GenericParamsContext

func NewGenericParamsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GenericParamsContext

func (s *GenericParamsContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *GenericParamsContext) AllCOMMA() []antlr.TerminalNode

func (s *GenericParamsContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *GenericParamsContext) COMMA(i int) antlr.TerminalNode

func (s *GenericParamsContext) GetParser() antlr.Parser

func (s *GenericParamsContext) GetRuleContext() antlr.RuleContext

func (s *GenericParamsContext) IDENTIFIER(i int) antlr.TerminalNode

func (*GenericParamsContext) IsGenericParamsContext()

func (s *GenericParamsContext) LBRACKET() antlr.TerminalNode

func (s *GenericParamsContext) RBRACKET() antlr.TerminalNode

func (s *GenericParamsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type HexLiteralContext struct {
	PrimaryContext
}

func NewHexLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *HexLiteralContext

func (s *HexLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *HexLiteralContext) GetRuleContext() antlr.RuleContext

func (s *HexLiteralContext) HEX_LIT() antlr.TerminalNode

type IArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsArgumentContext differentiates from other interfaces.
	IsArgumentContext()
}
    IArgumentContext is an interface to support dynamic dispatch.

type IArgumentListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllArgument() []IArgumentContext
	Argument(i int) IArgumentContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsArgumentListContext differentiates from other interfaces.
	IsArgumentListContext()
}
    IArgumentListContext is an interface to support dynamic dispatch.

type IAssignOpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASSIGN() antlr.TerminalNode
	ADD_ASSIGN() antlr.TerminalNode
	SUB_ASSIGN() antlr.TerminalNode
	MUL_ASSIGN() antlr.TerminalNode
	DIV_ASSIGN() antlr.TerminalNode
	MOD_ASSIGN() antlr.TerminalNode
	AND_ASSIGN() antlr.TerminalNode
	OR_ASSIGN() antlr.TerminalNode
	XOR_ASSIGN() antlr.TerminalNode
	SHL_ASSIGN() antlr.TerminalNode
	SHR_ASSIGN() antlr.TerminalNode

	// IsAssignOpContext differentiates from other interfaces.
	IsAssignOpContext()
}
    IAssignOpContext is an interface to support dynamic dispatch.

type IAssignmentStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AssignmentTarget() IAssignmentTargetContext
	AssignOp() IAssignOpContext
	Expression() IExpressionContext
	INC() antlr.TerminalNode
	DEC() antlr.TerminalNode

	// IsAssignmentStatementContext differentiates from other interfaces.
	IsAssignmentStatementContext()
}
    IAssignmentStatementContext is an interface to support dynamic dispatch.

type IAssignmentTargetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	DOT() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	RBRACKET() antlr.TerminalNode

	// IsAssignmentTargetContext differentiates from other interfaces.
	IsAssignmentTargetContext()
}
    IAssignmentTargetContext is an interface to support dynamic dispatch.

type IAttributeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsAttributeContext differentiates from other interfaces.
	IsAttributeContext()
}
    IAttributeContext is an interface to support dynamic dispatch.

type IBaseTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimitiveType() IPrimitiveTypeContext
	VOID() antlr.TerminalNode
	BOOL() antlr.TerminalNode
	STRING() antlr.TerminalNode
	BYTE() antlr.TerminalNode
	CHAR() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	GenericArgs() IGenericArgsContext
	IDENTIFIER() antlr.TerminalNode
	VECTOR() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	AllTypeRef() []ITypeRefContext
	TypeRef(i int) ITypeRefContext
	RBRACKET() antlr.TerminalNode
	MAP() antlr.TerminalNode
	Expression() IExpressionContext

	// IsBaseTypeContext differentiates from other interfaces.
	IsBaseTypeContext()
}
    IBaseTypeContext is an interface to support dynamic dispatch.

type IBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsBlockContext differentiates from other interfaces.
	IsBlockContext()
}
    IBlockContext is an interface to support dynamic dispatch.

type IBreakStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	BREAK() antlr.TerminalNode

	// IsBreakStatementContext differentiates from other interfaces.
	IsBreakStatementContext()
}
    IBreakStatementContext is an interface to support dynamic dispatch.

type ICallingConventionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CDECL() antlr.TerminalNode
	STDCALL() antlr.TerminalNode
	THISCALL() antlr.TerminalNode
	VECTORCALL() antlr.TerminalNode
	FASTCALL() antlr.TerminalNode

	// IsCallingConventionContext differentiates from other interfaces.
	IsCallingConventionContext()
}
    ICallingConventionContext is an interface to support dynamic dispatch.

type ICompilationUnitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NamespaceDecl() INamespaceDeclContext
	EOF() antlr.TerminalNode
	AllTopLevelDecl() []ITopLevelDeclContext
	TopLevelDecl(i int) ITopLevelDeclContext

	// IsCompilationUnitContext differentiates from other interfaces.
	IsCompilationUnitContext()
}
    ICompilationUnitContext is an interface to support dynamic dispatch.

type IConstDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CONST() antlr.TerminalNode
	AllConstSpec() []IConstSpecContext
	ConstSpec(i int) IConstSpecContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsConstDeclContext differentiates from other interfaces.
	IsConstDeclContext()
}
    IConstDeclContext is an interface to support dynamic dispatch.

type IConstSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	Semi() ISemiContext
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsConstSpecContext differentiates from other interfaces.
	IsConstSpecContext()
}
    IConstSpecContext is an interface to support dynamic dispatch.

type IContinueStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CONTINUE() antlr.TerminalNode

	// IsContinueStatementContext differentiates from other interfaces.
	IsContinueStatementContext()
}
    IContinueStatementContext is an interface to support dynamic dispatch.

type IDeferStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFER() antlr.TerminalNode
	Expression() IExpressionContext

	// IsDeferStatementContext differentiates from other interfaces.
	IsDeferStatementContext()
}
    IDeferStatementContext is an interface to support dynamic dispatch.

type IDeinitDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEINIT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	SelfParam() ISelfParamContext
	RPAREN() antlr.TerminalNode
	Block() IBlockContext

	// IsDeinitDeclContext differentiates from other interfaces.
	IsDeinitDeclContext()
}
    IDeinitDeclContext is an interface to support dynamic dispatch.

type IEnumDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ENUM() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	COLON() antlr.TerminalNode
	PrimitiveType() IPrimitiveTypeContext
	AllEnumMember() []IEnumMemberContext
	EnumMember(i int) IEnumMemberContext
	Semi() ISemiContext

	// IsEnumDeclContext differentiates from other interfaces.
	IsEnumDeclContext()
}
    IEnumDeclContext is an interface to support dynamic dispatch.

type IEnumMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	Semi() ISemiContext
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext

	// IsEnumMemberContext differentiates from other interfaces.
	IsEnumMemberContext()
}
    IEnumMemberContext is an interface to support dynamic dispatch.

type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}
    IExpressionContext is an interface to support dynamic dispatch.

type IExpressionListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExpressionListContext differentiates from other interfaces.
	IsExpressionListContext()
}
    IExpressionListContext is an interface to support dynamic dispatch.

type IExpressionStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsExpressionStatementContext differentiates from other interfaces.
	IsExpressionStatementContext()
}
    IExpressionStatementContext is an interface to support dynamic dispatch.

type IExternClassContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CLASS() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	ABSTRACT() antlr.TerminalNode
	ExternSymbol() IExternSymbolContext
	AllExternClassMember() []IExternClassMemberContext
	ExternClassMember(i int) IExternClassMemberContext

	// IsExternClassContext differentiates from other interfaces.
	IsExternClassContext()
}
    IExternClassContext is an interface to support dynamic dispatch.

type IExternClassMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExternVirtualMethod() IExternVirtualMethodContext
	Semi() ISemiContext
	ExternStaticMethod() IExternStaticMethodContext
	ExternConstructor() IExternConstructorContext
	ExternDestructor() IExternDestructorContext

	// IsExternClassMemberContext differentiates from other interfaces.
	IsExternClassMemberContext()
}
    IExternClassMemberContext is an interface to support dynamic dispatch.

type IExternConstructorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NEW() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ExternType() IExternTypeContext
	ExternParamList() IExternParamListContext

	// IsExternConstructorContext differentiates from other interfaces.
	IsExternConstructorContext()
}
    IExternConstructorContext is an interface to support dynamic dispatch.

type IExternDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXTERN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllExternMember() []IExternMemberContext
	ExternMember(i int) IExternMemberContext

	// IsExternDeclContext differentiates from other interfaces.
	IsExternDeclContext()
}
    IExternDeclContext is an interface to support dynamic dispatch.

type IExternDestructorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DELETE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	ExternMethodParam() IExternMethodParamContext
	RPAREN() antlr.TerminalNode
	VOID() antlr.TerminalNode

	// IsExternDestructorContext differentiates from other interfaces.
	IsExternDestructorContext()
}
    IExternDestructorContext is an interface to support dynamic dispatch.

type IExternFuncDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	CallingConvention() ICallingConventionContext
	ExternSymbol() IExternSymbolContext
	ExternParamList() IExternParamListContext
	ExternReturnType() IExternReturnTypeContext

	// IsExternFuncDeclContext differentiates from other interfaces.
	IsExternFuncDeclContext()
}
    IExternFuncDeclContext is an interface to support dynamic dispatch.

type IExternFunctionPtrTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ExternParamList() IExternParamListContext
	ExternReturnType() IExternReturnTypeContext

	// IsExternFunctionPtrTypeContext differentiates from other interfaces.
	IsExternFunctionPtrTypeContext()
}
    IExternFunctionPtrTypeContext is an interface to support dynamic dispatch.

type IExternMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExternFuncDecl() IExternFuncDeclContext
	Semi() ISemiContext
	ExternTypeAlias() IExternTypeAliasContext
	ExternNamespace() IExternNamespaceContext
	ExternClass() IExternClassContext

	// IsExternMemberContext differentiates from other interfaces.
	IsExternMemberContext()
}
    IExternMemberContext is an interface to support dynamic dispatch.

type IExternMethodParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELF() antlr.TerminalNode
	ExternType() IExternTypeContext

	// IsExternMethodParamContext differentiates from other interfaces.
	IsExternMethodParamContext()
}
    IExternMethodParamContext is an interface to support dynamic dispatch.

type IExternMethodParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExternMethodParam() IExternMethodParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AllExternParam() []IExternParamContext
	ExternParam(i int) IExternParamContext
	ELLIPSIS() antlr.TerminalNode

	// IsExternMethodParamListContext differentiates from other interfaces.
	IsExternMethodParamListContext()
}
    IExternMethodParamListContext is an interface to support dynamic dispatch.

type IExternNamespaceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NAMESPACE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	AllExternMember() []IExternMemberContext
	ExternMember(i int) IExternMemberContext

	// IsExternNamespaceContext differentiates from other interfaces.
	IsExternNamespaceContext()
}
    IExternNamespaceContext is an interface to support dynamic dispatch.

type IExternParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExternType() IExternTypeContext

	// IsExternParamContext differentiates from other interfaces.
	IsExternParamContext()
}
    IExternParamContext is an interface to support dynamic dispatch.

type IExternParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExternParam() []IExternParamContext
	ExternParam(i int) IExternParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	ELLIPSIS() antlr.TerminalNode

	// IsExternParamListContext differentiates from other interfaces.
	IsExternParamListContext()
}
    IExternParamListContext is an interface to support dynamic dispatch.

type IExternReturnTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ExternType() IExternTypeContext
	CONST() antlr.TerminalNode

	// IsExternReturnTypeContext differentiates from other interfaces.
	IsExternReturnTypeContext()
}
    IExternReturnTypeContext is an interface to support dynamic dispatch.

type IExternStaticMethodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STATIC() antlr.TerminalNode
	FUNC() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ExternSymbol() IExternSymbolContext
	ExternParamList() IExternParamListContext
	ExternReturnType() IExternReturnTypeContext

	// IsExternStaticMethodContext differentiates from other interfaces.
	IsExternStaticMethodContext()
}
    IExternStaticMethodContext is an interface to support dynamic dispatch.

type IExternSymbolContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LIT() antlr.TerminalNode

	// IsExternSymbolContext differentiates from other interfaces.
	IsExternSymbolContext()
}
    IExternSymbolContext is an interface to support dynamic dispatch.

type IExternTypeAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	ExternFunctionPtrType() IExternFunctionPtrTypeContext

	// IsExternTypeAliasContext differentiates from other interfaces.
	IsExternTypeAliasContext()
}
    IExternTypeAliasContext is an interface to support dynamic dispatch.

type IExternTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSTAR() []antlr.TerminalNode
	STAR(i int) antlr.TerminalNode
	ExternType() IExternTypeContext
	CONST() antlr.TerminalNode
	AMP() antlr.TerminalNode
	PrimitiveType() IPrimitiveTypeContext
	VOID() antlr.TerminalNode
	BOOL() antlr.TerminalNode
	STRING() antlr.TerminalNode
	BYTE() antlr.TerminalNode
	CHAR() antlr.TerminalNode
	USIZE() antlr.TerminalNode
	ISIZE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	IDENTIFIER() antlr.TerminalNode
	LBRACKET() antlr.TerminalNode
	Expression() IExpressionContext
	RBRACKET() antlr.TerminalNode

	// IsExternTypeContext differentiates from other interfaces.
	IsExternTypeContext()
}
    IExternTypeContext is an interface to support dynamic dispatch.

type IExternVirtualMethodContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VIRTUAL() antlr.TerminalNode
	FUNC() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	CallingConvention() ICallingConventionContext
	ExternMethodParamList() IExternMethodParamListContext
	ExternReturnType() IExternReturnTypeContext

	// IsExternVirtualMethodContext differentiates from other interfaces.
	IsExternVirtualMethodContext()
}
    IExternVirtualMethodContext is an interface to support dynamic dispatch.

type IFieldInitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	Expression() IExpressionContext

	// IsFieldInitContext differentiates from other interfaces.
	IsFieldInitContext()
}
    IFieldInitContext is an interface to support dynamic dispatch.

type IForHeaderContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ForInit() IForInitContext
	AllSEMI() []antlr.TerminalNode
	SEMI(i int) antlr.TerminalNode
	Expression() IExpressionContext
	ForPost() IForPostContext
	ForIterator() IForIteratorContext

	// IsForHeaderContext differentiates from other interfaces.
	IsForHeaderContext()
}
    IForHeaderContext is an interface to support dynamic dispatch.

type IForInitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LET() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsForInitContext differentiates from other interfaces.
	IsForInitContext()
}
    IForInitContext is an interface to support dynamic dispatch.

type IForIteratorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	IN() antlr.TerminalNode
	Expression() IExpressionContext
	COMMA() antlr.TerminalNode

	// IsForIteratorContext differentiates from other interfaces.
	IsForIteratorContext()
}
    IForIteratorContext is an interface to support dynamic dispatch.

type IForPostContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	AssignmentTarget() IAssignmentTargetContext
	AssignOp() IAssignOpContext
	INC() antlr.TerminalNode
	DEC() antlr.TerminalNode

	// IsForPostContext differentiates from other interfaces.
	IsForPostContext()
}
    IForPostContext is an interface to support dynamic dispatch.

type IForStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FOR() antlr.TerminalNode
	ForHeader() IForHeaderContext
	Block() IBlockContext

	// IsForStatementContext differentiates from other interfaces.
	IsForStatementContext()
}
    IForStatementContext is an interface to support dynamic dispatch.

type IFuncDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	Block() IBlockContext
	AllFuncModifier() []IFuncModifierContext
	FuncModifier(i int) IFuncModifierContext
	GenericParams() IGenericParamsContext
	ParamList() IParamListContext
	ReturnType() IReturnTypeContext

	// IsFuncDeclContext differentiates from other interfaces.
	IsFuncDeclContext()
}
    IFuncDeclContext is an interface to support dynamic dispatch.

type IFuncModifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ASYNC() antlr.TerminalNode
	GPU() antlr.TerminalNode

	// IsFuncModifierContext differentiates from other interfaces.
	IsFuncModifierContext()
}
    IFuncModifierContext is an interface to support dynamic dispatch.

type IFunctionTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FUNC() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ASYNC() antlr.TerminalNode
	TypeList() ITypeListContext
	TypeRef() ITypeRefContext

	// IsFunctionTypeContext differentiates from other interfaces.
	IsFunctionTypeContext()
}
    IFunctionTypeContext is an interface to support dynamic dispatch.

type IGenericArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	AllTypeRef() []ITypeRefContext
	TypeRef(i int) ITypeRefContext
	RBRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGenericArgsContext differentiates from other interfaces.
	IsGenericArgsContext()
}
    IGenericArgsContext is an interface to support dynamic dispatch.

type IGenericParamsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACKET() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	RBRACKET() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGenericParamsContext differentiates from other interfaces.
	IsGenericParamsContext()
}
    IGenericParamsContext is an interface to support dynamic dispatch.

type IIfStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIF() []antlr.TerminalNode
	IF(i int) antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllBlock() []IBlockContext
	Block(i int) IBlockContext
	AllELSE() []antlr.TerminalNode
	ELSE(i int) antlr.TerminalNode

	// IsIfStatementContext differentiates from other interfaces.
	IsIfStatementContext()
}
    IIfStatementContext is an interface to support dynamic dispatch.

type IImportAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	UNDERSCORE() antlr.TerminalNode
	DOT() antlr.TerminalNode

	// IsImportAliasContext differentiates from other interfaces.
	IsImportAliasContext()
}
    IImportAliasContext is an interface to support dynamic dispatch.

type IImportDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IMPORT() antlr.TerminalNode
	AllImportSpec() []IImportSpecContext
	ImportSpec(i int) IImportSpecContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode

	// IsImportDeclContext differentiates from other interfaces.
	IsImportDeclContext()
}
    IImportDeclContext is an interface to support dynamic dispatch.

type IImportSpecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LIT() antlr.TerminalNode
	Semi() ISemiContext
	ImportAlias() IImportAliasContext

	// IsImportSpecContext differentiates from other interfaces.
	IsImportSpecContext()
}
    IImportSpecContext is an interface to support dynamic dispatch.

type IInitializerBlockContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllFieldInit() []IFieldInitContext
	FieldInit(i int) IFieldInitContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	Semi() ISemiContext
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllMapEntry() []IMapEntryContext
	MapEntry(i int) IMapEntryContext

	// IsInitializerBlockContext differentiates from other interfaces.
	IsInitializerBlockContext()
}
    IInitializerBlockContext is an interface to support dynamic dispatch.

type IInterfaceDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INTERFACE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	GenericParams() IGenericParamsContext
	AllInterfaceField() []IInterfaceFieldContext
	InterfaceField(i int) IInterfaceFieldContext

	// IsInterfaceDeclContext differentiates from other interfaces.
	IsInterfaceDeclContext()
}
    IInterfaceDeclContext is an interface to support dynamic dispatch.

type IInterfaceFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext
	Semi() ISemiContext

	// IsInterfaceFieldContext differentiates from other interfaces.
	IsInterfaceFieldContext()
}
    IInterfaceFieldContext is an interface to support dynamic dispatch.

type ILambdaParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsLambdaParamContext differentiates from other interfaces.
	IsLambdaParamContext()
}
    ILambdaParamContext is an interface to support dynamic dispatch.

type ILambdaParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllLambdaParam() []ILambdaParamContext
	LambdaParam(i int) ILambdaParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsLambdaParamListContext differentiates from other interfaces.
	IsLambdaParamListContext()
}
    ILambdaParamListContext is an interface to support dynamic dispatch.

type ILetStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LET() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsLetStatementContext differentiates from other interfaces.
	IsLetStatementContext()
}
    ILetStatementContext is an interface to support dynamic dispatch.

type IMapEntryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	COLON() antlr.TerminalNode

	// IsMapEntryContext differentiates from other interfaces.
	IsMapEntryContext()
}
    IMapEntryContext is an interface to support dynamic dispatch.

type INamespaceDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NAMESPACE() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	Semi() ISemiContext
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsNamespaceDeclContext differentiates from other interfaces.
	IsNamespaceDeclContext()
}
    INamespaceDeclContext is an interface to support dynamic dispatch.

type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelfParam() ISelfParamContext
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ParamType() IParamTypeContext
	ELLIPSIS() antlr.TerminalNode

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}
    IParamContext is an interface to support dynamic dispatch.

type IParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllParam() []IParamContext
	Param(i int) IParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsParamListContext differentiates from other interfaces.
	IsParamListContext()
}
    IParamListContext is an interface to support dynamic dispatch.

type IParamTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AMP() antlr.TerminalNode
	MUT() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsParamTypeContext differentiates from other interfaces.
	IsParamTypeContext()
}
    IParamTypeContext is an interface to support dynamic dispatch.

type IPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPrimaryContext differentiates from other interfaces.
	IsPrimaryContext()
}
    IPrimaryContext is an interface to support dynamic dispatch.

type IPrimitiveTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INT8() antlr.TerminalNode
	INT16() antlr.TerminalNode
	INT32() antlr.TerminalNode
	INT64() antlr.TerminalNode
	UINT8() antlr.TerminalNode
	UINT16() antlr.TerminalNode
	UINT32() antlr.TerminalNode
	UINT64() antlr.TerminalNode
	USIZE() antlr.TerminalNode
	ISIZE() antlr.TerminalNode
	FLOAT32() antlr.TerminalNode
	FLOAT64() antlr.TerminalNode

	// IsPrimitiveTypeContext differentiates from other interfaces.
	IsPrimitiveTypeContext()
}
    IPrimitiveTypeContext is an interface to support dynamic dispatch.

type IQualifiedNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsQualifiedNameContext differentiates from other interfaces.
	IsQualifiedNameContext()
}
    IQualifiedNameContext is an interface to support dynamic dispatch.

type IReturnStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	RETURN() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsReturnStatementContext differentiates from other interfaces.
	IsReturnStatementContext()
}
    IReturnStatementContext is an interface to support dynamic dispatch.

type IReturnTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TupleType() ITupleTypeContext
	TypeRef() ITypeRefContext

	// IsReturnTypeContext differentiates from other interfaces.
	IsReturnTypeContext()
}
    IReturnTypeContext is an interface to support dynamic dispatch.

type ISelfParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELF() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	COLON() antlr.TerminalNode
	ParamType() IParamTypeContext
	AMP() antlr.TerminalNode
	MUT() antlr.TerminalNode

	// IsSelfParamContext differentiates from other interfaces.
	IsSelfParamContext()
}
    ISelfParamContext is an interface to support dynamic dispatch.

type ISemiContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSEMI() []antlr.TerminalNode
	SEMI(i int) antlr.TerminalNode

	// IsSemiContext differentiates from other interfaces.
	IsSemiContext()
}
    ISemiContext is an interface to support dynamic dispatch.

type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LetStatement() ILetStatementContext
	Semi() ISemiContext
	VarStatement() IVarStatementContext
	ConstDecl() IConstDeclContext
	ReturnStatement() IReturnStatementContext
	BreakStatement() IBreakStatementContext
	ContinueStatement() IContinueStatementContext
	DeferStatement() IDeferStatementContext
	IfStatement() IIfStatementContext
	ForStatement() IForStatementContext
	SwitchStatement() ISwitchStatementContext
	AssignmentStatement() IAssignmentStatementContext
	ExpressionStatement() IExpressionStatementContext

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}
    IStatementContext is an interface to support dynamic dispatch.

type ISwitchCaseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	ExpressionList() IExpressionListContext
	COLON() antlr.TerminalNode
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsSwitchCaseContext differentiates from other interfaces.
	IsSwitchCaseContext()
}
    ISwitchCaseContext is an interface to support dynamic dispatch.

type ISwitchDefaultContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DEFAULT() antlr.TerminalNode
	COLON() antlr.TerminalNode
	AllStatement() []IStatementContext
	Statement(i int) IStatementContext

	// IsSwitchDefaultContext differentiates from other interfaces.
	IsSwitchDefaultContext()
}
    ISwitchDefaultContext is an interface to support dynamic dispatch.

type ISwitchStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SWITCH() antlr.TerminalNode
	Expression() IExpressionContext
	LBRACE() antlr.TerminalNode
	RBRACE() antlr.TerminalNode
	AllSwitchCase() []ISwitchCaseContext
	SwitchCase(i int) ISwitchCaseContext
	SwitchDefault() ISwitchDefaultContext

	// IsSwitchStatementContext differentiates from other interfaces.
	IsSwitchStatementContext()
}
    ISwitchStatementContext is an interface to support dynamic dispatch.

type ITopLevelDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ImportDecl() IImportDeclContext
	ConstDecl() IConstDeclContext
	TopLevelVarDecl() ITopLevelVarDeclContext
	TopLevelLetDecl() ITopLevelLetDeclContext
	FuncDecl() IFuncDeclContext
	DeinitDecl() IDeinitDeclContext
	InterfaceDecl() IInterfaceDeclContext
	AllAttribute() []IAttributeContext
	Attribute(i int) IAttributeContext
	EnumDecl() IEnumDeclContext
	TypeAliasDecl() ITypeAliasDeclContext
	ExternDecl() IExternDeclContext
	Semi() ISemiContext

	// IsTopLevelDeclContext differentiates from other interfaces.
	IsTopLevelDeclContext()
}
    ITopLevelDeclContext is an interface to support dynamic dispatch.

type ITopLevelLetDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LET() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	Semi() ISemiContext
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext

	// IsTopLevelLetDeclContext differentiates from other interfaces.
	IsTopLevelLetDeclContext()
}
    ITopLevelLetDeclContext is an interface to support dynamic dispatch.

type ITopLevelVarDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VAR() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	Semi() ISemiContext
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext
	NULL() antlr.TerminalNode

	// IsTopLevelVarDeclContext differentiates from other interfaces.
	IsTopLevelVarDeclContext()
}
    ITopLevelVarDeclContext is an interface to support dynamic dispatch.

type ITupleTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	AllTypeRef() []ITypeRefContext
	TypeRef(i int) ITypeRefContext
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTupleTypeContext differentiates from other interfaces.
	IsTupleTypeContext()
}
    ITupleTypeContext is an interface to support dynamic dispatch.

type ITypeAliasDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TYPE() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	OPAQUE() antlr.TerminalNode
	Semi() ISemiContext
	TypeRef() ITypeRefContext

	// IsTypeAliasDeclContext differentiates from other interfaces.
	IsTypeAliasDeclContext()
}
    ITypeAliasDeclContext is an interface to support dynamic dispatch.

type ITypeListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTypeRef() []ITypeRefContext
	TypeRef(i int) ITypeRefContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTypeListContext differentiates from other interfaces.
	IsTypeListContext()
}
    ITypeListContext is an interface to support dynamic dispatch.

type ITypeRefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	FunctionType() IFunctionTypeContext
	BaseType() IBaseTypeContext

	// IsTypeRefContext differentiates from other interfaces.
	IsTypeRefContext()
}
    ITypeRefContext is an interface to support dynamic dispatch.

type IVarStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	VAR() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	ASSIGN() antlr.TerminalNode
	Expression() IExpressionContext
	COLON() antlr.TerminalNode
	TypeRef() ITypeRefContext
	NULL() antlr.TerminalNode

	// IsVarStatementContext differentiates from other interfaces.
	IsVarStatementContext()
}
    IVarStatementContext is an interface to support dynamic dispatch.

type IdentExprContext struct {
	PrimaryContext
}

func NewIdentExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentExprContext

func (s *IdentExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IdentExprContext) GetRuleContext() antlr.RuleContext

func (s *IdentExprContext) IDENTIFIER() antlr.TerminalNode

type IfStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyIfStatementContext() *IfStatementContext

func NewIfStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IfStatementContext

func (s *IfStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IfStatementContext) AllBlock() []IBlockContext

func (s *IfStatementContext) AllELSE() []antlr.TerminalNode

func (s *IfStatementContext) AllExpression() []IExpressionContext

func (s *IfStatementContext) AllIF() []antlr.TerminalNode

func (s *IfStatementContext) Block(i int) IBlockContext

func (s *IfStatementContext) ELSE(i int) antlr.TerminalNode

func (s *IfStatementContext) Expression(i int) IExpressionContext

func (s *IfStatementContext) GetParser() antlr.Parser

func (s *IfStatementContext) GetRuleContext() antlr.RuleContext

func (s *IfStatementContext) IF(i int) antlr.TerminalNode

func (*IfStatementContext) IsIfStatementContext()

func (s *IfStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImportAliasContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportAliasContext() *ImportAliasContext

func NewImportAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportAliasContext

func (s *ImportAliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportAliasContext) DOT() antlr.TerminalNode

func (s *ImportAliasContext) GetParser() antlr.Parser

func (s *ImportAliasContext) GetRuleContext() antlr.RuleContext

func (s *ImportAliasContext) IDENTIFIER() antlr.TerminalNode

func (*ImportAliasContext) IsImportAliasContext()

func (s *ImportAliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ImportAliasContext) UNDERSCORE() antlr.TerminalNode

type ImportDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportDeclContext() *ImportDeclContext

func NewImportDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportDeclContext

func (s *ImportDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportDeclContext) AllImportSpec() []IImportSpecContext

func (s *ImportDeclContext) GetParser() antlr.Parser

func (s *ImportDeclContext) GetRuleContext() antlr.RuleContext

func (s *ImportDeclContext) IMPORT() antlr.TerminalNode

func (s *ImportDeclContext) ImportSpec(i int) IImportSpecContext

func (*ImportDeclContext) IsImportDeclContext()

func (s *ImportDeclContext) LPAREN() antlr.TerminalNode

func (s *ImportDeclContext) RPAREN() antlr.TerminalNode

func (s *ImportDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ImportSpecContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyImportSpecContext() *ImportSpecContext

func NewImportSpecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ImportSpecContext

func (s *ImportSpecContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ImportSpecContext) GetParser() antlr.Parser

func (s *ImportSpecContext) GetRuleContext() antlr.RuleContext

func (s *ImportSpecContext) ImportAlias() IImportAliasContext

func (*ImportSpecContext) IsImportSpecContext()

func (s *ImportSpecContext) STRING_LIT() antlr.TerminalNode

func (s *ImportSpecContext) Semi() ISemiContext

func (s *ImportSpecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IndexExprContext struct {
	ExpressionContext
}

func NewIndexExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexExprContext

func (s *IndexExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IndexExprContext) AllExpression() []IExpressionContext

func (s *IndexExprContext) Expression(i int) IExpressionContext

func (s *IndexExprContext) GetRuleContext() antlr.RuleContext

func (s *IndexExprContext) LBRACKET() antlr.TerminalNode

func (s *IndexExprContext) RBRACKET() antlr.TerminalNode

type InitializerBlockContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInitializerBlockContext() *InitializerBlockContext

func NewInitializerBlockContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InitializerBlockContext

func (s *InitializerBlockContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InitializerBlockContext) AllCOMMA() []antlr.TerminalNode

func (s *InitializerBlockContext) AllExpression() []IExpressionContext

func (s *InitializerBlockContext) AllFieldInit() []IFieldInitContext

func (s *InitializerBlockContext) AllMapEntry() []IMapEntryContext

func (s *InitializerBlockContext) COMMA(i int) antlr.TerminalNode

func (s *InitializerBlockContext) Expression(i int) IExpressionContext

func (s *InitializerBlockContext) FieldInit(i int) IFieldInitContext

func (s *InitializerBlockContext) GetParser() antlr.Parser

func (s *InitializerBlockContext) GetRuleContext() antlr.RuleContext

func (*InitializerBlockContext) IsInitializerBlockContext()

func (s *InitializerBlockContext) LBRACE() antlr.TerminalNode

func (s *InitializerBlockContext) MapEntry(i int) IMapEntryContext

func (s *InitializerBlockContext) RBRACE() antlr.TerminalNode

func (s *InitializerBlockContext) Semi() ISemiContext

func (s *InitializerBlockContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type IntLiteralContext struct {
	PrimaryContext
}

func NewIntLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntLiteralContext

func (s *IntLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *IntLiteralContext) GetRuleContext() antlr.RuleContext

func (s *IntLiteralContext) INT_LIT() antlr.TerminalNode

type InterfaceDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInterfaceDeclContext() *InterfaceDeclContext

func NewInterfaceDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InterfaceDeclContext

func (s *InterfaceDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InterfaceDeclContext) AllInterfaceField() []IInterfaceFieldContext

func (s *InterfaceDeclContext) GenericParams() IGenericParamsContext

func (s *InterfaceDeclContext) GetParser() antlr.Parser

func (s *InterfaceDeclContext) GetRuleContext() antlr.RuleContext

func (s *InterfaceDeclContext) IDENTIFIER() antlr.TerminalNode

func (s *InterfaceDeclContext) INTERFACE() antlr.TerminalNode

func (s *InterfaceDeclContext) InterfaceField(i int) IInterfaceFieldContext

func (*InterfaceDeclContext) IsInterfaceDeclContext()

func (s *InterfaceDeclContext) LBRACE() antlr.TerminalNode

func (s *InterfaceDeclContext) RBRACE() antlr.TerminalNode

func (s *InterfaceDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type InterfaceFieldContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyInterfaceFieldContext() *InterfaceFieldContext

func NewInterfaceFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InterfaceFieldContext

func (s *InterfaceFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *InterfaceFieldContext) COLON() antlr.TerminalNode

func (s *InterfaceFieldContext) GetParser() antlr.Parser

func (s *InterfaceFieldContext) GetRuleContext() antlr.RuleContext

func (s *InterfaceFieldContext) IDENTIFIER() antlr.TerminalNode

func (*InterfaceFieldContext) IsInterfaceFieldContext()

func (s *InterfaceFieldContext) Semi() ISemiContext

func (s *InterfaceFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *InterfaceFieldContext) TypeRef() ITypeRefContext

type LambdaExprContext struct {
	PrimaryContext
}

func NewLambdaExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LambdaExprContext

func (s *LambdaExprContext) ARROW() antlr.TerminalNode

func (s *LambdaExprContext) ASYNC() antlr.TerminalNode

func (s *LambdaExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LambdaExprContext) Block() IBlockContext

func (s *LambdaExprContext) GetRuleContext() antlr.RuleContext

func (s *LambdaExprContext) LPAREN() antlr.TerminalNode

func (s *LambdaExprContext) LambdaParamList() ILambdaParamListContext

func (s *LambdaExprContext) RPAREN() antlr.TerminalNode

type LambdaParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLambdaParamContext() *LambdaParamContext

func NewLambdaParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LambdaParamContext

func (s *LambdaParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LambdaParamContext) COLON() antlr.TerminalNode

func (s *LambdaParamContext) GetParser() antlr.Parser

func (s *LambdaParamContext) GetRuleContext() antlr.RuleContext

func (s *LambdaParamContext) IDENTIFIER() antlr.TerminalNode

func (*LambdaParamContext) IsLambdaParamContext()

func (s *LambdaParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *LambdaParamContext) TypeRef() ITypeRefContext

type LambdaParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLambdaParamListContext() *LambdaParamListContext

func NewLambdaParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LambdaParamListContext

func (s *LambdaParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LambdaParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *LambdaParamListContext) AllLambdaParam() []ILambdaParamContext

func (s *LambdaParamListContext) COMMA(i int) antlr.TerminalNode

func (s *LambdaParamListContext) GetParser() antlr.Parser

func (s *LambdaParamListContext) GetRuleContext() antlr.RuleContext

func (*LambdaParamListContext) IsLambdaParamListContext()

func (s *LambdaParamListContext) LambdaParam(i int) ILambdaParamContext

func (s *LambdaParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type LetStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyLetStatementContext() *LetStatementContext

func NewLetStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LetStatementContext

func (s *LetStatementContext) ASSIGN() antlr.TerminalNode

func (s *LetStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LetStatementContext) AllCOMMA() []antlr.TerminalNode

func (s *LetStatementContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *LetStatementContext) COLON() antlr.TerminalNode

func (s *LetStatementContext) COMMA(i int) antlr.TerminalNode

func (s *LetStatementContext) Expression() IExpressionContext

func (s *LetStatementContext) GetParser() antlr.Parser

func (s *LetStatementContext) GetRuleContext() antlr.RuleContext

func (s *LetStatementContext) IDENTIFIER(i int) antlr.TerminalNode

func (*LetStatementContext) IsLetStatementContext()

func (s *LetStatementContext) LET() antlr.TerminalNode

func (s *LetStatementContext) LPAREN() antlr.TerminalNode

func (s *LetStatementContext) RPAREN() antlr.TerminalNode

func (s *LetStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *LetStatementContext) TypeRef() ITypeRefContext

type LogicalAndExprContext struct {
	ExpressionContext
}

func NewLogicalAndExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalAndExprContext

func (s *LogicalAndExprContext) AND() antlr.TerminalNode

func (s *LogicalAndExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LogicalAndExprContext) AllExpression() []IExpressionContext

func (s *LogicalAndExprContext) Expression(i int) IExpressionContext

func (s *LogicalAndExprContext) GetRuleContext() antlr.RuleContext

type LogicalNotContext struct {
	ExpressionContext
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext

func (s *LogicalNotContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LogicalNotContext) BANG() antlr.TerminalNode

func (s *LogicalNotContext) Expression() IExpressionContext

func (s *LogicalNotContext) GetRuleContext() antlr.RuleContext

type LogicalOrExprContext struct {
	ExpressionContext
}

func NewLogicalOrExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalOrExprContext

func (s *LogicalOrExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *LogicalOrExprContext) AllExpression() []IExpressionContext

func (s *LogicalOrExprContext) Expression(i int) IExpressionContext

func (s *LogicalOrExprContext) GetRuleContext() antlr.RuleContext

func (s *LogicalOrExprContext) OR() antlr.TerminalNode

type MapEntryContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyMapEntryContext() *MapEntryContext

func NewMapEntryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapEntryContext

func (s *MapEntryContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MapEntryContext) AllExpression() []IExpressionContext

func (s *MapEntryContext) COLON() antlr.TerminalNode

func (s *MapEntryContext) Expression(i int) IExpressionContext

func (s *MapEntryContext) GetParser() antlr.Parser

func (s *MapEntryContext) GetRuleContext() antlr.RuleContext

func (*MapEntryContext) IsMapEntryContext()

func (s *MapEntryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type MapLiteralContext struct {
	PrimaryContext
}

func NewMapLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MapLiteralContext

func (s *MapLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MapLiteralContext) AllTypeRef() []ITypeRefContext

func (s *MapLiteralContext) GetRuleContext() antlr.RuleContext

func (s *MapLiteralContext) InitializerBlock() IInitializerBlockContext

func (s *MapLiteralContext) LBRACKET() antlr.TerminalNode

func (s *MapLiteralContext) MAP() antlr.TerminalNode

func (s *MapLiteralContext) RBRACKET() antlr.TerminalNode

func (s *MapLiteralContext) TypeRef(i int) ITypeRefContext

type MemberAccessContext struct {
	ExpressionContext
}

func NewMemberAccessContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberAccessContext

func (s *MemberAccessContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MemberAccessContext) DOT() antlr.TerminalNode

func (s *MemberAccessContext) Expression() IExpressionContext

func (s *MemberAccessContext) GetRuleContext() antlr.RuleContext

func (s *MemberAccessContext) IDENTIFIER() antlr.TerminalNode

type MulExprContext struct {
	ExpressionContext
	// Has unexported fields.
}

func NewMulExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MulExprContext

func (s *MulExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *MulExprContext) AllExpression() []IExpressionContext

func (s *MulExprContext) Expression(i int) IExpressionContext

func (s *MulExprContext) GetOp() antlr.Token

func (s *MulExprContext) GetRuleContext() antlr.RuleContext

func (s *MulExprContext) PERCENT() antlr.TerminalNode

func (s *MulExprContext) SLASH() antlr.TerminalNode

func (s *MulExprContext) STAR() antlr.TerminalNode

func (s *MulExprContext) SetOp(v antlr.Token)

type NamespaceDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyNamespaceDeclContext() *NamespaceDeclContext

func NewNamespaceDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamespaceDeclContext

func (s *NamespaceDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NamespaceDeclContext) AllDOT() []antlr.TerminalNode

func (s *NamespaceDeclContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *NamespaceDeclContext) DOT(i int) antlr.TerminalNode

func (s *NamespaceDeclContext) GetParser() antlr.Parser

func (s *NamespaceDeclContext) GetRuleContext() antlr.RuleContext

func (s *NamespaceDeclContext) IDENTIFIER(i int) antlr.TerminalNode

func (*NamespaceDeclContext) IsNamespaceDeclContext()

func (s *NamespaceDeclContext) NAMESPACE() antlr.TerminalNode

func (s *NamespaceDeclContext) Semi() ISemiContext

func (s *NamespaceDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type NewArrayExprContext struct {
	PrimaryContext
}

func NewNewArrayExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NewArrayExprContext

func (s *NewArrayExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NewArrayExprContext) Expression() IExpressionContext

func (s *NewArrayExprContext) GetRuleContext() antlr.RuleContext

func (s *NewArrayExprContext) LBRACKET() antlr.TerminalNode

func (s *NewArrayExprContext) NEW() antlr.TerminalNode

func (s *NewArrayExprContext) RBRACKET() antlr.TerminalNode

func (s *NewArrayExprContext) TypeRef() ITypeRefContext

type NewExprContext struct {
	PrimaryContext
}

func NewNewExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NewExprContext

func (s *NewExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NewExprContext) GetRuleContext() antlr.RuleContext

func (s *NewExprContext) InitializerBlock() IInitializerBlockContext

func (s *NewExprContext) NEW() antlr.TerminalNode

func (s *NewExprContext) TypeRef() ITypeRefContext

type NullLiteralContext struct {
	PrimaryContext
}

func NewNullLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullLiteralContext

func (s *NullLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *NullLiteralContext) GetRuleContext() antlr.RuleContext

func (s *NullLiteralContext) NULL() antlr.TerminalNode

type ParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParamContext() *ParamContext

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParamContext) COLON() antlr.TerminalNode

func (s *ParamContext) ELLIPSIS() antlr.TerminalNode

func (s *ParamContext) GetParser() antlr.Parser

func (s *ParamContext) GetRuleContext() antlr.RuleContext

func (s *ParamContext) IDENTIFIER() antlr.TerminalNode

func (*ParamContext) IsParamContext()

func (s *ParamContext) ParamType() IParamTypeContext

func (s *ParamContext) SelfParam() ISelfParamContext

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParamListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParamListContext() *ParamListContext

func NewParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamListContext

func (s *ParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParamListContext) AllCOMMA() []antlr.TerminalNode

func (s *ParamListContext) AllParam() []IParamContext

func (s *ParamListContext) COMMA(i int) antlr.TerminalNode

func (s *ParamListContext) GetParser() antlr.Parser

func (s *ParamListContext) GetRuleContext() antlr.RuleContext

func (*ParamListContext) IsParamListContext()

func (s *ParamListContext) Param(i int) IParamContext

func (s *ParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ParamTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyParamTypeContext() *ParamTypeContext

func NewParamTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamTypeContext

func (s *ParamTypeContext) AMP() antlr.TerminalNode

func (s *ParamTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParamTypeContext) GetParser() antlr.Parser

func (s *ParamTypeContext) GetRuleContext() antlr.RuleContext

func (*ParamTypeContext) IsParamTypeContext()

func (s *ParamTypeContext) MUT() antlr.TerminalNode

func (s *ParamTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ParamTypeContext) TypeRef() ITypeRefContext

type ParenExprContext struct {
	PrimaryContext
}

func NewParenExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenExprContext

func (s *ParenExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ParenExprContext) Expression() IExpressionContext

func (s *ParenExprContext) GetRuleContext() antlr.RuleContext

func (s *ParenExprContext) LPAREN() antlr.TerminalNode

func (s *ParenExprContext) RPAREN() antlr.TerminalNode

type PostDecrementContext struct {
	ExpressionContext
}

func NewPostDecrementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PostDecrementContext

func (s *PostDecrementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostDecrementContext) DEC() antlr.TerminalNode

func (s *PostDecrementContext) Expression() IExpressionContext

func (s *PostDecrementContext) GetRuleContext() antlr.RuleContext

type PostIncrementContext struct {
	ExpressionContext
}

func NewPostIncrementContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PostIncrementContext

func (s *PostIncrementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PostIncrementContext) Expression() IExpressionContext

func (s *PostIncrementContext) GetRuleContext() antlr.RuleContext

func (s *PostIncrementContext) INC() antlr.TerminalNode

type PrimaryContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimaryContext() *PrimaryContext

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext

func (s *PrimaryContext) CopyAll(ctx *PrimaryContext)

func (s *PrimaryContext) GetParser() antlr.Parser

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext

func (*PrimaryContext) IsPrimaryContext()

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type PrimaryExprContext struct {
	ExpressionContext
}

func NewPrimaryExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimaryExprContext

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimaryExprContext) GetRuleContext() antlr.RuleContext

func (s *PrimaryExprContext) Primary() IPrimaryContext

type PrimitiveTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyPrimitiveTypeContext() *PrimitiveTypeContext

func NewPrimitiveTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimitiveTypeContext

func (s *PrimitiveTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimitiveTypeContext) FLOAT32() antlr.TerminalNode

func (s *PrimitiveTypeContext) FLOAT64() antlr.TerminalNode

func (s *PrimitiveTypeContext) GetParser() antlr.Parser

func (s *PrimitiveTypeContext) GetRuleContext() antlr.RuleContext

func (s *PrimitiveTypeContext) INT16() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT32() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT64() antlr.TerminalNode

func (s *PrimitiveTypeContext) INT8() antlr.TerminalNode

func (s *PrimitiveTypeContext) ISIZE() antlr.TerminalNode

func (*PrimitiveTypeContext) IsPrimitiveTypeContext()

func (s *PrimitiveTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *PrimitiveTypeContext) UINT16() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT32() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT64() antlr.TerminalNode

func (s *PrimitiveTypeContext) UINT8() antlr.TerminalNode

func (s *PrimitiveTypeContext) USIZE() antlr.TerminalNode

type PrimitiveTypeExprContext struct {
	PrimaryContext
}

func NewPrimitiveTypeExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimitiveTypeExprContext

func (s *PrimitiveTypeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *PrimitiveTypeExprContext) GetRuleContext() antlr.RuleContext

func (s *PrimitiveTypeExprContext) PrimitiveType() IPrimitiveTypeContext

type ProcessExprContext struct {
	PrimaryContext
}

func NewProcessExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ProcessExprContext

func (s *ProcessExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ProcessExprContext) AllLPAREN() []antlr.TerminalNode

func (s *ProcessExprContext) AllRPAREN() []antlr.TerminalNode

func (s *ProcessExprContext) ArgumentList() IArgumentListContext

func (s *ProcessExprContext) Block() IBlockContext

func (s *ProcessExprContext) FUNC() antlr.TerminalNode

func (s *ProcessExprContext) GetRuleContext() antlr.RuleContext

func (s *ProcessExprContext) LPAREN(i int) antlr.TerminalNode

func (s *ProcessExprContext) PROCESS() antlr.TerminalNode

func (s *ProcessExprContext) ParamList() IParamListContext

func (s *ProcessExprContext) RPAREN(i int) antlr.TerminalNode

func (s *ProcessExprContext) ReturnType() IReturnTypeContext

type QualifiedExprContext struct {
	PrimaryContext
}

func NewQualifiedExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QualifiedExprContext

func (s *QualifiedExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *QualifiedExprContext) GetRuleContext() antlr.RuleContext

func (s *QualifiedExprContext) QualifiedName() IQualifiedNameContext

type QualifiedNameContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyQualifiedNameContext() *QualifiedNameContext

func NewQualifiedNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedNameContext

func (s *QualifiedNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *QualifiedNameContext) AllDOT() []antlr.TerminalNode

func (s *QualifiedNameContext) AllIDENTIFIER() []antlr.TerminalNode

func (s *QualifiedNameContext) DOT(i int) antlr.TerminalNode

func (s *QualifiedNameContext) GetParser() antlr.Parser

func (s *QualifiedNameContext) GetRuleContext() antlr.RuleContext

func (s *QualifiedNameContext) IDENTIFIER(i int) antlr.TerminalNode

func (*QualifiedNameContext) IsQualifiedNameContext()

func (s *QualifiedNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type RangeExprContext struct {
	ExpressionContext
}

func NewRangeExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RangeExprContext

func (s *RangeExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RangeExprContext) AllExpression() []IExpressionContext

func (s *RangeExprContext) Expression(i int) IExpressionContext

func (s *RangeExprContext) GetRuleContext() antlr.RuleContext

func (s *RangeExprContext) RANGE() antlr.TerminalNode

type RelationalExprContext struct {
	ExpressionContext
	// Has unexported fields.
}

func NewRelationalExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RelationalExprContext

func (s *RelationalExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *RelationalExprContext) AllExpression() []IExpressionContext

func (s *RelationalExprContext) Expression(i int) IExpressionContext

func (s *RelationalExprContext) GE() antlr.TerminalNode

func (s *RelationalExprContext) GT() antlr.TerminalNode

func (s *RelationalExprContext) GetOp() antlr.Token

func (s *RelationalExprContext) GetRuleContext() antlr.RuleContext

func (s *RelationalExprContext) LE() antlr.TerminalNode

func (s *RelationalExprContext) LT() antlr.TerminalNode

func (s *RelationalExprContext) SetOp(v antlr.Token)

type ReturnStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyReturnStatementContext() *ReturnStatementContext

func NewReturnStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnStatementContext

func (s *ReturnStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReturnStatementContext) AllCOMMA() []antlr.TerminalNode

func (s *ReturnStatementContext) AllExpression() []IExpressionContext

func (s *ReturnStatementContext) COMMA(i int) antlr.TerminalNode

func (s *ReturnStatementContext) Expression(i int) IExpressionContext

func (s *ReturnStatementContext) GetParser() antlr.Parser

func (s *ReturnStatementContext) GetRuleContext() antlr.RuleContext

func (*ReturnStatementContext) IsReturnStatementContext()

func (s *ReturnStatementContext) LPAREN() antlr.TerminalNode

func (s *ReturnStatementContext) RETURN() antlr.TerminalNode

func (s *ReturnStatementContext) RPAREN() antlr.TerminalNode

func (s *ReturnStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ReturnTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyReturnTypeContext() *ReturnTypeContext

func NewReturnTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ReturnTypeContext

func (s *ReturnTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ReturnTypeContext) GetParser() antlr.Parser

func (s *ReturnTypeContext) GetRuleContext() antlr.RuleContext

func (*ReturnTypeContext) IsReturnTypeContext()

func (s *ReturnTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *ReturnTypeContext) TupleType() ITupleTypeContext

func (s *ReturnTypeContext) TypeRef() ITypeRefContext

type SelfParamContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySelfParamContext() *SelfParamContext

func NewSelfParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelfParamContext

func (s *SelfParamContext) AMP() antlr.TerminalNode

func (s *SelfParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SelfParamContext) COLON() antlr.TerminalNode

func (s *SelfParamContext) GetParser() antlr.Parser

func (s *SelfParamContext) GetRuleContext() antlr.RuleContext

func (s *SelfParamContext) IDENTIFIER() antlr.TerminalNode

func (*SelfParamContext) IsSelfParamContext()

func (s *SelfParamContext) MUT() antlr.TerminalNode

func (s *SelfParamContext) ParamType() IParamTypeContext

func (s *SelfParamContext) SELF() antlr.TerminalNode

func (s *SelfParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SemiContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySemiContext() *SemiContext

func NewSemiContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SemiContext

func (s *SemiContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SemiContext) AllSEMI() []antlr.TerminalNode

func (s *SemiContext) GetParser() antlr.Parser

func (s *SemiContext) GetRuleContext() antlr.RuleContext

func (*SemiContext) IsSemiContext()

func (s *SemiContext) SEMI(i int) antlr.TerminalNode

func (s *SemiContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type ShiftExprContext struct {
	ExpressionContext
	// Has unexported fields.
}

func NewShiftExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShiftExprContext

func (s *ShiftExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *ShiftExprContext) AllExpression() []IExpressionContext

func (s *ShiftExprContext) Expression(i int) IExpressionContext

func (s *ShiftExprContext) GetOp() antlr.Token

func (s *ShiftExprContext) GetRuleContext() antlr.RuleContext

func (s *ShiftExprContext) LSHIFT() antlr.TerminalNode

func (s *ShiftExprContext) RSHIFT() antlr.TerminalNode

func (s *ShiftExprContext) SetOp(v antlr.Token)

type SliceExprContext struct {
	ExpressionContext
}

func NewSliceExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SliceExprContext

func (s *SliceExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SliceExprContext) AllExpression() []IExpressionContext

func (s *SliceExprContext) Expression(i int) IExpressionContext

func (s *SliceExprContext) GetRuleContext() antlr.RuleContext

func (s *SliceExprContext) LBRACKET() antlr.TerminalNode

func (s *SliceExprContext) RANGE() antlr.TerminalNode

func (s *SliceExprContext) RBRACKET() antlr.TerminalNode

type StatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyStatementContext() *StatementContext

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext

func (s *StatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StatementContext) AssignmentStatement() IAssignmentStatementContext

func (s *StatementContext) BreakStatement() IBreakStatementContext

func (s *StatementContext) ConstDecl() IConstDeclContext

func (s *StatementContext) ContinueStatement() IContinueStatementContext

func (s *StatementContext) DeferStatement() IDeferStatementContext

func (s *StatementContext) ExpressionStatement() IExpressionStatementContext

func (s *StatementContext) ForStatement() IForStatementContext

func (s *StatementContext) GetParser() antlr.Parser

func (s *StatementContext) GetRuleContext() antlr.RuleContext

func (s *StatementContext) IfStatement() IIfStatementContext

func (*StatementContext) IsStatementContext()

func (s *StatementContext) LetStatement() ILetStatementContext

func (s *StatementContext) ReturnStatement() IReturnStatementContext

func (s *StatementContext) Semi() ISemiContext

func (s *StatementContext) SwitchStatement() ISwitchStatementContext

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *StatementContext) VarStatement() IVarStatementContext

type StringLiteralContext struct {
	PrimaryContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext

func (s *StringLiteralContext) STRING_LIT() antlr.TerminalNode

type SwitchCaseContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchCaseContext() *SwitchCaseContext

func NewSwitchCaseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchCaseContext

func (s *SwitchCaseContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchCaseContext) AllStatement() []IStatementContext

func (s *SwitchCaseContext) CASE() antlr.TerminalNode

func (s *SwitchCaseContext) COLON() antlr.TerminalNode

func (s *SwitchCaseContext) ExpressionList() IExpressionListContext

func (s *SwitchCaseContext) GetParser() antlr.Parser

func (s *SwitchCaseContext) GetRuleContext() antlr.RuleContext

func (*SwitchCaseContext) IsSwitchCaseContext()

func (s *SwitchCaseContext) Statement(i int) IStatementContext

func (s *SwitchCaseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchDefaultContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchDefaultContext() *SwitchDefaultContext

func NewSwitchDefaultContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchDefaultContext

func (s *SwitchDefaultContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchDefaultContext) AllStatement() []IStatementContext

func (s *SwitchDefaultContext) COLON() antlr.TerminalNode

func (s *SwitchDefaultContext) DEFAULT() antlr.TerminalNode

func (s *SwitchDefaultContext) GetParser() antlr.Parser

func (s *SwitchDefaultContext) GetRuleContext() antlr.RuleContext

func (*SwitchDefaultContext) IsSwitchDefaultContext()

func (s *SwitchDefaultContext) Statement(i int) IStatementContext

func (s *SwitchDefaultContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type SwitchStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptySwitchStatementContext() *SwitchStatementContext

func NewSwitchStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SwitchStatementContext

func (s *SwitchStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *SwitchStatementContext) AllSwitchCase() []ISwitchCaseContext

func (s *SwitchStatementContext) Expression() IExpressionContext

func (s *SwitchStatementContext) GetParser() antlr.Parser

func (s *SwitchStatementContext) GetRuleContext() antlr.RuleContext

func (*SwitchStatementContext) IsSwitchStatementContext()

func (s *SwitchStatementContext) LBRACE() antlr.TerminalNode

func (s *SwitchStatementContext) RBRACE() antlr.TerminalNode

func (s *SwitchStatementContext) SWITCH() antlr.TerminalNode

func (s *SwitchStatementContext) SwitchCase(i int) ISwitchCaseContext

func (s *SwitchStatementContext) SwitchDefault() ISwitchDefaultContext

func (s *SwitchStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TopLevelDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTopLevelDeclContext() *TopLevelDeclContext

func NewTopLevelDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelDeclContext

func (s *TopLevelDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TopLevelDeclContext) AllAttribute() []IAttributeContext

func (s *TopLevelDeclContext) Attribute(i int) IAttributeContext

func (s *TopLevelDeclContext) ConstDecl() IConstDeclContext

func (s *TopLevelDeclContext) DeinitDecl() IDeinitDeclContext

func (s *TopLevelDeclContext) EnumDecl() IEnumDeclContext

func (s *TopLevelDeclContext) ExternDecl() IExternDeclContext

func (s *TopLevelDeclContext) FuncDecl() IFuncDeclContext

func (s *TopLevelDeclContext) GetParser() antlr.Parser

func (s *TopLevelDeclContext) GetRuleContext() antlr.RuleContext

func (s *TopLevelDeclContext) ImportDecl() IImportDeclContext

func (s *TopLevelDeclContext) InterfaceDecl() IInterfaceDeclContext

func (*TopLevelDeclContext) IsTopLevelDeclContext()

func (s *TopLevelDeclContext) Semi() ISemiContext

func (s *TopLevelDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TopLevelDeclContext) TopLevelLetDecl() ITopLevelLetDeclContext

func (s *TopLevelDeclContext) TopLevelVarDecl() ITopLevelVarDeclContext

func (s *TopLevelDeclContext) TypeAliasDecl() ITypeAliasDeclContext

type TopLevelLetDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTopLevelLetDeclContext() *TopLevelLetDeclContext

func NewTopLevelLetDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelLetDeclContext

func (s *TopLevelLetDeclContext) ASSIGN() antlr.TerminalNode

func (s *TopLevelLetDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TopLevelLetDeclContext) COLON() antlr.TerminalNode

func (s *TopLevelLetDeclContext) Expression() IExpressionContext

func (s *TopLevelLetDeclContext) GetParser() antlr.Parser

func (s *TopLevelLetDeclContext) GetRuleContext() antlr.RuleContext

func (s *TopLevelLetDeclContext) IDENTIFIER() antlr.TerminalNode

func (*TopLevelLetDeclContext) IsTopLevelLetDeclContext()

func (s *TopLevelLetDeclContext) LET() antlr.TerminalNode

func (s *TopLevelLetDeclContext) Semi() ISemiContext

func (s *TopLevelLetDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TopLevelLetDeclContext) TypeRef() ITypeRefContext

type TopLevelVarDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTopLevelVarDeclContext() *TopLevelVarDeclContext

func NewTopLevelVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TopLevelVarDeclContext

func (s *TopLevelVarDeclContext) ASSIGN() antlr.TerminalNode

func (s *TopLevelVarDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TopLevelVarDeclContext) COLON() antlr.TerminalNode

func (s *TopLevelVarDeclContext) Expression() IExpressionContext

func (s *TopLevelVarDeclContext) GetParser() antlr.Parser

func (s *TopLevelVarDeclContext) GetRuleContext() antlr.RuleContext

func (s *TopLevelVarDeclContext) IDENTIFIER() antlr.TerminalNode

func (*TopLevelVarDeclContext) IsTopLevelVarDeclContext()

func (s *TopLevelVarDeclContext) NULL() antlr.TerminalNode

func (s *TopLevelVarDeclContext) Semi() ISemiContext

func (s *TopLevelVarDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TopLevelVarDeclContext) TypeRef() ITypeRefContext

func (s *TopLevelVarDeclContext) VAR() antlr.TerminalNode

type TrueLiteralContext struct {
	PrimaryContext
}

func NewTrueLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TrueLiteralContext

func (s *TrueLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TrueLiteralContext) GetRuleContext() antlr.RuleContext

func (s *TrueLiteralContext) TRUE() antlr.TerminalNode

type TupleLiteralContext struct {
	PrimaryContext
}

func NewTupleLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TupleLiteralContext

func (s *TupleLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleLiteralContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleLiteralContext) AllExpression() []IExpressionContext

func (s *TupleLiteralContext) COMMA(i int) antlr.TerminalNode

func (s *TupleLiteralContext) Expression(i int) IExpressionContext

func (s *TupleLiteralContext) GetRuleContext() antlr.RuleContext

func (s *TupleLiteralContext) LPAREN() antlr.TerminalNode

func (s *TupleLiteralContext) RPAREN() antlr.TerminalNode

type TupleTypeContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTupleTypeContext() *TupleTypeContext

func NewTupleTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TupleTypeContext

func (s *TupleTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TupleTypeContext) AllCOMMA() []antlr.TerminalNode

func (s *TupleTypeContext) AllTypeRef() []ITypeRefContext

func (s *TupleTypeContext) COMMA(i int) antlr.TerminalNode

func (s *TupleTypeContext) GetParser() antlr.Parser

func (s *TupleTypeContext) GetRuleContext() antlr.RuleContext

func (*TupleTypeContext) IsTupleTypeContext()

func (s *TupleTypeContext) LPAREN() antlr.TerminalNode

func (s *TupleTypeContext) RPAREN() antlr.TerminalNode

func (s *TupleTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TupleTypeContext) TypeRef(i int) ITypeRefContext

type TypeAliasDeclContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeAliasDeclContext() *TypeAliasDeclContext

func NewTypeAliasDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeAliasDeclContext

func (s *TypeAliasDeclContext) ASSIGN() antlr.TerminalNode

func (s *TypeAliasDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeAliasDeclContext) GetParser() antlr.Parser

func (s *TypeAliasDeclContext) GetRuleContext() antlr.RuleContext

func (s *TypeAliasDeclContext) IDENTIFIER() antlr.TerminalNode

func (*TypeAliasDeclContext) IsTypeAliasDeclContext()

func (s *TypeAliasDeclContext) OPAQUE() antlr.TerminalNode

func (s *TypeAliasDeclContext) Semi() ISemiContext

func (s *TypeAliasDeclContext) TYPE() antlr.TerminalNode

func (s *TypeAliasDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeAliasDeclContext) TypeRef() ITypeRefContext

type TypeListContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeListContext() *TypeListContext

func NewTypeListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeListContext

func (s *TypeListContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeListContext) AllCOMMA() []antlr.TerminalNode

func (s *TypeListContext) AllTypeRef() []ITypeRefContext

func (s *TypeListContext) COMMA(i int) antlr.TerminalNode

func (s *TypeListContext) GetParser() antlr.Parser

func (s *TypeListContext) GetRuleContext() antlr.RuleContext

func (*TypeListContext) IsTypeListContext()

func (s *TypeListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *TypeListContext) TypeRef(i int) ITypeRefContext

type TypeRefContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyTypeRefContext() *TypeRefContext

func NewTypeRefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeRefContext

func (s *TypeRefContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypeRefContext) BaseType() IBaseTypeContext

func (s *TypeRefContext) FunctionType() IFunctionTypeContext

func (s *TypeRefContext) GetParser() antlr.Parser

func (s *TypeRefContext) GetRuleContext() antlr.RuleContext

func (*TypeRefContext) IsTypeRefContext()

func (s *TypeRefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

type TypedInitExprContext struct {
	PrimaryContext
}

func NewTypedInitExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypedInitExprContext

func (s *TypedInitExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *TypedInitExprContext) GenericArgs() IGenericArgsContext

func (s *TypedInitExprContext) GetRuleContext() antlr.RuleContext

func (s *TypedInitExprContext) IDENTIFIER() antlr.TerminalNode

func (s *TypedInitExprContext) InitializerBlock() IInitializerBlockContext

func (s *TypedInitExprContext) QualifiedName() IQualifiedNameContext

type UnaryMinusContext struct {
	ExpressionContext
}

func NewUnaryMinusContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnaryMinusContext

func (s *UnaryMinusContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *UnaryMinusContext) Expression() IExpressionContext

func (s *UnaryMinusContext) GetRuleContext() antlr.RuleContext

func (s *UnaryMinusContext) MINUS() antlr.TerminalNode

type VarStatementContext struct {
	antlr.BaseParserRuleContext
	// Has unexported fields.
}

func NewEmptyVarStatementContext() *VarStatementContext

func NewVarStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarStatementContext

func (s *VarStatementContext) ASSIGN() antlr.TerminalNode

func (s *VarStatementContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VarStatementContext) COLON() antlr.TerminalNode

func (s *VarStatementContext) Expression() IExpressionContext

func (s *VarStatementContext) GetParser() antlr.Parser

func (s *VarStatementContext) GetRuleContext() antlr.RuleContext

func (s *VarStatementContext) IDENTIFIER() antlr.TerminalNode

func (*VarStatementContext) IsVarStatementContext()

func (s *VarStatementContext) NULL() antlr.TerminalNode

func (s *VarStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string

func (s *VarStatementContext) TypeRef() ITypeRefContext

func (s *VarStatementContext) VAR() antlr.TerminalNode

type VectorLiteralContext struct {
	PrimaryContext
}

func NewVectorLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VectorLiteralContext

func (s *VectorLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{}

func (s *VectorLiteralContext) GetRuleContext() antlr.RuleContext

func (s *VectorLiteralContext) InitializerBlock() IInitializerBlockContext

func (s *VectorLiteralContext) LBRACKET() antlr.TerminalNode

func (s *VectorLiteralContext) RBRACKET() antlr.TerminalNode

func (s *VectorLiteralContext) TypeRef() ITypeRefContext

func (s *VectorLiteralContext) VECTOR() antlr.TerminalNode

