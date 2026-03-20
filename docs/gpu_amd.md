# Arc AMD GPU Model

**Key Concept:** Functions marked `async func<gpu>` compile **entirely** to AMDGCN assembly or LLVM IR. Everything - control flow, loops, operations - becomes GPU instructions that run on AMD Radeon and Instinct GPUs.

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
| **GPU** | `pkg/amdgcngen` | AMDGCN Assembly or LLVM IR | High-performance parallel kernels |

```
Arc Source (async func<gpu> + await)
    ↓
Arc Parser → AST (unified for CPU + GPU)
    ↓
IR Generation → SSA IR with Target tags
    ↓
    ├─→ pkg/codegen → AMD64/ARM64 machine code
    └─→ pkg/amdgcngen → LLVM IR or AMDGCN Assembly
    ↓
hipcc/clang → Code Object (.co) with metadata
    ↓
Arc Linker → Embeds .co in executable, links to libamdhip64.so
    ↓
Runtime (HIP API) → ROCr Runtime → GPU Execution
```

---

## Two Compilation Strategies

Arc supports **two approaches** for AMD GPU code generation:

### Strategy 1: LLVM IR Path (Recommended)

```
Arc IR → LLVM IR (.ll) → hipcc → Code Object (.co)
```

**Advantages:**
- Leverage ROCm's optimizer
- Future-proof (architecture-independent)
- Easier debugging with LLVM tools

**Example Flow:**
```bash
# Arc generates LLVM IR
arc build --target=amd-gpu --gpu-backend=llvm kernel.arc
# Produces: kernel.ll

# hipcc compiles to code object
hipcc --genco --offload-arch=gfx90a kernel.ll -o kernel.co

# Arc embeds kernel.co in executable
```

### Strategy 2: Direct AMDGCN Assembly (Advanced)

```
Arc IR → AMDGCN Assembly (.s) → clang assembler → Code Object (.co)
```

**Advantages:**
- Maximum control over instruction selection
- Inline assembly support
- Performance tuning for specific hardware

**Example Flow:**
```bash
# Arc generates AMDGCN assembly
arc build --target=amd-gpu --gpu-backend=asm kernel.arc
# Produces: kernel.s with AMDGCN instructions

# Assemble and link
clang -x assembler -target amdgcn-amd-amdhsa -mcpu=gfx90a \
      -c kernel.s -o kernel.o
ld.lld -shared kernel.o -o kernel.co

# Arc embeds kernel.co in executable
```

---

## Control Flow Mapping

The `pkg/amdgcngen` backend maps Arc control flow directly to AMDGCN instructions.

| Arc Source | AMDGCN Assembly Output |
|-----------|------------------------|
| `if (cond)` | `v_cmp_gt_f32 vcc, v0, v1` + `s_cbranch_vccz LABEL_ELSE` |
| `for` loop | `LABEL_LOOP:` + body + `s_add_u32` + `s_branch LABEL_LOOP` |
| `while` loop | `LABEL_LOOP:` + condition + `s_cbranch_scc0 EXIT` + body |
| `break` | `s_branch LABEL_EXIT` |
| `continue` | `s_branch LABEL_LOOP` |
| `a + b` | `v_add_f32 v0, v1, v2` |
| `a * b` | `v_mul_f32 v0, v1, v2` |
| `arr[idx]` | `global_load_dword v0, v[1:2], off` |
| `arr[idx] = val` | `global_store_dword v[0:1], v2, off` |

### Example: if-else (LLVM IR Path)

```arc
async func example<gpu>(x: float32, threshold: float32) float32 {
    if x > threshold {
        return x * 2.0
    } else {
        return x * 0.5
    }
}
```

**Generated LLVM IR:**
```llvm
define amdgpu_kernel void @example(float %x, float %threshold, float* %result) {
entry:
  %cmp = fcmp ogt float %x, %threshold
  br i1 %cmp, label %then, label %else

then:
  %mul1 = fmul float %x, 2.0
  store float %mul1, float* %result
  br label %exit

else:
  %mul2 = fmul float %x, 0.5
  store float %mul2, float* %result
  br label %exit

exit:
  ret void
}
```

### Example: if-else (Direct Assembly)

**Generated AMDGCN Assembly:**
```asm
.text
.amdgcn_target "amdgcn-amd-amdhsa--gfx90a"
.globl example
.type example,@function

example:
    ; Compare x > threshold
    v_cmp_gt_f32 vcc, v0, v1
    s_cbranch_vccz .L_else
    
    ; Then branch: x * 2.0
    v_mul_f32 v2, v0, 2.0
    s_branch .L_exit
    
.L_else:
    ; Else branch: x * 0.5
    v_mul_f32 v2, v0, 0.5
    
.L_exit:
    ; Store result
    global_store_dword v[3:4], v2, off
    s_endpgm
```

---

## Type System Extension

Arc's type system is extended for AMD GPU hardware (GCN, CDNA, RDNA).

### Standard Types (CPU + GPU)

| Arc Type | AMDGCN Type | LLVM IR |
|----------|-------------|---------|
| `bool` | Predicate/VCC | `i1` |
| `int8` / `uint8` | 8-bit VGPR | `i8` |
| `int16` / `uint16` | 16-bit VGPR | `i16` |
| `int32` / `uint32` | 32-bit VGPR | `i32` |
| `int64` / `uint64` | 64-bit VGPR pair | `i64` |
| `float32` | 32-bit VGPR | `float` |
| `float64` | 64-bit VGPR pair | `double` |

### GPU-Specific Types

| Arc Type | AMDGCN Mapping | Description |
|----------|----------------|-------------|
| `float16` | 16-bit VGPR | IEEE 754 Half Precision |
| `bfloat16` | 16-bit VGPR | Brain Floating Point (CDNA) |
| `float8_e4m3` | 8-bit VGPR (storage) | FP8 (MI300 Training) |
| `float8_e5m2` | 8-bit VGPR (storage) | FP8 (MI300 Inference) |
| `vector2<T>` | `.v2` suffix | 2-element tuple (coalesced loads) |
| `vector4<T>` | `.v4` suffix | 4-element tuple (coalesced loads) |

### Memory Space & Opaque Types

| Arc Type | AMDGCN Mapping | Usage |
|----------|----------------|-------|
| `shared<T, N>` | LDS allocation | Local Data Share (on-chip memory) |
| `gpu.Barrier` | `s_barrier` | Workgroup synchronization |
| `gpu.Fragment<T,M,N,K>` | AGPRs (CDNA) | Matrix Core Fragment (MFMA) |

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

- Calls `hipMallocManaged()` - accessible from both CPU and GPU
- ROCm runtime automatically handles transfers
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
    
    hip.hipMalloc(&gpu_data, 1024 * sizeof<float32>)
    hip.hipMalloc(&gpu_result, 1024 * sizeof<float32>)
    
    hip.hipMemcpy(gpu_data, cpu_data, 1024 * sizeof<float32>, hipMemcpyHostToDevice)
    await process_explicit(gpu_data, gpu_result, 1024)
    hip.hipMemcpy(cpu_data, gpu_result, 1024 * sizeof<float32>, hipMemcpyDeviceToHost)
    
    hip.hipFree(gpu_data)
    hip.hipFree(gpu_result)
}
```

### Local Data Share / Shared Memory (On-Chip Cache)

```arc
async func matrix_tile<gpu>(A: *float32, B: *float32, C: *float32, N: usize) {
    // Allocate LDS (64KB per CU on CDNA)
    let tile_A: shared<float32, 256>
    let tile_B: shared<float32, 256>
    
    let tid = gpu.thread_id()
    
    // Cooperative load into LDS
    tile_A[tid] = A[/* ... */]
    tile_B[tid] = B[/* ... */]
    
    gpu.sync_threads()  // s_barrier - wait for all threads
    
    // Compute using fast LDS
    let mut sum: float32 = 0.0
    for i in 0..256 {
        sum += tile_A[i] * tile_B[i]
    }
    
    C[tid] = sum
}
```

**AMDGCN Assembly for LDS:**
```asm
; Allocate LDS
.amd_kernel_code_t
    workgroup_group_segment_byte_size = 2048  ; 256 floats * 4 bytes * 2 arrays
.end_amd_kernel_code_t

; Store to LDS
ds_write_b32 v0, v1         ; tile_A[tid] = value

; Load from LDS  
ds_read_b32 v2, v0          ; value = tile_A[i]

; Synchronize
s_barrier                   ; All threads in workgroup wait
```

---

## GPU Intrinsics: `extern gpu`

The `extern gpu` interface exposes AMD hardware features. Defined in `lib/gpu/intrinsics.arc`:

### 1. Indexing & Dimensions (Work-Item/Work-Group Model)

```arc
extern gpu {
    // Thread/Work-item IDs
    func thread_id() int32          // Flat work-item ID in work-group
    func thread_id_x() int32        // X dimension
    func thread_id_y() int32        // Y dimension
    func thread_id_z() int32        // Z dimension
    
    // Work-group IDs
    func block_id() int32           // Flat work-group ID
    func block_id_x() int32         // X dimension
    func block_id_y() int32         // Y dimension
    func block_id_z() int32         // Z dimension
    
    // Dimensions
    func block_dim() int32          // Work-group size
    func block_dim_x() int32
    func block_dim_y() int32
    func block_dim_z() int32
    func grid_dim() int32           // Number of work-groups
    
    // Wave/Wavefront IDs (AMD-specific)
    func lane_id() int32            // 0-63 within wavefront
    func wave_id() int32            // Wavefront ID within work-group
}
```

**AMDGCN Assembly Mapping:**
```asm
; thread_id_x()
v_mov_b32 v0, v0            ; Work-item ID is in v0 by default

; block_id_x()
s_mov_b32 s0, s0            ; Work-group ID in s0 (kernarg)

; lane_id()
v_mbcnt_lo_u32_b32 v1, -1, 0
v_mbcnt_hi_u32_b32 v1, -1, v1

; block_dim()
; Passed as kernel argument in SGPRs
```

**Usage:**
```arc
async func kernel<gpu>(data: *float32, n: usize) {
    let gid_x = gpu.block_id_x() * gpu.block_dim_x() + gpu.thread_id_x()
    let lane = gpu.lane_id()  // 0-63 within wavefront
    
    if gid_x < n {
        data[gid_x] *= 2.0
    }
}
```

### 2. Math Intrinsics

```arc
extern gpu {
    // Trigonometry
    func sin(float32) float32       // v_sin_f32
    func cos(float32) float32       // v_cos_f32
    func tan(float32) float32
    
    // Power & Roots
    func pow(float32, float32) float32
    func sqrt(float32) float32      // v_sqrt_f32
    func rsqrt(float32) float32     // v_rsq_f32 - fast 1/sqrt(x)
    func rcp(float32) float32       // v_rcp_f32 - reciprocal
    
    // Arithmetic
    func abs(float32) float32       // v_abs_f32
    func abs(int32) int32           // v_abs_i32
    func fma(float32, float32, float32) float32  // v_fma_f32 or v_mac_f32
    func saturate(float32) float32  // v_med3_f32 for clamp [0,1]
    
    // Logarithms
    func log2(float32) float32      // v_log_f32
    func exp2(float32) float32      // v_exp_f32
}
```

**AMDGCN Assembly Examples:**
```asm
; fma(a, b, c) = a*b + c
v_fma_f32 v0, v1, v2, v3

; rsqrt(x)
v_rsq_f32 v0, v1

; abs(x)
v_abs_f32 v0, v1
```

### 3. Bit Manipulation

```arc
extern gpu {
    func popc(uint32) int32         // v_bcnt_u32_b32 - count set bits
    func clz(uint32) int32          // v_ffbh_u32 - count leading zeros
    func ffs(uint32) int32          // v_ffbl_b32 - find first set bit
    func brev(uint32) uint32        // v_bfrev_b32 - bit reverse
    func bitfield_extract(uint32, uint32, uint32) uint32  // v_bfe_u32
}
```

### 4. Wavefront Control & Synchronization

```arc
extern gpu {
    // Synchronization
    func sync_threads()             // s_barrier - work-group barrier
    func sync_wave()                // s_waitcnt - wavefront barrier
    
    // Voting (Across Wavefront - 64 threads)
    func all(bool) bool             // All lanes true?
    func any(bool) bool             // Any lane true?
    func ballot(bool) uint64        // 64-bit mask of true lanes
    
    // Shuffling (Data Exchange within Wavefront)
    func shuffle<T>(T, int32) T             // ds_bpermute_b32 - get from lane N
    func shuffle_up<T>(T, int32) T          // DPP move
    func shuffle_down<T>(T, int32) T        // DPP move
    func shuffle_xor<T>(T, int32) T         // ds_permute_b32 - butterfly
}
```

**AMDGCN Assembly for Shuffles:**
```asm
; shuffle_xor(value, 16) - butterfly pattern
; Uses ds_bpermute for cross-lane exchange
v_lshlrev_b32 v1, 2, v_lane_id      ; lane_id * 4
v_xor_b32 v1, v1, 64                ; XOR with 16*4
ds_bpermute_b32 v0, v1, v0          ; Read from permuted lane
```

**Example: Wavefront Reduction (64 threads)**
```arc
async func wave_sum<gpu>(value: float32) float32 {
    let mut sum = value
    
    // Butterfly reduction across 64-thread wavefront
    sum += gpu.shuffle_xor(sum, 32)
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
    func atomic_cas<T>(*T, T, T) T      // Compare-and-swap
}
```

**AMDGCN Assembly:**
```asm
; Global memory atomic add
global_atomic_add v0, v[1:2], v3, off

; LDS atomic add
ds_add_u32 v0, v1

; Global atomic CAS
global_atomic_cmpswap v0, v[1:2], v[3:4], off
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

### 6. Memory Fences & Async Operations

```arc
extern gpu {
    // Memory Fences
    func thread_fence()             // s_waitcnt vmcnt(0) - global memory
    func thread_fence_block()       // s_waitcnt lgkmcnt(0) - LDS
    func thread_fence_system()      // System-wide fence
    
    // Load with Cache Hints
    func load_cached<T>(*T) T       // global_load with glc flag
    func load_streaming<T>(*T) T    // global_load with slc flag
    func load_volatile<T>(*T) T     // global_load with volatile
}
```

**AMDGCN Assembly:**
```asm
; Wait for all vector memory ops
s_waitcnt vmcnt(0)

; Wait for LDS operations
s_waitcnt lgkmcnt(0)

; Cached load
global_load_dword v0, v[1:2], off glc

; Streaming load (evict from cache)
global_load_dword v0, v[1:2], off slc
```

### 7. Matrix Core Operations (CDNA - MI100/MI200/MI300)

```arc
extern gpu {
    // Matrix Fused Multiply-Add (MFMA)
    func mfma_f32_32x32x1(gpu.Fragment<float32, 32, 32, 1>, 
                          float32, float32) gpu.Fragment<float32, 32, 32, 1>
    
    func mfma_f16_32x32x4(gpu.Fragment<float16, 32, 32, 4>,
                          vector4<float16>, vector4<float16>) gpu.Fragment<float32, 32, 32, 4>
    
    func mfma_bf16_32x32x4(gpu.Fragment<bfloat16, 32, 32, 4>,
                           vector4<bfloat16>, vector4<bfloat16>) gpu.Fragment<float32, 32, 32, 4>
    
    // FP8 (MI300 only)
    func mfma_f8_32x32x16(gpu.Fragment<float8_e4m3, 32, 32, 16>,
                          vector4<float8_e4m3>, vector4<float8_e4m3>) gpu.Fragment<float32, 32, 32, 16>
}
```

**AMDGCN Assembly (CDNA):**
```asm
; 32x32x1 FP32 matrix multiply using AGPRs
v_mfma_f32_32x32x1f32 a[0:31], v0, v1, a[0:31]

; 32x32x4 FP16 matrix multiply
v_mfma_f16_32x32x4f16 a[0:15], v[0:1], v[2:3], a[0:15]

; Move from AGPR to VGPR for output
v_accvgpr_read_b32 v0, a0
```

---

## Complete Examples

### Example 1: Vector Addition

```arc
async func vector_add<gpu>(a: *float32, b: *float32, c: *float32, n: usize) {
    let gid = gpu.block_id() * gpu.block_dim() + gpu.thread_id()
    
    if gid < n {
        c[gid] = a[gid] + b[gid]
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
    
    // Execute on GPU (256 threads per work-group)
    await vector_add<<<N/256, 256>>>(a, b, c, N)
    
    io.printf("c[0] = %f\n", c[0])
    
    gpu.unified_free(a)
    gpu.unified_free(b)
    gpu.unified_free(c)
}
```

**Generated LLVM IR:**
```llvm
define amdgpu_kernel void @vector_add(float addrspace(1)* %a, 
                                      float addrspace(1)* %b,
                                      float addrspace(1)* %c,
                                      i64 %n) {
entry:
  %bid = call i32 @llvm.amdgcn.workgroup.id.x()
  %tid = call i32 @llvm.amdgcn.workitem.id.x()
  %bdim = call i32 @llvm.amdgcn.dispatch.ptr().ngroups.x
  
  %bid64 = zext i32 %bid to i64
  %tid64 = zext i32 %tid to i64
  %bdim64 = zext i32 %bdim to i64
  
  %offset = mul i64 %bid64, %bdim64
  %gid = add i64 %offset, %tid64
  
  %cmp = icmp ult i64 %gid, %n
  br i1 %cmp, label %compute, label %exit

compute:
  %a_ptr = getelementptr float, float addrspace(1)* %a, i64 %gid
  %b_ptr = getelementptr float, float addrspace(1)* %b, i64 %gid
  %c_ptr = getelementptr float, float addrspace(1)* %c, i64 %gid
  
  %a_val = load float, float addrspace(1)* %a_ptr
  %b_val = load float, float addrspace(1)* %b_ptr
  %sum = fadd float %a_val, %b_val
  store float %sum, float addrspace(1)* %c_ptr
  br label %exit

exit:
  ret void
}
```

**Or Direct AMDGCN Assembly:**
```asm
.text
.amdgcn_target "amdgcn-amd-amdhsa--gfx90a"
.globl vector_add
.type vector_add,@function
.amdgpu_hsa_kernel vector_add

vector_add:
    .amd_kernel_code_t
        enable_sgpr_kernarg_segment_ptr = 1
        kernarg_segment_byte_size = 32
        workitem_vgpr_count = 8
        wavefront_sgpr_count = 16
    .end_amd_kernel_code_t

    ; Load kernel arguments (pointers + n)
    s_load_dwordx8 s[0:7], s[4:5], 0x0
    s_waitcnt lgkmcnt(0)
    
    ; Compute global ID
    s_mul_i32 s8, s6, s7           ; block_id * block_dim
    v_add_u32 v0, s8, v0           ; + thread_id
    
    ; Bounds check
    v_cmp_lt_u32 vcc, v0, s7       ; gid < n
    s_cbranch_vccz .L_exit
    
    ; Calculate addresses
    v_lshlrev_b32 v1, 2, v0        ; gid * 4 (float32)
    v_add_u32 v2, s0, v1           ; &a[gid]
    v_add_u32 v3, s2, v1           ; &b[gid]
    v_add_u32 v4, s4, v1           ; &c[gid]
    
    ; Load a[gid] and b[gid]
    global_load_dword v5, v[2:3], off
    global_load_dword v6, v[3:4], off
    s_waitcnt vmcnt(0)
    
    ; Compute sum
    v_add_f32 v7, v5, v6
    
    ; Store c[gid]
    global_store_dword v[4:5], v7, off

.L_exit:
    s_endpgm
```

### Example 2: Matrix Multiplication with LDS

```arc
const TILE_SIZE: usize = 16

async func matmul_tiled<gpu>(A: *float32, B: *float32, C: *float32, N: usize) {
    let tile_A: shared<float32, 256>  // 16x16 tile in LDS
    let tile_B: shared<float32, 256>
    
    let tx = gpu.thread_id_x()
    let ty = gpu.thread_id_y()
    let bx = gpu.block_id_x()
    let by = gpu.block_id_y()
    
    let row = by * TILE_SIZE + ty
    let col = bx * TILE_SIZE + tx
    
    let mut sum: float32 = 0.0
    
    // Tile over K dimension
    for t in 0..(N / TILE_SIZE) {
        // Cooperative load into LDS
        tile_A[ty * TILE_SIZE + tx] = A[row * N + (t * TILE_SIZE + tx)]
        tile_B[ty * TILE_SIZE + tx] = B[(t * TILE_SIZE + ty) * N + col]
        
        gpu.sync_threads()  // s_barrier
        
        // Compute partial sum using LDS
        for k in 0..TILE_SIZE {
            sum += tile_A[ty * TILE_SIZE + k] * tile_B[k * TILE_SIZE + tx]
        }
        
        gpu.sync_threads()
    }
    
    C[row * N + col] = sum
}
```

**Key AMDGCN Instructions:**
```asm
; Allocate LDS
.amd_kernel_code_t
    workgroup_group_segment_byte_size = 2048  ; 2 * 256 * 4 bytes
.end_amd_kernel_code_t

; Store to LDS
ds_write_b32 v_offset, v_data

; Synchronize
s_barrier

; Load from LDS
ds_read_b32 v_result, v_offset

; Multiply-add
v_fma_f32 v_sum, v_a, v_b, v_sum
```

### Example 3: Parallel Reduction

```arc
async func reduce_sum<gpu>(input: *float32, output: *float32, n: usize) {
    let sdata: shared<float32, 256>
    
    let tid = gpu.thread_id()
    let gid = gpu.block_id() * gpu.block_dim() + tid
    
    // Load into LDS
    sdata[tid] = if gid < n { input[gid] } else { 0.0 }
    gpu.sync_threads()
    
    // Reduction in LDS
    let mut stride = gpu.block_dim() / 2
    while stride > 0 {
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

## Code Object Format

AMD GPU code is packaged in **HSA Code Objects (.co files)** with ELF structure:

### Code Object Structure

```
┌─────────────────────────────────────┐
│   ELF Header                        │
│   - Machine: EM_AMDGPU              │
│   - ABI: AMDGPU_HSA_V5              │
│   - Architecture: gfx90a, etc.      │
├─────────────────────────────────────┤
│   .text Section                     │
│   - AMDGCN machine code             │
├─────────────────────────────────────┤
│   .rodata Section                   │
│   - Constant data                   │
├─────────────────────────────────────┤
│   .note Section                     │
│   - Metadata (MessagePack)          │
│     * Kernel entry points           │
│     * Register usage (VGPRs/SGPRs)  │
│     * LDS size                      │
│     * Kernel arguments layout       │
├─────────────────────────────────────┤
│   Symbol Table                      │
│   - Kernel function names           │
└─────────────────────────────────────┘
```

### Loading Code Objects at Runtime

```arc
// Arc runtime automatically handles this
hip.hipModuleLoad(&module, "kernel.co")
hip.hipModuleGetFunction(&kernel, module, "vector_add")
hip.hipLaunchKernel(kernel, gridDim, blockDim, args, ...)
```

---

## Two-Phase Compilation

### Build Time (Ahead-of-Time)

```bash
arc build --target=amd-gpu --gpu-arch=gfx90a my_program.arc
```

**What happens:**

1. **Parse Arc source** → AST
2. **Compile `async func<gpu>`:**
   - **Option A (LLVM):** Arc IR → LLVM IR → `hipcc --genco` → Code Object
   - **Option B (ASM):** Arc IR → AMDGCN Assembly → `clang + ld.lld` → Code Object
3. **Embed `.co` file** in executable (binary blob)
4. **Compile regular code** → AMD64/ARM64
5. **Link to libamdhip64.so** (HIP runtime)

**Output:** Standalone executable with embedded code object

### Runtime Execution

```arc
let result = await double_array(data, 1024)  // First call
```

**What happens:**

1. **Load embedded code object** from executable memory
2. **`hipModuleLoadData(co_bytes)`** → module handle
3. **Code object already contains native AMDGCN machine code** (no JIT needed)
4. **`hipModuleGetFunction()`** → get kernel entry point
5. **`hipLaunchKernel()`** → dispatch to GPU
6. **`await`** synchronizes with `hipDeviceSynchronize()`
7. **Return result** pointer

**No JIT compilation** - code objects contain final ISA for target architecture

---

## HIP API Interface

```arc
extern hip {
    // Initialization
    func hipInit(flags: uint32) int32
    func hipGetDeviceCount(count: *int32) int32
    func hipSetDevice(device: int32) int32
    
    // Memory Management
    func hipMalloc(ptr: **void, size: usize) int32
    func hipMallocManaged(ptr: **void, size: usize, flags: uint32) int32
    func hipFree(ptr: *void) int32
    
    func hipMemcpy(dst: *void, src: *void, size: usize, kind: int32) int32
    func hipMemcpyAsync(dst: *void, src: *void, size: usize, 
                        kind: int32, stream: *void) int32
    
    // Module Management
    func hipModuleLoad(module: **void, fname: *byte) int32
    func hipModuleLoadData(module: **void, image: *void) int32
    func hipModuleGetFunction(func: **void, module: *void, name: *byte) int32
    func hipModuleUnload(module: *void) int32
    
    // Kernel Execution
    func hipModuleLaunchKernel(
        func: *void,
        gridDimX: uint32, gridDimY: uint32, gridDimZ: uint32,
        blockDimX: uint32, blockDimY: uint32, blockDimZ: uint32,
        sharedMemBytes: uint32, stream: *void,
        kernelParams: **void, extra: **void
    ) int32
    
    func hipLaunchKernelGGL(
        func: *void, gridDim: dim3, blockDim: dim3,
        args: **void, sharedMem: usize, stream: *void
    ) int32
    
    // Synchronization
    func hipDeviceSynchronize() int32
    func hipStreamSynchronize(stream: *void) int32
    
    // Device Properties
    func hipGetDeviceProperties(props: *hipDeviceProp, device: int32) int32
}

// Helper types
struct dim3 {
    x: uint32
    y: uint32
    z: uint32
}

struct hipDeviceProp {
    name: [256]byte
    totalGlobalMem: usize
    sharedMemPerBlock: usize
    regsPerBlock: int32
    warpSize: int32
    maxThreadsPerBlock: int32
    maxThreadsDim: [3]int32
    maxGridSize: [3]int32
    clockRate: int32
    multiProcessorCount: int32
    // ... more fields
}
```

---

## Runtime Implementation

```arc
namespace gpu

let gpu_devices: array<int32, 8> = []
let device_count: int32 = 0
let initialized: bool = false
let current_device: int32 = 0

func init() {
    if initialized { return }
    
    hip.hipInit(0)
    hip.hipGetDeviceCount(&device_count)
    
    if device_count > 0 {
        hip.hipSetDevice(0)
        current_device = 0
    }
    
    initialized = true
}

func unified_malloc<T>(count: usize) *T {
    init()
    let ptr: *void = null
    hip.hipMallocManaged(&ptr, count * sizeof<T>, 
                         hipMemAttachGlobal)
    return cast<*T>(ptr)
}

func unified_free(ptr: *void) {
    hip.hipFree(ptr)
}

func device_count() int32 {
    init()
    return device_count
}

// Intrinsics are compiler-generated
func thread_id() int32 {
    // Replaced with: v_mov_b32 v0, v0 (work-item ID in v0)
}

func block_id() int32 {
    // Replaced with: s_load_dword s0, s[4:5], 0x0 (from kernarg)
}
```

---

## Architecture-Specific Compilation

Arc supports targeting specific AMD GPU architectures:

```bash
# Target specific GPU
arc build --target=amd-gpu --gpu-arch=gfx90a   # MI200 series
arc build --target=amd-gpu --gpu-arch=gfx940   # MI300 series  
arc build --target=amd-gpu --gpu-arch=gfx1030  # RDNA2 (RX 6000)
arc build --target=amd-gpu --gpu-arch=gfx1100  # RDNA3 (RX 7000)

# Auto-detect from system
arc build --target=amd-gpu --gpu-arch=native

# Multiple architectures (fat binary)
arc build --target=amd-gpu --gpu-arch=gfx90a,gfx940,gfx1100
```

**Architecture Feature Detection:**
```arc
func use_matrix_cores() {
    if gpu.has_mfma() {  // Check for CDNA matrix cores
        await matmul_mfma(A, B, C, N)
    } else {
        await matmul_standard(A, B, C, N)
    }
}
```

---

## Multi-Device Execution

```arc
func multi_gpu() {
    let num_gpus = gpu.device_count()
    
    // Execute on specific devices
    hip.hipSetDevice(0)
    let r0 = await process_chunk(data0, size)
    
    hip.hipSetDevice(1)
    let r1 = await process_chunk(data1, size)
    
    hip.hipSetDevice(2)
    let r2 = await process_chunk(data2, size)
    
    let combined = average(r0, r1, r2)
}
```

---

## Debugging and Profiling

### Inspect Generated Code

```bash
# View LLVM IR
arc build --target=amd-gpu --emit=llvm-ir kernel.arc
cat kernel.ll

# View AMDGCN assembly
arc build --target=amd-gpu --emit=asm kernel.arc
cat kernel.s

# Inspect code object
roc-obj-ls -v kernel.co
roc-obj-extract kernel.co .text
```

### Profile with ROCm Tools

```bash
# Profile kernel execution
rocprof --stats ./my_program

# Trace HIP calls
roctracer ./my_program

# Use ROCm debugger
rocgdb ./my_program
```

---

## Build and Run

```bash
# Build Arc program with AMD GPU support
arc build --target=amd-gpu --gpu-arch=gfx90a my_program.arc

# Arc automatically:
# 1. Compiles CPU code → AMD64/ARM64
# 2. Compiles GPU code → LLVM IR or AMDGCN assembly
# 3. Generates code object (.co) via hipcc
# 4. Embeds .co in executable
# 5. Links to libamdhip64.so
# 6. Produces standalone executable

# Run
./my_program
# Code object is already compiled - loads instantly
# No JIT compilation needed
```

**Requirements:**
- AMD GPU with ROCm support (GCN 3+, CDNA, RDNA)
- ROCm driver and runtime
- libamdhip64.so (HIP runtime library)

**NOT Required:**
- CUDA
- nvcc
- Full ROCm SDK (for running - only needed for compiling)

---

## Comparison: NVIDIA PTX vs AMD Code Objects

| Aspect | NVIDIA (PTX) | AMD (Code Object) |
|--------|--------------|-------------------|
| **Intermediate Format** | PTX assembly text | LLVM IR or AMDGCN assembly |
| **Final Format** | SASS (JIT compiled) | Native AMDGCN ISA (AOT) |
| **Compilation** | Driver JIT at runtime | Ahead-of-time with hipcc |
| **Portability** | One PTX for all GPUs | One .co per architecture |
| **Performance** | JIT overhead on first call | No JIT - instant load |
| **Size** | Compact text | Larger binary per arch |
| **Debugging** | PTX visible | Disassemble with roc-obj |

---

## Summary

1. **`async func<gpu>`** compiles to AMDGCN assembly or LLVM IR
2. **Two backend strategies:**
   - **LLVM IR → hipcc** (recommended, portable)
   - **Direct AMDGCN assembly** (advanced, maximum control)
3. **Control flow maps directly** - if/for/while → branches/predicates
4. **Extended type system** - `float16`, `bfloat16`, `shared<T,N>`
5. **`extern gpu` intrinsics** - expose AMD hardware features
6. **Wavefront model** - 64 threads (vs 32 in NVIDIA warps)
7. **Code objects** - ELF format with metadata, no JIT needed
8. **HIP API** - manages devices, execution, memory
9. **LDS (Local Data Share)** - on-chip shared memory
10. **MFMA instructions** - matrix cores on CDNA (MI100/MI200/MI300)

```arc
// You write:
async func double<gpu>(arr: *float32, n: usize) *float32 {
    let result = gpu.unified_malloc<float32>(n)
    let gid = gpu.block_id() * gpu.block_dim() + gpu.thread_id()
    if gid < n {
        result[gid] = arr[gid] * 2.0
    }
    return result
}

// Arc generates LLVM IR or AMDGCN assembly (AOT)
// hipcc compiles to code object (.co)
// Embedded in executable - no runtime JIT

let result = await double(data, 1024)  // Instant load and execute
```

**HIP provides the runtime. ROCm provides the compiler. Together they give you GPU execution through an open-source stack compatible with NVIDIA CUDA semantics.** 🚀