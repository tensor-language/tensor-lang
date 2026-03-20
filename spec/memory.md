# Memory

---

## const
```arc
const MAX: int32 = 100
const NAME = "arc"
const PI: float64 = 3.141592653589793
```

A compile-time immutable value. `const` declares a value that is fixed at compile time
and never changes at runtime. The compiler is free to inline it everywhere it appears,
fold it into expressions, and eliminate it from the binary entirely as a named location.
There is no memory address you can take, no mutation that can happen, no scope that can
affect it. Type annotation is optional — the compiler infers it from the literal. Constants
are the right choice for magic numbers, sizes, protocol values, OS flags, and any value
that is a fact about the program rather than state within it.

---

## const, grouped
```arc
const (
    MAX:  int32   = 100
    MIN:  int32   = 0
    NAME           = "arc"
    PI:   float64 = 3.141592653589793
)

const (
    O_RDONLY: int32 = 0
    O_WRONLY: int32 = 1
    O_RDWR:   int32 = 2
    O_CREAT:  int32 = 64
)
```

Multiple constants declared together in a single block delimited by parentheses. Each
constant is on its own line with the same syntax as a single declaration — name, optional
type annotation, and value. Grouped constants are the standard form when declaring a
related set of values together: OS flags, protocol constants, configuration defaults,
error codes. Keeping related constants in one block makes their relationship visible and
keeps the file organized. The same inference rules apply — type annotation is optional
when the compiler can derive it from the literal.

---

## let
```arc
let x: int32 = 42
let point = Point{x: 1, y: 2}
let name = "hello"
```

A stack-allocated value. `let` declares a value that lives on the current stack frame.
When the scope that contains it ends, it is gone — no heap involvement, no ref counting,
no cleanup step. Assignment always copies the value: if you assign a `let` variable to
another variable, the receiving variable gets its own independent copy. This means `let`
is inherently safe to pass around within a scope but cannot be shared across scope
boundaries or stored in longer-lived structures. It is the right default for local
working values, intermediate results, and anything whose lifetime is obviously bounded
by the current function.

---

## var
```arc
var client = Client{name: "test", port: 8080}
var server: net.Server = null
```

A heap-allocated, reference-counted value. `var` declares a value that lives on the heap
and is managed by the compiler through reference counting. Every time a `var` is assigned
to another variable, passed into a function that holds onto it, or stored in a structure,
the reference count increments. When a reference goes out of scope the count decrements.
When the count reaches zero, `deinit` fires and the memory is freed. Because the value
lives on the heap, `var` can be safely shared across scopes, stored in other types, and
passed across function boundaries without copying the underlying data. `var` is also the
only declaration that can be assigned `null` — the compiler rejects null for `let` and
`const`.

---

## var, nullable
```arc
var client: Client = null
var connection: net.Conn = null

if client == null {
    // handle missing client
}
```

`var` declarations can be explicitly initialized to `null` and checked against `null` at
runtime. This is only valid for `var` — not `let`, not `const`, not `new`. Null represents
the absence of a heap-allocated ref-counted value. The compiler enforces this boundary:
assigning null to a stack value is a compile error because a stack slot always holds a
real value. Null is the standard way to express optional ownership in arc, and explicit
null checks are the correct way to guard before use. There is no implicit null coercion
and no automatic null propagation.

---

## new
```arc
let node = new Node{}
let buf  = new [4096]byte
let ids  = new [256]uint32
```

A heap-allocated, manually managed value. `new` allocates on the heap with no reference
counting and no compiler assistance with lifetime. You own it entirely — the compiler
will not free it, will not track references to it, and will not call `deinit` when it
goes out of scope. `new` is designed for kernel drivers, interrupt handlers, spinlock-held
sections, and hot paths where the overhead of reference counting is unacceptable. It is
also the right tool when you need precise, deterministic control over exactly when memory
is freed. The returned value is a typed pointer — field access and indexing work directly
on it without any extra syntax. Always pair `new` with `defer delete` immediately after
allocation.

---

## new, field initialization
```arc
let node = new Node{value: 10, next: null}
let desc = new D3D11_BUFFER_DESC{byte_width: 36, usage: 0, bind_flags: 1}
```

`new` accepts an initializer block with named fields, identical in syntax to stack
initialization. Fields not listed in the initializer are zero-initialized by the
compiler. This means a `new T{}` with an empty initializer block always produces
a fully zeroed allocation — there is no uninitialized memory to reason about. Named
field initialization makes the intent explicit at the callsite and prevents mistakes
that come from positional struct initialization where field order matters and is easy
to get wrong.

---

## new, fixed array
```arc
let buf = new [4096]byte
let ids = new [256]uint32
let mat = new [16]float32
```

`new` can allocate a fixed-size array on the heap with manual ownership. The size is
a compile-time constant, the type is any arc type, and the result is a typed pointer to
the first element. All elements are zero-initialized on allocation. This is the pattern
for I/O buffers, DMA regions, scratch memory, and any fixed-size working area that is
too large for the stack or needs a lifetime that outlives the current scope. Like all
`new` allocations, the compiler will not free it — `defer delete` is mandatory.

---

## delete
```arc
let node = new Node{value: 42}
defer delete(node)

let buf = new [4096]byte
defer delete(buf)
```

The manual deallocation call that pairs with `new`. `delete` frees the heap memory that
`new` allocated. It does not call `deinit` — that is only for `var`. It does not
decrement a reference count — there is none. It simply releases the allocation back to
the allocator. `delete` must always be called manually. The standard pattern is to write
`defer delete(x)` on the line immediately after the `new` — this ensures the free
happens at scope exit regardless of how the function returns, including early returns
and error paths. Forgetting `delete` is a memory leak. Calling it twice is undefined
behavior.

---

## defer
```arc
let buf = new [4096]byte
defer delete(buf)

let fd = open("file.txt", O_RDONLY)
defer close(fd)

let db = sqlite3_open("test.db")
defer sqlite3_close(db)
```

Schedules a call to execute at the end of the current scope, regardless of how that
scope exits — normal return, early return, or any other exit path. `defer` is not
specific to memory management but it is the standard mechanism for pairing cleanup with
allocation. Deferred calls execute in reverse order of declaration, so the last resource
acquired is the first one released. Writing `defer delete(x)` immediately after `new`
is the idiomatic arc pattern — it keeps acquisition and release visually adjacent in
the source and eliminates an entire class of leaks caused by early returns or added
code paths that forget to clean up.

---

## deinit
```arc
interface Client {
    name: string
    port: int32
}

deinit(self c: Client) {
    io.close(c.port)
}

var c = Client{name: "test", port: 8080}
// ... when ref count hits 0, deinit fires automatically
```

A lifecycle hook that fires automatically when a `var` declaration's reference count
reaches zero. `deinit` is the place to release resources that the type owns — file
descriptors, sockets, GPU handles, anything that needs explicit cleanup when the value
is no longer referenced. You never call `deinit` manually — the compiler inserts the
call at the point where the last reference drops. `deinit` only fires for `var`
declarations. It never fires for `let`, `const`, or `new`. If a type has resources that
need cleanup and it will be used as a `var`, it should define a `deinit`. If it will
only ever be used as a `new`, cleanup is your responsibility through `delete` and manual
resource management.

---

## var vs new
```arc
var node = Node{}       // heap, ref counted — compiler manages lifetime
let node = new Node{}   // heap, manual — you manage lifetime
```

Both `var` and `new` allocate on the heap, but they represent fundamentally different
ownership contracts. `var` hands lifetime responsibility to the compiler — the ref count
tracks who holds references, `deinit` fires when the last one drops, and you cannot
forget to free it. `new` hands lifetime responsibility entirely to you — there is no
ref count, no `deinit`, no safety net. The tradeoff is overhead and control. `var` is
safe and automatic but carries ref count increment and decrement at every assignment and
scope exit. `new` is unsafe and manual but has zero overhead beyond the allocation itself,
making it the right choice for hot paths, kernel code, interrupt handlers, and anywhere
the cost of ref counting is measurable and unacceptable.

---

## Stack vs Heap at a Glance
```arc
const MAX = 100                            // compile-time, no memory
let point = Point{x: 1, y: 2}             // stack, freed at scope end
var client = Client{name: "x", port: 80}  // heap, ref counted, deinit on zero
var client: Client = null                  // heap, nullable
let node = new Node{value: 42}             // heap, manual, you free it
let buf  = new [4096]byte                  // heap, manual fixed array
defer delete(node)                         // always paired with new
defer delete(buf)
```

The four declaration keywords encode the full ownership story at the point of
declaration. `const` is a compile-time fact. `let` is a stack value with scope lifetime.
`var` is a shared heap value with compiler-managed lifetime. `new` is a raw heap
allocation with manually managed lifetime. Reading a declaration tells you immediately
who is responsible for the memory, how long it lives, whether it can be null, and
whether cleanup is automatic or manual. The type annotation never carries this
information — the keyword always does.