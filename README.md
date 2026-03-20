<h1 align="center">
  <img src="./github/tensor_logo.png" alt="Tensor Language" width="400px">
</h1>

<h4 align="center">Systems Programming Language<br>High Performance, Native Code, Modern Syntax</h4>

<p align="center">
    <img src="https://img.shields.io/badge/Version-2.0-blue" alt="Version">
    <img src="https://img.shields.io/badge/Targets-CPU%20%7C%20GPU%20%7C%20TPU-purple" alt="Targets">
    <img src="https://img.shields.io/badge/License-MIT-green" alt="License">
</p>

---

## What is Tensor?

Tensor is a modern systems programming language for building native applications, servers, CLI tools, AI model inference, and kernel drivers with high performance.

Tensor uses automatic reference counting for memory safety, provides zero-cost abstractions, and seamless C/C++ interoperability. The language compiles to efficient native code for x86-64, ARM64, and other CPU architectures.

**Hardware acceleration built-in:** When you need it, Tensor functions can also compile to GPUs and TPUs without leaving your codebase or learning new APIs.

---

## Language Overview

### AI Model Integration

```tensor
import "ai"
import "io"

func main() {
    // Load model from file
    let model = ai.load_model("models/model-7b.gguf")
    defer model.free()
    
    // Configure inference
    let config = ai.InferenceConfig{
        temperature: 0.7,
        top_p: 0.9,
        max_tokens: 512
    }
    
    // Run inference
    let prompt = "Explain quantum computing in simple terms:"
    let tokens = model.tokenize(prompt)
    
    for token in model.generate(tokens, config) {
        let text = model.decode(token)
        io.printf("%s", text)
    }
}
```

### Type System

```tensor
// Fixed-width integers
let i: int32 = -500
let u: uint64 = 10000

// Pointer-sized integers
let size: usize = 100      // Unsigned (array indexing, sizes)
let offset: isize = -4     // Signed (offsets)

// Floating point
let f: float32 = 3.14
let d: float64 = 2.71828

// Slices (view into memory, ptr + length, no allocation)
let view: []byte = buffer[0..64]

// Mutable references — &mut in type, & at call site
func increment(x: &mut int32) { x += 1 }
increment(&i)

// Interfaces (value types — stack allocated with let)
interface Point {
    x: int32
    y: int32
}

// Interfaces (reference types — heap allocated with var)
interface Client {
    name: string
    port: int32
}
```

### Memory Management

Tensor uses automatic reference counting for `var` declarations and manual allocation for low-level work:

```tensor
// Manual heap allocation
let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)

// Heap allocation (via FFI)
extern c {
    func malloc(usize) *void
    func free(*void) void
}

let ptr = malloc(1024)
defer free(ptr)  // Cleanup on scope exit

// Get address of variable for extern calls
let val = 42
let raw = memptr(&val)    // address of val as memory pointer
let sentinel = memptr(-1) // cast integer value to memory pointer
```

### Foreign Function Interface

Direct interop with C and C++:

```tensor
// C libraries
extern c {
    func printf(*byte, ...) int32
    func sqlite3_open(*byte, **sqlite3) int32
}

// C++ libraries
extern cpp {
    namespace DirectX {
        class ID3D11Device {
            virtual func CreateBuffer(
                self *ID3D11Device,
                *D3D11_BUFFER_DESC,
                **ID3D11Buffer
            ) HRESULT
        }
    }
}
```

### Async/Await

```tensor
async func fetch_data(url: string) string {
    let response = await http.get(url)
    return response.body
}

async func main() {
    let data = await fetch_data("https://api.example.com")
    io.print(data)
}
```

### Generics

Monomorphized at compile time. Note the use of square brackets `[...]` for type parameters.

```tensor
func swap[T](a: &mut T, b: &mut T) {
    let tmp: T = a
    a = b
    b = tmp
}

let x = 10
let y = 20
swap(&x, &y)

interface Box[T] {
    value: T
}

func get[T](self b: Box[T]) T {
    return b.value
}
```

---

## Example Programs

### HTTP Server

```tensor
namespace main

import "net"
import "io"

async func handle_request(conn: net.TcpStream) {
    let buffer: vector[byte] = {1, 2, 3}
    let bytes_read = await conn.read(&buffer)
    
    let response = "HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nOK"
    await conn.write(response.as_bytes())
}

async func main() {
    let listener = net.TcpListener.bind("0.0.0.0:8080")
    io.print("Server listening on port 8080")
    
    for {
        let (conn, addr) = await listener.accept()
        // Fire-and-forget concurrent processing
        process func(c: net.TcpStream) {
             handle_request(c)
        }(conn)
    }
}
```

### Database Application

```tensor
namespace main

extern c {
    interface sqlite3 {}
    interface sqlite3_stmt {}
    
    const SQLITE_OK: int32 = 0
    const SQLITE_ROW: int32 = 100
    
    func sqlite3_open(*byte, **sqlite3) int32
    func sqlite3_close(*sqlite3) int32
    func sqlite3_prepare_v2(*sqlite3, *byte, int32, **sqlite3_stmt, **byte) int32
    func sqlite3_step(*sqlite3_stmt) int32
    func sqlite3_column_text(*sqlite3_stmt, int32) *byte
    func printf(*byte, ...) int32
}

func main() {
    var db: sqlite3 = null
    
    if sqlite3_open("app.db", memptr(&db)) != SQLITE_OK {
        printf("Failed to open database\n")
        return
    }
    defer sqlite3_close(db)
    
    var stmt: sqlite3_stmt = null
    sqlite3_prepare_v2(db, "SELECT name FROM users", -1, memptr(&stmt), null)
    
    for sqlite3_step(stmt) == SQLITE_ROW {
        let name = sqlite3_column_text(stmt, 0)
        printf("User: %s\n", name)
    }
}
```

### Graphics Application (DirectX 11)

```tensor
namespace main

type HRESULT = int32
const D3D11_SDK_VERSION: uint32 = 7
const D3D_DRIVER_TYPE_HARDWARE: uint32 = 1

extern cpp {
    namespace DirectX {
        func D3D11CreateDevice(
            *void, uint32, *void, uint32, *uint32, uint32,
            uint32, **ID3D11Device, *uint32, **ID3D11DeviceContext
        ) HRESULT

        class ID3D11Device {
            virtual func Release(self *ID3D11Device) uint32
        }
        
        class ID3D11DeviceContext {
            virtual func Release(self *ID3D11DeviceContext) uint32
        }
    }
}

func main() {
    var device:  DirectX.ID3D11Device        = null
    var context: DirectX.ID3D11DeviceContext = null

    // Direct C++ Interop
    let hr = DirectX.D3D11CreateDevice(
        null, D3D_DRIVER_TYPE_HARDWARE, null, 0, null, 0,
        D3D11_SDK_VERSION, memptr(&device), null, memptr(&context)
    )

    if hr != 0 {
        return
    }
    defer device.Release()
    defer context.Release()
    
    // Use device...
}
```

### Kernel Module

```tensor
namespace driver

import "linux/kernel/driver"
import "linux/kernel/log"

func init_module() int32 {
    log.info("Driver loading")
    
    let dev = driver.CharDevice.new("custom_device", 0)
    
    dev.on_read(func(file: driver.File, buffer: []byte, size: uint64) int64 {
        let data = "Hello from kernel"
        mem_copy(memptr(&buffer[0]), data.as_bytes(), data.len())
        return int64(data.len())
    })
    
    return 0
}

func cleanup_module() {
    log.info("Driver unloading")
}
```

---

## Hardware Acceleration

Tensor can compile functions to run on specialized hardware when you need maximum performance. The target (cuda, metal, rocm, etc.) is set in `build.vs`, not in source code.

```tensor
namespace compute

import "ai"

// CPU version (default)
func process_data(data: []float32, size: usize) {
    for let i: usize = 0; i < size; i++ {
        data[i] = data[i] * 2.0
    }
}

// GPU version — target set in build.vs
gpu func process_gpu(data: []float32, size: usize) {
    // Note: thread_id() is an intrinsic valid only inside gpu func
    let idx = thread_id()
    if idx < size {
        data[idx] = data[idx] * 2.0
    }
}

// TPU version
gpu func process_tpu(data: Tensor) Tensor {
    return data.multiply(2.0)
}

// Train model on GPU
gpu func train(model: &mut ai.Model, data: Tensor) {
    for epoch in 0..100 {
        let loss = model.forward(data)
        model.backward(loss)
        model.step()
    }
}

func main() {
    let model = ai.load_model("models/model-13b.gguf")
    let data = ai.load_tensor("training_data.bin")
    
    await train(&model, data)
    
    io.printf("Training complete\n")
}
```

---

## Package Management

Tensor downloads packages via HTTPS to a local cache (`~/.tensor/`). No system package managers required.

### Source Code

```tensor
// main.vs
namespace main

import c "sqlite3"
import c "curl"
import "io"
import "ai"

func main() {
    // Use imported libraries
}
```

The compiler detects your platform and downloads the appropriate packages to `~/.tensor/cache/`.

---

## Supported Targets

### CPU Architectures

* x86-64 (Intel, AMD)
* ARM64 (Apple Silicon, ARM servers)
* RISC-V (in progress)

### Operating Systems

* Linux (Ubuntu, Debian, Arch, Fedora, Alpine, etc.)
* macOS (Intel and Apple Silicon)
* Windows (x64)
* FreeBSD

### Accelerators

* NVIDIA GPUs (`gpu.cuda`)
* AMD GPUs (`gpu.rocm`)
* Apple Silicon (`gpu.metal`)
* Intel GPUs (`gpu.oneapi`)
* Google TPUs (`tpu`)
* AWS Trainium (`aws.trainium`)

---

## Installation

```bash
git clone https://github.com/tensor-language/tensor-lang
cd tensor-lang/cmd
./build build
./test_runner
```

### Build a Program

```bash
./tensor main.vs -o main
./main
```

---

## Documentation

* **[Language Reference](docs/reference.md)** - Complete syntax and semantics
* **[Grammar Specification](docs/grammar_1.0.md)** - Language grammar
* **[Package Management](docs/package_manager.md)** - Dependency resolution
* **[Foreign Function Interface](docs/extern.md)** - C/C++ interop
* **[Kernel Drivers](docs/kernel_drivers.md)** - Systems programming
* **[Compiler Intrinsics](docs/intrinsics_1.2.md)** - Built-in functions

---

## Current Status

**Beta Release**

Working:

* Core language features
* C/C++ FFI
* Package management
* CPU compilation (x86-64, ARM64)
* GPU compilation (CUDA, Metal)
* Kernel driver support

In Development:

* Standard library
* AI model runtime
* TPU backend
* Additional GPU targets
* Tooling (LSP, debugger)

---

## License

Licensed under either of

* Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

at your option.

## Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in the work by you, as defined in the Apache-2.0 license, shall be dual licensed as above, without any additional terms or conditions.
