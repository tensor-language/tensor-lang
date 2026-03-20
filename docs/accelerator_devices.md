# Arc Accelerator Targets

Arc supports heterogeneous computing across GPUs, TPUs, QPUs (Quantum), and custom AI accelerators through target specifications. Functions marked with `async func<target>` compile entirely to platform-specific instruction sets, enabling high-performance execution on specialized hardware.

---

## Syntax

```arc
async func function_name<target>(parameters) return_type {
    // Code compiles to target-specific instructions
}
```

**Target specification format:** `<category.framework>` or `<platform.device>`

---

## GPU Targets

### `<gpu.cuda>` - NVIDIA CUDA Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** PTX assembly → **Runtime:** CUDA Driver API → **Final:** SASS (GPU machine code)

**Hardware Support:**
- NVIDIA GPUs: Maxwell, Pascal, Volta, Turing, Ampere, Ada Lovelace, Hopper architectures
- Compute Capability 5.0+

**What You Get:**
- Fine-grained SIMT parallelism via `gpu.thread_id()` (thousands of threads)
- Unified memory or explicit device memory management
- NVIDIA-specific features: Tensor Cores, warp-level primitives, shared memory
- JIT compilation: PTX→SASS on first call, cached thereafter

**Example:**
```arc
async func matmul<gpu.cuda>(A: *float32, B: *float32, C: *float32, N: usize) {
    let total = N * N
    let idx = gpu.thread_id()
    
    if idx < total {
        let row = idx / N
        let col = idx % N
        let mut sum: float32 = 0.0
        
        for k in 0..N {
            sum += A[row * N + k] * B[k * N + col]
        }
        
        C[row * N + col] = sum
    }
}

// First call: Arc loads embedded PTX, CUDA Driver JIT compiles to SASS
await matmul(A, B, C, 1024)
```

**Runtime Dependencies:**
- `libcuda.so` (Linux) or `nvcuda.dll` (Windows) - CUDA Driver
- NVIDIA GPU driver

---

### `<gpu.rocm>` - AMD ROCm Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** HIP/GCN → **Runtime:** ROCm → **Final:** AMD GPU ISA

**Hardware Support:**
- AMD GPUs: GCN 3.0+, RDNA, RDNA 2, RDNA 3, CDNA architectures
- Instinct MI series (MI100, MI200, MI300X)
- Radeon RX series

**What You Get:**
- HIP (Heterogeneous-computing Interface for Portability)
- AMD-specific features: Matrix Cores (MFMA), wavefront intrinsics
- Open-source software stack

**Example:**
```arc
async func vector_add<gpu.rocm>(a: *float32, b: *float32, c: *float32, n: usize) {
    let idx = gpu.thread_id()
    if idx < n {
        c[idx] = a[idx] + b[idx]
    }
}
await vector_add(a, b, c, 1024)
```

**Runtime Dependencies:**
- `libhsa-runtime64.so` - HSA runtime
- ROCm driver stack

---

### `<gpu.metal>` - Apple Metal Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** Metal Shading Language (MSL) → **Runtime:** Metal API → **Final:** Apple GPU ISA

**Hardware Support:**
- Apple Silicon: M1, M2, M3, M4 families
- Tile-Based Deferred Rendering (TBDR) architecture

**What You Get:**
- Unified memory architecture (CPU and GPU share memory)
- Metal Performance Shaders (MPS)
- Integration with macOS/iOS ecosystem

**Example:**
```arc
async func image_filter<gpu.metal>(input: *float32, output: *float32) {
    let idx = gpu.thread_id()
    // Metal-specific: unified memory, no explicit transfers needed
    output[idx] = input[idx] * 1.2
}
```

---

### `<gpu.oneapi>` - Intel oneAPI Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** SYCL (DPC++) → **Runtime:** Level Zero API → **Final:** Intel GPU ISA

**Hardware Support:**
- Intel Data Center GPU Max Series
- Intel Arc GPUs (Alchemist, Battlemage)

**What You Get:**
- SYCL 2020 standard programming model
- Cross-vendor compatibility via oneAPI
- Unified programming model across Intel CPU/GPU/FPGA

---

## TPU Targets

### `<tpu>` - Google TPU

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** HLO (High-Level Operations) → **Runtime:** PJRT C API → **Final:** TPU machine code

**Hardware Support:**
- Google Cloud TPU v2, v3, v4, v5e, v5p

**What You Get:**
- Coarse-grained SPMD parallelism via `tpu.replica_id()`
- Optimized for matrix operations and transformers
- Automatic buffer allocation and shape inference via XLA

**Example:**
```arc
async func matrix_multiply<tpu>(A: *float32, B: *float32, N: usize) *float32 {
    let idx = tpu.replica_id()
    // Compiler infers output buffer shape via HLO
    return matmul_op(A, B) 
}
```

---

## QPU Targets (Quantum)

Arc treats Quantum Processing Units (QPUs) as asynchronous accelerators. Code is compiled to **QIR (Quantum Intermediate Representation)**, an LLVM-based standard, enabling execution on various quantum backends.

### `<qpu.ibm>` - IBM Quantum Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** QIR (LLVM) → **Runtime:** Qiskit Runtime → **Final:** Microwave Pulses (Transmons)

**Hardware Support:**
- IBM Eagle, Osprey, Heron, Condor processors (Superconducting)

**What You Get:**
- Native gate access (`h`, `cx`, `rz`, `sx`)
- Access to IBM Quantum Cloud execution queues
- Error mitigation primitives

**Example:**
```arc
// Create a Bell State (Entanglement)
async func bell_state<qpu.ibm>(shots: usize) Distribution {
    // 1. Allocate Qubits (Opaque Handles)
    let q = qpu.alloc(2)
    let c = qpu.alloc_bit(2)
    
    // 2. Quantum Operations
    qpu.h(q[0])          // Hadamard gate
    qpu.cx(q[0], q[1])   // CNOT gate (entangle)
    
    // 3. Measurement
    c[0] = qpu.measure(q[0])
    c[1] = qpu.measure(q[1])
    
    // 4. Return results (waits for queue)
    return qpu.result(c) 
}

let dist = await bell_state(1024)
// Result: { "00": 0.51, "11": 0.49 }
```

### `<qpu.ionq>` - IonQ Trapped Ion Platform

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** QIR (LLVM) → **Runtime:** Azure Quantum / AWS Braket

**Hardware Support:**
- IonQ Harmony, Aria, Forte

**What You Get:**
- **All-to-All Connectivity:** Any qubit can interact with any other (no SWAP overhead)
- High gate fidelity, longer coherence times
- Ideal for optimization problems and deep circuits

### `<qpu.sim>` - Local Simulator

**Behavior:**
- Runs locally on CPU or GPU using tensor network simulation.
- Used for debugging logic before incurring cloud costs.
- Supports "perfect" (noiseless) simulation.

**Example:**
```arc
async func debug_circuit<qpu.sim>() {
    let q = qpu.alloc(1)
    qpu.rx(q[0], 3.14)
    
    // Simulator-only debug feature
    print(qpu.debug_state_vector()) 
}
```

---

## AWS Silicon Targets

### `<aws.trainium>` & `<aws.inferentia>`

**Compilation Pipeline:**
- **Source:** Arc → **Intermediate:** Neuron IR → **Runtime:** Neuron SDK → **Final:** AWS ISA

**Hardware Support:**
- AWS Trainium (Trn1/Trn2) for training
- AWS Inferentia (Inf1/Inf2) for inference

**What You Get:**
- Purpose-built silicon for deep learning
- Integration with AWS Neuron SDK
- Significant cost/performance benefits over generic GPUs on AWS

---

## Generic Targets

### `<gpu>` - Generic GPU
Portable subset of GPU features. Compiles to CUDA, HIP, or Metal depending on the host.

### `<qpu>` - Generic Quantum
Portable QIR "Base Profile" code. Compiles to the available quantum backend (IBM, IonQ, or Simulator).

---

## Comparison Table

| Target | Backend Format | Parallelism Model | Memory Model | Key Feature | Market |
|--------|---------------|-------------------|--------------|-------------|--------|
| `<gpu.cuda>` | PTX Assembly | SIMT (Threads) | Unified/Explicit | Mature Ecosystem | 95% AI/HPC |
| `<tpu>` | HLO / XLA | SPMD (Cores) | Auto-Managed | Matrix Scaling | Google Cloud |
| `<qpu.ibm>` | QIR / QASM | Probabilistic (Shots) | Opaque Handles | Error Mitigation | Quantum Research |
| `<qpu.ionq>` | QIR / Quil | Probabilistic (Shots) | Opaque Handles | All-to-All Connect | Quantum Optimization |
| `<aws.trainium>`| Neuron IR | Custom Cores | SDK-Managed | Cost Efficiency | AWS Training |

---

## Hybrid Execution (The "Killer Feature")

Arc allows mixing targets in a single workflow. This is essential for variational quantum algorithms (VQE/QAOA) or AI pipelines.

```arc
func main() {
    let mut params = [0.1, 0.5, 0.2]
    
    // Hybrid Quantum-Classical Loop (VQE)
    for i in 0..100 {
        // 1. Run Quantum Circuit (QPU)
        let energy = await quantum_ansatz<qpu.ibm>(params)
        
        // 2. Calculate Gradient (GPU)
        let grad = await compute_gradient<gpu.cuda>(energy, params)
        
        // 3. Update (CPU)
        params = update_weights(params, grad)
        
        if grad < 0.001 { break }
    }
}
```

---

## Build System Integration

### Specifying Targets

```bash
# Build for NVIDIA GPU and IBM Quantum
arc build --target=gpu.cuda,qpu.ibm my_algo.arc

# Build for Simulator only
arc build --target=qpu.sim my_algo.arc
```

## Future Targets

```arc
<gpu.vulkan>      // Vulkan compute (cross-vendor)
<qpu.quantinuum>  // H-Series Trapped Ion Quantum
<qpu.pascal>      // Neutral Atom Quantum
<fpga.xilinx>     // Xilinx FPGAs
<asic.cerebras>   // Wafer-scale engine
```

---

## Summary

Arc's accelerator target system provides:

1.  **Explicit compilation targets** (`async func<target>`)
2.  **Native Code Generation** (PTX, HLO, QIR)
3.  **Heterogeneous Flexibility** (Mix CPU + GPU + QPU)
4.  **Quantum Readiness** (First-class support for the QIR standard)

**The pattern:** `async func<category.framework>` tells Arc exactly where and how to execute your code, giving you control over the entire compute spectrum—from bits to qubits.