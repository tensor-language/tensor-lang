Here's the doc with the entire Linkage section removed:

---

# Foreign Function Interface (`extern`)

Arc provides seamless interoperability with C and C++ libraries through the
`extern` keyword. Arc's ABI is already compatible with C/C++ calling
conventions, so `extern` primarily serves as a linker hint and mangling
directive.

> **Note:** Stars (`*`, `&`) only appear inside `extern` blocks.
> They faithfully describe the C/C++ interface. All regular arc code is
> star-free — the compiler handles indirection automatically.

## Quick Reference

```arc
// C functions (no name mangling)
extern c {
    func printf(*byte, ...) int32
    func malloc(usize) *void
}

// C++ functions (name mangling + vtables)
extern cpp {
    namespace DirectX {
        class ID3D11Device {
            virtual func Release(self *ID3D11Device) uint32
        }
        func CreateDevice(...) HRESULT
    }
}
```

---

## `extern c`

Use `extern c` for C libraries. No name mangling — symbols are used as-is.

### Basic Usage

```arc
namespace mymodule

extern c {
    func printf(*byte, ...) int32
    func sprintf(*byte, *byte, ...) int32
    func malloc(usize) *void
    func free(*void) void
    func memcpy(*void, *void, usize) *void
    
    func sin(float64) float64
    func cos(float64) float64
    func sqrt(float64) float64
}

func main() {
    let ptr = malloc(1024)
    defer free(ptr)
    
    printf("sqrt(2) = %f\n", sqrt(2.0))
}
```

### Symbol Renaming

Use a string literal after the function name to specify the actual C symbol:

```arc
extern c {
    func print "printf" (*byte, ...) int32
    func alloc "malloc" (usize) *void
    func c_free "free" (*void) void
    
    func sleep "_sleep" (uint32) void      // Windows
}

print("Hello %s\n", "World")
let ptr = alloc(100)
c_free(ptr)
```

### Interfaces for C Types

Arc interfaces are binary-compatible with C structs. Define matching layouts:

```arc
interface timeval {
    tv_sec: int64
    tv_usec: int64
}

interface stat {
    st_dev: uint64
    st_ino: uint64
    st_mode: uint32
    st_nlink: uint64
    st_uid: uint32
    st_gid: uint32
    st_size: int64
}

extern c {
    func fopen(*byte, *byte) *FILE
    func fclose(*FILE) int32
    func fread(*void, usize, usize, *FILE) usize
    
    func gettimeofday(*timeval, *void) int32
    func stat_file "stat" (*byte, *stat) int32
}
```

For C types with no accessible fields, declare an empty interface outside the extern block:

```arc
// Empty body = opaque layout, can only use as pointer
interface FILE {}
interface sqlite3 {}
interface sqlite3_stmt {}

extern c {
    func fopen(*byte, *byte) *FILE
    func sqlite3_open(*byte, **sqlite3) int32
}
```

### Variadic Functions

```arc
extern c {
    func printf(*byte, ...) int32
    func scanf(*byte, ...) int32
    func ioctl(int32, uint64, ...) int32
}
```

### Constants

Constants are declared at the top level, not inside `extern` blocks:

```arc
const O_RDONLY: int32 = 0
const O_WRONLY: int32 = 1
const O_RDWR: int32 = 2
const O_CREAT: int32 = 64

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

### Callbacks / Function Pointers

```arc
extern c {
    type Comparator = func(*void, *void) int32
    type SignalHandler = func(int32) void
    type ThreadFunc = func(*void) *void
    
    func qsort(*void, usize, usize, Comparator) void
    func signal(int32, SignalHandler) SignalHandler
    func pthread_create(*void, *void, ThreadFunc, *void) int32
    func atexit(func() void) int32
}

func compare_ints(a: *void, b: *void) int32 {
    let x = int32(a)
    let y = int32(b)
    return x - y
}

func on_signal(sig: int32) void {
    printf("Signal: %d\n", sig)
}

func main() {
    // memptr(&data) — get address of data for extern call
    qsort(memptr(&data), 5, sizeof(int32), compare_ints)
    signal(SIGINT, on_signal)
}
```

---

## `extern cpp`

Use `extern cpp` for C++ libraries. Handles name mangling and vtable calls.

### Basic Usage

```arc
namespace graphics

extern cpp {
    func CreateDevice(*void, uint32) *Device
    func DestroyDevice(*Device) void
}

graphics.CreateDevice(adapter, flags)
```

### Namespaces

**Nested blocks** — best for deep hierarchies:

```arc
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

graphics.DirectX.CreateDevice(...)
graphics.DirectX.DXGI.GetAdapter(0)
graphics.DirectX.DXGI.Debugging.ReportLiveObjects()
```

**Dot notation** — best for sparse declarations:

```arc
extern cpp {
    namespace std.chrono {
        func now() TimePoint
    }
    
    namespace std.filesystem {
        func exists(*Path) bool
        func create_directory(*Path) bool
    }
    
    namespace boost.asio.ip {
        func make_address(*byte) Address
    }
}

mymodule.std.chrono.now()
mymodule.std.filesystem.exists(path)
mymodule.boost.asio.ip.make_address("127.0.0.1")
```

**Mix both styles** as needed:

```arc
extern cpp {
    namespace DirectX {
        func CreateDevice(...) HRESULT
        
        namespace DXGI {
            func GetAdapter(uint32) *IDXGIAdapter
        }
    }
    
    namespace DirectX.D3D12 {
        func CreateCommandQueue(...) HRESULT
    }
}
```

### Classes and Virtual Methods

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

Generated code for virtual calls:

```asm
; device.Release()
mov rax, [rcx]        ; Load vtable pointer
call [rax + 16]       ; vtable[2] (offset 2 * 8 bytes)
```

### Constructors and Destructors

```arc
extern cpp {
    class Widget {
        new(int32, int32) *Widget
        new() *Widget
        new(*byte) *Widget
        
        delete(self *Widget) void
        
        virtual func Process(self *Widget) void
    }
}

let w = Widget.new(100, 200)
defer w.delete()
w.Process()
```

### Static Methods

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

let f = Factory.Create("config.json")
defer f.delete()

let version = Factory.GetVersion()
```

### Abstract Classes

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
    
    renderer.Init()
    render_frame(renderer, scene)
    renderer.Shutdown()
}
```

### Const Methods

```arc
extern cpp {
    class Buffer {
        new(usize) *Buffer
        delete(self *Buffer) void
        
        virtual func Size(self *const Buffer) const usize
        virtual func Data(self *const Buffer) const *byte
        virtual func IsEmpty(self *const Buffer) const bool
        
        virtual func Resize(self *Buffer, usize) void
        virtual func Clear(self *Buffer) void
        virtual func Write(self *Buffer, *byte, usize) void
    }
}
```

### Reference Parameters

```arc
extern cpp {
    class Container {
        virtual func Find(self *Container, &const Key) *Value
        virtual func Swap(self *Container, &int32, &int32) void
        virtual func At(self *Container, usize) &Item
        virtual func Front(self *const Container) const &Item
    }
    
    namespace std {
        func swap(&int32, &int32) void
    }
}
```

### Function Overloading

```arc
extern cpp {
    namespace Math {
        func Clamp(int32, int32, int32) int32
        func Clamp(float32, float32, float32) float32
        func Clamp(float64, float64, float64) float64
        
        func Abs(int32) int32
        func Abs(float32) float32
        func Abs(float64) float64
    }
}

let i = Math.Clamp(x, 0, 100)      // int32 version
let f = Math.Clamp(x, 0.0, 1.0)    // float32 version
```

If auto-mangling fails use explicit symbol:

```arc
extern cpp {
    namespace Math {
        func ClampInt   "?Clamp@Math@@YAHHHH@Z"  (int32, int32, int32) int32
        func ClampFloat "?Clamp@Math@@YAMMM@Z"   (float32, float32, float32) float32
    }
}
```

### Template Instantiations

```arc
extern cpp {
    namespace std {
        class IntVector "std::vector<int>" {
            new() *IntVector
            delete(self *IntVector) void
            
            virtual func push_back(self *IntVector, int32) void
            virtual func pop_back(self *IntVector) void
            virtual func size(self *const IntVector) const usize
            virtual func at(self *IntVector, usize) &int32
            virtual func data(self *IntVector) *int32
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
            virtual func clear(self *String) void
        }
    }
}

let vec = std.IntVector.new()
defer vec.delete()

vec.push_back(10)
vec.push_back(20)
printf("Size: %zu\n", vec.size())
```

### Callbacks / Function Pointers

```arc
extern cpp {
    type EventCallback = func(*Event, *void) void
    
    class EventSystem {
        new() *EventSystem
        delete(self *EventSystem) void
        
        virtual func Subscribe(self *EventSystem, *byte, EventCallback, *void) void
        virtual func Unsubscribe(self *EventSystem, *byte, EventCallback) void
    }
}

func on_event(event: Event, user_data: memptr) {
    printf("Event received!\n")
}

func main() {
    let events = EventSystem.new()
    defer events.delete()
    events.Subscribe("click", on_event, null)
}
```

---

## Interface Attributes

### Alignment

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

### Packed

```arc
@packed
interface FileHeader {
    magic: uint32      // offset 0
    version: uint16    // offset 4
    flags: uint16      // offset 6
    size: uint32       // offset 8
}
// Total: 12 bytes (not 16)

@packed
interface NetworkPacket {
    type: uint8
    length: uint16
    sequence: uint32
}
```

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

| Convention | Description | Platform |
| --- | --- | --- |
| `cdecl` | Caller cleans stack | All |
| `stdcall` | Callee cleans stack | Windows x86 |
| `thiscall` | `this` in ECX/RCX | MSVC C++ |
| `vectorcall` | SIMD registers | Windows |
| `fastcall` | First args in registers | Legacy |

---

## Complete Examples

### DirectX 11

```arc
namespace graphics.d3d11

type HRESULT = int32
type GUID = [16]byte

interface D3D11_BUFFER_DESC {
    byte_width: uint32
    usage: uint32
    bind_flags: uint32
    cpu_access_flags: uint32
    misc_flags: uint32
    structure_byte_stride: uint32
}

interface D3D11_SUBRESOURCE_DATA {
    sys_mem: *void
    sys_mem_pitch: uint32
    sys_mem_slice_pitch: uint32
}

const D3D11_SDK_VERSION: uint32 = 7
const D3D_DRIVER_TYPE_HARDWARE: uint32 = 1
const D3D11_USAGE_DEFAULT: uint32 = 0
const D3D11_BIND_VERTEX_BUFFER: uint32 = 1

interface IDXGIAdapter {}
interface ID3D11DeviceContext {}
interface ID3D11Texture2D {}

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
    var device: DirectX.ID3D11Device = null
    var context: DirectX.ID3D11DeviceContext = null
    
    // memptr(&x) — get address of pointer for output params
    let hr = DirectX.D3D11CreateDevice(
        null,
        D3D_DRIVER_TYPE_HARDWARE,
        null, 0, null, 0,
        D3D11_SDK_VERSION,
        memptr(&device),   // **ID3D11Device
        null,
        memptr(&context)   // **ID3D11DeviceContext
    )
    
    if hr != 0 {
        printf("Failed: %d\n", hr)
        return
    }
    defer device.Release()
    defer context.Release()
    
    let desc = D3D11_BUFFER_DESC{
        byte_width: 36,
        usage: D3D11_USAGE_DEFAULT,
        bind_flags: D3D11_BIND_VERTEX_BUFFER,
        cpu_access_flags: 0,
        misc_flags: 0,
        structure_byte_stride: 0
    }
    
    var buffer: DirectX.ID3D11Buffer = null
    
    hr = device.CreateBuffer(memptr(&desc), init_data, memptr(&buffer))
    
    if hr == 0 {
        defer buffer.Release()
        printf("Buffer created!\n")
    }
}
```

### SQLite (C)

```arc
namespace main

interface sqlite3 {}
interface sqlite3_stmt {}

const SQLITE_OK: int32 = 0
const SQLITE_ROW: int32 = 100
const SQLITE_DONE: int32 = 101

extern c {
    func sqlite3_open(*byte, **sqlite3) int32
    func sqlite3_close(*sqlite3) int32
    func sqlite3_errmsg(*sqlite3) *byte
    
    func sqlite3_prepare_v2(*sqlite3, *byte, int32, **sqlite3_stmt, **byte) int32
    func sqlite3_step(*sqlite3_stmt) int32
    func sqlite3_finalize(*sqlite3_stmt) int32
    func sqlite3_reset(*sqlite3_stmt) int32
    
    func sqlite3_bind_int(*sqlite3_stmt, int32, int32) int32
    func sqlite3_bind_text(*sqlite3_stmt, int32, *byte, int32, *void) int32
    
    func sqlite3_column_int(*sqlite3_stmt, int32) int32
    func sqlite3_column_text(*sqlite3_stmt, int32) *byte
    
    func printf(*byte, ...) int32
}

func main() {
    var db: sqlite3 = null
    
    // memptr(&db) — pass address of pointer for output param
    if sqlite3_open("test.db", memptr(&db)) != SQLITE_OK {
        printf("Failed to open: %s\n", sqlite3_errmsg(db))
        return
    }
    defer sqlite3_close(db)
    
    var stmt: sqlite3_stmt = null
    let sql = "SELECT id, name FROM users"
    
    if sqlite3_prepare_v2(db, sql, -1, memptr(&stmt), null) != SQLITE_OK {
        printf("Failed to prepare: %s\n", sqlite3_errmsg(db))
        return
    }
    defer sqlite3_finalize(stmt)
    
    for sqlite3_step(stmt) == SQLITE_ROW {
        let id = sqlite3_column_int(stmt, 0)
        let name = sqlite3_column_text(stmt, 1)
        printf("User %d: %s\n", id, name)
    }
}
```

### Mixed C and C++

```arc
namespace main

extern c {
    func printf(*byte, ...) int32
    func malloc(usize) *void
    func free(*void) void
}

extern cpp {
    namespace Physics {
        class World {
            new() *World
            delete(self *World) void
            
            virtual func Step(self *World, float32) void
            virtual func GetGravity(self *const World) const Vec3
        }
        
        namespace Collision {
            abstract class Shape {}
            
            class BoxShape {
                new(float32, float32, float32) *BoxShape
                delete(self *BoxShape) void
            }
            
            class SphereShape {
                new(float32) *SphereShape
                delete(self *SphereShape) void
            }
            
            func TestCollision(*Shape, *Shape) bool
        }
    }
    
    namespace Audio.Effects {
        class Reverb {
            new() *Reverb
            delete(self *Reverb) void
            
            virtual func SetDecay(self *Reverb, float32) void
            virtual func Process(self *Reverb, *float32, usize) void
        }
    }
}

func main() {
    let world = Physics.World.new()
    defer world.delete()
    
    let box = Physics.Collision.BoxShape.new(1.0, 2.0, 1.0)
    defer box.delete()
    
    let sphere = Physics.Collision.SphereShape.new(0.5)
    defer sphere.delete()
    
    if Physics.Collision.TestCollision(box, sphere) {
        printf("Collision detected!\n")
    }
    
    let reverb = Audio.Effects.Reverb.new()
    defer reverb.delete()
    
    reverb.SetDecay(0.8)
    world.Step(0.016)
}
```

---

## Comparison: `extern c` vs `extern cpp`

| Feature | `extern c` | `extern cpp` |
| --- | --- | --- |
| Name mangling | None | C++ ABI mangling |
| Namespaces | N/A | Nested access paths |
| Classes | N/A | Full support |
| Virtual methods | N/A | `virtual` keyword |
| Constructors | N/A | `new(...) *T` |
| Destructors | N/A | `delete(self *T)` |
| Static methods | N/A | `static func` |
| Abstract types | N/A | `abstract class` |
| Overloading | N/A | Via mangling |
| Symbol override | `func name "symbol"` | `func name "symbol"` |
| Function pointers | ✅ | ✅ |

---

## Quick Reference

```arc
// C
extern c {
    func name(types...) ReturnType
    func arc_name "c_symbol" (types...) ReturnType
    type Callback = func(types...) ReturnType
    
    stdcall func WinApiFunc(...) int32
}

// C++
extern cpp {
    namespace Outer {
        func OuterFunc() void
        
        namespace Inner {
            func InnerFunc() void
        }
    }
    
    namespace Outer.Inner.Deep {
        func DeepFunc() void
    }
    
    abstract class IInterface {
        virtual func Method(self *IInterface) void
    }
    
    class ClassName "optional::mangled::name" {
        new(types...) *ClassName
        delete(self *ClassName) void
        
        static func StaticMethod() *ClassName
        
        virtual func Method(self *ClassName, types...) ReturnType
        virtual func ConstMethod(self *const ClassName) const ReturnType
        virtual func RefParam(self *ClassName, &Type, &const Type) void
    }
    
    type Callback = func(*Event) void
    
    vectorcall func SimdFunc(Vec4, Vec4) Vec4
}

// Opaque C types — declared outside extern, not inside
interface FILE {}
interface sqlite3 {}

// Constants — declared at top level, not inside extern
const SQLITE_OK: int32 = 0
const O_RDONLY: int32 = 0

// Interface attributes
@align(16)
interface Aligned { ... }

@packed
interface Packed { ... }

// memptr usage
memptr(-1)      // cast value to memory pointer
memptr(&val)    // get address of val as memory pointer
```