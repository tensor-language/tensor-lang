# Arc NVIDIA GPU Model

**Key Concept:** Functions marked `async func<gpu>` compile **entirely** to PTX (Parallel Thread Execution). Everything - control flow, loops, operations - becomes PTX instructions.

---

## Compiler Architecture: One Parser, Two Backends

Arc uses a **Unified Frontend** with a **Split Backend**. The same parser handles both CPU and GPU code.

### Compilation Pipeline

#### Pass 1: Frontend (Parser & Type Checker)
- **Input:** Source code (`.arc`)
- Parses standard syntax (`struct`, `func`, `if`, `let`)
- Validates types (including GPU-specific types like `float16`)
- Resolves `gpu.*` calls using `extern gpu` intrinsics
- **Validation:** Functions tagged `<gpu>` can only call other GPU functions or `extern gpu` intrinsics

#### Pass 2: IR Generation
- Generates High-Level IR (SSA) for entire program
- Functions defined as `async func<gpu>` are tagged with `Target: GPU`

#### Pass 3: Backend Dispatch

| Target | Package | Output | Role |
|--------|---------|--------|------|
| **CPU** | `pkg/codegen` | Machine Code (ELF/PE) | Application logic, memory management, kernel launching |
| **GPU** | `pkg/ptxgen` | PTX Assembly (Text) | High-performance parallel kernels |

```
Arc Source (async func<gpu> + await)
    ↓
Arc Parser → AST (unified for CPU + GPU)
    ↓
IR Generation → SSA IR with Target tags
    ↓
    ├─→ pkg/codegen → AMD64/ARM64 machine code
    └─→ pkg/ptxgen → PTX assembly (embedded in executable)
    ↓
Arc Linker → Links to libcuda.so
    ↓
Runtime (CUDA Driver API) → NVIDIA Driver JIT → GPU Execution
```

---

## Control Flow Mapping

The `pkg/ptxgen` backend maps Arc control flow directly to PTX. **No new syntax required.**

| Arc Source | PTX Output |
|-----------|------------|
| `if (cond)` | **Short block:** `setp.lt.u32 %p1, ...` + `@%p1 add.f32 ...` (predication) |
| `if (cond)` | **Long block:** `setp.lt.u32 %p1, ...` + `@!%p1 bra LABEL_ELSE;` (branching) |
| `for` loop | `LABEL_LOOP:` + loop body + `bra LABEL_LOOP;` |
| `while` loop | `LABEL_LOOP:` + condition + `@!%p1 bra EXIT;` + body + `bra LABEL_LOOP;` |
| `break` | `bra LABEL_EXIT;` |
| `continue` | `bra LABEL_LOOP;` |
| `a + b` | `add.f32 %r1, %r2, %r3;` |
| `a * b` | `mul.f32 %r1, %r2, %r3;` |
| `arr[idx]` | `ld.global.f32 %r1, [%ptr];` |
| `arr[idx] = val` | `st.global.f32 [%ptr], %r1;` |

### Example: if-else

```arc
async func example<gpu>(x: float32, threshold: float32) float32 {
    if x > threshold {
        return x * 2.0
    } else {
        return x * 0.5
    }
}
```

**Generated PTX:**
```ptx
setp.gt.f32 %p1, %x, %threshold;
@%p1  mul.f32 %result, %x, 2.0;
@!%p1 mul.f32 %result, %x, 0.5;
```

---

## Type System Extension

Arc's type system is extended for GPU hardware support (Volta, Ampere, Hopper).

### Standard Types (CPU + GPU)

| Arc Type | PTX Type |
|----------|----------|
| `bool` | `.pred` |
| `int8` / `uint8` | `.s8` / `.u8` |
| `int16` / `uint16` | `.s16` / `.u16` |
| `int32` / `uint32` | `.s32` / `.u32` |
| `int64` / `uint64` | `.s64` / `.u64` |
| `float32` | `.f32` |
| `float64` | `.f64` |

### GPU-Specific Types

| Arc Type | PTX Mapping | Description |
|----------|-------------|-------------|
| `float16` | `.f16` | IEEE 754 Half Precision |
| `bfloat16` | `.bf16` | Brain Floating Point (AI/ML) |
| `float8_e4m3` | `.b8` (storage) | FP8 (Hopper Inference) |
| `float8_e5m2` | `.b8` (storage) | FP8 (Hopper Training) |
| `vector2<T>` | `.v2` | 2-element tuple (coalesced loads) |
| `vector4<T>` | `.v4` | 4-element tuple (coalesced loads) |

### Memory Space & Opaque Types

| Arc Type | PTX Mapping | Usage |
|----------|-------------|-------|
| `shared<T, N>` | `.shared` | L1 Shared Memory allocation |
| `gpu.Barrier` | `.b64` | Hardware `mbarrier` object |
| `gpu.TensorMap` | `.u64` | TMA Descriptor Pointer (H100) |
| `gpu.Fragment<T,M,N,K>` | `.reg` array | Tensor Core Matrix Fragment |

---

## Memory Model

### Unified Memory (Recommended)

```arc
async func process<gpu>(arr: *float32, n: usize) *float32 {
    let result = gpu.unified_malloc<float32>(n)  // CPU+GPU accessible
    let idx = gpu.thread_id()
    if idx < n {
        result[idx] = arr[idx] * 2.0
    }
    return result
}
```

- Calls `cuMemAllocManaged()` - accessible from both CPU and GPU
- CUDA driver automatically handles transfers
- Simpler but may have performance overhead

### Explicit Memory (Performance Critical)

```arc
async func process_explicit<gpu>(gpu_arr: *float32, gpu_result: *float32, n: usize) {
    let idx = gpu.thread_id()
    if idx < n {
        gpu_result[idx] = gpu_arr[idx] * 2.0
    }
}

func use_explicit() {
    let cpu_data = alloca<float32>(1024)
    let gpu_data: *float32 = null
    let gpu_result: *float32 = null
    
    cuda.cuMemAlloc(&gpu_data, 1024 * sizeof<float32>)
    cuda.cuMemAlloc(&gpu_result, 1024 * sizeof<float32>)
    
    cuda.cuMemcpyHtoD(gpu_data, cpu_data, 1024 * sizeof<float32>)
    await process_explicit(gpu_data, gpu_result, 1024)
    cuda.cuMemcpyDtoH(cpu_data, gpu_result, 1024 * sizeof<float32>)
    
    cuda.cuMemFree(gpu_data)
    cuda.cuMemFree(gpu_result)
}
```

### Shared Memory (On-Chip Cache)

```arc
async func matrix_tile<gpu>(A: *float32, B: *float32, C: *float32, N: usize) {
    // Allocate shared memory (64KB per SM)
    let tile_A: shared<float32, 256>
    let tile_B: shared<float32, 256>
    
    let tid = gpu.thread_id()
    
    // Cooperative load into shared memory
    tile_A[tid] = A[/* ... */]
    tile_B[tid] = B[/* ... */]
    
    gpu.sync_threads()  // Wait for all threads
    
    // Compute using fast shared memory
    let mut sum: float32 = 0.0
    for i in 0..256 {
        sum += tile_A[i] * tile_B[i]
    }
    
    C[tid] = sum
}
```

---

## GPU Intrinsics: `extern gpu`

The `extern gpu` interface exposes hardware features without new syntax. Defined in `lib/gpu/intrinsics.arc`:

### 1. Indexing & Dimensions

```arc
extern gpu {
    func thread_id() int32          // %tid.x
    func thread_id_y() int32        // %tid.y
    func thread_id_z() int32        // %tid.z
    func block_id() int32           // %ctaid.x
    func block_id_y() int32         // %ctaid.y
    func block_id_z() int32         // %ctaid.z
    func lane_id() int32            // %laneid (0-31 within warp)
    func warp_id() int32            // %warpid
    func grid_dim() int32           // %nctaid.x
    func block_dim() int32          // %ntid.x
}
```

**Usage:**
```arc
async func kernel<gpu>(data: *float32, n: usize) {
    let idx = gpu.thread_id()
    let bid = gpu.block_id()
    let tid_in_block = gpu.thread_id()
    
    let global_idx = bid * gpu.block_dim() + tid_in_block
    
    if global_idx < n {
        data[global_idx] *= 2.0
    }
}
```

### 2. Math Intrinsics

```arc
extern gpu {
    // Trigonometry
    func sin(float32) float32       // sin.approx.f32
    func cos(float32) float32       // cos.approx.f32
    func tan(float32) float32
    
    // Power & Roots
    func pow(float32, float32) float32
    func sqrt(float32) float32
    func rsqrt(float32) float32     // 1.0 / sqrt(x) - fast inverse sqrt
    func rcp(float32) float32       // 1.0 / x - reciprocal
    
    // Arithmetic
    func abs(float32) float32
    func abs(int32) int32
    func fma(float32, float32, float32) float32  // (a*b)+c - fused multiply-add
    func saturate(float32) float32  // Clamp to [0.0, 1.0]
    
    // Logarithms
    func log2(float32) float32
    func exp2(float32) float32
}
```

### 3. Bit Manipulation

```arc
extern gpu {
    func popc(uint32) int32         // Population count (count set bits)
    func clz(uint32) int32          // Count leading zeros
    func ffs(uint32) int32          // Find first set bit
    func brev(uint32) uint32        // Bit reverse
}
```

### 4. Warp Control & Synchronization

```arc
extern gpu {
    // Synchronization
    func sync_threads()             // bar.sync 0 (block-level barrier)
    func sync_warp()                // bar.warp.sync (warp-level barrier)
    func active_mask() uint32       // activemask (which lanes are active)
    
    // Voting (Predicate Logic across Warp)
    func all(bool) bool             // vote.all.pred - all lanes true?
    func any(bool) bool             // vote.any.pred - any lane true?
    func ballot(bool) uint32        // vote.ballot.b32 - bitmask of true lanes
    
    // Shuffling (Register Exchange within Warp)
    func shuffle<T>(T, int32) T             // shfl.sync.idx - get from lane N
    func shuffle_up<T>(T, int32) T          // shfl.sync.up - get from lane-delta
    func shuffle_down<T>(T, int32) T        // shfl.sync.down - get from lane+delta
    func shuffle_xor<T>(T, int32) T         // shfl.sync.bfly - butterfly pattern
}
```

**Example: Warp Reduction**
```arc
async func warp_sum<gpu>(value: float32) float32 {
    let mut sum = value
    
    // Butterfly reduction across 32-thread warp
    sum += gpu.shuffle_xor(sum, 16)
    sum += gpu.shuffle_xor(sum, 8)
    sum += gpu.shuffle_xor(sum, 4)
    sum += gpu.shuffle_xor(sum, 2)
    sum += gpu.shuffle_xor(sum, 1)
    
    return sum  // All lanes have the same sum
}
```

### 5. Atomic Operations

```arc
extern gpu {
    func atomic_add<T>(*T, T) T         // Returns OLD value
    func atomic_sub<T>(*T, T) T
    func atomic_min<T>(*T, T) T
    func atomic_max<T>(*T, T) T
    func atomic_exchange<T>(*T, T) T
    func atomic_cas<T>(*T, T, T) T      // (ptr, compare, val) - compare-and-swap
}
```

**Example: Histogram**
```arc
async func histogram<gpu>(data: *uint8, hist: *uint32, n: usize) {
    let idx = gpu.thread_id()
    
    if idx < n {
        let bin = data[idx]
        gpu.atomic_add(&hist[bin], 1)
    }
}
```

### 6. Memory & Async Pipeline

```arc
extern gpu {
    // Fences
    func thread_fence()             // membar.gl - global memory fence
    func thread_fence_block()       // membar.cta - block-level fence
    func thread_fence_system()      // membar.sys - system-wide fence
    
    // Load Caching Hints
    func load_cached<T>(*T) T       // ld.global.ca - cache in L1+L2
    func load_streaming<T>(*T) T    // ld.global.cs - cache streaming (evict first)
    func load_volatile<T>(*T) T     // ld.volatile - bypass cache
    
    // Async Copy (Global → Shared, Ampere+)
    func async_copy(*void, *void, usize)    // cp.async.ca.shared.global
    func async_commit_group()               // cp.async.commit_group
    func async_wait_group(int32)            // cp.async.wait_group
}
```

### 7. Tensor Memory Accelerator (Hopper H100+)

```arc
extern gpu {
    func tma_copy(gpu.TensorMap, *void, gpu.Barrier)   // cp.async.bulk.tensor
    func barrier_wait(gpu.Barrier)                      // mbarrier.test_wait
}
```

---

## Complete Examples

### Example 1: Vector Addition

```arc
async func vector_add<gpu>(a: *float32, b: *float32, c: *float32, n: usize) {
    let idx = gpu.thread_id()
    
    if idx < n {
        c[idx] = a[idx] + b[idx]
    }
}

func main() {
    const N: usize = 1024
    
    let a = gpu.unified_malloc<float32>(N)
    let b = gpu.unified_malloc<float32>(N)
    let c = gpu.unified_malloc<float32>(N)
    
    // Initialize on CPU
    for i in 0..N {
        a[i] = cast<float32>(i)
        b[i] = cast<float32>(i * 2)
    }
    
    // Execute on GPU
    await vector_add(a, b, c, N)
    
    io.printf("c[0] = %f\n", c[0])
    
    gpu.unified_free(a)
    gpu.unified_free(b)
    gpu.unified_free(c)
}
```

### Example 2: Matrix Multiplication with Shared Memory

```arc
const TILE_SIZE: usize = 16

async func matmul_tiled<gpu>(A: *float32, B: *float32, C: *float32, N: usize) {
    let tile_A: shared<float32, 256>  // 16x16 tile
    let tile_B: shared<float32, 256>
    
    let tx = gpu.thread_id()
    let ty = gpu.thread_id_y()
    let bx = gpu.block_id()
    let by = gpu.block_id_y()
    
    let row = by * TILE_SIZE + ty
    let col = bx * TILE_SIZE + tx
    
    let mut sum: float32 = 0.0
    
    // Tile over K dimension
    for t in 0..(N / TILE_SIZE) {
        // Load tiles into shared memory
        tile_A[ty * TILE_SIZE + tx] = A[row * N + (t * TILE_SIZE + tx)]
        tile_B[ty * TILE_SIZE + tx] = B[(t * TILE_SIZE + ty) * N + col]
        
        gpu.sync_threads()
        
        // Compute partial sum
        for k in 0..TILE_SIZE {
            sum += tile_A[ty * TILE_SIZE + k] * tile_B[k * TILE_SIZE + tx]
        }
        
        gpu.sync_threads()
    }
    
    C[row * N + col] = sum
}
```

### Example 3: Parallel Reduction

```arc
async func reduce_sum<gpu>(input: *float32, output: *float32, n: usize) {
    let sdata: shared<float32, 256>
    
    let tid = gpu.thread_id()
    let idx = gpu.block_id() * gpu.block_dim() + tid
    
    // Load data into shared memory
    sdata[tid] = if idx < n { input[idx] } else { 0.0 }
    gpu.sync_threads()
    
    // Reduction in shared memory
    let mut stride = gpu.block_dim() / 2
    for stride > 0 {
        if tid < stride {
            sdata[tid] += sdata[tid + stride]
        }
        gpu.sync_threads()
        stride /= 2
    }
    
    // Thread 0 writes result
    if tid == 0 {
        gpu.atomic_add(&output[0], sdata[0])
    }
}
```

---

## Two-Phase Compilation

### Build Time (AOT)

```bash
arc build --target=gpu my_program.arc
```

**What happens:**
1. Parse Arc source → AST
2. Compile `async func<gpu>` → PTX assembly (text)
3. **Embed PTX text in executable** (like embedding strings)
4. Compile regular code → AMD64/ARM64
5. Link to `libcuda.so`

**Output:** Standalone executable with embedded PTX

### Runtime Execution (JIT)

```arc
let result = await double_array(data, 1024)  // First call
```

**What happens:**
1. Load embedded PTX string from executable
2. `cuModuleLoadData(ptx_string)` → module handle
3. **NVIDIA Driver JIT compiles PTX → SASS** (GPU machine code)
4. Driver caches compiled SASS for reuse
5. `cuLaunchKernel()` executes on GPU
6. `await` synchronizes with `cuCtxSynchronize()`
7. Return result pointer

**Subsequent calls reuse cached SASS**

---

## CUDA Driver API Interface

```arc
extern cuda {
    func cuInit(flags: uint32) int32
    func cuDeviceGet(device: *int32, ordinal: int32) int32
    func cuCtxCreate(ctx: **void, flags: uint32, device: int32) int32
    func cuModuleLoadData(module: **void, ptx: *byte) int32
    func cuModuleGetFunction(func: **void, module: *void, name: *byte) int32
    func cuLaunchKernel(
        func: *void,
        gridDimX: uint32, gridDimY: uint32, gridDimZ: uint32,
        blockDimX: uint32, blockDimY: uint32, blockDimZ: uint32,
        sharedMemBytes: uint32, stream: *void,
        kernelParams: **void, extra: **void
    ) int32
    func cuCtxSynchronize() int32
    func cuMemAllocManaged(ptr: **void, size: usize, flags: uint32) int32
    func cuMemFree(ptr: *void) int32
}
```

---

## Runtime Implementation

```arc
namespace gpu

let gpu_contexts: array<*void, 8> = []
let device_count: int32 = 0
let initialized: bool = false

func init() {
    if initialized { return }
    
    cuda.cuInit(0)
    cuda.cuDeviceGetCount(&device_count)
    
    for i in 0..device_count {
        let device: int32 = 0
        cuda.cuDeviceGet(&device, i)
        cuda.cuCtxCreate(&gpu_contexts[i], 0, device)
    }
    
    initialized = true
}

func unified_malloc<T>(count: usize) *T {
    init()
    let ptr: *void = null
    cuda.cuMemAllocManaged(&ptr, count * sizeof<T>, 1)
    return cast<*T>(ptr)
}

func unified_free(ptr: *void) {
    cuda.cuMemFree(ptr)
}

func thread_id() int32 {
    // Compiler intrinsic - replaced with PTX: mov.u32 %tid, %tid.x
}
```

---

## Multi-Device Execution

```arc
func multi_gpu() {
    let num_gpus = gpu.device_count()
    
    // Execute on specific devices
    let r0 = await(0) process_chunk(data0, size)
    let r1 = await(1) process_chunk(data1, size)
    let r2 = await(2) process_chunk(data2, size)
    
    let combined = average(r0, r1, r2)
}
```

---

## Build and Run

```bash
# Build
arc build --target=gpu my_program.arc

# Arc automatically:
# 1. Compiles CPU code → AMD64/ARM64
# 2. Compiles GPU code → PTX assembly (AOT)
# 3. Embeds PTX text in executable
# 4. Links to libcuda.so
# 5. Produces standalone executable

# Run
./my_program
# First call: Driver JIT compiles PTX → SASS
# Subsequent calls: Reuse cached SASS
```

**Requirements:**
- NVIDIA GPU with driver
- CUDA Driver library (`libcuda.so` or `nvcuda.dll`)

**NOT Required:**
- CUDA Toolkit
- nvcc compiler
- gcc/g++ toolchain

---

## Summary

1. **`async func<gpu>`** compiles entirely to PTX
2. **One parser, two backends** - same syntax for CPU and GPU
3. **Control flow maps directly** - if/for/while → PTX branches/predicates
4. **Extended type system** - `float16`, `bfloat16`, `shared<T,N>`
5. **`extern gpu` intrinsics** - expose hardware features without new syntax
6. **Unified memory is simplest** - CPU and GPU share pointers
7. **Two-phase compilation** - Arc→PTX (AOT), PTX→SASS (JIT on first call)
8. **CUDA Driver API** - manages devices, execution, memory
9. **`gpu.thread_id()`** - fine-grained SIMT parallelism (thousands of threads)

```arc
// You write:
async func double<gpu>(arr: *float32, n: usize) *float32 {
    let result = gpu.unified_malloc<float32>(n)
    let idx = gpu.thread_id()
    if idx < n {
        result[idx] = arr[idx] * 2.0
    }
    return result
}

// Compiler generates PTX (AOT):
// - Kernel with thread indexing
// - Load, multiply, store operations

// First await: Driver JIT compiles PTX → SASS
let result = await double(data, 1024)
```

**CUDA Driver API provides the runtime. NVIDIA Driver provides the JIT compiler. Together they give you GPU execution through a stable C API.** 🚀