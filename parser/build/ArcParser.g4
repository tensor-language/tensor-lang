parser grammar ArcParser;

options { tokenVocab = ArcLexer; }

// ═══════════════════════════════════════════════
//  Top Level
// ═══════════════════════════════════════════════

compilationUnit
    : namespaceDecl topLevelDecl* EOF
    ;

namespaceDecl
    : NAMESPACE IDENTIFIER (DOT IDENTIFIER)* semi
    ;

topLevelDecl
    : importDecl
    | constDecl
    | topLevelVarDecl
    | topLevelLetDecl
    | funcDecl
    | deinitDecl
    | attribute* interfaceDecl
    | enumDecl
    | typeAliasDecl
    | externDecl
    | semi
    ;

semi
    : SEMI+
    ;

// ═══════════════════════════════════════════════
//  Imports
// ═══════════════════════════════════════════════

importDecl
    : IMPORT importSpec
    | IMPORT LPAREN importSpec+ RPAREN
    ;

importSpec
    : importAlias? STRING_LIT semi
    ;

importAlias
    : IDENTIFIER
    | UNDERSCORE
    | DOT
    ;

// ═══════════════════════════════════════════════
//  Constants
// ═══════════════════════════════════════════════

constDecl
    : CONST constSpec
    | CONST LPAREN constSpec+ RPAREN
    ;

constSpec
    : IDENTIFIER (COLON typeRef)? ASSIGN expression semi
    ;

// ═══════════════════════════════════════════════
//  Variables (Top Level)
// ═══════════════════════════════════════════════

topLevelVarDecl
    : VAR IDENTIFIER (COLON typeRef)? ASSIGN expression semi
    | VAR IDENTIFIER COLON typeRef ASSIGN NULL semi
    ;

topLevelLetDecl
    : LET IDENTIFIER (COLON typeRef)? ASSIGN expression semi
    ;

// ═══════════════════════════════════════════════
//  Functions
// ═══════════════════════════════════════════════

funcDecl
    : funcModifier* FUNC IDENTIFIER genericParams? LPAREN paramList? RPAREN returnType? block
    ;

funcModifier
    : ASYNC
    | GPU
    ;

deinitDecl
    : DEINIT LPAREN selfParam RPAREN block
    ;

paramList
    : param (COMMA param)* COMMA?
    ;

param
    : selfParam
    | IDENTIFIER COLON paramType
    | ELLIPSIS
    ;

selfParam
    : SELF IDENTIFIER COLON paramType
    | SELF AMP MUT IDENTIFIER COLON paramType
    ;

paramType
    : AMP MUT typeRef
    | typeRef
    ;

returnType
    : tupleType
    | typeRef
    ;

tupleType
    : LPAREN typeRef (COMMA typeRef)+ RPAREN
    ;

genericParams
    : LBRACKET IDENTIFIER (COMMA IDENTIFIER)* RBRACKET
    ;

genericArgs
    : LBRACKET typeRef (COMMA typeRef)* RBRACKET
    ;

// ═══════════════════════════════════════════════
//  Interfaces
// ═══════════════════════════════════════════════

interfaceDecl
    : INTERFACE IDENTIFIER genericParams? LBRACE interfaceField* RBRACE
    ;

interfaceField
    : IDENTIFIER COLON typeRef semi
    ;

// ═══════════════════════════════════════════════
//  Enums
// ═══════════════════════════════════════════════

enumDecl
    : ENUM IDENTIFIER (COLON primitiveType)? LBRACE enumMember+ RBRACE semi?
    ;

enumMember
    : IDENTIFIER (ASSIGN expression)? semi
    ;

// ═══════════════════════════════════════════════
//  Type Alias / Opaque
// ═══════════════════════════════════════════════

typeAliasDecl
    : TYPE IDENTIFIER ASSIGN OPAQUE semi
    | TYPE IDENTIFIER ASSIGN typeRef semi
    ;

// ═══════════════════════════════════════════════
//  Attributes
// ═══════════════════════════════════════════════

attribute
    : AT IDENTIFIER (LPAREN expression RPAREN)?
    ;

// ═══════════════════════════════════════════════
//  Types
// ═══════════════════════════════════════════════

typeRef
    : functionType
    | baseType
    ;

functionType
    : ASYNC? FUNC LPAREN typeList? RPAREN typeRef?
    ;

baseType
    : primitiveType
    | VOID
    | BOOL
    | STRING
    | BYTE
    | CHAR
    | qualifiedName genericArgs?
    | IDENTIFIER genericArgs?
    | VECTOR LBRACKET typeRef RBRACKET
    | MAP LBRACKET typeRef RBRACKET typeRef
    | LBRACKET RBRACKET typeRef
    | LBRACKET expression RBRACKET typeRef
    ;

primitiveType
    : INT8 | INT16 | INT32 | INT64
    | UINT8 | UINT16 | UINT32 | UINT64
    | USIZE | ISIZE
    | FLOAT32 | FLOAT64
    ;

typeList
    : typeRef (COMMA typeRef)*
    ;

// ═══════════════════════════════════════════════
//  Extern Blocks
// ═══════════════════════════════════════════════

externDecl
    : EXTERN IDENTIFIER LBRACE externMember* RBRACE
    ;

externMember
    : externFuncDecl       semi
    | externTypeAlias      semi
    | externNamespace
    | externClass
    ;

// ── Extern Functions ─────────────────────────

externFuncDecl
    : callingConvention? FUNC IDENTIFIER externSymbol? LPAREN externParamList? RPAREN externReturnType?
    ;

callingConvention
    : CDECL | STDCALL | THISCALL | VECTORCALL | FASTCALL
    ;

externSymbol
    : STRING_LIT
    ;

externParamList
    : externParam (COMMA externParam)* (COMMA ELLIPSIS)?
    | ELLIPSIS
    ;

externParam
    : externType
    ;

externReturnType
    : CONST? externType
    ;

externType
    : STAR STAR externType
    | STAR CONST? externType
    | AMP CONST? externType
    | primitiveType
    | VOID
    | BOOL
    | STRING
    | BYTE
    | CHAR
    | USIZE | ISIZE
    | qualifiedName
    | IDENTIFIER
    | LBRACKET expression RBRACKET externType
    ;

// ── Extern Namespaces ────────────────────────

externNamespace
    : NAMESPACE IDENTIFIER (DOT IDENTIFIER)* LBRACE externMember* RBRACE
    ;

// ── Extern Classes ───────────────────────────

externClass
    : ABSTRACT? CLASS IDENTIFIER externSymbol? LBRACE externClassMember* RBRACE
    ;

externClassMember
    : externVirtualMethod   semi
    | externStaticMethod    semi
    | externConstructor     semi
    | externDestructor      semi
    ;

externVirtualMethod
    : callingConvention? VIRTUAL FUNC IDENTIFIER
      LPAREN externMethodParamList? RPAREN externReturnType?
    ;

externStaticMethod
    : STATIC FUNC IDENTIFIER externSymbol?
      LPAREN externParamList? RPAREN externReturnType?
    ;

externConstructor
    : NEW LPAREN externParamList? RPAREN externType
    ;

externDestructor
    : DELETE LPAREN externMethodParam RPAREN VOID?
    ;

externMethodParamList
    : externMethodParam (COMMA externParam)* (COMMA ELLIPSIS)?
    ;

externMethodParam
    : SELF externType
    ;

// ── Extern Type Alias ────────────────────────

externTypeAlias
    : TYPE IDENTIFIER ASSIGN externFunctionPtrType
    ;

externFunctionPtrType
    : FUNC LPAREN externParamList? RPAREN externReturnType?
    ;

// ═══════════════════════════════════════════════
//  Statements
// ═══════════════════════════════════════════════

block
    : LBRACE statement* RBRACE
    ;

statement
    : letStatement        semi
    | varStatement        semi
    | constDecl           // constDecl has its own semi rule
    | returnStatement     semi
    | breakStatement      semi
    | continueStatement   semi
    | deferStatement      semi
    | ifStatement
    | forStatement
    | switchStatement
    | assignmentStatement semi
    | expressionStatement semi
    | semi
    ;

letStatement
    : LET LPAREN IDENTIFIER (COMMA IDENTIFIER)+ RPAREN ASSIGN expression
    | LET IDENTIFIER (COLON typeRef)? ASSIGN expression
    ;

varStatement
    : VAR IDENTIFIER (COLON typeRef)? ASSIGN expression
    | VAR IDENTIFIER COLON typeRef ASSIGN NULL
    ;

returnStatement
    : RETURN expression?
    | RETURN LPAREN expression (COMMA expression)+ RPAREN
    ;

breakStatement
    : BREAK
    ;

continueStatement
    : CONTINUE
    ;

deferStatement
    : DEFER expression
    ;

ifStatement
    : IF expression block (ELSE IF expression block)* (ELSE block)?
    ;

forStatement
    : FOR forHeader block
    ;

forHeader
    : forInit SEMI expression SEMI forPost
    | forIterator
    | expression
    |
    ;

forInit
    : LET IDENTIFIER (COLON typeRef)? ASSIGN expression
    | expression
    ;

forPost
    : expression
    | assignmentTarget assignOp expression
    | expression (INC | DEC)
    ;

forIterator
    : IDENTIFIER (COMMA IDENTIFIER)? IN expression
    ;

switchStatement
    : SWITCH expression LBRACE switchCase* switchDefault? RBRACE
    ;

switchCase
    : CASE expressionList COLON statement*
    ;

switchDefault
    : DEFAULT COLON statement*
    ;

expressionList
    : expression (COMMA expression)*
    ;

assignmentStatement
    : assignmentTarget assignOp expression
    | expression (INC | DEC)
    ;

assignmentTarget
    : expression DOT IDENTIFIER
    | expression LBRACKET expression RBRACKET
    | IDENTIFIER
    ;

assignOp
    : ASSIGN
    | ADD_ASSIGN | SUB_ASSIGN | MUL_ASSIGN | DIV_ASSIGN | MOD_ASSIGN
    | AND_ASSIGN | OR_ASSIGN | XOR_ASSIGN
    | SHL_ASSIGN | SHR_ASSIGN
    ;

expressionStatement
    : expression
    ;

// ═══════════════════════════════════════════════
//  Expressions
// ═══════════════════════════════════════════════

expression
    : primary                                   # PrimaryExpr

    // ── Postfix ──
    | expression DOT IDENTIFIER                 # MemberAccess
    | expression LBRACKET expression RBRACKET   # IndexExpr
    | expression LBRACKET expression RANGE expression RBRACKET # SliceExpr
    | expression LPAREN argumentList? RPAREN    # CallExpr
    | expression INC                            # PostIncrement
    | expression DEC                            # PostDecrement

    // ── Unary ──
    | MINUS expression                          # UnaryMinus
    | BANG expression                           # LogicalNot
    | TILDE expression                          # BitwiseNot
    | AMP expression                            # AddressOf
    | AWAIT expression                          # AwaitExpr

    // ── Binary ──
    | expression op=(STAR | SLASH | PERCENT) expression  # MulExpr
    | expression op=(PLUS | MINUS) expression            # AddExpr
    | expression op=(LSHIFT | RSHIFT) expression         # ShiftExpr
    | expression op=(LT | GT | LE | GE) expression       # RelationalExpr
    | expression op=(EQ | NEQ) expression                # EqualityExpr
    | expression AMP expression                          # BitwiseAndExpr
    | expression CARET expression                        # BitwiseXorExpr
    | expression PIPE expression                         # BitwiseOrExpr
    | expression AND expression                          # LogicalAndExpr
    | expression OR expression                           # LogicalOrExpr

    // ── Range ──
    | expression RANGE expression               # RangeExpr
    ;

primary
    : INT_LIT                   # IntLiteral
    | HEX_LIT                   # HexLiteral
    | FLOAT_LIT                 # FloatLiteral
    | STRING_LIT                # StringLiteral
    | CHAR_LIT                  # CharLiteral
    | TRUE                      # TrueLiteral
    | FALSE                     # FalseLiteral
    | NULL                      # NullLiteral

    // Typed Initializer: Point{...} or routes{...}
    | (qualifiedName | IDENTIFIER) genericArgs? initializerBlock       # TypedInitExpr

    // Bare Initializer: {...} (for array/struct inference)
    | initializerBlock          # BareInitExpr

    | VECTOR LBRACKET typeRef RBRACKET initializerBlock                # VectorLiteral
    | MAP LBRACKET typeRef RBRACKET typeRef initializerBlock           # MapLiteral

    | qualifiedName             # QualifiedExpr
    | IDENTIFIER                # IdentExpr
    | primitiveType             # PrimitiveTypeExpr

    | LPAREN expression RPAREN                          # ParenExpr
    | LPAREN expression (COMMA expression)+ RPAREN      # TupleLiteral

    | NEW typeRef initializerBlock                      # NewExpr
    | NEW LBRACKET expression RBRACKET typeRef          # NewArrayExpr

    | DELETE LPAREN expression RPAREN                   # DeleteExpr

    | ASYNC? LPAREN lambdaParamList? RPAREN ARROW block # LambdaExpr

    | PROCESS FUNC LPAREN paramList? RPAREN returnType? block
      LPAREN argumentList? RPAREN                       # ProcessExpr
    ;

// ── Initializers ─────────────────────────────

initializerBlock
    : LBRACE RBRACE
    | LBRACE fieldInit (COMMA fieldInit)* COMMA? semi? RBRACE
    | LBRACE expression (COMMA expression)* COMMA? semi? RBRACE
    | LBRACE mapEntry (COMMA mapEntry)* COMMA? semi? RBRACE
    ;

fieldInit
    : IDENTIFIER COLON expression
    ;

mapEntry
    : expression COLON expression
    ;

// ── Arguments ────────────────────────────────

argumentList
    : argument (COMMA argument)*
    ;

argument
    : expression
    ;

// ── Lambda Parameters ────────────────────────

lambdaParamList
    : lambdaParam (COMMA lambdaParam)*
    ;

lambdaParam
    : IDENTIFIER COLON typeRef
    ;

// ═══════════════════════════════════════════════
//  Shared
// ═══════════════════════════════════════════════

qualifiedName
    : IDENTIFIER (DOT IDENTIFIER)+
    ;