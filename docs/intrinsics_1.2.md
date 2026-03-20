# arc Compiler Intrinsics (Version 2.3)

> **Note**: All intrinsics are compiler-handled functions, not parser keywords.
> They are mapped in the compiler's intrinsic registry, not built into the grammar.

---

## memptr — Memory Pointer

Two forms. One name. The C boundary lives here.
```arc
// memptr(&val) — get address of variable as memory pointer
// memptr(x)    — cast integer value to memory pointer

// 1. Get address of a variable
let x = 100
let p = memptr(&x)

// 2. Cast integer value to pointer
let transient = memptr(-1)
sqlite3_bind_text(stmt, 1, val, len, transient)

// 3. Pointer arithmetic (via usize)
let addr     = usize(some_ptr)
let next_ptr = memptr(addr + 16)

// 4. Pass address to extern output params
var db: sqlite3 = null
sqlite3_open("test.db", memptr(&db))
```

- Only valid inside `extern c` blocks or when calling extern functions
- No star syntax outside extern blocks — `memptr` is the C boundary
- Pointer arithmetic always goes through `usize` first

---

## new / delete — Manual Heap Allocation

`new` and `delete` are for manual memory management.
No ref counting, no deinit, no compiler assistance — you own it, you free it.
Designed for kernel drivers, hot paths, and anywhere `var` ref counting is too heavy.
```arc
// allocate single type, zero initialized
let node = new Node{}

// allocate with field initialization
let node = new Node{value: 10, next: null}

// allocate fixed-size array
let buf  = new [4096]byte
let ids  = new [256]uint32

// always pair with defer delete
let node = new Node{}
defer delete(node)

let buf = new [4096]byte
defer delete(buf)
```

The difference between `var` and `new`:
```arc
var node = Node{}       // heap, ref counted — compiler manages lifetime
let node = new Node{}   // heap, manual — you manage lifetime
```

- `new` returns a typed pointer — field access just works
- `delete` must be called manually — use `defer` to avoid leaks
- `new` never triggers `deinit` — that is only for `var`
- Safe to use in interrupt handlers and spinlock-held sections
- No ref count overhead — safe in hot paths

---

## mem — Raw Memory Operations

Operations on raw memory. Work on any memory — stack, heap, mapped regions.
No ownership implied, no lifetime tracked.

### mem_zero — Fill with zeros
```arc
// mem_zero(ptr, count)
mem_zero(buf, sizeof(buf))
mem_zero(addr, sizeof(SockAddrIn))

// zero init after alloc
let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)

// zero init struct before syscall
let stat_buf = new Stat{}
mem_zero(stat_buf, sizeof(Stat))
defer delete(stat_buf)
```

### mem_copy — Copy memory (non-overlapping)
```arc
// mem_copy(dest, src, count)
// UNSAFE if regions overlap — use mem_move for overlapping regions
mem_copy(dest, src, 1024)

// copy a struct
let src  = new Point{x: 10, y: 20}
let dest = new Point{}
mem_copy(dest, src, sizeof(Point))
defer delete(src)
defer delete(dest)

// copy fixed array
let src  = new [1024]byte
let dest = new [1024]byte
mem_copy(dest, src, sizeof(src))
defer delete(src)
defer delete(dest)
```

### mem_move — Copy memory (overlap safe)
```arc
// mem_move(dest, src, count)
// safe for overlapping regions, slightly slower than mem_copy
mem_move(dest, src, 1024)

// shift buffer contents
mem_move(buffer, memptr(usize(buffer) + 100), 900)
```

### mem_compare — Compare memory regions
```arc
// mem_compare(ptr1, ptr2, count) -> int32
// returns 0 if equal, <0 if ptr1 < ptr2, >0 if ptr1 > ptr2
let diff = mem_compare(ptr1, ptr2, 1024)

if mem_compare(a, b, sizeof(Point)) == 0 {
    // structs are equal
}
```

---

## Compile-Time Builtins

Zero runtime cost. Resolved entirely at compile time.
No namespace — they are not operations, they are facts about types.

### sizeof — Size of type in bytes
```arc
// sizeof(T) -> usize
let sz     = sizeof(int32)         // 4
let st_sz  = sizeof(SockAddrIn)    // struct size with padding
let arr_sz = sizeof(buf)           // size of fixed array

// common pattern
let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)
```

### alignof — Alignment requirement of type
```arc
// alignof(T) -> usize
let align        = alignof(float64)   // 8
let struct_align = alignof(Point)
```

### bitcast — Reinterpret bits without conversion
```arc
// bitcast(T, value) -> T
// zero runtime cost — compiler reinterprets the bits, no conversion

// IEEE754 float to integer bits
let f    = float32(1.0)
let bits = bitcast(uint32, f)     // 0x3F800000

// type punning
let i = int32(-1)
let u = bitcast(uint32, i)        // 0xFFFFFFFF

// integer back to float bits
let raw  = uint32(0x3F800000)
let back = bitcast(float32, raw)  // 1.0
```

---

## len — String and Slice Length

One name. The type decides the cost.
```arc
// arc string — free, compiles to s.len field access
let s = "hello"
let n = len(s)        // 5, O(1), no scan

// slice — free, compiles to slice.len field access
let data: []byte = buffer[0..64]
let n = len(data)     // 64, O(1), no scan

// vector — free, compiles to vec.len field access
let items: vector[int32] = {1, 2, 3}
let n = len(items)    // 3, O(1), no scan

// C null-terminated string — memory scan, costs O(n)
let cstr: *byte = some_c_func()
let n = len(cstr)     // walks memory until \0
```

- `len` on any arc type is always free — it reads a stored length field
- `len` on `*byte` compiles to a `strlen` memory scan — cost is proportional to string length
- The type is the warning — if you pass a `*byte` you are doing a scan

---

## syscall — Direct System Call

Direct kernel interface. Up to 6 arguments.
Universal vocabulary — every systems developer already knows it.
```arc
// syscall(number, arg1..arg6)
let msg  = "Hello\n"
let wlen = 6

// SYS_WRITE
let result = syscall(SYS_WRITE, STDOUT, msg, wlen)

// SYS_OPEN
let fd = syscall(SYS_OPEN, path, O_RDONLY, 0)

// SYS_MMAP
let addr = syscall(SYS_MMAP, null, 4096,
                   PROT_READ | PROT_WRITE,
                   MAP_PRIVATE | MAP_ANONYMOUS, -1, 0)

// SYS_READ into manual buffer
let buf        = new [4096]byte
let bytes_read = syscall(SYS_READ, fd, buf, 4096)
defer delete(buf)
```

---

## Variadic Arguments

For implementing variadic functions that mirror C's `va_list` pattern.
```arc
func printf(fmt: string, ...) {
    let args = va_start(fmt)
    defer va_end(args)

    let val = va_arg(args, int32)
    let str = va_arg(args, string)
}
```

- Always pair `va_start` with `defer va_end` — never skip cleanup
- `va_arg` advances the argument list — call once per argument in order

---

## gpu func — Accelerator Kernels
```arc
// gpu func — all params are gpu bound
// compiler maps to build target (cuda, metal, xla)
gpu func kernel(data: float32, n: usize) {
    let idx   = thread_id()
    data[idx] = data[idx] * 2.0
}

async func main() {
    let result = await kernel(data, n)
}
```

---

## Quick Reference
```arc
// addressing
memptr(&val)              // address of val as memory pointer
memptr(-1)                // integer as memory pointer
memptr(usize(p) + 16)     // pointer arithmetic

// manual heap
let node = new Node{}           // allocate single, you own it
let node = new Node{value: 42}  // allocate with fields
let buf  = new [4096]byte       // allocate fixed array
let ids  = new [256]uint32      // allocate fixed array of any type
defer delete(node)               // free it, always use defer
defer delete(buf)                // free array, same pattern

// raw memory ops
mem_zero(ptr, sizeof(ptr))       // zero fill
mem_copy(dst, src, n)            // copy, no overlap
mem_move(dst, src, n)            // copy, overlap safe
mem_compare(a, b, n)             // compare, 0 = equal

// compile-time
sizeof(T)                        // size in bytes
sizeof(buf)                      // size of fixed array
alignof(T)                       // alignment requirement
bitcast(T, val)                  // reinterpret bits, zero cost

// length
len(s)                           // arc string/slice/vector — free
len(cstr)                        // *byte C string — O(n) scan

// syscall
syscall(SYS_WRITE, fd, buf, n)  // direct kernel interface
```

---

## Usage Notes

**Memory Safety**
- `mem_copy` is faster but UNSAFE for overlapping regions
- `mem_move` handles overlaps safely but is slightly slower
- Always `defer delete` immediately after `new` — never skip it
- `new` never calls `deinit` — that is only for `var`

**Performance**
- `sizeof`, `alignof`, `bitcast` are compile-time — zero runtime cost
- `new` / `delete` have no ref count overhead — safe in hot paths
- `mem_zero`, `mem_copy`, `mem_move` compile to optimized CPU instructions
- `len` on arc types is always free — only `*byte` costs a scan

**Kernel and Driver Patterns**
```arc
// MMIO register access
let regs = memptr(0x3F200000)

// DMA buffer
let dma = new DmaBuffer{}
mem_zero(dma, sizeof(DmaBuffer))
defer delete(dma)
syscall(SYS_IOCTL, fd, DMA_MAP, memptr(&dma))

// syscall buffer pattern
let buf        = new [4096]byte
let bytes_read = syscall(SYS_READ, fd, buf, sizeof(buf))
defer delete(buf)

// C struct interop
let addr = new SockAddrIn{}
mem_zero(addr, sizeof(SockAddrIn))
defer delete(addr)
connect(fd, memptr(&addr), uint32(sizeof(SockAddrIn)))
```