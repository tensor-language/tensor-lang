# Arc Async Smart Threads - Complete Architecture

## What Are Arc Async Smart Threads?

Arc's `async` functions compile to **real OS threads** with **self-growing stacks**. Unlike traditional threading models:

- **No manual stack sizing** - stacks grow automatically from 64KB to 1MB as needed
- **Thousands of concurrent threads** - low initial memory (64KB per thread)
- **Zero runtime** - direct syscalls, no scheduler overhead in V1
- **Cross-platform** - Linux, macOS, Windows with same semantics

```arc
async compute(x: int) int {
    x * 2  // Runs on real OS thread with auto-growing stack
}

async func main() {
    let result = await compute(42)  // Block until thread completes
}
```

---

## High-Level Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                   Arc Async Execution Model                    │
└────────────────────────────────────────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   async compute(x: int)       │
              │       x * 2                   │
              └───────────────┬───────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   Compiler Translation        │
              └───────────────┬───────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌───────────────┐    ┌────────────────┐    ┌──────────────┐
│ Stack Manager │    │ Thread Spawner │    │ Synchronizer │
│  (V1 + V2)    │    │     (V1)       │    │    (V2)      │
└───────┬───────┘    └────────┬───────┘    └──────┬───────┘
        │                     │                    │
        │                     │                    │
┌───────▼──────────────────────▼────────────────────▼───────┐
│                    OS Thread (clone3/pthread/CreateThread) │
│  ┌──────────────────────────────────────────────────────┐ │
│  │                    Stack Layout                      │ │
│  │  High Addr  ┌─────────────────┐ ◄─ RSP (grows down)│ │
│  │             │  Committed (64KB)│ ✓ Read/Write       │ │
│  │             ├─────────────────┤                     │ │
│  │             │  Guard Pages     │ ✗ No Access        │ │
│  │  Low Addr   └─────────────────┘ (triggers handler)  │ │
│  │             ◄────── 1MB Virtual ──────►             │ │
│  └──────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────┘
```

---

## Three Core Systems

### 1. Stack Manager - Auto-Growing Stacks

**Purpose:** Allow threads to start small (64KB) but grow to 1MB without crashes.

**How it works:**
1. Reserve 1MB virtual address space (costs ~0 physical memory)
2. Commit only top 64KB as read/write (physical memory allocated)
3. Rest is "guard zone" - no permissions
4. When thread accesses guard zone → **page fault** → kernel calls handler
5. Handler commits more pages (makes them read/write)
6. Thread resumes, access succeeds

```
Thread Execution                    Kernel
─────────────────                   ──────
push rax         ───────────────► Page fault! (guard zone)
(access 0x...64KB)                  ├─ Lookup signal table
                                    ├─ Find registered handler
                 ◄────────────────  └─ Call handler()
                                    
Handler:
  mprotect(page, RW) ─────────────► Commit 64KB more pages
  return          ────────────────► Retry instruction
                                    
push rax         ───────────────► Success! (now committed)
```

**Platform Implementations:**

| Platform | Reserve Mechanism | Commit Mechanism | Handler Type |
|:---------|:------------------|:-----------------|:-------------|
| **Linux** | `mmap(PROT_NONE)` | `mprotect(RW)` | `SIGSEGV` signal |
| **macOS** | `mmap(PROT_NONE)` | `mprotect(RW)` | Mach exception port |
| **Windows** | `VirtualAlloc(MEM_RESERVE)` | `VirtualAlloc(MEM_COMMIT)` | VEH (Vectored Exception) |

---

### 2. Thread Spawner - Creating Threads

**Purpose:** Create OS threads with custom stacks.

**How it works:**
1. Allocate stack (using Stack Manager)
2. Spawn OS thread pointing to that stack
3. Thread starts executing async function
4. On completion: store result, wake waiters, exit

```
Parent Thread                           Child Thread
─────────────                          ────────────
async_spawn(compute, 42)
├─ mmap(1MB, PROT_NONE)
├─ mprotect(64KB, RW)
├─ clone3(stack=new_stack) ──────────► [New thread born]
│                                       ├─ RSP = new_stack + 1MB
│                                       ├─ call compute(42)
│                                       ├─ result = 84
│                                       ├─ store result in AsyncResult
│                                       ├─ futex_wake() (V2)
│                                       └─ exit(0)
└─ return thread_id
```

**Platform Implementations:**

| Platform | Thread Creation | Stack Setup | Calling Convention |
|:---------|:----------------|:------------|:-------------------|
| **Linux** | `clone3()` syscall | Pass stack in `clone_args.stack` | SysV AMD64 |
| **macOS** | `pthread_create()` | `pthread_attr_setstack()` | SysV AMD64 |
| **Windows** | `CreateThread()` | `STACK_SIZE_PARAM_IS_A_RESERVATION` flag | Microsoft x64 |

---

### 3. Synchronizer - Awaiting Completion (V2)

**Purpose:** Make `await` block until thread finishes, retrieve result.

**How it works:**
1. Parent allocates `AsyncResult` struct in memory
2. Passes pointer to child thread
3. Child stores result when done, marks "ready", wakes parent
4. Parent sleeps on futex until "ready" flag is set

```
struct AsyncResult {
    ready: int32    // 0 = running, 1 = done (atomic)
    value: int64    // return value
}
```

**Execution flow:**

```
Parent                          Child                    Kernel
──────                          ─────                    ──────
result = alloc(AsyncResult)
result.ready = 0
spawn(child, result_ptr) ─────► start
                                execute...
await:                          done, result = 84
├─ while !result.ready:         store result.value = 84
│   futex_wait(&result.ready) ─► atomic_store(result.ready = 1)
│   (sleep)                      futex_wake(&result.ready) ──► Wake parent
│                                exit(0)
◄─ (wake up)
├─ load result.value
└─ return 84
```

**Why not just spin-wait?**
```c
// BAD: Burns CPU
while (!result.ready) { /* spin */ }

// GOOD: Sleeps in kernel, 0% CPU
futex_wait(&result.ready, 0);
```

**Platform Implementations:**

| Platform | Sleep Primitive | Wake Primitive | Syscall # |
|:---------|:----------------|:---------------|:----------|
| **Linux** | `futex(FUTEX_WAIT)` | `futex(FUTEX_WAKE)` | 202 |
| **macOS** | `__ulock_wait()` | `__ulock_wake()` | 0x2000203 / 0x2000204 |
| **Windows** | `WaitOnAddress()` | `WakeByAddressSingle()` | kernel32.dll |

---

## Complete Data Flow Example

```arc
async compute(x: int) -> int {
    let arr = [0; 100000]  // 800KB allocation
    x * 2
}

func main() {
    let result = await compute(42)
    print(result)
}
```

**Step-by-step execution:**

```
1. Main thread executes compute(42)
   ├─ mmap(1MB, PROT_NONE)          [Stack Manager]
   │  └─ Virtual: 0x7000_0000 - 0x7010_0000
   ├─ mprotect(0x700F_0000, 64KB, RW) [Stack Manager]
   │  └─ Committed: 0x700F_0000 - 0x7010_0000
   ├─ malloc(AsyncResult, 16 bytes)  [Synchronizer]
   │  └─ result @ 0x5555_1000 = {ready: 0, value: 0}
   └─ clone3(stack=0x7010_0000, func=compute, arg=42, result=0x5555_1000)
      └─ New thread TID=12345

2. Child thread (TID=12345) starts
   ├─ RSP = 0x7010_0000 (top of stack)
   ├─ call compute(42)
   │  ├─ allocate [0; 100000] → need 800KB
   │  ├─ RSP -= 800000 → RSP = 0x700D_D000
   │  └─ Access 0x700D_D000... ⚠️ GUARD ZONE!

3. Page Fault! 
   Kernel: "Thread 12345 accessed 0x700D_D000 (no permission)"
   ├─ Check signal table: SIGSEGV → handler @ 0x4000_1234
   └─ Call handler(signum=11, fault_addr=0x700D_D000, ...)

4. SIGSEGV Handler executes             [Stack Manager]
   ├─ Check: is 0x700D_D000 in stack bounds? YES
   ├─ Round down to page: 0x700D_D000 → 0x700D_D000
   ├─ mprotect(0x700D_D000, 64KB, RW)
   │  └─ Committed now: 0x700D_D000 - 0x7010_0000 (128KB total)
   └─ return → kernel retries instruction

5. Thread resumes
   ├─ RSP -= 800000 → SUCCESS (page now writable)
   ├─ compute: 42 * 2 = 84
   └─ return 84

6. Thread exit path                     [Synchronizer]
   ├─ store result.value = 84 @ 0x5555_1008
   ├─ atomic_store(result.ready = 1) @ 0x5555_1000
   ├─ futex(0x5555_1000, FUTEX_WAKE, 1)
   │  └─ Kernel: "Wake 1 thread waiting on 0x5555_1000"
   └─ exit(0)

7. Main thread (in await)               [Synchronizer]
   ├─ futex_wait(0x5555_1000, 0)
   │  └─ (sleeping, 0% CPU)
   ├─ ← Kernel wakes thread
   ├─ Check result.ready @ 0x5555_1000 → 1 ✓
   ├─ load result.value @ 0x5555_1008 → 84
   ├─ munmap(0x7000_0000, 1MB)  [cleanup child stack]
   └─ return 84

8. print(84)
```

---

## Platform-Specific Details

### Linux Implementation

```go
// 1. Stack Manager - SIGSEGV Handler (registered once at program start)
func (c *compiler) emitSigactionSetup() {
    handlerLabel := c.asm.NewLabel()
    
    // Signal handler code
    c.asm.Label(handlerLabel)
    // Args: RDI=signum, RSI=siginfo_t*, RDX=ucontext_t*
    
    // Get fault address from siginfo_t->si_addr
    c.asm.Mov(RegOp(RAX), NewMem(RSI, 16), 64)
    
    // Get RSP from ucontext
    c.asm.Mov(RegOp(RCX), NewMem(RDX, 160), 64)
    
    // Check if fault in stack bounds [RSP-1MB, RSP]
    c.asm.Mov(RegOp(RDI), RegOp(RCX), 64)
    c.asm.Sub(RegOp(RDI), ImmOp(1<<20))  // stack_base
    c.asm.Cmp(RAX, RDI)
    outOfBounds := c.asm.JccRel(CondL, 0)
    
    // Fault in our stack - grow it
    c.asm.And(RegOp(RAX), ImmOp(^0xFFF))  // page align
    
    // mprotect(page, 64KB, PROT_READ|PROT_WRITE)
    c.asm.Mov(RegOp(RDI), RegOp(RAX), 64)
    c.asm.Mov(RegOp(RSI), ImmOp(64<<10), 64)
    c.asm.Mov(RegOp(RDX), ImmOp(3), 64)  // PROT_READ|PROT_WRITE
    c.asm.Mov(RegOp(RAX), ImmOp(10), 64) // SYS_mprotect
    c.asm.Syscall()
    c.asm.Ret()
    
    // Out of bounds - re-raise signal
    c.asm.Label(outOfBounds)
    // ... restore default handler and re-raise
    
    // Install handler: rt_sigaction(SIGSEGV, &act, NULL, 8)
    c.asm.Sub(RegOp(RSP), ImmOp(152))  // sizeof(sigaction)
    c.asm.Mov(NewMem(RSP, 0), ImmOp(handlerAddr), 64)
    c.asm.Mov(NewMem(RSP, 8), ImmOp(0x04|0x08000000), 64) // SA_SIGINFO|SA_ONSTACK
    
    c.asm.Mov(RegOp(RDI), ImmOp(11), 64)  // SIGSEGV
    c.asm.Mov(RegOp(RSI), RegOp(RSP), 64)
    c.asm.Xor(RegOp(RDX), RegOp(RDX))
    c.asm.Mov(RegOp(R10), ImmOp(8), 64)
    c.asm.Mov(RegOp(RAX), ImmOp(13), 64)  // SYS_rt_sigaction
    c.asm.Syscall()
    
    c.asm.Add(RegOp(RSP), ImmOp(152))
}

// 2. Thread Spawner
func (c *compiler) compileAsyncSpawn(inst *ir.AsyncSpawnInst) error {
    // Allocate stack: mmap(NULL, 1MB, PROT_NONE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0)
    c.asm.Xor(RegOp(RDI), RegOp(RDI))
    c.asm.Mov(RegOp(RSI), ImmOp(1<<20), 64)
    c.asm.Xor(RegOp(RDX), RegOp(RDX))  // PROT_NONE
    c.asm.Mov(RegOp(R10), ImmOp(0x22), 64)  // MAP_PRIVATE|MAP_ANONYMOUS
    c.asm.Mov(RegOp(R8), ImmOp(-1), 64)
    c.asm.Xor(RegOp(R9), RegOp(R9))
    c.asm.Mov(RegOp(RAX), ImmOp(9), 64)  // SYS_mmap
    c.asm.Syscall()
    c.asm.Push(RAX)  // save stack base
    
    // Commit top 64KB: mprotect(base + 1MB - 64KB, 64KB, PROT_READ|PROT_WRITE)
    c.asm.Mov(RegOp(RDI), RegOp(RAX), 64)
    c.asm.Add(RegOp(RDI), ImmOp((1<<20)-(64<<10)))
    c.asm.Mov(RegOp(RSI), ImmOp(64<<10), 64)
    c.asm.Mov(RegOp(RDX), ImmOp(3), 64)
    c.asm.Mov(RegOp(RAX), ImmOp(10), 64)  // SYS_mprotect
    c.asm.Syscall()
    
    // Allocate AsyncResult: malloc(16)
    c.asm.Mov(RegOp(RDI), ImmOp(16), 64)
    c.emitCall("malloc")
    c.asm.Mov(NewMem(RAX, 0), ImmOp(0), 32)  // ready = 0
    c.asm.Push(RAX)  // save result pointer
    
    // Setup clone3 args
    c.asm.Sub(RegOp(RSP), ImmOp(88))
    c.asm.Mov(NewMem(RSP, 0), ImmOp(0x00010F00), 64)  // CLONE_VM|CLONE_THREAD|...
    c.asm.Mov(RegOp(RAX), NewMem(RSP, 88+8), 64)  // stack base
    c.asm.Mov(NewMem(RSP, 40), RegOp(RAX), 64)
    c.asm.Mov(NewMem(RSP, 48), ImmOp((1<<20)-16), 64)  // stack size
    
    // clone3(&args, 88)
    c.asm.Mov(RegOp(RDI), RegOp(RSP), 64)
    c.asm.Mov(RegOp(RSI), ImmOp(88), 64)
    c.asm.Mov(RegOp(RAX), ImmOp(435), 64)  // SYS_clone3
    c.asm.Syscall()
    
    // Check if child (RAX == 0)
    c.asm.Test(RAX, RAX)
    childJmp := c.asm.JccRel(CondEq, 0)
    
    // Parent: cleanup and return
    c.asm.Add(RegOp(RSP), ImmOp(88+8+8))
    c.store(RAX, inst)  // return thread ID
    parentDone := c.asm.JmpRel(0)
    
    // Child: execute function and exit
    c.asm.PatchJump(childJmp)
    c.asm.Pop(RDI)  // result pointer (arg 1)
    c.asm.Pop(RSI)  // user arg (arg 2)
    c.load(RAX, inst.Func)
    c.asm.CallReg(RAX)  // call async function
    
    // Store result and wake
    c.asm.Pop(RDI)  // result pointer
    c.asm.Mov(NewMem(RDI, 8), RegOp(RAX), 64)  // result.value = RAX
    c.asm.Mov(NewMem(RDI, 0), ImmOp(1), 32)    // result.ready = 1
    
    // futex(&result.ready, FUTEX_WAKE, 1, ...)
    c.asm.Mov(RegOp(RSI), ImmOp(1), 64)  // FUTEX_WAKE
    c.asm.Mov(RegOp(RDX), ImmOp(1), 64)  // wake 1 thread
    c.asm.Xor(RegOp(R10), RegOp(R10))
    c.asm.Xor(RegOp(R8), RegOp(R8))
    c.asm.Xor(RegOp(R9), RegOp(R9))
    c.asm.Mov(RegOp(RAX), ImmOp(202), 64)  // SYS_futex
    c.asm.Syscall()
    
    // exit(0)
    c.asm.Xor(RegOp(RDI), RegOp(RDI))
    c.asm.Mov(RegOp(RAX), ImmOp(60), 64)  // SYS_exit
    c.asm.Syscall()
    
    c.asm.PatchJump(parentDone)
    return nil
}

// 3. Synchronizer - await
func (c *compiler) compileAwait(inst *ir.AwaitInst) error {
    c.load(RDI, inst.Handle)  // RDI = &result
    
    waitLoop := c.asm.NewLabel()
    c.asm.Label(waitLoop)
    
    // Check if ready
    c.asm.Mov(RegOp(RAX), NewMem(RDI, 0), 32)
    c.asm.Test(RAX, RAX)
    doneLabel := c.asm.JccRel(CondNE, 0)
    
    // futex(&result.ready, FUTEX_WAIT, 0, NULL, NULL, 0)
    c.asm.Xor(RegOp(RSI), RegOp(RSI))  // FUTEX_WAIT
    c.asm.Xor(RegOp(RDX), RegOp(RDX))  // expected = 0
    c.asm.Xor(RegOp(R10), RegOp(R10))  // timeout = NULL
    c.asm.Xor(RegOp(R8), RegOp(R8))
    c.asm.Xor(RegOp(R9), RegOp(R9))
    c.asm.Mov(RegOp(RAX), ImmOp(202), 64)  // SYS_futex
    c.asm.Syscall()
    
    c.asm.JmpRel(waitLoop)
    
    // Done
    c.asm.Label(doneLabel)
    c.asm.Mov(RegOp(RAX), NewMem(RDI, 8), 64)  // load result.value
    c.store(RAX, inst)
    return nil
}
```

### macOS Implementation

**Key differences from Linux:**
- Use Mach exception ports instead of signals
- syscall numbers are different (0x2000000 base)
- Use `pthread_create` with `pthread_attr_setstack` instead of `clone3`

```go
// 1. Stack Manager - Mach Exception Handler
func (c *compiler) emitMachExceptionHandler() {
    // Register exception port for EXC_BAD_ACCESS
    // Similar logic to SIGSEGV but uses:
    //   - mach_task_self() to get task port
    //   - task_set_exception_ports() to register handler
    //   - Handler uses vm_protect() instead of mprotect()
}

// 2. Thread Spawner - Uses pthread API
func (c *compiler) compileAsyncSpawnDarwin(inst *ir.AsyncSpawnInst) error {
    // Stack allocation same as Linux
    // ... mmap syscall 0x20000C5 ...
    
    // Use pthread_create with custom stack:
    //   pthread_attr_t attr;
    //   pthread_attr_init(&attr);
    //   pthread_attr_setstack(&attr, stack_base, 1MB);
    //   pthread_create(&tid, &attr, start_routine, arg);
}

// 3. Synchronizer - Uses __ulock_wait
func (c *compiler) compileAwaitDarwin(inst *ir.AwaitInst) error {
    // Same logic as Linux but:
    //   - syscall 0x2000203 for __ulock_wait
    //   - syscall 0x2000204 for __ulock_wake
    //   - RDI = UL_COMPARE_AND_WAIT (1)
}
```

### Windows Implementation

**Key differences:**
- Must use DLL functions (kernel32.dll), not raw syscalls
- Different calling convention (RCX, RDX, R8, R9 for args)
- Stack growth partially automatic via guard pages

```go
// 1. Stack Manager - VEH (Vectored Exception Handler)
func (c *compiler) emitVEHHandler() {
    // AddVectoredExceptionHandler(1, handler)
    handlerLabel := c.asm.NewLabel()
    
    c.asm.Label(handlerLabel)
    // Args: RCX = EXCEPTION_POINTERS*
    
    // Check exception code: if EXCEPTION_GUARD_PAGE
    c.asm.Mov(RegOp(RAX), NewMem(RCX, 0), 64)  // ExceptionRecord
    c.asm.Mov(RegOp(RAX), NewMem(RAX, 0), 32)  // ExceptionCode
    c.asm.Cmp(RAX, ImmOp(0x80000001))  // EXCEPTION_GUARD_PAGE
    notGuard := c.asm.JccRel(CondNE, 0)
    
    // Get fault address
    c.asm.Mov(RegOp(RAX), NewMem(RCX, 0), 64)
    c.asm.Mov(RegOp(RDX), NewMem(RAX, 16), 64)  // ExceptionInformation[1]
    
    // VirtualAlloc(addr, 64KB, MEM_COMMIT, PAGE_READWRITE)
    c.asm.And(RegOp(RDX), ImmOp(^0xFFFF))  // 64KB align
    c.asm.Mov(RegOp(RCX), RegOp(RDX), 64)
    c.asm.Mov(RegOp(RDX), ImmOp(64<<10), 64)
    c.asm.Mov(RegOp(R8), ImmOp(0x1000), 64)   // MEM_COMMIT
    c.asm.Mov(RegOp(R9), ImmOp(0x04), 64)     // PAGE_READWRITE
    c.emitWindowsCall("VirtualAlloc")
    
    // Return EXCEPTION_CONTINUE_EXECUTION
    c.asm.Xor(RegOp(RAX), RegOp(RAX))
    c.asm.Ret()
    
    c.asm.Label(notGuard)
    c.asm.Mov(RegOp(RAX), ImmOp(-1), 64)  // EXCEPTION_CONTINUE_SEARCH
    c.asm.Ret()
    
    // Install: AddVectoredExceptionHandler(1, handler)
    c.asm.Mov(RegOp(RCX), ImmOp(1), 64)  // first=1
    c.asm.Mov(RegOp(RDX), ImmOp(handlerAddr), 64)
    c.emitWindowsCall("AddVectoredExceptionHandler")
}

// 2. Thread Spawner
func (c *compiler) compileAsyncSpawnWindows(inst *ir.AsyncSpawnInst) error {
    // VirtualAlloc(NULL, 1MB, MEM_RESERVE, PAGE_NOACCESS)
    c.asm.Xor(RegOp(RCX), RegOp(RCX))
    c.asm.Mov(RegOp(RDX), ImmOp(1<<20), 64)
    c.asm.Mov(RegOp(R8), ImmOp(0x2000), 64)   // MEM_RESERVE
    c.asm.Mov(RegOp(R9), ImmOp(0x01), 64)     // PAGE_NOACCESS
    c.emitWindowsCall("VirtualAlloc")
    c.asm.Push(RAX)
    
    // VirtualAlloc(base+1MB-64KB, 64KB, MEM_COMMIT, PAGE_READWRITE)
    c.asm.Mov(RegOp(RCX), RegOp(RAX), 64)
    c.asm.Add(RegOp(RCX), ImmOp((1<<20)-(64<<10)))
    c.asm.Mov(RegOp(RDX), ImmOp(64<<10), 64)
    c.asm.Mov(RegOp(R8), ImmOp(0x1000), 64)
    c.asm.Mov(RegOp(R9), ImmOp(0x04), 64)
    c.emitWindowsCall("VirtualAlloc")
    
    // CreateThread(NULL, 0, start_routine, arg, STACK_SIZE_PARAM_IS_A_RESERVATION, NULL)
    c.asm.Xor(RegOp(RCX), RegOp(RCX))  // lpThreadAttributes
    c.asm.Xor(RegOp(RDX), RegOp(RDX))  // dwStackSize
    c.load(R8, inst.Func)               // lpStartAddress
    c.load(R9, inst.Arg)                // lpParameter
    c.asm.Sub(RegOp(RSP), ImmOp(32))   // shadow space
    c.asm.Mov(NewMem(RSP, 32), ImmOp(0x00010000), 64)  // STACK_SIZE_PARAM_IS_A_RESERVATION
    c.asm.Mov(NewMem(RSP, 40), ImmOp(0), 64)  // lpThreadId
    c.emitWindowsCall("CreateThread")
    c.asm.Add(RegOp(RSP), ImmOp(32))
    
    c.store(RAX, inst)
    return nil
}

// 3. Synchronizer - WaitOnAddress
func (c *compiler) compileAwaitWindows(inst *ir.AwaitInst) error {
    c.load(RCX, inst.Handle)
    
    waitLoop := c.asm.NewLabel()
    c.asm.Label(waitLoop)
    
    // Check if ready
    c.asm.Mov(RegOp(RAX), NewMem(RCX, 0), 32)
    c.asm.Test(RAX, RAX)
    doneLabel := c.asm.JccRel(CondNE, 0)
    
    // WaitOnAddress(&result.ready, &expected, 4, INFINITE)
    c.asm.Push(RCX)  // save result ptr
    c.asm.Sub(RegOp(RSP), ImmOp(32+8))  // shadow space + expected
    c.asm.Mov(NewMem(RSP, 32), ImmOp(0), 32)  // expected = 0
    
    c.asm.Lea(RegOp(RDX), NewMem(RSP, 32))  // &expected
    c.asm.Mov(RegOp(R8), ImmOp(4), 64)      // sizeof(int)
    c.asm.Mov(RegOp(R9), ImmOp(0xFFFFFFFF), 64)  // INFINITE
    c.emitWindowsCall("WaitOnAddress")
    
    c.asm.Add(RegOp(RSP), ImmOp(32+8))
    c.asm.Pop(RCX)
    c.asm.JmpRel(waitLoop)
    
    c.asm.Label(doneLabel)
    c.asm.Mov(RegOp(RAX), NewMem(RCX, 8), 64)
    c.store(RAX, inst)
    return nil
}
```

---

## System Comparison Table

### Stack Growth Mechanism

| Platform | Reserve | Commit | Fault Type | Handler Registration | Grow Call |
|:---------|:--------|:-------|:-----------|:--------------------|:----------|
| **Linux** | `mmap` + `PROT_NONE` | `mprotect` + `RW` | `SIGSEGV` signal | `rt_sigaction()` | `mprotect()` |
| **macOS** | `mmap` + `PROT_NONE` | `mprotect` + `RW` | Mach exception | `task_set_exception_ports()` | `vm_protect()` |
| **Windows** | `VirtualAlloc` + `MEM_RESERVE` | `VirtualAlloc` + `MEM_COMMIT` | Guard page exception | `AddVectoredExceptionHandler()` | `VirtualAlloc` + `MEM_COMMIT` |

### Thread Creation

| Platform | API | Stack Setup | Result |
|:---------|:----|:------------|:-------|
| **Linux** | `clone3()` syscall | Pass in `clone_args.stack` | Thread ID (TID) |
| **macOS** | `pthread_create()` | `pthread_attr_setstack()` | `pthread_t` handle |
| **Windows** | `CreateThread()` | `STACK_SIZE_PARAM_IS_A_RESERVATION` flag | Thread handle |

### Thread Synchronization

| Platform | Sleep Call | Wake Call | Type |
|:---------|:-----------|:----------|:-----|
| **Linux** | `futex(FUTEX_WAIT)` | `futex(FUTEX_WAKE)` | Syscall (202) |
| **macOS** | `__ulock_wait()` | `__ulock_wake()` | Syscall (0x2000203/4) |
| **Windows** | `WaitOnAddress()` | `WakeByAddressSingle()` | DLL function |

---

## Why This Design?

### 1. Real Threads vs Green Threads

**Green threads** (Go, early Java):
- Managed by runtime scheduler
- Context switch in userspace
- Blocked thread = blocked runtime

**OS threads** (Arc's approach):
- Kernel schedules fairly across cores
- Blocked thread sleeps in kernel (0% CPU)
- No runtime overhead

### 2. Auto-Growing Stacks vs Fixed Size

**Fixed size** (pthreads default = 8MB):
- 1000 threads = 8GB virtual memory
- Most never use >64KB

**Auto-growing** (Arc):
- 1000 threads = 64MB committed, 1GB virtual
- Grows only when needed

### 3. Futex vs Spin-Wait

**Spin-wait:**
```c
while (!ready) { /* burns 100% CPU */ }
```

**Futex:**
```c
futex_wait(&ready, 0);  // sleeps in kernel, 0% CPU
```

---

## Performance Characteristics

| Operation | Cost | Notes |
|:----------|:-----|:------|
| **Thread spawn** | ~5-10µs | One-time syscall overhead |
| **Stack page fault** | ~1-5µs | Only on first access to new page |
| **futex wake** | ~500ns | Fast-path if waiter in L1 cache |
| **futex wait** | ~1µs | Context switch to kernel |
| **Typical await** | ~2µs | Wake + resume |

**Scaling:**
- 1,000 threads: 64MB committed memory
- 10,000 threads: 640MB committed memory
- 100,000 threads: Likely hit OS thread limits first (~32k on Linux)

---

## Implementation Checklist

### V1: Fire-and-Forget (No await)
- [ ] `mmap`/`VirtualAlloc` stack allocation
- [ ] `mprotect`/`VirtualAlloc` initial commit
- [ ] SIGSEGV/Mach/VEH handler registration
- [ ] Handler: grow stack on page fault
- [ ] `clone3`/`pthread_create`/`CreateThread` spawner
- [ ] Thread exit cleanup

### V2: Add await
- [ ] `AsyncResult` struct allocation
- [ ] Store result on thread exit
- [ ] `futex`/`__ulock`/`WaitOnAddress` wait loop
- [ ] `futex_wake`/`__ulock_wake`/`WakeByAddressSingle` on exit
- [ ] Stack cleanup after await

### V3: Optimizations (Future)
- [ ] Per-function stack sizing (compiler analysis)
- [ ] Worker pool (reuse threads)
- [ ] Task migration scheduler

---

## Summary

Arc async smart threads combine three orthogonal systems:

1. **Stack Manager**: Auto-grows stacks via page fault handlers
2. **Thread Spawner**: Creates OS threads with custom stacks
3. **Synchronizer**: Blocks threads efficiently via futex-like primitives

All three work together:
- **`async`** keyword → Stack Manager + Thread Spawner
- **Stack overflow** → Stack Manager handler grows stack
- **`await`** keyword → Synchronizer sleeps thread until result ready

**Cross-platform** via OS-specific syscalls/APIs, but same semantics on all platforms.