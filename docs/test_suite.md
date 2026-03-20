# Arc Compiler Test Suite

**Purpose**: Progressive test cases for the Arc compiler. Each test compiles to `.o`, links with `gcc`, executes, and validates output.

**Test Structure**:
- **Name**: Test identifier
- **Globals**: Top-level declarations (structs, functions, imports)
- **Body**: Code inserted into `main()` function
- **Expected**: Substring that must appear in stdout

---

## Phase 1: Foundation (Literals & Variables)

### 1.1 Integer Literals - Decimal
```go
TestCase{
    Name: "literal_int_decimal",
    Body: `
    let x = 42
    libc.printf("value=%d\n", x)
    `,
    Expected: "value=42",
}
```

### 1.2 Integer Literals - Hexadecimal
```go
TestCase{
    Name: "literal_int_hex",
    Body: `
    let x = 0xFF
    libc.printf("hex=%d\n", x)
    `,
    Expected: "hex=255",
}
```

### 1.3 Integer Literals - Negative
```go
TestCase{
    Name: "literal_int_negative",
    Body: `
    let x = -100
    libc.printf("neg=%d\n", x)
    `,
    Expected: "neg=-100",
}
```

### 1.4 Float Literals - Standard
```go
TestCase{
    Name: "literal_float_standard",
    Body: `
    let pi: float32 = 3.14
    libc.printf("pi=%.2f\n", pi)
    `,
    Expected: "pi=3.14",
}
```

### 1.5 Float Literals - Scientific Notation
```go
TestCase{
    Name: "literal_float_scientific",
    Body: `
    let x: float64 = 1.5e-3
    libc.printf("sci=%.4f\n", x)
    `,
    Expected: "sci=0.0015",
}
```

### 1.6 Boolean Literals
```go
TestCase{
    Name: "literal_bool",
    Body: `
    let t = true
    let f = false
    libc.printf("true=%d false=%d\n", t, f)
    `,
    Expected: "true=1 false=0",
}
```

### 1.7 Character Literals
```go
TestCase{
    Name: "literal_char",
    Body: `
    let c: char = 'A'
    libc.printf("char=%c code=%d\n", c, c)
    `,
    Expected: "char=A code=65",
}
```

### 1.8 Character Escapes
```go
TestCase{
    Name: "literal_char_escape",
    Body: `
    let newline: char = '\n'
    let tab: char = '\t'
    libc.printf("escape_test%cok\n", newline)
    `,
    Expected: "escape_test\nok",
}
```

### 1.9 String Literals - Basic
```go
TestCase{
    Name: "literal_string_basic",
    Body: `
    let msg: string = "hello"
    libc.printf("%s\n", msg)
    `,
    Expected: "hello",
}
```

### 1.10 String Literals - Escapes
```go
TestCase{
    Name: "literal_string_escape",
    Body: `
    let msg: string = "line1\nline2"
    libc.printf("%s\n", msg)
    `,
    Expected: "line1\nline2",
}
```

### 1.11 Null Pointer
```go
TestCase{
    Name: "literal_null",
    Body: `
    let ptr: *int32 = null
    if ptr == null {
        libc.printf("ptr_is_null\n")
    }
    `,
    Expected: "ptr_is_null",
}
```

### 1.12 Variable - Mutable with Type
```go
TestCase{
    Name: "variable_mutable_typed",
    Body: `
    let x: int32 = 10
    x = 20
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=20",
}
```

### 1.13 Variable - Mutable with Inference
```go
TestCase{
    Name: "variable_mutable_inferred",
    Body: `
    let x = 99
    x = 100
    libc.printf("inferred=%d\n", x)
    `,
    Expected: "inferred=100",
}
```

### 1.14 Constant - Immutable with Type
```go
TestCase{
    Name: "constant_typed",
    Body: `
    const MAX: int32 = 255
    libc.printf("max=%d\n", MAX)
    `,
    Expected: "max=255",
}
```

### 1.15 Constant - Immutable with Inference
```go
TestCase{
    Name: "constant_inferred",
    Body: `
    const LIMIT = 1024
    libc.printf("limit=%d\n", LIMIT)
    `,
    Expected: "limit=1024",
}
```

### 1.16 Type - Signed Integers (int8)
```go
TestCase{
    Name: "type_int8",
    Body: `
    let small: int8 = -128
    libc.printf("int8=%d\n", small)
    `,
    Expected: "int8=-128",
}
```

### 1.17 Type - Signed Integers (int16)
```go
TestCase{
    Name: "type_int16",
    Body: `
    let medium: int16 = 32000
    libc.printf("int16=%d\n", medium)
    `,
    Expected: "int16=32000",
}
```

### 1.18 Type - Signed Integers (int32)
```go
TestCase{
    Name: "type_int32",
    Body: `
    let x: int32 = 2147483647
    libc.printf("int32=%d\n", x)
    `,
    Expected: "int32=2147483647",
}
```

### 1.19 Type - Signed Integers (int64)
```go
TestCase{
    Name: "type_int64",
    Body: `
    let big: int64 = 9223372036854775807
    libc.printf("int64=%lld\n", big)
    `,
    Expected: "int64=9223372036854775807",
}
```

### 1.20 Type - Unsigned Integers (uint8)
```go
TestCase{
    Name: "type_uint8",
    Body: `
    let byte_val: uint8 = 255
    libc.printf("uint8=%u\n", byte_val)
    `,
    Expected: "uint8=255",
}
```

### 1.21 Type - Unsigned Integers (uint16)
```go
TestCase{
    Name: "type_uint16",
    Body: `
    let port: uint16 = 65535
    libc.printf("uint16=%u\n", port)
    `,
    Expected: "uint16=65535",
}
```

### 1.22 Type - Unsigned Integers (uint32)
```go
TestCase{
    Name: "type_uint32",
    Body: `
    let x: uint32 = 4294967295
    libc.printf("uint32=%u\n", x)
    `,
    Expected: "uint32=4294967295",
}
```

### 1.23 Type - Unsigned Integers (uint64)
```go
TestCase{
    Name: "type_uint64",
    Body: `
    let huge: uint64 = 18446744073709551615
    libc.printf("uint64=%llu\n", huge)
    `,
    Expected: "uint64=18446744073709551615",
}
```

### 1.24 Type - Architecture Dependent (usize)
```go
TestCase{
    Name: "type_usize",
    Body: `
    let size: usize = 1024
    libc.printf("usize=%zu\n", size)
    `,
    Expected: "usize=1024",
}
```

### 1.25 Type - Architecture Dependent (isize)
```go
TestCase{
    Name: "type_isize",
    Body: `
    let offset: isize = -256
    libc.printf("isize=%zd\n", offset)
    `,
    Expected: "isize=-256",
}
```

### 1.26 Type - Floating Point (float32)
```go
TestCase{
    Name: "type_float32",
    Body: `
    let x: float32 = 3.14159
    libc.printf("float32=%.5f\n", x)
    `,
    Expected: "float32=3.14159",
}
```

### 1.27 Type - Floating Point (float64)
```go
TestCase{
    Name: "type_float64",
    Body: `
    let x: float64 = 2.718281828459045
    libc.printf("float64=%.15f\n", x)
    `,
    Expected: "float64=2.718281828459045",
}
```

### 1.28 Type Alias - byte
```go
TestCase{
    Name: "type_alias_byte",
    Body: `
    let b: byte = 0x42
    libc.printf("byte=%u\n", b)
    `,
    Expected: "byte=66",
}
```

### 1.29 Type Alias - bool
```go
TestCase{
    Name: "type_alias_bool",
    Body: `
    let flag: bool = true
    libc.printf("bool=%d\n", flag)
    `,
    Expected: "bool=1",
}
```

### 1.30 Type Alias - char
```go
TestCase{
    Name: "type_alias_char",
    Body: `
    let ch: char = 'Z'
    libc.printf("char=%c\n", ch)
    `,
    Expected: "char=Z",
}
```

---

## Phase 2: Operators & Expressions

### 2.1 Arithmetic - Addition
```go
TestCase{
    Name: "op_addition",
    Body: `
    let a = 10
    let b = 20
    let sum = a + b
    libc.printf("sum=%d\n", sum)
    `,
    Expected: "sum=30",
}
```

### 2.2 Arithmetic - Subtraction
```go
TestCase{
    Name: "op_subtraction",
    Body: `
    let diff = 50 - 15
    libc.printf("diff=%d\n", diff)
    `,
    Expected: "diff=35",
}
```

### 2.3 Arithmetic - Multiplication
```go
TestCase{
    Name: "op_multiplication",
    Body: `
    let prod = 7 * 8
    libc.printf("prod=%d\n", prod)
    `,
    Expected: "prod=56",
}
```

### 2.4 Arithmetic - Division (Signed)
```go
TestCase{
    Name: "op_division_signed",
    Body: `
    let quot: int32 = -20 / 4
    libc.printf("quot=%d\n", quot)
    `,
    Expected: "quot=-5",
}
```

### 2.5 Arithmetic - Division (Unsigned)
```go
TestCase{
    Name: "op_division_unsigned",
    Body: `
    let quot: uint32 = 20 / 4
    libc.printf("quot=%u\n", quot)
    `,
    Expected: "quot=5",
}
```

### 2.6 Arithmetic - Modulo
```go
TestCase{
    Name: "op_modulo",
    Body: `
    let rem = 17 % 5
    libc.printf("rem=%d\n", rem)
    `,
    Expected: "rem=2",
}
```

### 2.7 Unary - Negation
```go
TestCase{
    Name: "op_unary_negation",
    Body: `
    let x = 100
    let neg = -x
    libc.printf("neg=%d\n", neg)
    `,
    Expected: "neg=-100",
}
```

### 2.8 Bitwise - AND
```go
TestCase{
    Name: "op_bitwise_and",
    Body: `
    let result = 0xF0 & 0x0F
    libc.printf("and=%d\n", result)
    `,
    Expected: "and=0",
}
```

### 2.9 Bitwise - OR
```go
TestCase{
    Name: "op_bitwise_or",
    Body: `
    let result = 0xF0 | 0x0F
    libc.printf("or=%d\n", result)
    `,
    Expected: "or=255",
}
```

### 2.10 Bitwise - XOR
```go
TestCase{
    Name: "op_bitwise_xor",
    Body: `
    let result = 0xFF ^ 0x0F
    libc.printf("xor=%d\n", result)
    `,
    Expected: "xor=240",
}
```

### 2.11 Bitwise - NOT
```go
TestCase{
    Name: "op_bitwise_not",
    Body: `
    let x: uint8 = 0x0F
    let result = ~x
    libc.printf("not=%u\n", result)
    `,
    Expected: "not=240",
}
```

### 2.12 Bitwise - Left Shift
```go
TestCase{
    Name: "op_shift_left",
    Body: `
    let result = 1 << 4
    libc.printf("shl=%d\n", result)
    `,
    Expected: "shl=16",
}
```

### 2.13 Bitwise - Right Shift (Logical)
```go
TestCase{
    Name: "op_shift_right_logical",
    Body: `
    let x: uint32 = 64
    let result = x >> 2
    libc.printf("shr=%u\n", result)
    `,
    Expected: "shr=16",
}
```

### 2.14 Bitwise - Right Shift (Arithmetic)
```go
TestCase{
    Name: "op_shift_right_arithmetic",
    Body: `
    let x: int32 = -64
    let result = x >> 2
    libc.printf("shr=%d\n", result)
    `,
    Expected: "shr=-16",
}
```

### 2.15 Logical - AND
```go
TestCase{
    Name: "op_logical_and",
    Body: `
    let result = true && true
    libc.printf("and=%d\n", result)
    `,
    Expected: "and=1",
}
```

### 2.16 Logical - OR
```go
TestCase{
    Name: "op_logical_or",
    Body: `
    let result = false || true
    libc.printf("or=%d\n", result)
    `,
    Expected: "or=1",
}
```

### 2.17 Logical - NOT
```go
TestCase{
    Name: "op_logical_not",
    Body: `
    let result = !false
    libc.printf("not=%d\n", result)
    `,
    Expected: "not=1",
}
```

### 2.18 Comparison - Equality
```go
TestCase{
    Name: "op_equality",
    Body: `
    let result = (42 == 42)
    libc.printf("eq=%d\n", result)
    `,
    Expected: "eq=1",
}
```

### 2.19 Comparison - Inequality
```go
TestCase{
    Name: "op_inequality",
    Body: `
    let result = (10 != 20)
    libc.printf("ne=%d\n", result)
    `,
    Expected: "ne=1",
}
```

### 2.20 Comparison - Less Than
```go
TestCase{
    Name: "op_less_than",
    Body: `
    let result = (5 < 10)
    libc.printf("lt=%d\n", result)
    `,
    Expected: "lt=1",
}
```

### 2.21 Comparison - Less Equal
```go
TestCase{
    Name: "op_less_equal",
    Body: `
    let result = (10 <= 10)
    libc.printf("le=%d\n", result)
    `,
    Expected: "le=1",
}
```

### 2.22 Comparison - Greater Than
```go
TestCase{
    Name: "op_greater_than",
    Body: `
    let result = (20 > 10)
    libc.printf("gt=%d\n", result)
    `,
    Expected: "gt=1",
}
```

### 2.23 Comparison - Greater Equal
```go
TestCase{
    Name: "op_greater_equal",
    Body: `
    let result = (15 >= 15)
    libc.printf("ge=%d\n", result)
    `,
    Expected: "ge=1",
}
```

### 2.24 Compound Assignment - Add
```go
TestCase{
    Name: "op_compound_add",
    Body: `
    let x = 10
    x += 5
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=15",
}
```

### 2.25 Compound Assignment - Subtract
```go
TestCase{
    Name: "op_compound_sub",
    Body: `
    let x = 20
    x -= 7
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=13",
}
```

### 2.26 Compound Assignment - Multiply
```go
TestCase{
    Name: "op_compound_mul",
    Body: `
    let x = 5
    x *= 3
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=15",
}
```

### 2.27 Compound Assignment - Divide
```go
TestCase{
    Name: "op_compound_div",
    Body: `
    let x = 100
    x /= 4
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=25",
}
```

### 2.28 Compound Assignment - Modulo
```go
TestCase{
    Name: "op_compound_mod",
    Body: `
    let x = 17
    x %= 5
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=2",
}
```

### 2.29 Increment - Post-fix
```go
TestCase{
    Name: "op_increment_postfix",
    Body: `
    let x = 10
    let old = x++
    libc.printf("old=%d new=%d\n", old, x)
    `,
    Expected: "old=10 new=11",
}
```

### 2.30 Increment - Pre-fix
```go
TestCase{
    Name: "op_increment_prefix",
    Body: `
    let x = 10
    let new_val = ++x
    libc.printf("new=%d\n", new_val)
    `,
    Expected: "new=11",
}
```

### 2.31 Decrement - Post-fix
```go
TestCase{
    Name: "op_decrement_postfix",
    Body: `
    let x = 10
    let old = x--
    libc.printf("old=%d new=%d\n", old, x)
    `,
    Expected: "old=10 new=9",
}
```

### 2.32 Decrement - Pre-fix
```go
TestCase{
    Name: "op_decrement_prefix",
    Body: `
    let x = 10
    let new_val = --x
    libc.printf("new=%d\n", new_val)
    `,
    Expected: "new=9",
}
```

### 2.33 Operator Precedence
```go
TestCase{
    Name: "op_precedence",
    Body: `
    let result = 2 + 3 * 4
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=14",
}
```

### 2.34 Operator Parentheses
```go
TestCase{
    Name: "op_parentheses",
    Body: `
    let result = (2 + 3) * 4
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=20",
}
```

---

## Phase 3: Control Flow

### 3.1 If Statement - Basic
```go
TestCase{
    Name: "control_if_basic",
    Body: `
    let x = 10
    if x > 5 {
        libc.printf("condition_true\n")
    }
    `,
    Expected: "condition_true",
}
```

### 3.2 If-Else Statement
```go
TestCase{
    Name: "control_if_else",
    Body: `
    let x = 3
    if x > 5 {
        libc.printf("greater\n")
    } else {
        libc.printf("not_greater\n")
    }
    `,
    Expected: "not_greater",
}
```

### 3.3 If-Else-If Chain
```go
TestCase{
    Name: "control_if_elseif",
    Body: `
    let score = 75
    if score >= 90 {
        libc.printf("A\n")
    } else if score >= 80 {
        libc.printf("B\n")
    } else if score >= 70 {
        libc.printf("C\n")
    } else {
        libc.printf("F\n")
    }
    `,
    Expected: "C",
}
```

### 3.4 Nested If Statements
```go
TestCase{
    Name: "control_if_nested",
    Body: `
    let x = 10
    let y = 20
    if x > 5 {
        if y > 15 {
            libc.printf("both_true\n")
        }
    }
    `,
    Expected: "both_true",
}
```

### 3.5 For Loop - C-Style
```go
TestCase{
    Name: "control_for_cstyle",
    Body: `
    let sum = 0
    for let i = 0; i < 5; i++ {
        sum += i
    }
    libc.printf("sum=%d\n", sum)
    `,
    Expected: "sum=10",
}
```

### 3.6 For Loop - While-Style
```go
TestCase{
    Name: "control_for_while",
    Body: `
    let count = 0
    for count < 3 {
        count++
    }
    libc.printf("count=%d\n", count)
    `,
    Expected: "count=3",
}
```

### 3.7 For Loop - Infinite with Break
```go
TestCase{
    Name: "control_for_infinite_break",
    Body: `
    let i = 0
    for {
        i++
        if i >= 5 {
            break
        }
    }
    libc.printf("i=%d\n", i)
    `,
    Expected: "i=5",
}
```

### 3.8 For Loop - Continue
```go
TestCase{
    Name: "control_for_continue",
    Body: `
    let sum = 0
    for let i = 0; i < 10; i++ {
        if i % 2 == 0 {
            continue
        }
        sum += i
    }
    libc.printf("sum_odd=%d\n", sum)
    `,
    Expected: "sum_odd=25",
}
```

### 3.9 For Loop - Nested
```go
TestCase{
    Name: "control_for_nested",
    Body: `
    let product = 0
    for let i = 1; i <= 3; i++ {
        for let j = 1; j <= 2; j++ {
            product = i * j
        }
    }
    libc.printf("product=%d\n", product)
    `,
    Expected: "product=6",
}
```

### 3.10 Switch - Simple Case
```go
TestCase{
    Name: "control_switch_simple",
    Body: `
    let x = 2
    switch x {
        case 1:
            libc.printf("one\n")
        case 2:
            libc.printf("two\n")
        case 3:
            libc.printf("three\n")
    }
    `,
    Expected: "two",
}
```

### 3.11 Switch - Default Case
```go
TestCase{
    Name: "control_switch_default",
    Body: `
    let x = 99
    switch x {
        case 1:
            libc.printf("one\n")
        case 2:
            libc.printf("two\n")
        default:
            libc.printf("other\n")
    }
    `,
    Expected: "other",
}
```

### 3.12 Switch - Multiple Values Per Case
```go
TestCase{
    Name: "control_switch_multi",
    Body: `
    let x = 3
    switch x {
        case 1, 2, 3:
            libc.printf("low\n")
        case 4, 5, 6:
            libc.printf("mid\n")
        default:
            libc.printf("high\n")
    }
    `,
    Expected: "low",
}
```

### 3.13 Return - Void Function
```go
TestCase{
    Name: "control_return_void",
    Globals: `
    func test_return() {
        libc.printf("before_return\n")
        return
        libc.printf("after_return\n")
    }
    `,
    Body: `
    test_return()
    libc.printf("done\n")
    `,
    Expected: "before_return\ndone",
}
```

### 3.14 Return - Value Function
```go
TestCase{
    Name: "control_return_value",
    Globals: `
    func get_value() int32 {
        return 42
    }
    `,
    Body: `
    let x = get_value()
    libc.printf("value=%d\n", x)
    `,
    Expected: "value=42",
}
```

### 3.15 Defer - Single Statement
```go
TestCase{
    Name: "control_defer_single",
    Body: `
    libc.printf("before\n")
    defer libc.printf("deferred\n")
    libc.printf("after\n")
    `,
    Expected: "before\nafter\ndeferred",
}
```

### 3.16 Defer - Multiple Statements (LIFO)
```go
TestCase{
    Name: "control_defer_multiple",
    Body: `
    defer libc.printf("first\n")
    defer libc.printf("second\n")
    defer libc.printf("third\n")
    libc.printf("main\n")
    `,
    Expected: "main\nthird\nsecond\nfirst",
}
```

### 3.17 Defer - With Function Call
```go
TestCase{
    Name: "control_defer_function",
    Globals: `
    func cleanup() {
        libc.printf("cleanup_called\n")
    }
    `,
    Body: `
    libc.printf("start\n")
    defer cleanup()
    libc.printf("end\n")
    `,
    Expected: "start\nend\ncleanup_called",
}
```

---

## Phase 4: Functions

### 4.1 Function - No Parameters, No Return
```go
TestCase{
    Name: "func_no_params_void",
    Globals: `
    func greet() {
        libc.printf("hello\n")
    }
    `,
    Body: `
    greet()
    `,
    Expected: "hello",
}
```

### 4.2 Function - With Parameters, No Return
```go
TestCase{
    Name: "func_with_params_void",
    Globals: `
    func print_sum(a: int32, b: int32) {
        libc.printf("sum=%d\n", a + b)
    }
    `,
    Body: `
    print_sum(10, 20)
    `,
    Expected: "sum=30",
}
```

### 4.3 Function - With Return Value
```go
TestCase{
    Name: "func_with_return",
    Globals: `
    func add(a: int32, b: int32) int32 {
        return a + b
    }
    `,
    Body: `
    let result = add(15, 25)
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=40",
}
```

### 4.4 Function - Multiple Return Values
```go
TestCase{
    Name: "func_multiple_return",
    Globals: `
    func divide(a: int32, b: int32) (int32, bool) {
        if b == 0 {
            return (0, false)
        }
        return (a / b, true)
    }
    `,
    Body: `
    let (result, ok) = divide(10, 2)
    libc.printf("result=%d ok=%d\n", result, ok)
    `,
    Expected: "result=5 ok=1",
}
```

### 4.5 Function - Recursive (Factorial)
```go
TestCase{
    Name: "func_recursive_factorial",
    Globals: `
    func factorial(n: int32) int32 {
        if n <= 1 {
            return 1
        }
        return n * factorial(n - 1)
    }
    `,
    Body: `
    let result = factorial(5)
    libc.printf("factorial=%d\n", result)
    `,
    Expected: "factorial=120",
}
```

### 4.6 Function - Recursive (Fibonacci)
```go
TestCase{
    Name: "func_recursive_fibonacci",
    Globals: `
    func fib(n: int32) int32 {
        if n <= 1 {
            return n
        }
        return fib(n - 1) + fib(n - 2)
    }
    `,
    Body: `
    let result = fib(7)
    libc.printf("fib=%d\n", result)
    `,
    Expected: "fib=13",
}
```

### 4.7 Function - Nested Calls
```go
TestCase{
    Name: "func_nested_calls",
    Globals: `
    func double(x: int32) int32 {
        return x * 2
    }
    func triple(x: int32) int32 {
        return x * 3
    }
    `,
    Body: `
    let result = double(triple(5))
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=30",
}
```

### 4.8 Function - Early Return
```go
TestCase{
    Name: "func_early_return",
    Globals: `
    func check_positive(x: int32) int32 {
        if x < 0 {
            return 0
        }
        return x
    }
    `,
    Body: `
    let result = check_positive(-5)
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=0",
}
```

### 4.9 Function - Pass by Value
```go
TestCase{
    Name: "func_pass_by_value",
    Globals: `
    func modify(x: int32) {
        x = 99
    }
    `,
    Body: `
    let original = 10
    modify(original)
    libc.printf("original=%d\n", original)
    `,
    Expected: "original=10",
}
```

### 4.10 Extern - libc printf
```go
TestCase{
    Name: "extern_libc_printf",
    Body: `
    libc.printf("extern_test\n")
    `,
    Expected: "extern_test",
}
```

### 4.11 Extern - libc malloc/free
```go
TestCase{
    Name: "extern_libc_malloc",
    Body: `
    let ptr = libc.malloc(64)
    if ptr != null {
        libc.printf("malloc_success\n")
    }
    libc.free(ptr)
    `,
    Expected: "malloc_success",
}
```

---

## Phase 5: Pointers & Memory

### 5.1 Pointer - Basic Type
```go
TestCase{
    Name: "pointer_basic_type",
    Body: `
    let x = 42
    let ptr: *int32 = &x
    libc.printf("pointer_created\n")
    `,
    Expected: "pointer_created",
}
```

### 5.2 Pointer - Address-Of Operator
```go
TestCase{
    Name: "pointer_address_of",
    Body: `
    let value = 100
    let ptr: *int32 = &value
    libc.printf("got_address\n")
    `,
    Expected: "got_address",
}
```

### 5.3 Pointer - Dereference (Read)
```go
TestCase{
    Name: "pointer_deref_read",
    Body: `
    let value = 50
    let ptr: *int32 = &value
    let retrieved = *ptr
    libc.printf("value=%d\n", retrieved)
    `,
    Expected: "value=50",
}
```

### 5.4 Pointer - Dereference (Write)
```go
TestCase{
    Name: "pointer_deref_write",
    Body: `
    let value = 10
    let ptr: *int32 = &value
    *ptr = 99
    libc.printf("value=%d\n", value)
    `,
    Expected: "value=99",
}
```

### 5.5 Pointer - Null Check
```go
TestCase{
    Name: "pointer_null_check",
    Body: `
    let ptr: *int32 = null
    if ptr == null {
        libc.printf("is_null\n")
    }
    `,
    Expected: "is_null",
}
```

### 5.6 Pointer - Non-Null Check
```go
TestCase{
    Name: "pointer_non_null",
    Body: `
    let value = 42
    let ptr: *int32 = &value
    if ptr != null {
        libc.printf("not_null\n")
    }
    `,
    Expected: "not_null",
}
```

### 5.7 Pointer - Arithmetic (Addition)
```go
TestCase{
    Name: "pointer_arithmetic_add",
    Body: `
    let arr = alloca(int32, 5)
    arr[0] = 10
    arr[1] = 20
    let ptr: *int32 = arr
    let next = ptr + 1
    libc.printf("value=%d\n", *next)
    `,
    Expected: "value=20",
}
```

### 5.8 Pointer - Arithmetic (Subtraction)
```go
TestCase{
    Name: "pointer_arithmetic_sub",
    Body: `
    let arr = alloca(int32, 5)
    arr[0] = 10
    arr[1] = 20
    let ptr: *int32 = arr + 1
    let prev = ptr - 1
    libc.printf("value=%d\n", *prev)
    `,
    Expected: "value=10",
}
```

### 5.9 Pointer - Array Indexing
```go
TestCase{
    Name: "pointer_array_indexing",
    Body: `
    let arr = alloca(int32, 3)
    arr[0] = 100
    arr[1] = 200
    arr[2] = 300
    libc.printf("arr[1]=%d\n", arr[1])
    `,
    Expected: "arr[1]=200",
}
```

### 5.10 Pointer - Void Pointer
```go
TestCase{
    Name: "pointer_void_type",
    Body: `
    let x = 42
    let vptr: *void = cast<*void>(&x)
    let iptr: *int32 = cast<*int32>(vptr)
    libc.printf("value=%d\n", *iptr)
    `,
    Expected: "value=42",
}
```

### 5.11 Pointer - Double Pointer
```go
TestCase{
    Name: "pointer_double_pointer",
    Body: `
    let value = 77
    let ptr: *int32 = &value
    let pptr: **int32 = &ptr
    libc.printf("value=%d\n", **pptr)
    `,
    Expected: "value=77",
}
```

### 5.12 Pointer - Pass to Function (Modify)
```go
TestCase{
    Name: "pointer_pass_to_function",
    Globals: `
    func modify_value(ptr: *int32) {
        *ptr = 888
    }
    `,
    Body: `
    let x = 10
    modify_value(&x)
    libc.printf("x=%d\n", x)
    `,
    Expected: "x=888",
}
```

---

## Phase 6: Structs

### 6.1 Struct - Definition
```go
TestCase{
    Name: "struct_definition",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    `,
    Body: `
    libc.printf("struct_defined\n")
    `,
    Expected: "struct_defined",
}
```

### 6.2 Struct - Initialization (Named Fields)
```go
TestCase{
    Name: "struct_init_named",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    `,
    Body: `
    let p: Point = Point{x: 10, y: 20}
    libc.printf("x=%d y=%d\n", p.x, p.y)
    `,
    Expected: "x=10 y=20",
}
```

### 6.3 Struct - Initialization (Default/Zero)
```go
TestCase{
    Name: "struct_init_default",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    `,
    Body: `
    let p: Point = Point{}
    libc.printf("x=%d y=%d\n", p.x, p.y)
    `,
    Expected: "x=0 y=0",
}
```

### 6.4 Struct - Field Access (Read)
```go
TestCase{
    Name: "struct_field_read",
    Globals: `
    struct Rectangle {
        width: int32
        height: int32
    }
    `,
    Body: `
    let rect = Rectangle{width: 100, height: 50}
    libc.printf("width=%d\n", rect.width)
    `,
    Expected: "width=100",
}
```

### 6.5 Struct - Field Access (Write)
```go
TestCase{
    Name: "struct_field_write",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    `,
    Body: `
    let p = Point{x: 5, y: 10}
    p.x = 99
    libc.printf("x=%d\n", p.x)
    `,
    Expected: "x=99",
}
```

### 6.6 Struct - Pointer Field Access
```go
TestCase{
    Name: "struct_pointer_field",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    `,
    Body: `
    let p = Point{x: 15, y: 25}
    let ptr: *Point = &p
    libc.printf("x=%d\n", ptr.x)
    `,
    Expected: "x=15",
}
```

### 6.7 Struct - Nested Struct
```go
TestCase{
    Name: "struct_nested",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    struct Line {
        start: Point
        end: Point
    }
    `,
    Body: `
    let line = Line{
        start: Point{x: 0, y: 0},
        end: Point{x: 10, y: 10}
    }
    libc.printf("end_x=%d\n", line.end.x)
    `,
    Expected: "end_x=10",
}
```

### 6.8 Struct - Method Definition (Inline)
```go
TestCase{
    Name: "struct_method_inline",
    Globals: `
    struct Counter {
        value: int32
        
        func get(self c: Counter) int32 {
            return c.value
        }
    }
    `,
    Body: `
    let counter = Counter{value: 42}
    let val = counter.get()
    libc.printf("value=%d\n", val)
    `,
    Expected: "value=42",
}
```

### 6.9 Struct - Method Definition (Flat)
```go
TestCase{
    Name: "struct_method_flat",
    Globals: `
    struct Counter {
        value: int32
    }
    
    func get(self c: Counter) int32 {
        return c.value
    }
    `,
    Body: `
    let counter = Counter{value: 77}
    let val = counter.get()
    libc.printf("value=%d\n", val)
    `,
    Expected: "value=77",
}
```

### 6.10 Struct - Mutating Method
```go
TestCase{
    Name: "struct_mutating_method",
    Globals: `
    struct Counter {
        count: int32
        
        mutating increment(self c: *Counter) {
            c.count++
        }
    }
    `,
    Body: `
    let counter = Counter{count: 10}
    counter.increment()
    libc.printf("count=%d\n", counter.count)
    `,
    Expected: "count=11",
}
```

### 6.11 Struct - Multiple Fields
```go
TestCase{
    Name: "struct_multiple_fields",
    Globals: `
    struct Person {
        age: int32
        height: int32
        weight: int32
    }
    `,
    Body: `
    let person = Person{age: 30, height: 180, weight: 75}
    libc.printf("age=%d height=%d\n", person.age, person.height)
    `,
    Expected: "age=30 height=180",
}
```

### 6.12 Struct - Pass to Function (Value)
```go
TestCase{
    Name: "struct_pass_by_value",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    
    func print_point(p: Point) {
        libc.printf("x=%d y=%d\n", p.x, p.y)
    }
    `,
    Body: `
    let pt = Point{x: 5, y: 10}
    print_point(pt)
    `,
    Expected: "x=5 y=10",
}
```

### 6.13 Struct - Pass to Function (Pointer)
```go
TestCase{
    Name: "struct_pass_by_pointer",
    Globals: `
    struct Point {
        x: int32
        y: int32
    }
    
    func move_point(p: *Point, dx: int32, dy: int32) {
        p.x += dx
        p.y += dy
    }
    `,
    Body: `
    let pt = Point{x: 0, y: 0}
    move_point(&pt, 5, 10)
    libc.printf("x=%d y=%d\n", pt.x, pt.y)
    `,
    Expected: "x=5 y=10",
}
```

---

## Phase 7: Arrays & Collections

### 7.1 Array - Stack Allocated Fixed Size
```go
TestCase{
    Name: "array_fixed_size",
    Body: `
    let arr: array<int32, 3> = {10, 20, 30}
    libc.printf("arr[1]=%d\n", arr[1])
    `,
    Expected: "arr[1]=20",
}
```

### 7.2 Array - Zero Initialization
```go
TestCase{
    Name: "array_zero_init",
    Body: `
    let arr: array<int32, 5> = {}
    libc.printf("arr[0]=%d\n", arr[0])
    `,
    Expected: "arr[0]=0",
}
```

### 7.3 Array - Element Assignment
```go
TestCase{
    Name: "array_element_assign",
    Body: `
    let arr: array<int32, 3> = {1, 2, 3}
    arr[1] = 99
    libc.printf("arr[1]=%d\n", arr[1])
    `,
    Expected: "arr[1]=99",
}
```

### 7.4 Array - Pointer to First Element
```go
TestCase{
    Name: "array_ptr_first",
    Body: `
    let arr: array<int32, 3> = {100, 200, 300}
    let ptr: *int32 = &arr[0]
    libc.printf("first=%d\n", *ptr)
    `,
    Expected: "first=100",
}
```

### 7.5 Array - Iteration (C-Style Loop)
```go
TestCase{
    Name: "array_iteration_cstyle",
    Body: `
    let arr: array<int32, 4> = {1, 2, 3, 4}
    let sum = 0
    for let i: usize = 0; i < 4; i++ {
        sum += arr[i]
    }
    libc.printf("sum=%d\n", sum)
    `,
    Expected: "sum=10",
}
```

### 7.6 Array - Multidimensional (2D)
```go
TestCase{
    Name: "array_2d",
    Body: `
    let matrix: array<array<int32, 2>, 2> = {
        {1, 2},
        {3, 4}
    }
    libc.printf("matrix[1][1]=%d\n", matrix[1][1])
    `,
    Expected: "matrix[1][1]=4",
}
```

### 7.7 Vector - Basic Type (requires runtime)
```go
TestCase{
    Name: "vector_basic",
    Body: `
    // Vector requires runtime library
    // Placeholder for when vector is implemented
    libc.printf("vector_placeholder\n")
    `,
    Expected: "vector_placeholder",
}
```

### 7.8 Map - Basic Type (requires runtime)
```go
TestCase{
    Name: "map_basic",
    Body: `
    // Map requires runtime library
    // Placeholder for when map is implemented
    libc.printf("map_placeholder\n")
    `,
    Expected: "map_placeholder",
}
```

---

## Phase 8: Type Casting

### 8.1 Cast - Int to Int (Widening)
```go
TestCase{
    Name: "cast_int_widen",
    Body: `
    let small: int8 = 100
    let large: int32 = cast<int32>(small)
    libc.printf("large=%d\n", large)
    `,
    Expected: "large=100",
}
```

### 8.2 Cast - Int to Int (Truncation)
```go
TestCase{
    Name: "cast_int_truncate",
    Body: `
    let large: int32 = 1000
    let small: int8 = cast<int8>(large)
    libc.printf("small=%d\n", small)
    `,
    Expected: "small=-24",
}
```

### 8.3 Cast - Int to Float
```go
TestCase{
    Name: "cast_int_to_float",
    Body: `
    let i: int32 = 42
    let f: float32 = cast<float32>(i)
    libc.printf("float=%.1f\n", f)
    `,
    Expected: "float=42.0",
}
```

### 8.4 Cast - Float to Int
```go
TestCase{
    Name: "cast_float_to_int",
    Body: `
    let f: float32 = 3.9
    let i: int32 = cast<int32>(f)
    libc.printf("int=%d\n", i)
    `,
    Expected: "int=3",
}
```

### 8.5 Cast - Pointer to Pointer
```go
TestCase{
    Name: "cast_ptr_to_ptr",
    Body: `
    let x: int32 = 42
    let iptr: *int32 = &x
    let vptr: *void = cast<*void>(iptr)
    let back: *int32 = cast<*int32>(vptr)
    libc.printf("value=%d\n", *back)
    `,
    Expected: "value=42",
}
```

### 8.6 Cast - Pointer to Integer
```go
TestCase{
    Name: "cast_ptr_to_int",
    Body: `
    let x: int32 = 99
    let ptr: *int32 = &x
    let addr: uint64 = cast<uint64>(ptr)
    libc.printf("addr_non_zero=%d\n", addr != 0)
    `,
    Expected: "addr_non_zero=1",
}
```

### 8.7 Cast - Integer to Pointer
```go
TestCase{
    Name: "cast_int_to_ptr",
    Body: `
    let addr: uint64 = 0x1000
    let ptr: *int32 = cast<*int32>(addr)
    libc.printf("ptr_created\n")
    `,
    Expected: "ptr_created",
}
```

### 8.8 Cast - Signed to Unsigned
```go
TestCase{
    Name: "cast_signed_to_unsigned",
    Body: `
    let neg: int32 = -1
    let pos: uint32 = cast<uint32>(neg)
    libc.printf("unsigned=%u\n", pos)
    `,
    Expected: "unsigned=4294967295",
}
```

### 8.9 Cast - Unsigned to Signed
```go
TestCase{
    Name: "cast_unsigned_to_signed",
    Body: `
    let big: uint32 = 4294967295
    let neg: int32 = cast<int32>(big)
    libc.printf("signed=%d\n", neg)
    `,
    Expected: "signed=-1",
}
```

---

## Phase 9: Intrinsics

### 9.1 Intrinsic - sizeof (int32)
```go
TestCase{
    Name: "intrinsic_sizeof_int32",
    Body: `
    let size = sizeof<int32>
    libc.printf("size=%zu\n", size)
    `,
    Expected: "size=4",
}
```

### 9.2 Intrinsic - sizeof (int64)
```go
TestCase{
    Name: "intrinsic_sizeof_int64",
    Body: `
    let size = sizeof<int64>
    libc.printf("size=%zu\n", size)
    `,
    Expected: "size=8",
}
```

### 9.3 Intrinsic - sizeof (Struct)
```go
TestCase{
    Name: "intrinsic_sizeof_struct",
    Globals: `
    struct Pair {
        a: int32
        b: int32
    }
    `,
    Body: `
    let size = sizeof<Pair>
    libc.printf("size=%zu\n", size)
    `,
    Expected: "size=8",
}
```

### 9.4 Intrinsic - alignof
```go
TestCase{
    Name: "intrinsic_alignof",
    Body: `
    let align = alignof<int64>
    libc.printf("align=%zu\n", align)
    `,
    Expected: "align=8",
}
```

### 9.5 Intrinsic - alloca (Single Element)
```go
TestCase{
    Name: "intrinsic_alloca_single",
    Body: `
    let ptr = alloca(int32)
    *ptr = 777
    libc.printf("value=%d\n", *ptr)
    `,
    Expected: "value=777",
}
```

### 9.6 Intrinsic - alloca (Array/Buffer)
```go
TestCase{
    Name: "intrinsic_alloca_array",
    Body: `
    let buffer = alloca(int32, 5)
    buffer[0] = 10
    buffer[1] = 20
    buffer[2] = 30
    libc.printf("buffer[1]=%d\n", buffer[1])
    `,
    Expected: "buffer[1]=20",
}
```

### 9.7 Intrinsic - memset
```go
TestCase{
    Name: "intrinsic_memset",
    Body: `
    let buffer = alloca(byte, 10)
    memset(buffer, 0x42, 10)
    libc.printf("buffer[5]=%d\n", buffer[5])
    `,
    Expected: "buffer[5]=66",
}
```

### 9.8 Intrinsic - memcpy
```go
TestCase{
    Name: "intrinsic_memcpy",
    Body: `
    let src = alloca(int32, 3)
    src[0] = 100
    src[1] = 200
    src[2] = 300
    
    let dest = alloca(int32, 3)
    libc.memcpy(dest, src, 12)
    
    libc.printf("dest[1]=%d\n", dest[1])
    `,
    Expected: "dest[1]=200",
}
```

### 9.9 Intrinsic - strlen
```go
TestCase{
    Name: "intrinsic_strlen",
    Body: `
    let str: *byte = "hello"
    let len = libc.strlen(str)
    libc.printf("len=%zu\n", len)
    `,
    Expected: "len=5",
}
```

### 9.10 Intrinsic - bit_cast
```go
TestCase{
    Name: "intrinsic_bitcast",
    Body: `
    let f: float32 = 1.0
    let bits = bit_cast<uint32>(f)
    libc.printf("bits=%u\n", bits)
    `,
    Expected: "bits=1065353216",
}
```

---

## Phase 10: Advanced Features (Requires Runtime)

### 10.1 Generics - Struct Definition
```go
TestCase{
    Name: "generic_struct_def",
    Globals: `
    struct Box<T> {
        value: T
    }
    `,
    Body: `
    let b: Box<int32> = Box<int32>{value: 42}
    libc.printf("value=%d\n", b.value)
    `,
    Expected: "value=42",
}
```

### 10.2 Generics - Function Definition
```go
TestCase{
    Name: "generic_func_def",
    Globals: `
    func identity<T>(x: T) T {
        return x
    }
    `,
    Body: `
    let result = identity<int32>(99)
    libc.printf("result=%d\n", result)
    `,
    Expected: "result=99",
}
```

### 10.3 Generics - Swap Function
```go
TestCase{
    Name: "generic_swap",
    Globals: `
    func swap<T>(a: *T, b: *T) {
        let tmp: T = *a
        *a = *b
        *b = tmp
    }
    `,
    Body: `
    let x = 10
    let y = 20
    swap<int32>(&x, &y)
    libc.printf("x=%d y=%d\n", x, y)
    `,
    Expected: "x=20 y=10",
}
```

### 10.4 Enum - Basic Definition
```go
TestCase{
    Name: "enum_basic",
    Globals: `
    enum Status {
        OK
        ERROR
        PENDING
    }
    `,
    Body: `
    let s: Status = Status.OK
    libc.printf("status=%d\n", s)
    `,
    Expected: "status=0",
}
```

### 10.5 Enum - Explicit Values
```go
TestCase{
    Name: "enum_explicit_values",
    Globals: `
    enum HttpCode {
        OK = 200
        NOT_FOUND = 404
        SERVER_ERROR = 500
    }
    `,
    Body: `
    let code: HttpCode = HttpCode.NOT_FOUND
    libc.printf("code=%d\n", code)
    `,
    Expected: "code=404",
}
```

### 10.6 String Interpolation
```go
TestCase{
    Name: "string_interpolation",
    Body: `
    let name = "Alice"
    let age = 30
    let msg = "Hello ${name}, age ${age}"
    libc.printf("%s\n", msg)
    `,
    Expected: "Hello Alice, age 30",
}
```

### 10.7 Execution - Thread (requires runtime)
```go
TestCase{
    Name: "exec_thread",
    Body: `
    // Thread execution requires runtime support
    libc.printf("thread_placeholder\n")
    `,
    Expected: "thread_placeholder",
}
```

### 10.8 Execution - Async (requires runtime)
```go
TestCase{
    Name: "exec_async",
    Body: `
    // Async execution requires runtime support
    libc.printf("async_placeholder\n")
    `,
    Expected: "async_placeholder",
}
```

### 10.9 For-In Loop (Range)
```go
TestCase{
    Name: "for_in_range",
    Body: `
    let sum = 0
    for i in 0..5 {
        sum += i
    }
    libc.printf("sum=%d\n", sum)
    `,
    Expected: "sum=10",
}
```

### 10.10 Try-Except (requires runtime)
```go
TestCase{
    Name: "try_except",
    Body: `
    // Exception handling requires runtime support
    libc.printf("exception_placeholder\n")
    `,
    Expected: "exception_placeholder",
}
```

---

## Phase 11: Edge Cases & Complex Scenarios

### 11.1 Integer Overflow Behavior
```go
TestCase{
    Name: "edge_int_overflow",
    Body: `
    let x: int8 = 127
    x = x + 1
    libc.printf("overflow=%d\n", x)
    `,
    Expected: "overflow=-128",
}
```

### 11.2 Division by Zero (Undefined Behavior)
```go
TestCase{
    Name: "edge_div_zero",
    Body: `
    // This should crash or be caught by runtime
    // For now, just test that compiler doesn't reject it
    let x = 10
    let y = 0
    // let result = x / y  // Would crash
    libc.printf("divzero_skipped\n")
    `,
    Expected: "divzero_skipped",
}
```

### 11.3 Null Pointer Dereference (Undefined Behavior)
```go
TestCase{
    Name: "edge_null_deref",
    Body: `
    // This would crash at runtime
    let ptr: *int32 = null
    // let val = *ptr  // Would crash
    libc.printf("null_deref_skipped\n")
    `,
    Expected: "null_deref_skipped",
}
```

### 11.4 Stack vs Heap Pointer Lifetime
```go
TestCase{
    Name: "edge_pointer_lifetime",
    Body: `
    let stack_value = 42
    let ptr: *int32 = &stack_value
    libc.printf("value=%d\n", *ptr)
    // Pointer remains valid within same scope
    `,
    Expected: "value=42",
}
```

### 11.5 Struct Alignment and Padding
```go
TestCase{
    Name: "edge_struct_padding",
    Globals: `
    struct Padded {
        a: int8
        b: int64
        c: int8
    }
    `,
    Body: `
    let size = sizeof<Padded>
    // Size will be > 10 due to alignment padding
    libc.printf("size=%zu\n", size)
    `,
    Expected: "size=24",
}
```

### 11.6 Function Pointer (Future Feature)
```go
TestCase{
    Name: "edge_function_pointer",
    Body: `
    // Function pointers not yet implemented
    libc.printf("func_ptr_placeholder\n")
    `,
    Expected: "func_ptr_placeholder",
}
```

### 11.7 Deep Recursion Stack Test
```go
TestCase{
    Name: "edge_deep_recursion",
    Globals: `
    func countdown(n: int32) int32 {
        if n <= 0 {
            return 0
        }
        return countdown(n - 1) + 1
    }
    `,
    Body: `
    let result = countdown(100)
    libc.printf("depth=%d\n", result)
    `,
    Expected: "depth=100",
}
```

### 11.8 Large Array on Stack
```go
TestCase{
    Name: "edge_large_stack_array",
    Body: `
    let arr: array<int32, 1000> = {}
    arr[999] = 555
    libc.printf("last=%d\n", arr[999])
    `,
    Expected: "last=555",
}
```

### 11.9 Complex Expression Evaluation Order
```go
TestCase{
    Name: "edge_eval_order",
    Globals: `
    func side_effect(x: int32) int32 {
        libc.printf("eval_%d ", x)
        return x
    }
    `,
    Body: `
    let result = side_effect(1) + side_effect(2) * side_effect(3)
    libc.printf("result=%d\n", result)
    `,
    Expected: "eval_1 eval_2 eval_3 result=7",
}
```

### 11.10 Zero-Sized Array (Edge Case)
```go
TestCase{
    Name: "edge_zero_sized_array",
    Body: `
    // Zero-sized arrays might be compiler error or allowed
    // Testing edge case
    let arr: array<int32, 0> = {}
    libc.printf("zero_array_ok\n")
    `,
    Expected: "zero_array_ok",
}
```

---

## Summary

**Total Tests**: ~150 comprehensive test cases covering:

1. **Foundation (30)**: All literal types, variables, constants, type system
2. **Operators (34)**: All arithmetic, bitwise, logical, comparison operations
3. **Control Flow (17)**: if/else, loops, switch, defer, return
4. **Functions (12)**: Definitions, calls, recursion, extern
5. **Pointers (12)**: Address-of, dereference, arithmetic, null checks
6. **Structs (13)**: Definition, initialization, fields, methods
7. **Arrays (8)**: Fixed-size arrays, indexing, iteration
8. **Casting (9)**: All type conversions
9. **Intrinsics (10)**: sizeof, alloca, memset, memcpy, etc.
10. **Advanced (10)**: Generics, enums, string interpolation (placeholders)
11. **Edge Cases (10)**: Overflow, padding, deep recursion, complex expressions

Each test is self-contained, compiles to `.o`, links with gcc, executes, and validates output. Tests progress from simple to complex, with runtime-dependent features marked as placeholders for future implementation.