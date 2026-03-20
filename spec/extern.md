# Extern

---

## The C Boundary

Arc is star-free in regular code. The `*` and `&` pointer syntax only appears inside
`extern` blocks. This is intentional — pointer arithmetic, raw addresses, and manual
indirection are C concepts that belong at the C boundary, not scattered through arc
code. The `extern` block is where that world is described precisely and contained. Once
declared, extern functions and types are called from arc code like any other function,
with no special syntax at the callsite. The compiler handles the translation.

---

## extern c

```arc
extern c {
    func printf(*byte, ...) int32
    func malloc(usize) *void
    func free(*void) void
    func sqrt(float64) float64
}

func main() {
    let ptr = malloc(1024)
    defer free(ptr)
    printf("sqrt(2) = %f\n", sqrt(2.0))
}
```

Declares bindings to C library functions. No name mangling — symbols are used exactly
as they appear in the C header. The function signatures use C-compatible types with
pointer syntax allowed only inside the block. Every function declared in an `extern c`
block is callable from arc code directly by name after the block. Arc's ABI is already
compatible with C calling conventions so no conversion layer is inserted — the call
goes directly to the C symbol. `extern c` is the right block for any C library: system
calls wrappers, libc, platform APIs, third-party C libraries.

---

## extern cpp

```arc
extern cpp {
    namespace DirectX {
        class ID3D11Device {
            virtual func Release(self *ID3D11Device) uint32
        }
        func CreateDevice(...) HRESULT
    }
}

graphics.DirectX.CreateDevice(...)
```

Declares bindings to C++ library functions and types. Unlike `extern c`, C++ symbols
are name-mangled by the compiler according to the C++ ABI. `extern cpp` handles that
mangling automatically based on the declared namespaces, class names, and function
signatures. It also handles vtable layout for virtual methods and constructor/destructor
calling conventions. The full C++ object model — namespaces, classes, virtual dispatch,
overloading, templates — can be described inside an `extern cpp` block. Arc's ABI is
compatible with the C++ calling convention so no wrapper layer is needed.

---

## Symbol Renaming

```arc
extern c {
    func print  "printf"  (*byte, ...) int32
    func alloc  "malloc"  (usize) *void
    func c_free "free"    (*void) void
    func sleep  "_sleep"  (uint32) void
}

print("Hello %s\n", "World")
let ptr = alloc(100)
c_free(ptr)
```

A string literal after the function name specifies the actual symbol name in the
compiled library. The name before the string is what arc code uses to call the function.
This allows arc-idiomatic names at the callsite while binding to the real C symbol under
the hood. Useful when the C name conflicts with an arc keyword, when the platform
decorates symbol names with underscores or prefixes, or simply when the C name is
awkward and a cleaner arc name is preferred. Works in both `extern c` and `extern cpp`
blocks.

---

## opaque Type

```arc
type FILE         = opaque
type sqlite3      = opaque
type sqlite3_stmt = opaque

extern c {
    func fopen(*byte, *byte) *FILE
    func fclose(*FILE) int32
    func sqlite3_open(*byte, **sqlite3) int32
}
```

Opaque types give a C type a name in arc without describing its layout. Declared at the
top level outside the extern block with `type X = opaque`. Used as a pointer inside
extern blocks — you can form `*FILE` and `**sqlite3`, pass them across the C boundary,
and store them in variables, but you cannot access fields or initialize them directly.
The compiler enforces that no field access is ever attempted on an opaque type. This is
the correct pattern for any C type whose internals are private — file handles, database
handles, platform objects, any type the C library treats as an opaque handle.

---

## Interface for C Struct

```arc
interface timeval {
    tv_sec:  int64
    tv_usec: int64
}

interface stat {
    st_dev:   uint64
    st_ino:   uint64
    st_mode:  uint32
    st_nlink: uint64
    st_uid:   uint32
    st_gid:   uint32
    st_size:  int64
}

extern c {
    func gettimeofday(*timeval, *void) int32
    func stat_file "stat" (*byte, *stat) int32
}
```

Arc interfaces are binary-compatible with C structs when field types and order match
exactly. Declare an interface outside the extern block with fields that mirror the C
struct layout in the same order. The compiler lays out the fields with natural alignment
matching C's default struct layout rules. Pass a pointer to the interface into the C
function and the C code reads and writes the fields directly. Field names do not need
to match the C names — only the types and order matter for binary compatibility. Use
`@packed` on the interface when the C struct uses `__attribute__((packed))` or
`#pragma pack`.

---

## Constants

```arc
const O_RDONLY: int32 = 0
const O_WRONLY: int32 = 1
const O_RDWR:   int32 = 2
const O_CREAT:  int32 = 64

extern c {
    func open(*byte, int32, ...) int32
    func close(int32) int32
    func read(int32, *void, usize) isize
    func write(int32, *void, usize) isize
}

func main() {
    let fd = open("test.txt", O_RDONLY)
    if fd >= 0 {
        defer close(fd)
    }
}
```

Constants that correspond to C `#define` values are declared at the top level outside
the extern block using regular arc `const` declarations. There is no mechanism to import
C preprocessor defines directly — they are transcribed as typed arc constants. The type
should match what the C function expects at the parameter site. Constants declared this
way are compile-time values with zero runtime cost, identical in behavior to any other
arc `const`.

---

## Variadic Functions

```arc
extern c {
    func printf(*byte, ...) int32
    func scanf(*byte,  ...) int32
    func ioctl(int32, uint64, ...) int32
}
```

C variadic functions are declared with `...` as the final parameter. Arc passes
additional arguments to variadic C functions using the platform's C variadic calling
convention — no arc-level type checking is performed on the variadic arguments since
C has no type information for them at the call boundary. The fixed parameters before
`...` are type-checked normally. Variadic declarations work in both `extern c` and
`extern cpp` blocks.

---

## Function Pointers

```arc
extern c {
    type Comparator    = func(*void, *void) int32
    type SignalHandler = func(int32) void
    type ThreadFunc    = func(*void) *void

    func qsort(*void, usize, usize, Comparator) void
    func signal(int32, SignalHandler) SignalHandler
    func pthread_create(*void, *void, ThreadFunc, *void) int32
}

func compare_ints(a: *void, b: *void) int32 {
    let x = int32(a)
    let y = int32(b)
    return x - y
}

func main() {
    qsort(memptr(&data), 5, sizeof(int32), compare_ints)
    signal(SIGINT, on_signal)
}
```

C function pointer types are declared inside extern blocks with the `type` keyword.
The declared type alias can then be used as a parameter type in function declarations
within the same or other extern blocks. Arc functions whose signatures match the declared
type can be passed directly as callbacks — the compiler verifies the signature matches
at the callsite. `memptr(&val)` is used to get the address of arc data for passing to
C functions that expect a raw pointer.

---

## memptr at the C Boundary

```arc
extern c {
    func sqlite3_open(*byte, **sqlite3) int32
    func connect(int32, *void, uint32) int32
}

var db: sqlite3 = null
sqlite3_open("test.db", memptr(&db))

let addr = new SockAddrIn{}
mem_zero(addr, sizeof(SockAddrIn))
defer delete(addr)
connect(fd, memptr(&addr), uint32(sizeof(SockAddrIn)))
```

`memptr` is the bridge between arc variables and C pointer parameters. Two forms exist:
`memptr(&val)` gets the address of an arc variable as a raw memory pointer, and
`memptr(x)` casts an integer value to a memory pointer. At the C boundary both forms
produce the raw pointer value that C functions expect. `memptr` is only valid when
calling extern functions — it is the explicit marker that you are crossing the C
boundary with a raw address. Pointer arithmetic goes through `usize` first:
`memptr(usize(p) + 16)`.

---

## Namespaces in extern cpp

```arc
// nested blocks — best for deep hierarchies
extern cpp {
    namespace DirectX {
        func CreateDevice(...) HRESULT

        namespace DXGI {
            func GetAdapter(uint32) *IDXGIAdapter

            namespace Debugging {
                func ReportLiveObjects(...) void
            }
        }
    }
}

// dot notation — best for sparse declarations
extern cpp {
    namespace std.chrono {
        func now() TimePoint
    }

    namespace std.filesystem {
        func exists(*Path) bool
        func create_directory(*Path) bool
    }
}

// both styles can be mixed
extern cpp {
    namespace DirectX {
        func CreateDevice(...) HRESULT
    }

    namespace DirectX.D3D12 {
        func CreateCommandQueue(...) HRESULT
    }
}
```

C++ namespaces are declared inside `extern cpp` using either nested blocks or dot
notation. Nested blocks mirror the C++ source hierarchy and are clearest for deep or
dense namespace trees. Dot notation is more concise for declaring one or two functions
in a deep path without nesting many blocks. Both styles produce identical bindings —
the choice is purely readability. Mixed usage in the same extern block is valid. Arc
code accesses the declared functions through the full namespace path from the current
module: `graphics.DirectX.DXGI.GetAdapter(0)`.

---

## Classes in extern cpp

```arc
extern cpp {
    namespace DirectX {
        class ID3D11Device {
            virtual func QueryInterface(self *ID3D11Device, *GUID, **void) HRESULT
            virtual func AddRef(self *ID3D11Device) uint32
            virtual func Release(self *ID3D11Device) uint32
            virtual func CreateBuffer(
                self *ID3D11Device,
                *D3D11_BUFFER_DESC,
                *D3D11_SUBRESOURCE_DATA,
                **ID3D11Buffer
            ) HRESULT
        }
    }
}
```

C++ classes are declared with the `class` keyword inside `extern cpp` blocks. Virtual
methods are declared with `virtual func` and take an explicit `self` pointer as the
first parameter — this mirrors the hidden `this` pointer in the C++ ABI. The compiler
generates correct vtable dispatch for virtual calls: loading the vtable pointer from
the object, then calling through the correct vtable slot at the right byte offset.
Non-virtual methods and data members of C++ classes are not currently accessible —
only virtual dispatch, constructors, destructors, and static methods are exposed through
the extern cpp class declaration.

---

## Constructors and Destructors

```arc
extern cpp {
    class Widget {
        new()           *Widget
        new(int32, int32) *Widget
        new(*byte)      *Widget

        delete(self *Widget) void

        virtual func Process(self *Widget) void
    }
}

let w = Widget.new(100, 200)
defer w.delete()
w.Process()
```

C++ constructors are declared as `new` with their parameter list and a return type of
`*ClassName`. Multiple constructors with different signatures express C++ constructor
overloading. The destructor is declared as `delete` with a self pointer parameter.
Calling `new` on a C++ class invokes the constructor through the C++ ABI and returns
a pointer to the constructed object. Calling `delete` invokes the destructor and frees
the object. Always pair `new` with `defer delete` immediately after construction —
the same pattern as arc's own manual allocation.

---

## Static Methods

```arc
extern cpp {
    class Factory {
        static func Create(*byte) *Factory
        static func GetInstance() *Factory
        static func GetVersion() int32

        virtual func Process(self *Factory) void
        delete(self *Factory) void
    }
}

let f       = Factory.Create("config.json")
let version = Factory.GetVersion()
defer f.delete()
```

C++ static methods are declared with the `static` keyword. Static methods belong to
the class rather than an instance — they are called through the class name rather than
through a pointer to an object. No `self` parameter. Static methods can coexist with
virtual methods and constructors in the same class declaration. Arc calls them with
dot notation on the class name: `Factory.Create(...)`, `Factory.GetVersion()`.

---

## Abstract Classes

```arc
extern cpp {
    abstract class IRenderer {
        virtual func Init(self *IRenderer) bool
        virtual func Draw(self *IRenderer, *Scene) void
        virtual func Shutdown(self *IRenderer) void
    }

    class D3D11Renderer {
        new() *D3D11Renderer
        delete(self *D3D11Renderer) void

        virtual func Init(self *D3D11Renderer) bool
        virtual func Draw(self *D3D11Renderer, *Scene) void
        virtual func Shutdown(self *D3D11Renderer) void
    }
}

func render_frame(renderer: *IRenderer, scene: *Scene) {
    renderer.Draw(scene)
}

func main() {
    let renderer = D3D11Renderer.new()
    defer renderer.delete()
    render_frame(renderer, scene)
}
```

C++ abstract classes — those with pure virtual methods — are declared with the
`abstract` keyword. An abstract class cannot be instantiated directly and has no
constructor declaration. It exists to describe an interface through which concrete
implementations are called. Arc functions that accept `*IRenderer` can receive any
concrete class that implements the same vtable layout. The compiler generates the
correct virtual dispatch through the abstract class pointer. This is the standard
pattern for C++ polymorphism at the extern boundary: declare the abstract interface,
declare the concrete class, use the abstract pointer type in arc code that works with
either.

---

## Const Methods

```arc
extern cpp {
    class Buffer {
        new(usize) *Buffer
        delete(self *Buffer) void

        virtual func Size(self *const Buffer)  const usize
        virtual func Data(self *const Buffer)  const *byte
        virtual func IsEmpty(self *const Buffer) const bool

        virtual func Resize(self *Buffer, usize) void
        virtual func Clear(self *Buffer) void
        virtual func Write(self *Buffer, *byte, usize) void
    }
}
```

C++ const methods — those declared with `const` after the parameter list — are
declared in arc with `*const ClassName` as the self parameter type and `const` before
the return type. This matches the C++ ABI where const methods receive a pointer to a
const object. The distinction matters at the binary level: some C++ compilers and
optimizers rely on const correctness in the vtable layout. Declaring const methods
correctly ensures arc calls them through the right vtable slot with the right pointer
type.

---

## Reference Parameters

```arc
extern cpp {
    class Container {
        virtual func Find(self *Container, &const Key) *Value
        virtual func Swap(self *Container, &int32, &int32) void
        virtual func At(self *Container, usize) &Item
    }

    namespace std {
        func swap(&int32, &int32) void
    }
}
```

C++ reference parameters are declared with `&` before the type inside extern cpp
blocks. This is the only place `&` appears as a type qualifier in arc — in regular arc
code `&` at the callsite means mutable reference, but inside extern blocks it describes
the C++ reference type in the binary interface. `&const Type` is a const reference.
`&Type` is a non-const reference. Reference return types (`&Item`) indicate the C++
function returns a reference to an existing object rather than a new value. The compiler
maps these to the correct C++ reference calling convention.

---

## Function Overloading

```arc
extern cpp {
    namespace Math {
        func Clamp(int32, int32, int32) int32
        func Clamp(float32, float32, float32) float32
        func Clamp(float64, float64, float64) float64
    }
}

let i = Math.Clamp(x, 0, 100)       // resolves to int32 version
let f = Math.Clamp(x, 0.0, 1.0)     // resolves to float32 version
```

C++ function overloading is expressed by declaring multiple functions with the same
name and different parameter types in the same namespace or class. Arc resolves the
correct overload at the callsite from the argument types, then generates the correctly
mangled C++ symbol for that specific overload. When auto-mangling fails or produces
the wrong symbol, an explicit symbol string overrides it:

```arc
extern cpp {
    namespace Math {
        func ClampInt   "?Clamp@Math@@YAHHHH@Z"  (int32,   int32,   int32)   int32
        func ClampFloat "?Clamp@Math@@YAMMM@Z"   (float32, float32, float32) float32
    }
}
```

---

## Template Instantiations

```arc
extern cpp {
    namespace std {
        class IntVector "std::vector<int>" {
            new() *IntVector
            delete(self *IntVector) void

            virtual func push_back(self *IntVector, int32) void
            virtual func size(self *const IntVector) const usize
            virtual func at(self *IntVector, usize) &int32
            virtual func clear(self *IntVector) void
            virtual func empty(self *const IntVector) const bool
        }

        class String "std::string" {
            new() *String
            new(*byte) *String
            delete(self *String) void

            virtual func c_str(self *const String) const *byte
            virtual func size(self *const String) const usize
            virtual func empty(self *const String) const bool
        }
    }
}

let vec = std.IntVector.new()
defer vec.delete()
vec.push_back(10)
```

C++ template instantiations are bound by giving the arc class declaration an explicit
C++ symbol name as a string after the class name. The string is the fully qualified
C++ template instantiation name — `"std::vector<int>"`, `"std::string"`. Arc does not
generate template instantiations itself; it binds to ones that already exist in the
compiled C++ library. Each distinct instantiation is a separate arc class declaration
with its own name. The virtual method declarations describe the actual compiled vtable
of that specific instantiation.

---

## Calling Conventions

```arc
extern c {
    func printf(*byte, ...) int32

    stdcall func MessageBoxA(*void, *byte, *byte, uint32) int32
    stdcall func GetLastError() uint32
    stdcall func ExitProcess(uint32) void
}

extern cpp {
    class SomeClass {
        thiscall virtual func Method(self *SomeClass, int32) void
    }

    vectorcall func XMVectorAdd(XMVECTOR, XMVECTOR) XMVECTOR
}
```

The default calling convention is `cdecl` for `extern c` and the platform C++ ABI
convention for `extern cpp`. When a C or C++ function uses a non-default convention,
the convention keyword appears before `func` in the declaration. The compiler generates
the correct call sequence for the declared convention.

| Convention    | Description                  | Platform        |
| ------------- | ---------------------------- | --------------- |
| `cdecl`       | Caller cleans stack          | All             |
| `stdcall`     | Callee cleans stack          | Windows x86     |
| `thiscall`    | `this` in ECX/RCX            | MSVC C++        |
| `vectorcall`  | SIMD registers               | Windows         |
| `fastcall`    | First args in registers      | Legacy          |

---

## @align on Interface

```arc
@align(16)
interface XMVECTOR {
    x: float32
    y: float32
    z: float32
    w: float32
}

@align(4096)
interface PageAligned {
    data: [4096]byte
}
```

When a C or C++ type requires a specific memory alignment — SIMD vector types, DMA
buffers, cache-line aligned structures — declare the corresponding arc interface with
`@align(N)` where N is the required alignment in bytes and must be a power of two.
The compiler inserts the necessary alignment padding wherever the type is allocated.
`@align` applies regardless of allocation method — stack, heap ref counted, or manual
heap. DirectX math types, SSE/AVX vector types, and hardware buffer descriptors are
the most common cases.

---

## @packed on Interface

```arc
@packed
interface FileHeader {
    magic:   uint32    // offset 0
    version: uint16    // offset 4
    flags:   uint16    // offset 6
    size:    uint32    // offset 8
}                      // total: 12 bytes, no padding

@packed
interface NetworkPacket {
    type:     uint8
    length:   uint16
    sequence: uint32
}
```

Removes all padding between fields, placing them at consecutive byte offsets with no
alignment gaps. Without `@packed` the compiler inserts padding between fields to satisfy
natural alignment requirements. With `@packed` fields are placed exactly where the
previous field ends regardless of alignment. Used for wire protocols, file format
headers, and any layout where byte offsets are specified by an external standard and
must be exact. Accessing misaligned fields on architectures that require alignment can
cause a bus error — only use `@packed` when you control exactly how instances are
allocated and accessed, and when the C or C++ side uses the same packed layout.

---

## Complete Example — SQLite

```arc
namespace main

type sqlite3      = opaque
type sqlite3_stmt = opaque

const SQLITE_OK:   int32 = 0
const SQLITE_ROW:  int32 = 100
const SQLITE_DONE: int32 = 101

extern c {
    func sqlite3_open(*byte, **sqlite3) int32
    func sqlite3_close(*sqlite3) int32
    func sqlite3_errmsg(*sqlite3) *byte

    func sqlite3_prepare_v2(*sqlite3, *byte, int32, **sqlite3_stmt, **byte) int32
    func sqlite3_step(*sqlite3_stmt) int32
    func sqlite3_finalize(*sqlite3_stmt) int32

    func sqlite3_column_int(*sqlite3_stmt, int32) int32
    func sqlite3_column_text(*sqlite3_stmt, int32) *byte

    func printf(*byte, ...) int32
}

func main() {
    var db: sqlite3 = null

    if sqlite3_open("test.db", memptr(&db)) != SQLITE_OK {
        printf("failed to open: %s\n", sqlite3_errmsg(db))
        return
    }
    defer sqlite3_close(db)

    var stmt: sqlite3_stmt = null
    let sql = "SELECT id, name FROM users"

    if sqlite3_prepare_v2(db, sql, -1, memptr(&stmt), null) != SQLITE_OK {
        printf("failed to prepare: %s\n", sqlite3_errmsg(db))
        return
    }
    defer sqlite3_finalize(stmt)

    for sqlite3_step(stmt) == SQLITE_ROW {
        let id   = sqlite3_column_int(stmt, 0)
        let name = sqlite3_column_text(stmt, 1)
        printf("user %d: %s\n", id, name)
    }
}
```

---

## Complete Example — DirectX 11

```arc
namespace graphics.d3d11

type HRESULT = int32
type GUID    = [16]byte

interface D3D11_BUFFER_DESC {
    byte_width:             uint32
    usage:                  uint32
    bind_flags:             uint32
    cpu_access_flags:       uint32
    misc_flags:             uint32
    structure_byte_stride:  uint32
}

type IDXGIAdapter       = opaque
type ID3D11DeviceContext = opaque

const D3D11_SDK_VERSION:       uint32 = 7
const D3D_DRIVER_TYPE_HARDWARE: uint32 = 1
const D3D11_USAGE_DEFAULT:     uint32 = 0
const D3D11_BIND_VERTEX_BUFFER: uint32 = 1

extern cpp {
    namespace DirectX {
        func D3D11CreateDevice(
            *IDXGIAdapter,
            uint32,
            *void,
            uint32,
            *uint32,
            uint32,
            uint32,
            **ID3D11Device,
            *uint32,
            **ID3D11DeviceContext
        ) HRESULT

        class ID3D11Device {
            virtual func QueryInterface(self *ID3D11Device, *GUID, **void) HRESULT
            virtual func AddRef(self *ID3D11Device) uint32
            virtual func Release(self *ID3D11Device) uint32
            virtual func CreateBuffer(
                self *ID3D11Device,
                *D3D11_BUFFER_DESC,
                *D3D11_SUBRESOURCE_DATA,
                **ID3D11Buffer
            ) HRESULT
        }

        class ID3D11Buffer {
            virtual func QueryInterface(self *ID3D11Buffer, *GUID, **void) HRESULT
            virtual func AddRef(self *ID3D11Buffer) uint32
            virtual func Release(self *ID3D11Buffer) uint32
        }
    }
}

extern c {
    func printf(*byte, ...) int32
}

func main() {
    var device:  DirectX.ID3D11Device        = null
    var context: DirectX.ID3D11DeviceContext = null

    let hr = DirectX.D3D11CreateDevice(
        null,
        D3D_DRIVER_TYPE_HARDWARE,
        null, 0, null, 0,
        D3D11_SDK_VERSION,
        memptr(&device),
        null,
        memptr(&context)
    )

    if hr != 0 {
        printf("failed: %d\n", hr)
        return
    }
    defer device.Release()
    defer context.Release()

    let desc = D3D11_BUFFER_DESC{
        byte_width:   36,
        usage:        D3D11_USAGE_DEFAULT,
        bind_flags:   D3D11_BIND_VERTEX_BUFFER,
    }

    var buffer: DirectX.ID3D11Buffer = null
    hr = device.CreateBuffer(memptr(&desc), init_data, memptr(&buffer))

    if hr == 0 {
        defer buffer.Release()
        printf("buffer created\n")
    }
}
```