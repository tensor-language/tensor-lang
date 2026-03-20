# Intrinsics

---

## What Intrinsics Are

Intrinsics are compiler-handled functions. They are not keywords, not parser grammar,
and not part of the standard library. They are mapped in the compiler's intrinsic
registry and resolved before code generation. Calling an intrinsic looks identical to
calling a regular function at the source level — the difference is that the compiler
replaces the call with a direct instruction sequence, a compile-time computation, or
a platform-specific operation with no function call overhead. Some intrinsics have zero
runtime cost because they are resolved entirely at compile time. Others compile to a
single instruction. None of them go through a normal call stack.

---

## memptr — Address Of
```arc
let x = 100
let p = memptr(&x)
```

Gets the address of an arc variable as a raw memory pointer. `memptr(&val)` is the
form used when you need to pass the address of an arc variable to a C function that
expects a pointer — output parameters, struct pointers, buffer addresses. The `&`
inside `memptr` is not the mutable reference operator used in regular arc code — it is
specifically the address-of operation at the C boundary. `memptr` is only valid when
calling extern functions or performing explicit low-level operations. In regular arc
code there are no addresses, no pointer arithmetic, and no `&val` syntax — `memptr`
is the explicit marker that you are crossing that boundary.

---

## memptr — Integer to Pointer
```arc
let transient = memptr(-1)
sqlite3_bind_text(stmt, 1, val, len, transient)
```

Casts an integer value directly to a memory pointer. Used when a C API expects a
specific sentinel pointer value that is not a real address — SQLite's `SQLITE_TRANSIENT`
is `-1` cast to a function pointer, Win32 APIs use `(HANDLE)-1` for special handle
values. The integer is reinterpreted as a pointer-sized value with no conversion. This
form does not take the address of anything — it manufactures a pointer from a raw
integer. Valid only at the C boundary.

---

## memptr — Pointer Arithmetic
```arc
let addr     = usize(some_ptr)
let next_ptr = memptr(addr + 16)
```

Pointer arithmetic is performed by casting a pointer to `usize`, doing integer math,
then casting back with `memptr`. There is no direct pointer addition in arc — the
round-trip through `usize` is the required pattern and is intentional. It keeps
pointer arithmetic explicit and visible rather than allowing it to happen silently on
typed pointers. The `usize` type is the correct carrier for pointer-sized arithmetic
— it is guaranteed to be wide enough to hold a pointer value on the current platform.

---

## memptr — Output Parameters
```arc
var db: sqlite3 = null
sqlite3_open("test.db", memptr(&db))

var device: ID3D11Device = null
D3D11CreateDevice(..., memptr(&device), ...)
```

The most common use of `memptr` in practice. C functions that return values through
pointer parameters — output parameters typed as `T**` in C — receive the address of
an arc `var` variable via `memptr(&variable)`. The C function writes through the
pointer and the arc variable is updated. This pattern appears everywhere in C APIs that
allocate handles, return objects, or fill in multiple output values. The `var`
declaration is initialized to null before the call and checked after.

---

## new — Single Allocation
```arc
let node = new Node{}
let node = new Node{value: 10, next: null}
```

Allocates a single instance of a type on the heap with no reference counting and no
compiler-managed lifetime. Fields not listed in the initializer are zero-initialized.
`new` returns a typed pointer — field access works directly on it without any extra
syntax. The compiler does not track the allocation, will not free it, and will not call
`deinit` when it goes out of scope. You own it entirely. Designed for kernel drivers,
interrupt handlers, spinlock-held sections, and hot paths where the overhead of
reference counting is measurable and unacceptable. Always pair with `defer delete`
immediately after allocation.

---

## new — Array Allocation
```arc
let buf  = new [4096]byte
let ids  = new [256]uint32
let mat  = new [16]float32
```

Allocates a fixed-size array on the heap with manual ownership. The size is a
compile-time constant. All elements are zero-initialized. The result is a typed pointer
to the first element — indexing works directly. This is the standard pattern for I/O
buffers, DMA regions, scratch memory, and any fixed-size working area that is too large
for the stack or needs a lifetime that outlives the current scope. Same ownership rules
as single allocation — the compiler will not free it, `defer delete` is mandatory.

---

## delete
```arc
let node = new Node{value: 42}
defer delete(node)

let buf = new [4096]byte
defer delete(buf)
```

Frees a heap allocation made with `new`. Does not call `deinit` — that is only for
`var`. Does not decrement a reference count — there is none. Releases the allocation
back to the allocator immediately. Must always be called manually — the standard pattern
is `defer delete(x)` written on the line immediately after `new`. Deferred execution
ensures the free happens at scope exit regardless of how many return paths the function
has. Forgetting `delete` is a memory leak. Calling it twice is undefined behavior.

---

## mem_zero
```arc
mem_zero(buf, sizeof(buf))
mem_zero(addr, sizeof(SockAddrIn))

let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)
```

Fills a region of memory with zero bytes. Takes a pointer and a byte count. Works on
any memory — stack, heap, mapped regions. No ownership implied, no lifetime tracked.
The standard pattern after `new` for buffers and structs that will be passed to C
functions: allocate, zero, use, defer delete. Also the correct way to zero a struct
before a syscall to avoid leaking stack contents to the kernel. Compiles to optimized
CPU instructions — typically a `rep stosd` or equivalent SIMD sequence on x86.

---

## mem_copy
```arc
mem_copy(dest, src, 1024)

let src  = new Point{x: 10, y: 20}
let dest = new Point{}
mem_copy(dest, src, sizeof(Point))
defer delete(src)
defer delete(dest)
```

Copies a region of memory from source to destination. Takes a destination pointer,
a source pointer, and a byte count. The source and destination regions must not overlap
— behavior is undefined if they do. Faster than `mem_move` for non-overlapping regions
because the compiler can use unconstrained copy instructions. Use `mem_move` when there
is any possibility of overlap. Compiles to optimized CPU copy instructions — typically
`rep movsd` or a SIMD sequence. No ownership transfer implied — both pointers remain
valid after the copy.

---

## mem_move
```arc
mem_move(dest, src, 1024)

// shift buffer contents left by 100 bytes
mem_move(buffer, memptr(usize(buffer) + 100), 900)
```

Copies a region of memory from source to destination with correct behavior when the
regions overlap. Slower than `mem_copy` because it must handle the overlap case —
internally it may copy backwards depending on the relative positions of source and
destination. Use this when shifting data within a buffer, when the relationship between
source and destination is not guaranteed to be non-overlapping, or when you are not
certain. If you know regions do not overlap, `mem_copy` is faster.

---

## mem_compare
```arc
let diff = mem_compare(ptr1, ptr2, 1024)

if mem_compare(a, b, sizeof(Point)) == 0 {
    // structs are byte-equal
}
```

Compares two regions of memory byte by byte. Returns 0 if the regions are identical,
a negative value if the first differing byte in the first region is less than the
corresponding byte in the second region, and a positive value if it is greater. The
return value is `int32`. This is a raw byte comparison — it does not understand type
semantics, padding bytes, or floating point equality. Two structs that are logically
equal but have different values in their padding bytes will not compare as equal.
Use only when byte-level equality is what you actually need.

---

## sizeof
```arc
let sz     = sizeof(int32)         // 4
let st_sz  = sizeof(SockAddrIn)    // struct size with padding
let arr_sz = sizeof(buf)           // total size of fixed array

let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)
```

Returns the size of a type or variable in bytes as a `usize`. Resolved entirely at
compile time — zero runtime cost. For primitive types it returns the fixed byte width.
For interfaces it returns the total size including any padding inserted for alignment.
For fixed arrays it returns the total byte size of all elements. For slices and vectors
it returns the size of the header struct, not the data. `sizeof` is the correct way to
pass sizes to `mem_zero`, `mem_copy`, `mem_move`, `mem_compare`, and C functions that
expect a byte count. Never hardcode sizes — always use `sizeof`.

---

## alignof
```arc
let align        = alignof(float64)    // 8
let struct_align = alignof(Point)      // largest field alignment
let vec_align    = alignof(XMVECTOR)   // 16 if @align(16) declared
```

Returns the alignment requirement of a type in bytes as a `usize`. Resolved entirely
at compile time — zero runtime cost. For primitives it returns the natural alignment:
1 for byte, 2 for int16, 4 for int32, 8 for int64 and float64. For interfaces it
returns the alignment of the most strictly aligned field. For interfaces decorated with
`@align(N)`, it returns N. Used when manually managing aligned allocations, passing
alignment requirements to custom allocators, or verifying that a layout meets the
requirements of a C or C++ API.

---

## bitcast
```arc
// IEEE 754 float to integer bits
let f    = float32(1.0)
let bits = bitcast(uint32, f)       // 0x3F800000

// type punning
let i = int32(-1)
let u = bitcast(uint32, i)          // 0xFFFFFFFF

// back to float
let raw  = uint32(0x3F800000)
let back = bitcast(float32, raw)    // 1.0
```

Reinterprets the bits of a value as a different type without any conversion. Zero
runtime cost — the compiler reinterprets the existing bit pattern at the new type with
no instruction emitted. Both types must be the same size. `bitcast` is the correct tool
for IEEE 754 bit manipulation, type punning between integer and float representations,
and any case where you need to inspect or construct the raw bit pattern of a value.
It is not a numeric conversion — `bitcast(uint32, float32(1.0))` gives you the IEEE
754 encoding `0x3F800000`, not the integer `1`.

---

## len — Arc Types
```arc
let s     = "hello"
let n     = len(s)          // 5, O(1), reads stored length field

let data: []byte = buf[0..64]
let n     = len(data)        // 64, O(1), reads stored length field

let items: vector[int32] = {1, 2, 3}
let n     = len(items)       // 3, O(1), reads stored length field
```

Returns the number of elements or bytes in an arc string, slice, or vector. For all
arc-native types, `len` is free — it reads a length field stored alongside the data
pointer and the compiler may inline it to a single field load. There is no scan, no
iteration, no traversal. The length is always known because arc strings, slices, and
vectors all carry their length with them. Use `len` freely in loop conditions and size
calculations — it has no performance cost on arc types.

---

## len — C String
```arc
let cstr: *byte = some_c_func()
let n = len(cstr)               // walks memory until \0, O(n)
```

When `len` is called on a `*byte` — a C null-terminated string — it performs a memory
scan from the pointer until it finds a zero byte, equivalent to `strlen`. The cost is
proportional to the string length. This is not free. The type is the warning: `*byte`
means a C string, and `len` on a C string costs a scan. If you find yourself calling
`len` on a `*byte` in a loop or hot path, convert to an arc `string` once and use the
free version. The distinction between `string` and `*byte` in arc is precisely this
difference in behavior.

---

## syscall
```arc
let msg    = "Hello\n"
let result = syscall(SYS_WRITE, STDOUT, msg, 6)

let fd     = syscall(SYS_OPEN, path, O_RDONLY, 0)

let addr   = syscall(SYS_MMAP, null, 4096,
                     PROT_READ | PROT_WRITE,
                     MAP_PRIVATE | MAP_ANONYMOUS, -1, 0)

let buf        = new [4096]byte
let bytes_read = syscall(SYS_READ, fd, buf, 4096)
defer delete(buf)
```

Issues a direct system call to the kernel. Up to six arguments after the syscall
number. The syscall number and argument conventions are platform-specific — on Linux
x86-64 they match the syscall ABI directly. `syscall` bypasses libc entirely and goes
straight to the kernel interface. This is the correct tool for kernel drivers, sandboxed
environments where libc is unavailable, performance-critical paths where libc wrapper
overhead matters, and any situation where you need precise control over exactly what
the kernel receives. The return value is the raw kernel return — negative values
typically indicate errors encoded as `-errno`.

---

## va_start / va_arg / va_end
```arc
func log(fmt: string, ...) {
    let args = va_start(fmt)
    defer va_end(args)

    let code = va_arg(args, int32)
    let msg  = va_arg(args, string)
    printf("%s [%d]: %s\n", fmt, code, msg)
}
```

Intrinsics for implementing variadic functions that accept a variable number of
arguments. `va_start` initializes the argument list from the last named parameter and
returns a handle to the variadic argument state. `va_arg` advances through the argument
list one step at a time, returning each argument cast to the requested type — call once
per argument in the order they were passed. `va_end` cleans up the argument state.
Always pair `va_start` with `defer va_end` — never skip cleanup. `va_arg` calls must
match the types the caller actually passed — the compiler cannot verify this at the
call boundary, incorrect types produce undefined behavior.

---

## gpu func / thread_id
```arc
gpu func kernel(data: float32, n: usize) {
    let idx   = thread_id()
    data[idx] = data[idx] * 2.0
}

async func main() {
    let result = await kernel(data, n)
}
```

`gpu func` marks a function for compilation to an accelerator backend — CUDA, Metal,
XLA, or others set in the build configuration. The compiler routes `gpu func`
declarations to the appropriate device compiler rather than the host compiler. All
parameters are implicitly device-bound — data transfer between host and device memory
is handled automatically. Inside a `gpu func`, `thread_id()` is an intrinsic that
returns the index of the current parallel execution lane as a `usize`. Each lane runs
the same function body over its own index simultaneously. `thread_id()` is only valid
inside a `gpu func` — using it in host code is a compile error. GPU functions are
called with `await` from async host code and cannot call regular arc functions or use
`extern` blocks.

---

## Quick Reference
```arc
// memptr — the C boundary
memptr(&val)                    // address of val as memory pointer
memptr(-1)                      // integer as memory pointer
memptr(usize(p) + 16)           // pointer arithmetic via usize

// manual heap — always pair new with defer delete
let node = new Node{}           // single allocation, zero initialized
let node = new Node{value: 42}  // single allocation, field initialized
let buf  = new [4096]byte       // array allocation, zero initialized
defer delete(node)              // free single
defer delete(buf)               // free array

// raw memory — no ownership, no lifetime
mem_zero(ptr, sizeof(ptr))      // zero fill
mem_copy(dst, src, n)           // copy, regions must not overlap
mem_move(dst, src, n)           // copy, overlap safe, slightly slower
mem_compare(a, b, n)            // compare, 0 = equal

// compile-time — zero runtime cost
sizeof(T)                       // size of type in bytes
sizeof(buf)                     // size of fixed array in bytes
alignof(T)                      // alignment requirement in bytes
bitcast(T, val)                 // reinterpret bits, no conversion, no cost

// length — type decides the cost
len(s)                          // arc string / slice / vector — free, O(1)
len(cstr)                       // *byte C string — memory scan, O(n)

// syscall — direct kernel interface
syscall(SYS_WRITE, fd, buf, n)  // up to 6 args, raw kernel return

// variadic — always pair va_start with defer va_end
let args = va_start(fmt)
defer va_end(args)
let val  = va_arg(args, int32)

// gpu
thread_id()                     // current lane index, only inside gpu func
```