# Multi-Architecture Code Generation Design

## Overview

This document describes the architecture for three GPU/TPU code generation packages that provide a unified Go API for generating hardware-specific compute code.

## Why Different Output Formats?

The three packages have different output formats because of varying levels of vendor openness:
```
AMD (GCN/RDNA):     Machine Code (bytes)
  └─ ISA fully documented → Direct instruction encoding possible

NVIDIA (PTX):       Assembly Text (string)  
  └─ Hardware ISA (SASS) undocumented → PTX as portable IR
  └─ Driver JIT-compiles PTX → SASS at runtime

Google (TPU):       High-Level Operations (string)
  └─ ISA completely proprietary → Only XLA/HLO interface available
  └─ XLA compiler handles all low-level code generation
```

**Key insight:** AMD publishes complete instruction encodings, NVIDIA provides an intermediate representation (PTX), and Google only has high-level graph descriptions (HLO). Your code generation must match what each vendor allows you to target.

```
┌─────────────────────────────────────────────────────┐
│          Your Go Application                        │
│     (ML Framework, Compute Library, etc.)           │
└─────────────────────────────────────────────────────┘
                       │
         ┌─────────────┼─────────────┐
         │             │             │
         ▼             ▼             ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ ptx-codegen  │ │ gcn-codegen  │ │ hlo-codegen  │
│   (NVIDIA)   │ │    (AMD)     │ │  (Google)    │
└──────────────┘ └──────────────┘ └──────────────┘
         │             │             │
         ▼             ▼             ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  PTX String  │ │ Machine Code │ │  HLO String  │
└──────────────┘ └──────────────┘ └──────────────┘
         │             │             │
         ▼             ▼             ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ CUDA Driver  │ │ ROCm Driver  │ │ XLA Runtime  │
└──────────────┘ └──────────────┘ └──────────────┘
```

## Design Philosophy

### Key Principles

1. **Backend-Specific IR** - Each package has its own IR tailored to the target hardware
2. **Consistent API** - Similar builder patterns across all packages
3. **Type Safety** - Leverage Go's type system to catch errors at compile time
4. **No Shared Abstraction** - Abstraction happens in higher-level packages if needed
5. **Output Format Matches Platform** - PTX → string, GCN → bytes, HLO → string

### Why No Shared IR?

```
❌ Bad: Force-fit all architectures into one IR
   Common IR → Translate → PTX/GCN/HLO
   (Loses platform-specific features)

✅ Good: Let each package expose native capabilities
   Native IR → Native Output
   (Maximum control and performance)
```

## Package Architecture

### Common Structure (All Three Packages)

```
{package}-codegen/
├── README.md
├── ir/                    # IR node definitions
│   ├── node.go           # Base IR node interface
│   ├── ops.go            # Operations (Add, Mul, Load, etc.)
│   ├── types.go          # Type system
│   └── values.go         # Values and constants
│
├── builder/               # IR Builder API (Main Entry Point)
│   ├── builder.go        # Builder struct and methods
│   ├── control_flow.go   # If/Loop/Branch helpers
│   └── memory.go         # Memory operations
│
├── emit/                  # Code emission
│   ├── emitter.go        # IR → Output format
│   └── optimize.go       # Basic optimizations (optional)
│
├── runtime/               # Runtime integration (optional)
│   └── loader.go         # Load and execute generated code
│
└── examples/
    ├── vector_add.go
    ├── matmul.go
    └── conv2d.go
```

## 1. PTX Code Generation (NVIDIA)

### Package: `ptx-codegen`

**Purpose**: Generate NVIDIA PTX (Parallel Thread Execution) assembly text

**Output**: PTX string (portable across NVIDIA GPUs)

### IR Nodes

```go
// ir/node.go
package ir

type Node interface {
    Type() Type
    String() string
}

// PTX-specific register types
type Register interface {
    Node
    RegType() RegisterType  // .b32, .b64, .f32, .f64, .pred
    RegNum() int
}

type RegisterType int

const (
    RegB32 RegisterType = iota  // 32-bit integer
    RegB64                      // 64-bit integer
    RegF32                      // 32-bit float
    RegF64                      // 64-bit float
    RegPred                     // Predicate (boolean)
)
```

```go
// ir/ops.go
package ir

// PTX operations
type AddOp struct {
    Dst  Register
    Src1 Register
    Src2 Register
    Type RegisterType
}

type LoadOp struct {
    Dst     Register
    Address Register
    Space   AddressSpace  // .global, .shared, .local, .param
    Type    RegisterType
}

type StoreOp struct {
    Address Register
    Value   Register
    Space   AddressSpace
    Type    RegisterType
}

type MulOp struct {
    Dst  Register
    Src1 Register
    Src2 Register
    Type RegisterType
    Mode string  // .lo, .hi, .wide
}
```

### Builder API

```go
// builder/builder.go
package builder

type Builder struct {
    nodes    []ir.Node
    regs     map[string]ir.Register
    regCount map[ir.RegisterType]int
    params   []Parameter
}

func New() *Builder {
    return &Builder{
        regs:     make(map[string]ir.Register),
        regCount: make(map[ir.RegisterType]int),
    }
}

// Allocate a new register
func (b *Builder) NewReg(name string, t ir.RegisterType) ir.Register {
    num := b.regCount[t]
    b.regCount[t]++
    
    reg := &ptxRegister{
        name:    name,
        regType: t,
        regNum:  num,
    }
    b.regs[name] = reg
    return reg
}

// Add operation
func (b *Builder) Add(dst, src1, src2 ir.Register) {
    b.nodes = append(b.nodes, &ir.AddOp{
        Dst:  dst,
        Src1: src1,
        Src2: src2,
        Type: dst.RegType(),
    })
}

// Multiply operation
func (b *Builder) Mul(dst, src1, src2 ir.Register) {
    b.nodes = append(b.nodes, &ir.MulOp{
        Dst:  dst,
        Src1: src1,
        Src2: src2,
        Type: dst.RegType(),
    })
}

// Load from global memory
func (b *Builder) LoadGlobal(dst, addr ir.Register) {
    b.nodes = append(b.nodes, &ir.LoadOp{
        Dst:     dst,
        Address: addr,
        Space:   ir.SpaceGlobal,
        Type:    dst.RegType(),
    })
}

// Store to global memory
func (b *Builder) StoreGlobal(addr, value ir.Register) {
    b.nodes = append(b.nodes, &ir.StoreOp{
        Address: addr,
        Value:   value,
        Space:   ir.SpaceGlobal,
        Type:    value.RegType(),
    })
}

// Define kernel parameters
func (b *Builder) AddParam(name string, t ir.Type) ir.Register {
    param := Parameter{Name: name, Type: t}
    b.params = append(b.params, param)
    
    // Return register that will hold the parameter
    reg := b.NewReg(name, t.RegType())
    return reg
}

// Emit PTX string
func (b *Builder) Emit() string {
    e := emit.NewEmitter()
    return e.EmitPTX(b.nodes, b.params, b.regCount)
}
```

### Emitter

```go
// emit/emitter.go
package emit

import "strings"

type Emitter struct {
    buf strings.Builder
}

func (e *Emitter) EmitPTX(nodes []ir.Node, params []Parameter, regCount map[ir.RegisterType]int) string {
    e.buf.WriteString(".version 8.0\n")
    e.buf.WriteString(".target sm_80\n")
    e.buf.WriteString(".address_size 64\n\n")
    
    // Emit kernel entry
    e.buf.WriteString(".visible .entry kernel(\n")
    for i, p := range params {
        e.buf.WriteString(fmt.Sprintf("    .param .%s %s", p.Type.PTXType(), p.Name))
        if i < len(params)-1 {
            e.buf.WriteString(",\n")
        }
    }
    e.buf.WriteString("\n) {\n")
    
    // Declare registers
    for regType, count := range regCount {
        if count > 0 {
            e.buf.WriteString(fmt.Sprintf("    .reg .%s %%r<%d>;\n", 
                regType.PTXName(), count))
        }
    }
    e.buf.WriteString("\n")
    
    // Emit operations
    for _, node := range nodes {
        e.emitNode(node)
    }
    
    e.buf.WriteString("    ret;\n")
    e.buf.WriteString("}\n")
    
    return e.buf.String()
}

func (e *Emitter) emitNode(node ir.Node) {
    switch n := node.(type) {
    case *ir.AddOp:
        e.buf.WriteString(fmt.Sprintf("    add.%s %s, %s, %s;\n",
            n.Type.PTXName(), n.Dst, n.Src1, n.Src2))
    case *ir.MulOp:
        e.buf.WriteString(fmt.Sprintf("    mul.%s %s, %s, %s;\n",
            n.Type.PTXName(), n.Dst, n.Src1, n.Src2))
    case *ir.LoadOp:
        e.buf.WriteString(fmt.Sprintf("    ld.%s.%s %s, [%s];\n",
            n.Space.PTXName(), n.Type.PTXName(), n.Dst, n.Address))
    case *ir.StoreOp:
        e.buf.WriteString(fmt.Sprintf("    st.%s.%s [%s], %s;\n",
            n.Space.PTXName(), n.Type.PTXName(), n.Address, n.Value))
    }
}
```

### Usage Example

```go
package main

import (
    "fmt"
    "ptx-codegen/builder"
    "ptx-codegen/ir"
)

func main() {
    b := builder.New()
    
    // Define kernel parameters
    ptrA := b.AddParam("a", ir.TypeU64)
    ptrB := b.AddParam("b", ir.TypeU64)
    ptrC := b.AddParam("c", ir.TypeU64)
    
    // Allocate registers
    valA := b.NewReg("valA", ir.RegF32)
    valB := b.NewReg("valB", ir.RegF32)
    valC := b.NewReg("valC", ir.RegF32)
    
    // Generate code: C[i] = A[i] + B[i]
    b.LoadGlobal(valA, ptrA)
    b.LoadGlobal(valB, ptrB)
    b.Add(valC, valA, valB)
    b.StoreGlobal(ptrC, valC)
    
    // Emit PTX
    ptx := b.Emit()
    fmt.Println(ptx)
}
```

**Output**:
```ptx
.version 8.0
.target sm_80
.address_size 64

.visible .entry kernel(
    .param .u64 a,
    .param .u64 b,
    .param .u64 c
) {
    .reg .f32 %r<3>;
    .reg .u64 %rd<3>;
    
    ld.param.u64 %rd0, [a];
    ld.param.u64 %rd1, [b];
    ld.param.u64 %rd2, [c];
    ld.global.f32 %r0, [%rd0];
    ld.global.f32 %r1, [%rd1];
    add.f32 %r2, %r0, %r1;
    st.global.f32 [%rd2], %r2;
    ret;
}
```

## 2. GCN Code Generation (AMD)

### Package: `gcn-codegen`

**Purpose**: Generate AMD GCN/RDNA machine code

**Output**: Binary machine code (architecture-specific)

### IR Nodes

```go
// ir/node.go
package ir

type Node interface {
    Type() Type
    Encode(arch Architecture) ([]byte, error)
}

// AMD-specific register types
type VGPROperand struct {
    Index uint8
}

type SGPROperand struct {
    Index uint8
}

type VGPRPair struct {
    Base uint8  // Base register (uses Base and Base+1)
}

type SGPRRange struct {
    Base  uint8
    Count uint8
}
```

```go
// ir/ops.go
package ir

// AMD GCN operations
type VAddF32 struct {
    Dst  VGPROperand
    Src0 Operand
    Src1 Operand
}

type VMulF32 struct {
    Dst  VGPROperand
    Src0 Operand
    Src1 Operand
}

type BufferLoad struct {
    Dst    VGPROperand
    VAddr  VGPROperand
    SBase  SGPRRange  // 4 SGPRs for buffer descriptor
    Offset int32
}

type BufferStore struct {
    Src    VGPROperand
    VAddr  VGPROperand
    SBase  SGPRRange
    Offset int32
}

type SLoadB256 struct {
    DstBase uint8      // Loads into DstBase:DstBase+7
    SrcBase uint8      // Address in SrcBase:SrcBase+1
    Offset  int32
}
```

### Builder API

```go
// builder/builder.go
package builder

type Builder struct {
    nodes    []ir.Node
    arch     Architecture
    vgprUsed int
    sgprUsed int
}

func New(arch Architecture) *Builder {
    return &Builder{
        arch: arch,
    }
}

// Allocate VGPR
func (b *Builder) NewVGPR() ir.VGPROperand {
    vgpr := ir.VGPROperand{Index: uint8(b.vgprUsed)}
    b.vgprUsed++
    return vgpr
}

// Allocate SGPR range
func (b *Builder) NewSGPRRange(count uint8) ir.SGPRRange {
    base := uint8(b.sgprUsed)
    b.sgprUsed += int(count)
    return ir.SGPRRange{Base: base, Count: count}
}

// Vector add
func (b *Builder) VAddF32(dst, src0, src1 ir.VGPROperand) {
    b.nodes = append(b.nodes, &ir.VAddF32{
        Dst:  dst,
        Src0: src0,
        Src1: src1,
    })
}

// Vector multiply
func (b *Builder) VMulF32(dst, src0, src1 ir.VGPROperand) {
    b.nodes = append(b.nodes, &ir.VMulF32{
        Dst:  dst,
        Src0: src0,
        Src1: src1,
    })
}

// Load from buffer
func (b *Builder) BufferLoad(dst, vaddr ir.VGPROperand, sbase ir.SGPRRange) {
    b.nodes = append(b.nodes, &ir.BufferLoad{
        Dst:    dst,
        VAddr:  vaddr,
        SBase:  sbase,
        Offset: 0,
    })
}

// Store to buffer
func (b *Builder) BufferStore(src, vaddr ir.VGPROperand, sbase ir.SGPRRange) {
    b.nodes = append(b.nodes, &ir.BufferStore{
        Src:    src,
        VAddr:  vaddr,
        SBase:  sbase,
        Offset: 0,
    })
}

// Load kernel arguments
func (b *Builder) LoadArgs() {
    b.nodes = append(b.nodes, &ir.SLoadB256{
        DstBase: 4,      // Load into s[4:11]
        SrcBase: 0,      // From s[0:1]
        Offset:  0,
    })
    b.nodes = append(b.nodes, &ir.SWaitCnt{LGKMCnt: 0})
}

// End kernel
func (b *Builder) EndProgram() {
    b.nodes = append(b.nodes, &ir.SEndPgm{})
}

// Emit machine code
func (b *Builder) Emit() ([]byte, error) {
    var code []byte
    for _, node := range b.nodes {
        bytes, err := node.Encode(b.arch)
        if err != nil {
            return nil, err
        }
        code = append(code, bytes...)
    }
    return code, nil
}
```

### Usage Example

```go
package main

import (
    "fmt"
    "gcn-codegen/builder"
    "gcn-codegen/arch"
)

func main() {
    b := builder.New(arch.RDNA3())
    
    // Allocate registers
    v0 := b.NewVGPR()  // Thread offset
    v1 := b.NewVGPR()  // Value from A
    v2 := b.NewVGPR()  // Value from B
    v3 := b.NewVGPR()  // Result
    
    bufA := b.NewSGPRRange(4)  // Buffer descriptor for A
    bufB := b.NewSGPRRange(4)  // Buffer descriptor for B
    bufC := b.NewSGPRRange(4)  // Buffer descriptor for C
    
    // Load kernel arguments
    b.LoadArgs()
    
    // Generate code: C[i] = A[i] + B[i]
    b.BufferLoad(v1, v0, bufA)
    b.BufferLoad(v2, v0, bufB)
    b.Wait(0)  // Wait for loads
    b.VAddF32(v3, v1, v2)
    b.BufferStore(v3, v0, bufC)
    b.Wait(0)  // Wait for store
    b.EndProgram()
    
    // Emit machine code
    code, err := b.Emit()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Generated %d bytes of machine code\n", len(code))
    // code is now ready to wrap in HSACO and load via ROCm
}
```

## 3. HLO Code Generation (Google TPU)

### Package: `hlo-codegen`

**Purpose**: Generate XLA HLO (High-Level Optimizer) text format

**Output**: HLO string (for XLA compilation to TPU)

### IR Nodes

```go
// ir/node.go
package ir

type Node interface {
    Type() *TensorType
    String() string
}

// HLO-specific types
type TensorType struct {
    Shape    []int64
    ElemType PrimitiveType
}

type PrimitiveType int

const (
    F32 PrimitiveType = iota
    F64
    S32
    S64
    PRED
)
```

```go
// ir/ops.go
package ir

// HLO operations (tensor operations)
type AddOp struct {
    LHS *Tensor
    RHS *Tensor
}

type MulOp struct {
    LHS *Tensor
    RHS *Tensor
}

type DotOp struct {
    LHS *Tensor
    RHS *Tensor
    // Contracting dimensions
    LHSContracting []int64
    RHSContracting []int64
}

type ConvOp struct {
    Input  *Tensor
    Kernel *Tensor
    Window *WindowDimension
}

type ReduceOp struct {
    Input      *Tensor
    InitValue  *Tensor
    Dimensions []int64
    Reducer    string  // "add", "max", "min"
}

type ReshapeOp struct {
    Operand  *Tensor
    NewShape []int64
}
```

### Builder API

```go
// builder/builder.go
package builder

type Builder struct {
    nodes      []*Tensor
    nextID     int
    parameters map[string]*Tensor
}

func New() *Builder {
    return &Builder{
        parameters: make(map[string]*Tensor),
    }
}

// Create parameter (input tensor)
func (b *Builder) Parameter(name string, shape []int64, dtype ir.PrimitiveType) *Tensor {
    tensor := &Tensor{
        id:    b.nextID,
        op:    &ir.ParameterOp{Name: name},
        dtype: &ir.TensorType{Shape: shape, ElemType: dtype},
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    b.parameters[name] = tensor
    return tensor
}

// Create constant
func (b *Builder) Constant(value interface{}, shape []int64) *Tensor {
    tensor := &Tensor{
        id:    b.nextID,
        op:    &ir.ConstantOp{Value: value},
        dtype: inferType(value, shape),
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Element-wise add
func (b *Builder) Add(lhs, rhs *Tensor) *Tensor {
    tensor := &Tensor{
        id:    b.nextID,
        op:    &ir.AddOp{LHS: lhs, RHS: rhs},
        dtype: lhs.dtype,  // Assumes same type
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Element-wise multiply
func (b *Builder) Mul(lhs, rhs *Tensor) *Tensor {
    tensor := &Tensor{
        id:    b.nextID,
        op:    &ir.MulOp{LHS: lhs, RHS: rhs},
        dtype: lhs.dtype,
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Matrix multiply (dot product)
func (b *Builder) Dot(lhs, rhs *Tensor) *Tensor {
    tensor := &Tensor{
        id: b.nextID,
        op: &ir.DotOp{
            LHS:            lhs,
            RHS:            rhs,
            LHSContracting: []int64{len(lhs.dtype.Shape) - 1},
            RHSContracting: []int64{0},
        },
        dtype: computeDotType(lhs.dtype, rhs.dtype),
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Reshape
func (b *Builder) Reshape(operand *Tensor, newShape []int64) *Tensor {
    tensor := &Tensor{
        id: b.nextID,
        op: &ir.ReshapeOp{
            Operand:  operand,
            NewShape: newShape,
        },
        dtype: &ir.TensorType{
            Shape:    newShape,
            ElemType: operand.dtype.ElemType,
        },
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Reduce (sum, max, etc.)
func (b *Builder) Reduce(input *Tensor, dimensions []int64, reducer string) *Tensor {
    initValue := b.Constant(0.0, []int64{})
    
    tensor := &Tensor{
        id: b.nextID,
        op: &ir.ReduceOp{
            Input:      input,
            InitValue:  initValue,
            Dimensions: dimensions,
            Reducer:    reducer,
        },
        dtype: computeReduceType(input.dtype, dimensions),
    }
    b.nextID++
    b.nodes = append(b.nodes, tensor)
    return tensor
}

// Emit HLO text
func (b *Builder) Emit(root *Tensor) string {
    e := emit.NewEmitter()
    return e.EmitHLO(b.nodes, b.parameters, root)
}
```

### Emitter

```go
// emit/emitter.go
package emit

import "strings"

type Emitter struct {
    buf strings.Builder
}

func (e *Emitter) EmitHLO(nodes []*Tensor, params map[string]*Tensor, root *Tensor) string {
    e.buf.WriteString("HloModule computation\n\n")
    e.buf.WriteString("ENTRY main {\n")
    
    // Emit parameters
    for name, tensor := range params {
        e.buf.WriteString(fmt.Sprintf("  %%p%d = %s parameter(%d)\n",
            tensor.id, tensor.dtype.HLOString(), tensor.id))
    }
    e.buf.WriteString("\n")
    
    // Emit operations
    for _, tensor := range nodes {
        if _, isParam := tensor.op.(*ir.ParameterOp); isParam {
            continue  // Already emitted
        }
        e.emitTensor(tensor)
    }
    
    e.buf.WriteString(fmt.Sprintf("  ROOT %%result = %s %%t%d\n",
        root.dtype.HLOString(), root.id))
    e.buf.WriteString("}\n")
    
    return e.buf.String()
}

func (e *Emitter) emitTensor(t *Tensor) {
    switch op := t.op.(type) {
    case *ir.AddOp:
        e.buf.WriteString(fmt.Sprintf("  %%t%d = %s add(%%t%d, %%t%d)\n",
            t.id, t.dtype.HLOString(), op.LHS.id, op.RHS.id))
    case *ir.MulOp:
        e.buf.WriteString(fmt.Sprintf("  %%t%d = %s multiply(%%t%d, %%t%d)\n",
            t.id, t.dtype.HLOString(), op.LHS.id, op.RHS.id))
    case *ir.DotOp:
        e.buf.WriteString(fmt.Sprintf("  %%t%d = %s dot(%%t%d, %%t%d), "+
            "lhs_contracting_dims={%s}, rhs_contracting_dims={%s}\n",
            t.id, t.dtype.HLOString(), op.LHS.id, op.RHS.id,
            formatDims(op.LHSContracting), formatDims(op.RHSContracting)))
    case *ir.ReshapeOp:
        e.buf.WriteString(fmt.Sprintf("  %%t%d = %s reshape(%%t%d)\n",
            t.id, t.dtype.HLOString(), op.Operand.id))
    }
}
```

### Usage Example

```go
package main

import (
    "fmt"
    "hlo-codegen/builder"
    "hlo-codegen/ir"
)

func main() {
    b := builder.New()
    
    // Define parameters
    A := b.Parameter("A", []int64{128, 256}, ir.F32)
    B := b.Parameter("B", []int64{256, 512}, ir.F32)
    
    // Matrix multiply: C = A × B
    C := b.Dot(A, B)
    
    // Add bias
    bias := b.Constant(0.5, []int64{128, 512})
    result := b.Add(C, bias)
    
    // Emit HLO
    hlo := b.Emit(result)
    fmt.Println(hlo)
}
```

**Output**:
```hlo
HloModule computation

ENTRY main {
  %p0 = f32[128,256] parameter(0)
  %p1 = f32[256,512] parameter(1)
  
  %t2 = f32[128,512] dot(%p0, %p1), lhs_contracting_dims={1}, rhs_contracting_dims={0}
  %t3 = f32[128,512] constant(0.5)
  %t4 = f32[128,512] add(%t2, %t3)
  
  ROOT %result = f32[128,512] %t4
}
```

## Unified Simple API (Optional)

If you want a simple unified interface for basic operations:

```go
// simple/compute.go
package simple

type Compute interface {
    Add(a, b Value) Value
    Mul(a, b Value) Value
    Load(addr Value) Value
    Store(addr, value Value)
    Emit() interface{}  // Returns string or []byte
}

// Factory functions
func NewPTX() Compute { return &ptxCompute{...} }
func NewGCN(arch) Compute { return &gcnCompute{...} }
func NewHLO() Compute { return &hloCompute{...} }

// Usage
func genericKernel(c Compute) {
    a := c.Load(...)
    b := c.Load(...)
    result := c.Add(a, b)
    c.Store(..., result)
}

// Use with any backend
genericKernel(simple.NewPTX())
genericKernel(simple.NewGCN(arch.RDNA3()))
genericKernel(simple.NewHLO())
```

## Summary

### Package Outputs

| Package | Output Format | Size | Usage |
|---------|---------------|------|-------|
| `ptx-codegen` | PTX string | ~1-10 KB | Load via CUDA driver |
| `gcn-codegen` | Machine code bytes | ~100-1000 bytes | Wrap in HSACO, load via ROCm |
| `hlo-codegen` | HLO string | ~1-10 KB | Pass to XLA compiler |

### API Consistency

All three packages share similar patterns:
- ✅ Builder-based IR construction
- ✅ Type-safe operations
- ✅ `Emit()` method returns output
- ✅ Similar method names where possible

### When to Abstract

**Don't abstract immediately** - let each package mature independently.

**Later, create abstraction** in a separate package if you need:
- Backend-agnostic ML framework
- Multi-target compilation
- Runtime backend selection

```go
// Future: unified-compute package
package unified

type Backend int
const (
    NVIDIA Backend = iota
    AMD
    TPU
)

func Compile(ir *IR, backend Backend) interface{} {
    switch
    backend {
    case NVIDIA:
        return ptx.Compile(ir)
    case AMD:
        return gcn.Compile(ir)
    case TPU:
        return hlo.Compile(ir)
    }
}
```

## Conclusion

These three packages provide Go with native code generation capabilities for major compute platforms, each optimized for its target architecture while maintaining consistent API patterns.