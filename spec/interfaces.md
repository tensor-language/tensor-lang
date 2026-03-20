# Interfaces

---

## Basic Interface
```arc
interface Point {
    x: float32
    y: float32
}

interface Client {
    name: string
    port: int32
}
```

The single type declaration in arc. There is no `struct`, no `class`, no `record` —
`interface` is the one keyword for describing a named layout of fields. An interface
declares the shape of a type: what fields it has and what types those fields are. It
carries no memory model decision — whether an interface lives on the stack, the heap
with ref counting, or the heap with manual ownership is decided at the callsite with
`let`, `var`, or `new`. The interface declaration itself is purely a description of
layout.

---

## Field Types
```arc
interface Server {
    host:       string
    port:       int32
    max_conns:  uint32
    running:    bool
    timeout:    float64
}
```

Fields are declared one per line as `name: type`. Any arc type is valid as a field type
— primitives, strings, other interfaces, slices, vectors, fixed arrays, function types,
and qualified types from other namespaces. Field order in the declaration determines
field order in memory. There is no access control — all fields are accessible wherever
the type is in scope. The compiler lays out fields in declaration order with natural
alignment unless the interface carries a `@packed` attribute.

---

## Nested Interface
```arc
interface Address {
    host: string
    port: int32
}

interface Connection {
    addr:    Address
    timeout: float32
    active:  bool
}

let conn = Connection{
    addr:    Address{host: "localhost", port: 8080},
    timeout: 30.0,
    active:  true
}
```

An interface field can itself be of an interface type. The nested interface is stored
inline in the parent — not as a pointer, not as a heap allocation. The memory layout
is flat and contiguous. Nested interface fields are initialized with a nested initializer
block at the callsite. Accessing a nested field uses chained dot notation:
`conn.addr.host`. The decision of how to allocate the outer interface still applies to
the whole thing — a `let Connection` puts the entire layout including the nested Address
on the stack.

---

## Interface Instantiation
```arc
// stack — copied on assignment
let p = Point{x: 1.0, y: 2.0}

// heap, ref counted — shared, nullable
var c = Client{name: "test", port: 8080}

// heap, manual — you own it
let node = new Node{value: 42}
```

The same interface type can be instantiated three different ways depending on what the
callsite needs. `let` puts it on the stack and copies it on assignment. `var` puts it
on the heap with ref counting, enables sharing across scopes, and allows null. `new`
puts it on the heap with no ref counting and no compiler assistance — you manage the
lifetime manually with `delete`. The interface declaration does not change. The keyword
at the callsite is the entire ownership decision. This means you can use the same type
in a hot path with `new` and in regular code with `var` without defining separate types.

---

## Field Initialization
```arc
interface Config {
    host:    string
    port:    int32
    retries: uint32
    debug:   bool
}

let cfg = Config{host: "localhost", port: 8080, retries: 3, debug: false}
```

Interface instances are initialized with a braced block of named field assignments.
Fields are assigned by name, not by position — order does not matter and any field
omitted from the initializer is zero-initialized by the compiler. Zero initialization
means numeric fields become 0, bool fields become false, string fields become empty,
and var fields become null. This guarantees there is no uninitialized memory in an
arc interface instance. Named initialization also means adding a new field to an
interface does not silently break existing initialization sites — unspecified fields
just get zero.

---

## Methods
```arc
interface Rectangle {
    width:  float32
    height: float32
}

func area(self r: Rectangle) float32 {
    return r.width * r.height
}

func scale(self r: Rectangle, factor: float32) Rectangle {
    return Rectangle{width: r.width * factor, height: r.height * factor}
}

let rect = Rectangle{width: 10.0, height: 5.0}
let a = rect.area()
let scaled = rect.scale(2.0)
```

Methods are regular top-level functions whose first parameter is `self`. The `self`
keyword binds the function to the interface type of that parameter and enables dot-
notation calls on instances. Methods are not declared inside the interface block — the
interface only declares data. This separation means the interface declaration stays
focused on layout and methods can be added anywhere in the package without touching
the type. A method that does not need to modify the receiver takes it by value via
`self name: Type`. The receiver is a copy and mutations inside the method do not affect
the caller's instance.

---

## Mutable Methods
```arc
interface Counter {
    value: int32
}

func increment(self &mut c: Counter) {
    c.value += 1
}

func reset(self &mut c: Counter) {
    c.value = 0
}

var c = Counter{value: 0}
c.increment()
c.increment()
c.reset()
```

When a method needs to modify the receiver, the self parameter is declared as `&mut`.
This follows the same mutable reference rules as any other parameter — the mutation is
explicit in both the declaration and at the callsite. For `var` receivers, dot-notation
handles the implicit `&` because the compiler knows the variable is heap-allocated. For
`let` receivers the `&` must be passed explicitly. A method that reads should take by
value. A method that writes should take by `&mut`. The distinction is always visible in
the signature.

---

## Function Field
```arc
interface Button {
    label:   string
    onClick: func(int32, int32) void
}

interface Stream {
    port:     int32
    onData:   async func([]byte) void
    onClose:  func() void
}
```

A field can store a function — making behavior a first-class part of the type alongside
data. The field type declares the full signature: parameter types, return type, and
whether the callable is async. Function fields default to null when uninitialized. They
are assigned lambdas or named functions at the callsite. This is the arc pattern for
event handlers, callbacks, and injectable behavior — it avoids virtual dispatch and
keeps the wiring explicit at the point where the interface is configured rather than
buried in an inheritance hierarchy.

---

## Generic Interface
```arc
interface Box[T] {
    value: T
}

interface Pair[K, V] {
    key:   K
    value: V
}

interface Result[T, E] {
    data:    T
    error:   E
    success: bool
}
```

An interface parameterized over one or more types. The type parameters are listed in
square brackets after the interface name and can appear anywhere a type is valid in the
field list. The compiler generates a concrete layout for each distinct combination of
type arguments the interface is instantiated with. Generic interfaces carry no runtime
overhead — the generics are resolved entirely at compile time. Type parameters are
positional and named — `Pair[K, V]` and `Pair[V, K]` are the same declaration with
different parameter names but different meaning in the field layout.

---

## Generic Interface Methods
```arc
interface Box[T] {
    value: T
}

func get[T](self b: Box[T]) T {
    return b.value
}

func set[T](self &mut b: Box[T], val: T) {
    b.value = val
}

let b = Box[int32]{value: 42}
let v = b.get()
b.set(100)
```

Methods on generic interfaces use the same self-receiver pattern as regular methods,
combined with the same `[T]` type parameter syntax used by generic functions. The type
parameter is declared on the method itself — `func get[T]` — because methods are
top-level functions, not declared inside the interface. The interface only declares
data. The type parameter in the method signature must match the parameter used in the
receiver type — `Box[T]` in the self parameter connects the method's `T` to the
interface's `T`. The compiler infers the concrete type from the receiver at the
callsite — `b.get()` on a `Box[int32]` resolves `T` to `int32` without explicit
annotation. Generic methods can use the type parameter anywhere in the signature: as
parameter types, return types, or both.

---

## Opaque Type

```arc
type FILE         = opaque
type sqlite3      = opaque
type sqlite3_stmt = opaque
```

An opaque type declaration gives a C type a name in arc without describing its layout.
When a C library works with types whose internal layout is private — file handles,
database handles, platform-specific objects — you declare them as `opaque` to satisfy
the type system at the C boundary without pretending to know the fields. An opaque type
has no accessible fields, no initializer, and no size visible to arc code. It can only
be used as a pointer inside `extern c` blocks. The `type x = opaque` form is honest
about what it is — not an interface with zero fields, but a type that exists solely so
you can form a pointer to it and pass that pointer across the C boundary. Declared at
the top level outside the extern block, used as `*FILE` or `**sqlite3` inside it.

---

## @align
```arc
@align(16)
interface XMVECTOR {
    x: float32
    y: float32
    z: float32
    w: float32
}

@align(64)
interface CacheLine {
    data: [64]byte
}
```

Forces the interface to have at least the specified byte alignment in memory. The
alignment value must be a power of two. Used when C or C++ code, hardware, or SIMD
instructions require a type to start at a specific memory boundary. DirectX math types,
SSE/AVX vector types, and DMA buffers are common cases. The compiler inserts padding
before the type's allocation to satisfy the alignment requirement. `@align` applies to
every instance of the interface regardless of how it is allocated — stack, heap ref
counted, or manual heap.

---

## @packed
```arc
@packed
interface FileHeader {
    magic:   uint32
    version: uint16
    flags:   uint16
    size:    uint32
}

@packed
interface NetworkPacket {
    type:     uint8
    length:   uint16
    sequence: uint32
}
```

Removes all padding between fields, packing them at consecutive byte offsets with no
alignment gaps. Without `@packed`, the compiler inserts padding between fields to satisfy
each field's natural alignment requirement. With `@packed`, fields are placed exactly
where the previous field ends regardless of alignment. Used for wire protocols, file
formats, and any layout where the byte offset of every field must be exact and matches
an external specification. Accessing misaligned fields on architectures that require
alignment can cause a bus error — only use `@packed` when you control exactly how
instances are allocated and accessed.