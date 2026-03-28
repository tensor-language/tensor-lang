
## Language Overview

### The Memory Model: Safe by Default, Unsafe by Choice

Tensor completely eliminates the ambiguous `*` syntax found in C/C++. Memory intent is declared clearly reading left-to-right.

```tensor
// 1. Safe, read-only reference
func read_config(config: ref Config) {
    io.print(config.name)
}

// 2. Safe, mutable reference
func update_score(score: mut int32) {
    score += 10
}

// 3. Unsafe, unmanaged pointer (for C++ interop and raw memory)
func process_raw(buffer: ptr byte) {
    // ...
}
```

### Strict Namespace Rules (`.` vs `::`)

To keep code auditable, Tensor enforces a strict visual boundary between standard logic and memory math:

  * **The Dot (`.`)**: Used for standard module and namespace resolution (e.g., `io.print`, `ggml.init`).
  * **The Double Colon (`::`)**: Strictly reserved for low-level memory operations and casting. If you see `::`, raw memory is being manipulated.

<!-- end list -->

```tensor
let val = 42

// Safe code uses dots or nothing
let safe_val: ref int32 = ref val

// Unsafe memory operations trigger the :: syntax
let raw_ptr: ptr int32 = ptr::addr_of(val)
let next_ptr: ptr int32 = ptr::offset(raw_ptr, 1)
let empty_ptr: ptr void = ptr::null()

// Memory casting
let float_ptr = ptr::cast[ptr float32](raw_ptr)
```

### Native C/C++ Header Import

You never have to write `extern` blocks or wrapper libraries. Tensor's compiler natively parses Clang/C++ ASTs. Just import the header and use the functions natively. Tensor automatically maps C pointers to `ptr` types and C structs to `interface` types.

#### POSIX Socket Client (Native C)

```tensor
import "io"
import "sys/socket.h"
import "netinet/in.h"
import "arpa/inet.h"
import "unistd.h"

func main() {
    // 1. Call standard C functions directly
    let sock = socket(AF_INET, SOCK_STREAM, 0)
    if sock < 0 { return }
    defer close(sock)

    // 2. Instantiate the C struct natively (Tensor treats it as an interface)
    var server = sockaddr_in{
        sin_family: AF_INET,
        sin_port: htons(8080),
        sin_addr: in_addr{ s_addr: inet_addr("127.0.0.1".c_str()) }
    }

    // 3. Use Tensor's memory namespace to get pointers for the C function
    let raw_addr: ptr sockaddr_in = ptr::addr_of(server)
    let base_addr: ptr sockaddr = ptr::cast[ptr sockaddr](raw_addr)

    // 4. Pass the pointer directly to C's connect()
    connect(sock, base_addr, sizeof(sockaddr_in))
}
```

#### AI Model Inference (Native C++)

```tensor
import "io"
import "ggml.h" // Natively import C++ library

func load_weights(weights: []float32) {
    // Tensor natively understands the C++ struct from the header
    let params = ggml_init_params{
        mem_size: 1024 * 1024 * 512,
        mem_buffer: ptr::null(),
        no_alloc: false
    }
    
    // Pass by safe reference (ref) to the C++ function
    let ctx: ptr ggml_context = ggml_init(ref params)
    defer ggml_free(ctx)
    
    // Create a tensor
    let t: ptr ggml_tensor = ggml_new_tensor_1d(ctx, GGML_TYPE_F32, weights.len())
    
    // Get the raw pointer to our Tensor slice data using ::
    let raw_weights: ptr void = ptr::addr_of(weights[0])
    
    // Pass raw memory directly to C++
    ggml_set_data(t, raw_weights)
    
    io.print("Weights loaded using native header parsing!")
}
```

### Hardware Acceleration

Tensor can compile functions to run on specialized hardware when you need maximum performance. The target (`gpu.cuda`, `gpu.metal`, `tpu`) is set in the build file, not via messy pragmas in your source code.

```tensor
// CPU version (default)
func process_data(data: []float32, size: usize) {
    for let i: usize = 0; i < size; i++ {
        data[i] = data[i] * 2.0
    }
}

// GPU version
gpu func process_gpu(data: ptr float32, size: usize) {
    // thread_id() is an intrinsic valid only inside gpu func
    let idx = thread_id()
    if idx < size {
        let target_ptr = ptr::offset(data, idx)
        let current_val = ptr::read(target_ptr)
        ptr::write(target_ptr, current_val * 2.0)
    }
}
```
