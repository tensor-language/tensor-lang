# Types

---

## int8
```arc
let x: int8 = -128
let y: int8 = 127
```

A signed 8-bit integer. Holds values from -128 to 127. Rarely used for arithmetic —
its main purpose is interoperating with C APIs that traffic in signed bytes, or packing
data tightly in memory-mapped structures. The compiler will not silently widen it for
you — if you pass an int8 where an int32 is expected, you cast explicitly.

---

## int16
```arc
let x: int16 = -32768
let y: int16 = 32767
```

A signed 16-bit integer. Holds values from -32,768 to 32,767. Used when matching
C structs that contain short fields, or in audio and signal processing where 16-bit
samples are the native unit. Not a general-purpose arithmetic type — prefer int32 for
most integer work.

---

## int32
```arc
let x: int32 = -2147483648
let y: int32 = 2147483647
```

A signed 32-bit integer. The workhorse signed integer type. Holds values from roughly
-2.1 billion to 2.1 billion. This is the default choice for counts, indices, and most
arithmetic where you know values fit in 32 bits. Matches C's int on all major platforms.
Most C library functions return int32.

---

## int64
```arc
let x: int64 = -9223372036854775808
let y: int64 =  9223372036854775807
```

A signed 64-bit integer. Holds values up to roughly 9.2 quintillion in either direction.
Use when values may exceed int32 range — file sizes, timestamps, byte offsets into large
files, database row IDs. On 64-bit platforms this is a native register-width type, so
there is no performance penalty over int32.

---

## uint8
```arc
let x: uint8 = 0
let y: uint8 = 255
```

An unsigned 8-bit integer. Holds values from 0 to 255. Identical in layout to `byte` —
byte is just the preferred alias when the value represents raw data. Use uint8 when the
value is genuinely a small unsigned number rather than an opaque chunk of memory.
Common in pixel data, flag fields, and packed C structs.

---

## uint16
```arc
let x: uint16 = 0
let y: uint16 = 65535
```

An unsigned 16-bit integer. Holds values from 0 to 65,535. Used in network protocol
headers (port numbers are uint16), audio sample formats, and C struct fields typed as
unsigned short. Arithmetic wraps at 65,535 — the compiler will not warn you when it
rolls over.

---

## uint32
```arc
let x: uint32 = 0
let y: uint32 = 4294967295
```

An unsigned 32-bit integer. Holds values from 0 to roughly 4.2 billion. Common in
graphics (RGBA packed as uint32), checksums, hash values, flags fields, and C APIs
that return unsigned int. Matches C's unsigned int on all major platforms. Arithmetic
wraps at 4,294,967,295.

---

## uint64
```arc
let x: uint64 = 0
let y: uint64 = 18446744073709551615
```

An unsigned 64-bit integer. Holds values from 0 to roughly 18.4 quintillion. Used for
large counts, memory addresses treated as integers, 64-bit bitmasks, and file system
sizes. On 64-bit platforms this is a native register-width type. Arithmetic wraps at
2^64 - 1.

---

## usize
```arc
let len:    usize = 100
let offset: usize = buffer.length
let count:  usize = 4096
```

An unsigned pointer-sized integer. On a 64-bit platform it is 64 bits wide. On a 32-bit
platform it is 32 bits wide. This is the correct type for array lengths, memory sizes,
loop counters over collections, and anything that describes how many bytes or elements
exist. Using usize instead of uint64 keeps code portable across pointer widths and matches
what the compiler expects for indexing and sizeof results.

---

## isize
```arc
let delta:  isize = -4
let offset: isize = start - end
```

A signed pointer-sized integer. On a 64-bit platform it is 64 bits wide. On a 32-bit
platform it is 32 bits wide. The correct type for pointer offsets and differences between
addresses or indices that may be negative. Matches C's ptrdiff_t and ssize_t. Use isize
when you need to express a displacement or relative position rather than an absolute size.

---

## float32
```arc
let x: float32 = 3.14
let y: float32 = -0.001
let z: float32 = 1.0e10
```

A 32-bit IEEE 754 floating-point number. Provides roughly 7 decimal digits of precision.
This is the standard type in graphics, audio, physics engines, and GPU work where memory
bandwidth and SIMD throughput matter more than precision. float32 is what DirectX, OpenGL,
and most game math libraries expect. When precision requirements are relaxed, float32 is
significantly cheaper per element than float64.

---

## float64
```arc
let x: float64 = 3.141592653589793
let y: float64 = 2.718281828459045
```

A 64-bit IEEE 754 floating-point number. Provides roughly 15 decimal digits of precision.
The correct choice for scientific computation, financial calculations, geographic
coordinates, or any domain where accumulated rounding error in float32 would produce
wrong results. Most C standard library math functions (sin, cos, sqrt) operate on float64.
Doubles the memory footprint of float32 for the same element count.

---

## byte
```arc
let b:      byte = 255
let header: byte = 0xFF
```

An alias for uint8. Semantically signals that this value is raw, opaque data rather than
a small number — a chunk of memory, a wire byte, a pixel component, a character in a byte
stream. The compiler treats byte and uint8 as identical types; the distinction is purely
for human readers. Slices of bytes (`[]byte`) and arrays of bytes (`[N]byte`) are the
standard currency for I/O, network data, and C interop buffers.

---

## bool
```arc
let flag:    bool = true
let enabled: bool = false
```

An alias for uint8 with the values 1 (true) and 0 (false). Used for all conditional
logic. The compiler enforces that if statements and loop conditions receive a bool — there
is no implicit integer-to-bool coercion. When interfacing with C functions that return
int as a boolean, cast explicitly: `let ok: bool = bool(result)`. Stored as a single byte
in memory, same layout as C's bool in stdbool.h.

---

## char
```arc
let c:     char = 'a'
let digit: char = '5'
let euro:  char = '€'
```

An alias for uint32 representing a single Unicode code point. Not a byte, not a UTF-8
unit — a full Unicode scalar value. A char can represent any character in the Unicode
standard since uint32 covers the full code point range of 0 to 0x10FFFF. Character
literals use single quotes. Escape sequences (`'\n'`, `'\t'`, `'\\'`, `'\''`, `'\0'`)
produce the corresponding code point value as a uint32.

---

## string
```arc
let s:     string = "Hello, World!"
let empty: string = ""
let msg:   string = "Line one\nLine two"
```

A high-level string type stored internally as a pointer plus a length in bytes. The
content is UTF-8 encoded. Because the length is stored alongside the pointer, `len(s)`
is always a free O(1) operation — no memory scan required. Strings are immutable value
types: assigning a string to another variable copies the pointer and length, not the
underlying bytes. String literals support escape sequences (`\n`, `\t`, `\\`, `\"`,
`\0`). When passing string data to C functions that expect a null-terminated char pointer,
explicit conversion is required — arc strings are not implicitly null-terminated.

---

## Qualified Type
```arc
let client: net.Socket   = ...
let config: json.Config  = ...
let addr:   http.Address = ...
```

Any type can be referenced through its namespace using dot notation. This is not a
distinct type in the type system — it is the standard way to refer to a type that was
declared inside a namespace other than the current one. The full path uniquely identifies
the type across packages. There is no ambiguity resolution or implicit import: if you
reference `net.Socket`, the net package must be imported. Qualified types work everywhere
a plain type name works — variable declarations, function parameters, return types,
interface fields, and generic parameters.

---

## [N]T — Fixed-Size Array
```arc
let buf: [4096]byte  = ...
let ids: [256]uint32 = ...
let rgb: [3]float32  = ...
```

A fixed-size array whose length is known at compile time and is part of the type itself.
`[4096]byte` and `[1024]byte` are different, incompatible types. Fixed arrays can live
on the stack (declared with `let`) or on the heap with manual ownership (allocated with
`new`). The size N must be a compile-time constant. Indexing with a single integer
(`buf[2]`) produces a single element. Indexing with a range (`buf[0..4]`) produces a
`[]T` slice view into the array without any allocation. `sizeof(buf)` returns the total
byte size of the entire array at compile time.

---

## []T — Slice
```arc
let view:   []byte = buffer        // whole buffer as a slice view
let header: []byte = buffer[0..4]  // range into a buffer
let chunk:  []byte = existing_slice
let data:   []byte = some_vector   // whole vector as a slice view
```

A non-owning view into an existing sequence of elements, stored as a pointer and a
length. A slice never allocates memory and never owns what it points at — it is purely
a window into memory that already exists somewhere else. The source can be a fixed array,
a vector, another slice, or any contiguous region. Because the length is stored in the
slice itself, `len(view)` is always O(1). Slices are the standard way to pass a portion
of a buffer into a function without copying. The underlying memory must outlive the slice
— the compiler does not track this for you.

---

## vector[T] — Dynamic Array
```arc
let items: vector[int32]  = {1, 2, 3, 4, 5}    // annotation required, bare literal
let names: vector[string] = {"alice", "bob"}     // annotation required, bare literal

let nums  = vector[int32]{1, 2, 3}              // inferred from typed literal
let words = vector[string]{"hello", "world"}    // inferred from typed literal
```

An owned, heap-allocated, growable sequence of elements. Unlike a fixed array or slice,
a vector manages its own memory — it allocates, grows, and frees its backing buffer
automatically. Elements are accessed by index. A range index (`items[1..4]`) produces
a `[]T` slice view into the vector's backing memory without copying. `len(items)` is
always O(1). Vectors are the right choice when the number of elements is not known at
compile time and you need ownership over the collection.

When the element type is written directly in the literal, the variable annotation is not
needed. Typed vector literals can also be passed inline as function arguments without
declaring a variable first:
```arc
process(vector[int32]{1, 2, 3})
log_all(vector[string]{"error", "warn"})
```

---

## map[K]V — Map
```arc
let scores: map[string]int32 = {"alice": 100, "bob": 95}  // annotation required, bare literal

let hits  = map[string]int32{"alice": 100, "bob": 95}     // inferred from typed literal
let flags = map[string]bool{"debug": true, "trace": false} // inferred from typed literal
```

An owned, heap-allocated collection of key-value pairs. Keys must be a comparable type.
Values can be any arc type. Like typed vector literals, a typed map literal carries its
full type in the literal itself — no annotation is needed on the variable, and the literal
can be passed directly as a function argument without an intermediate variable:
```arc
configure(map[string]bool{"debug": true, "verbose": false})
seed(map[string]int32{"retries": 3, "timeout": 30})
```

Bare map literals with no type written in the literal still require an annotation on the
variable — the compiler has nothing to infer from:
```arc
let scores: map[string]int32 = {"alice": 100}   // annotation required
```

---

## Type Inference
```arc
let x     = 100                // int32
let name  = "hello"            // string
let flag  = true               // bool
let ratio = 3.14               // float64
let point = Point{x: 1, y: 2}  // Point
let nums  = vector[int32]{1, 2, 3}               // vector[int32]
let hits  = map[string]int32{"alice": 100}        // map[string]int32
```

The compiler infers the type of a declaration from its right-hand side. You never need
to write the type annotation when the value makes it unambiguous. Integer literals infer
as `int32`. Floating-point literals infer as `float64`. String literals infer as
`string`. Boolean literals infer as `bool`. Interface literals infer as the named
interface type. Typed collection literals carry their full type in the literal itself,
so inference works the same way.

When you need a type that differs from the default, annotate explicitly:
```arc
let x: int64   = 100   // would have been int32 without annotation
let y: float32 = 3.14  // would have been float64 without annotation
let z: uint8   = 255   // would have been int32 without annotation
```

Typed collection literals can also be passed directly as inline function arguments
without declaring a variable first — the type is unambiguous at the callsite:
```arc
render(vector[float32]{1.0, 0.5, 0.0})
connect(map[string]string{"host": "localhost", "port": "8080"})
```

**What cannot be inferred:**

Function parameters are never inferred — every parameter must have an explicit type,
enforced by the compiler with no exception:
```arc
func add(a: int32, b: int32) int32 { ... }   // always required
```

Nullable `var` declarations initialized to null cannot be inferred — null carries no
type information, so the annotation is required:
```arc
var server: net.Server = null   // annotation required
```

Bare collection literals with no type written in the literal cannot be inferred — the
compiler has nothing to derive the element or key-value types from:
```arc
let items: vector[int32]     = {1, 2, 3}         // annotation required
let scores: map[string]int32 = {"alice": 100}     // annotation required
```