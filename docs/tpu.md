# Arc Google TPU Model

**Key Concept:** Functions marked `async func<tpu>` compile **entirely** to StableHLO (Stable High-Level Operations). Everything - control flow, loops, operations - becomes StableHLO operations.

---

## How `<tpu>` Functions Work

### The Entire Function Becomes StableHLO

```arc
async func example<tpu>(data: *float32, n: usize) *float32 {
    // EVERYTHING in here → StableHLO operations
    let idx = tpu.replica_id()
    
    if idx < n {  // ← StableHLO CONDITIONAL/SELECT
        return data[idx] * 2.0  // ← StableHLO MULTIPLY
    }
}
```

### Control Flow Mapping

| Your Code | Compiles To |
|-----------|-------------|
| `if ... else` | StableHLO `stablehlo.case` or `stablehlo.select` |
| `for` loop | StableHLO `stablehlo.while` |
| `while` loop | StableHLO `stablehlo.while` |
| Arithmetic | StableHLO `stablehlo.add`, `stablehlo.multiply`, etc. |
| Array access | StableHLO `stablehlo.dynamic_slice` |

**Not Supported:**
- System calls, recursive functions, runtime memory allocation

**StableHLO Specification:** https://openxla.org/stablehlo/spec

---

## Memory Model: Buffer Inference

**Critical:** StableHLO computations don't allocate memory. They describe operations on buffers.

### What Really Happens

```arc
async func add<tpu>(a: *float32, b: *float32, n: usize) *float32 {
    // NO malloc in StableHLO! Just operations on inputs
    let idx = tpu.replica_id()
    if idx < n {
        return a[idx] + b[idx]  // Compiler infers output is float32[n]
    }
}
```

**Flow:**
1. **Compile time:** Compiler analyzes function, infers output shape `float32[n]`
2. **Compile time:** Generates StableHLO operating on pre-allocated buffers
3. **Runtime (await):** PJRT allocates input buffers on TPU HBM
4. **Runtime:** PJRT copies `a` and `b` to TPU
5. **Runtime:** PJRT allocates output buffer based on shape inference
6. **Runtime:** StableHLO executes on pre-allocated buffers
7. **Runtime:** PJRT copies result back to CPU

### Syntactic Sugar (if provided)

```arc
async func example<tpu>(data: *float32, n: usize) *float32 {
    let result = tpu.malloc<float32>(n)  // ← Compiler magic
    // Compiler converts this to shape annotations in StableHLO
    // Actual allocation happens at runtime by PJRT
}
```

**Behind the scenes:**
- Compiler infers buffer shapes at compile time
- Encodes shape metadata in StableHLO
- PJRT handles actual HBM allocation at runtime
- **No malloc operations in the StableHLO itself**

---

## Arc's Native Linking

Arc links directly to `libtpu.so`:

```
Arc Source (async func<tpu>)
    ↓
Arc Parser → AST
    ↓
StableHLO Backend → Binary MLIR (embedded in executable)
    ↓
Arc Linker → Links to libtpu.so
    ↓
Runtime (PJRT C API) → XLA Compiler → TPU Execution
```

---

## Why PJRT?

**XLA is a compiler (StableHLO → machine code). PJRT is the runtime (device management, execution, memory).**

```
Arc generates StableHLO → PJRT Runtime → XLA Compiler → TPU
                          ↓
                          Handles all the complexity
```

### What PJRT Does For You

If you used XLA directly, you'd manually implement:

1. Device initialization and configuration
2. Stream management for async operations
3. Buffer lifecycle (allocation, deallocation, transfers)
4. Multi-device coordination
5. Error handling and recovery
6. Memory management strategies
7. Shape inference for output buffers

PJRT handles all of this through a stable C API.

### Arc's Approach

```arc
extern pjrt {
    func PJRT_Client_Compile(
        client: *void,
        stablehlo_bytes: *byte,
        size: usize,
        executable: **void
    ) int32
    
    func PJRT_Executable_Execute(
        executable: *void,
        inputs: **void,
        num: usize,
        outputs: **void
    ) int32
    
    func PJRT_Buffer_ToHostBuffer(buffer: *void, dst: *void, size: usize) int32
}


// Your StableHLO as a text string:
let stablehlo_text = """
module {
  func.func @add(%arg0: tensor<f32>, %arg1: tensor<f32>) -> tensor<f32> {
    %0 = stablehlo.add %arg0, %arg1 : tensor<f32>
    return %0 : tensor<f32>
  }
}
"""

// Pass to PJRT:
PJRT_Program program;
program.code = stablehlo_text.ptr;     // ← Pointer to text string
program.code_size = stablehlo_text.len; // ← Length of string  
program.format = "mlir";                // ← Format identifier
program.format_size = 4;                // ← Length of "mlir"

PJRT_Client_Compile_Args compile_args;
compile_args.program = &program;
// ... set other args ...

pjrt_api.PJRT_Client_Compile(&compile_args);

```

**PJRT provides XLA compilation + runtime in a clean C API.**

---

## Two-Phase Compilation

### Build Time (AOT)

```bash
arc build --target=tpu my_program.arc
```

1. Compile `async func<tpu>` → StableHLO operations
2. Perform shape inference on StableHLO
3. Serialize StableHLO → binary MLIR bytecode
4. Embed StableHLO bytes in executable
5. Link to `libtpu.so`

### Runtime Execution (JIT)

```arc
let result = await matrix_multiply(A, B, N)  // First call
```

1. Load embedded StableHLO bytes from executable
2. `PJRT_Client_Compile(stablehlo_bytes)` → executable handle
3. **XLA JIT compiles StableHLO → TPU machine code**
4. PJRT caches compiled TPU code
5. PJRT allocates input buffers on TPU HBM
6. PJRT copies CPU data → TPU buffers
7. `PJRT_Executable_Execute()`
8. PJRT infers output shape, allocates output buffer
9. StableHLO executes on TPU
10. `await` synchronizes on completion
11. PJRT copies result TPU → CPU

**Subsequent calls reuse cached TPU executable**

---

## Complete Example

```arc
async func matrix_multiply<tpu>(A: *float32, B: *float32, N: usize) *float32 {
    // Compiler infers output shape: float32[N * N]
    let total = N * N
    let idx = tpu.replica_id()
    
    if idx < total {
        let row = idx / N
        let col = idx % N
        let mut sum: float32 = 0.0
        
        for k in 0..N {  // ← Becomes StableHLO while or optimized to dot
            sum += A[row * N + k] * B[k * N + col]
        }
        
        return sum  // Writes to output[idx]
    }
}

func main() {
    const N: usize = 1024
    
    let A = alloca<float32>(N * N)
    let B = alloca<float32>(N * N)
    
    // Initialize on CPU
    for i in 0..(N * N) {
        A[i] = 1.0
        B[i] = 2.0
    }
    
    // Execute on TPU - await handles all transfers
    let C = await matrix_multiply(A, B, N)
    
    io.printf("C[0] = %f\n", C[0])
}
```

---

## Pattern Matching

**Pattern matching can eliminate control flow and generate optimized StableHLO.**

### Best Case: Pure Array Operations

```arc
async func add<tpu>(a: [float32; N], b: [float32; N]) [float32; N] {
    return a + b  // Pure array operation, no indexing
}
```

**Compiler generates:** `stablehlo.add(parameter(0), parameter(1))` - Single operation!

### Manual Indexing

```arc
async func add_indexed<tpu>(a: *float32, b: *float32, n: usize) *float32 {
    let idx = tpu.replica_id()
    if idx < n {
        return a[idx] + b[idx]
    }
}
```

**Compiler generates:**
```
idx = stablehlo.replica_id()
bound_check = stablehlo.compare(idx, n, LT)
a_slice = stablehlo.dynamic_slice(a, idx)
b_slice = stablehlo.dynamic_slice(b, idx)
sum = stablehlo.add(a_slice, b_slice)
result = stablehlo.select(bound_check, sum, zero)
```

**Takeaway:** Write high-level array operations when possible. Manual indexing works but generates more StableHLO ops.

---

## SPMD Parallelism

```arc
async func parallel_work<tpu>(data: *float32, n: usize) *float32 {
    let core_id = tpu.replica_id()  // Which TPU core am I?
    // Each core runs same StableHLO but knows its index
    // Core 0 processes data[0], Core 1 processes data[1], etc.
}
```

**Maps to:** StableHLO's `stablehlo.replica_id()` operation

**Characteristics:**
- Coarse-grained: one index per TPU core
- Typically 8 cores per TPU v3 chip
- Not fine-grained like GPU threads

---

## Multi-Device Execution

```arc
func multi_tpu() {
    let num_tpus = tpu.device_count()
    
    let result0 = await(0) train_batch(data0, model, batch_size)
    let result1 = await(1) train_batch(data1, model, batch_size)
    let result2 = await(2) train_batch(data2, model, batch_size)
    
    let combined = average(result0, result1, result2)
}
```

---

## Runtime Implementation

```arc
namespace tpu

let tpu_client: *void = null
let initialized: bool = false

func init() {
    if initialized { return }
    
    let api = pjrt.GetPjrtApi()  // Load from libtpu.so
    pjrt.PJRT_Client_Create(api, &tpu_client, "tpu")
    
    initialized = true
}

func replica_id() usize {
    // Maps to StableHLO replica_id() operation
}
```

---

## Key Differences: TPU vs GPU

| Aspect | GPU/CUDA | TPU |
|--------|----------|-----|
| Backend Format | PTX (text assembly) | **StableHLO (binary MLIR)** |
| Abstraction | Low-level threads | High-level array ops |
| Control Flow | PTX branches | StableHLO conditional/select/while |
| Compilation | AOT (Arc→PTX), JIT (PTX→SASS) | AOT (Arc→StableHLO), JIT (StableHLO→TPU code) |
| Runtime API | CUDA Driver API | **PJRT C API** |
| Compiler | NVIDIA Driver JIT | **XLA (via PJRT)** |
| Parallelism | `gpu.thread_id()` (fine-grained) | `tpu.replica_id()` (coarse SPMD) |
| Memory Model | Explicit malloc/copy | **Buffer shapes inferred at compile time** |
| Threads | Thousands per kernel | Typically 8 cores (coarse parallelism) |

---

## Build and Run

```bash
# Build
arc build --target=tpu my_program.arc

# Run
./my_program
# First call: PJRT JIT compiles StableHLO → TPU code
# Subsequent calls: Reuse cached TPU executable
```

**Requirements:**
- TPU device (or Cloud TPU)
- PJRT TPU plugin (`libtpu.so`)

**NOT Required:**
- C++ compiler
- Protocol buffer compiler (`protoc`)
- Bazel
- TensorFlow/JAX
- XLA C++ headers

---

## StableHLO builder and codegen examples in golang 

- https://github.com/gomlx/go-xla

## Summary

1. `async func<tpu>` compiles entirely to **StableHLO** (see spec: https://openxla.org/stablehlo/spec)
2. builder api c++, https://openxla.org/stablehlo/generated/StablehloBuilder
2. Memory is pre-allocated - StableHLO operates on buffers
3. Two-phase: Arc→StableHLO (AOT), StableHLO→TPU code (JIT on first call)
4. `await` triggers JIT compilation, execution, and synchronization
5. Pattern matching optimizes to efficient array operations
6. PJRT manages devices, execution, memory, buffer allocation
7. XLA compiles StableHLO → TPU machine code
8. `tpu.replica_id()` enables SPMD parallelism across TPU cores
9. Shape inference determines buffer sizes at compile time

```arc
// You write:
async func add<tpu>(a: *float32, b: *float32, n: usize) *float32 {
    let idx = tpu.replica_id()
    if idx < n {
        return a[idx] + b[idx]
    }
}

// Compiler generates StableHLO (AOT):
// - Shape inference: output is tensor<?xf32> with dynamic dimension n
// - Operations: replica_id, compare, dynamic_slice, add, select

// First await: PJRT+XLA JIT compile StableHLO → TPU machine code
// PJRT allocates buffers, executes on TPU, returns result
let result = await add(data_a, data_b, 1024)
```

**PJRT provides the runtime. XLA provides the compiler. StableHLO provides the portable intermediate representation. Together they give you TPU execution through a clean C API.** 🚀