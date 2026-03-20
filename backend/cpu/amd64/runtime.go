package amd64

import (
	"github.com/arc-language/arc-lang/builder/ir"
	"github.com/arc-language/arc-lang/builder/types"
	"github.com/arc-language/arc-lang/backend/cpu/amd64/linux/layout"
)

// Runtime handles generating the assembly for OS-specific features
type Runtime struct {
	asm *Assembler
}

func NewRuntime(asm *Assembler) *Runtime {
	return &Runtime{asm: asm}
}

// EmitInitialization injects the setup code required at the start of execution.
func (r *Runtime) EmitInitialization() {
	// Preserve registers (RDI/RSI hold argc/argv)
	r.asm.Push(RDI)
	r.asm.Push(RSI)

	// --- 1. Setup Alt Stack for Main Thread ---
	// mmap(NULL, 8192, RW, ANON|PRIVATE, -1, 0)
	r.asm.Mov(RegOp(RDI), ImmOp(0), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(8192), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(R10), ImmOp(layout.MapPrivate|layout.MapAnonymous), 64)
	r.asm.Mov(RegOp(R8), ImmOp(-1), 64)
	r.asm.Xor(RegOp(R9), RegOp(R9))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMmap), 64)
	r.asm.Syscall()

	// Construct stack_t on stack: { ss_sp, ss_flags, ss_size }
	r.asm.Sub(RegOp(RSP), ImmOp(32))
	r.asm.Mov(NewMem(RSP, 0), RegOp(RAX), 64)   // ss_sp
	r.asm.Mov(NewMem(RSP, 8), ImmOp(0), 32)     // ss_flags
	r.asm.Mov(NewMem(RSP, 16), ImmOp(8192), 64) // ss_size

	// sigaltstack(&ss, NULL)
	r.asm.Mov(RegOp(RDI), RegOp(RSP), 64)
	r.asm.Xor(RegOp(RSI), RegOp(RSI))
	r.asm.Mov(RegOp(RAX), ImmOp(131), 64) // SYS_sigaltstack
	r.asm.Syscall()

	r.asm.Add(RegOp(RSP), ImmOp(32))

	// --- 2. Install Handler (rt_sigaction) ---
	
	// Jump over the handler code so it doesn't execute during init
	handlerLabel := r.asm.JmpRel(0)
	
	// -- SIGNAL HANDLER START --
	r.asm.PatchJump(handlerLabel)
	handlerStartOffset := r.asm.Len()
	
	// Emit handler and get the offset of the 'restorer' stub
	restorerOffset := r.emitSegfaultHandler()
	
	// -- SIGNAL HANDLER END --
	
	// Resume Main execution path
	startLabel := r.asm.NewLabel()
	r.asm.PatchInt32(handlerLabel, int32(startLabel - (handlerLabel + 4)))

	// Construct struct sigaction on stack
	// struct kernel_sigaction { handler, flags, restorer, mask }
	r.asm.Sub(RegOp(RSP), ImmOp(160))

	// A. Calculate Runtime Handler Address (RIP-Relative)
	// LEA RAX, [RIP + Disp]
	r.asm.emitByte(0x48); r.asm.emitByte(0x8D); r.asm.emitByte(0x05)
	disp := int32(handlerStartOffset - (r.asm.Len() + 4))
	r.asm.emitInt32(disp)
	
	r.asm.Mov(NewMem(RSP, 0), RegOp(RAX), 64) // sa_handler

	// B. Calculate Runtime Restorer Address (RIP-Relative)
	// LEA RAX, [RIP + Disp]
	r.asm.emitByte(0x48); r.asm.emitByte(0x8D); r.asm.emitByte(0x05)
	dispRestorer := int32(restorerOffset - (r.asm.Len() + 4))
	r.asm.emitInt32(dispRestorer)
	
	r.asm.Mov(NewMem(RSP, 16), RegOp(RAX), 64) // sa_restorer

	// C. Flags (Must include SA_RESTORER 0x04000000)
	// This flag tells the kernel "Return Address is valid, don't use vDSO"
	const SaRestorer = 0x04000000
	flags := int64(layout.SaSigInfo | layout.SaOnStack | layout.SaRestart | SaRestorer)
	r.asm.Mov(NewMem(RSP, 8), ImmOp(flags), 64)
	
	r.asm.Mov(NewMem(RSP, 24), ImmOp(0), 64) // sa_mask

	// rt_sigaction(SIGSEGV, &act, NULL, 8)
	r.asm.Mov(RegOp(RDI), ImmOp(layout.SigSegv), 64)
	r.asm.Mov(RegOp(RSI), RegOp(RSP), 64)
	r.asm.Xor(RegOp(RDX), RegOp(RDX))
	r.asm.Mov(RegOp(R10), ImmOp(8), 64) // sigsetsize
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysRtSigAction), 64)
	r.asm.Syscall()

	r.asm.Add(RegOp(RSP), ImmOp(160))
	r.asm.Pop(RSI)
	r.asm.Pop(RDI)
}

// emitSegfaultHandler generates the assembly and returns the offset of the restorer stub.
func (r *Runtime) emitSegfaultHandler() int {
	// 1. ENTRY ENDBR64 (CET Compatibility)
	r.asm.emitByte(0xF3); r.asm.emitByte(0x0F); r.asm.emitByte(0x1E); r.asm.emitByte(0xFA)

	// Inputs: RDI=signum, RSI=siginfo_t*, RDX=ucontext_t*

	// 2. Get Fault Address (siginfo->si_addr)
	r.asm.Mov(RegOp(RAX), NewMem(RSI, 16), 64)

	// 3. Page Align (RAX = RAX & -4096)
	r.asm.And(RegOp(RAX), ImmOp(-4096))

	// 4. mprotect(PageAddr, 64KB, RW)
	// We map a large chunk (64KB) to cover recursion depth.
	r.asm.Mov(RegOp(RDI), RegOp(RAX), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(64*1024), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMprotect), 64)
	r.asm.Syscall()

	// 5. Check Result
	r.asm.Test(RAX, RAX)
	okLabel := r.asm.JccRel(CondEq, 0)

	// Fail: Exit(139)
	r.asm.Mov(RegOp(RDI), ImmOp(139), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysExit), 64)
	r.asm.Syscall()

	// Success: Return from Handler
	r.asm.PatchJump(okLabel)
	// CRITICAL: We MUST 'RET' to pop the kernel's return address from the stack.
	// This jumps to 'restorerOffset' below.
	r.asm.Ret()
	
	// --- Restorer Stub ---
	restorerOffset := r.asm.Len()
	
	// RESTORER ENDBR64 (CET) - The 'ret' above targets this
	r.asm.emitByte(0xF3); r.asm.emitByte(0x0F); r.asm.emitByte(0x1E); r.asm.emitByte(0xFA)
	
	r.asm.Mov(RegOp(RAX), ImmOp(15), 64) // SYS_rt_sigreturn
	r.asm.Syscall()
	
	return restorerOffset
}

// EmitAsyncTaskCreate generates code to spawn a new thread (worker).
// It supports marshalling arbitrary arguments to the child thread using a Heap Struct.
func (r *Runtime) EmitAsyncTaskCreate(fn *ir.Function, args []ir.Value, stackSize int, loader func(Register, ir.Value)) {
	// 1. Save Callee-Saved Regs
	r.asm.Push(R12)
	r.asm.Push(R13) // We use R13 to hold StackBase securely

    // Calculate total size for AsyncHandle + Arguments
    totalSize := 16 + len(args)*8

	// 2. Allocate AsyncResult + Args Struct
	r.asm.Mov(RegOp(RDI), ImmOp(0), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(int64(totalSize)), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(R10), ImmOp(layout.MapPrivate|layout.MapAnonymous), 64)
	r.asm.Mov(RegOp(R8), ImmOp(-1), 64)
	r.asm.Xor(RegOp(R9), RegOp(R9))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMmap), 64)
	r.asm.Syscall()
	
	// Initialize Header
	r.asm.Mov(NewMem(RAX, 0), ImmOp(0), 64) // Status = 0
	r.asm.Mov(NewMem(RAX, 8), ImmOp(0), 64) // Result = 0
	r.asm.Mov(RegOp(R12), RegOp(RAX), 64)   // Save Handle to R12

    // 2b. Write Arguments into Struct
    for i, arg := range args {
        loader(RAX, arg)
        offset := 16 + (i * 8)
        r.asm.Mov(NewMem(R12, offset), RegOp(RAX), 64)
    }

	// 3. Allocate Stack Region (1MB, Guarded)
	r.asm.Xor(RegOp(RDI), RegOp(RDI))
	r.asm.Mov(RegOp(RSI), ImmOp(1024*1024), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtNone), 64) 
	r.asm.Mov(RegOp(R10), ImmOp(layout.MapPrivate|layout.MapAnonymous|layout.MapStack), 64)
	r.asm.Mov(RegOp(R8), ImmOp(-1), 64)
	r.asm.Xor(RegOp(R9), RegOp(R9))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMmap), 64)
	r.asm.Syscall()
	
    // Save Stack Base to R13 (Preserved register)
	r.asm.Mov(RegOp(R13), RegOp(RAX), 64)

	// 4. Commit Bottom 4KB (Backing store for Alt Stack)
	r.asm.Mov(RegOp(RDI), RegOp(R13), 64) // Base from R13
	r.asm.Mov(RegOp(RSI), ImmOp(4096), 64) 
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMprotect), 64)
	r.asm.Syscall()

	// 5. Commit Top 64KB (Initial Stack)
	r.asm.Mov(RegOp(RDI), RegOp(R13), 64) // Base from R13
	r.asm.Add(RegOp(RDI), ImmOp(1024*1024 - 64*1024))
	r.asm.Mov(RegOp(RSI), ImmOp(64*1024), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMprotect), 64)
	r.asm.Syscall()

	// 6. Get TLS
	r.asm.Sub(RegOp(RSP), ImmOp(16)) 
	r.asm.Mov(RegOp(RDI), ImmOp(0x1003), 64) 
	r.asm.Mov(RegOp(RSI), RegOp(RSP), 64)    
	r.asm.Mov(RegOp(RAX), ImmOp(158), 64)    
	r.asm.Syscall()
	r.asm.Mov(RegOp(R11), NewMem(RSP, 0), 64) 
    r.asm.Add(RegOp(RSP), ImmOp(16)) // Restore stack immediately

	// 7. Clone
	r.asm.Mov(RegOp(RDI), ImmOp(0x1D0F00), 64)
	r.asm.Mov(RegOp(RSI), RegOp(R13), 64)      // Base from R13
	r.asm.Add(RegOp(RSI), ImmOp(1024*1024))    // Top
	r.asm.Lea(RDX, NewMem(RSP, -8))            // Parent TID (dummy addr)
	r.asm.Xor(RegOp(R10), RegOp(R10))
	r.asm.Mov(RegOp(R8), RegOp(R11), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(56), 64)
	r.asm.Syscall()

	r.asm.Test(RAX, RAX)
	childLabel := r.asm.JccRel(CondEq, 0)

	// --- PARENT ---
    // Restore preserved regs
    r.asm.Pop(R13)
	r.asm.Mov(RegOp(RAX), RegOp(R12), 64) // Return Handle
	r.asm.Pop(R12) 
	doneLabel := r.asm.JmpRel(0)

	// --- CHILD ---
	r.asm.PatchJump(childLabel)
    // Child inherits R12 (AsyncResult) and R13 (StackBase)
    
    // Ensure 16-byte alignment
	r.asm.And(RegOp(RSP), ImmOp(-16)) 
	
	// Child Alt Stack (New Alloc) - 8KB
	r.asm.Mov(RegOp(RDI), ImmOp(0), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(8192), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(layout.ProtRead|layout.ProtWrite), 64)
	r.asm.Mov(RegOp(R10), ImmOp(layout.MapPrivate|layout.MapAnonymous), 64)
	r.asm.Mov(RegOp(R8), ImmOp(-1), 64)
	r.asm.Xor(RegOp(R9), RegOp(R9))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysMmap), 64)
	r.asm.Syscall()
	
	// sigaltstack
	r.asm.Sub(RegOp(RSP), ImmOp(32))
	r.asm.Mov(NewMem(RSP, 0), RegOp(RAX), 64)
	r.asm.Mov(NewMem(RSP, 8), ImmOp(0), 32)
	r.asm.Mov(NewMem(RSP, 16), ImmOp(8192), 64)
	
	r.asm.Mov(RegOp(RDI), RegOp(RSP), 64)
	r.asm.Xor(RegOp(RSI), RegOp(RSI))
	r.asm.Mov(RegOp(RAX), ImmOp(131), 64)
	r.asm.Syscall()
	r.asm.Add(RegOp(RSP), ImmOp(32))
	
    // 8. Restore Arguments from Struct (R12)
    gprRegs := []Register{RDI, RSI, RDX, RCX, R8, R9}
    xmmIdxC := 0
    gprIdxC := 0
    
    for i, arg := range args {
        offset := 16 + (i * 8)
        
        if types.IsFloat(arg.Type()) {
             if xmmIdxC < 8 {
                 r.asm.MovsdLoad(Register(xmmIdxC), NewMem(R12, offset))
                 xmmIdxC++
             }
        } else {
             if gprIdxC < 6 {
                 r.asm.Mov(RegOp(gprRegs[gprIdxC]), NewMem(R12, offset), 64)
                 gprIdxC++
             }
        }
    }

	// Run User Function
	r.asm.LeaRel(RAX, fn.Name())
	r.asm.CallReg(RAX)
	
	// Result & Wake
	r.asm.Mov(NewMem(R12, 8), RegOp(RAX), 64)
	r.asm.Mov(NewMem(R12, 0), ImmOp(1), 32)
	
	r.asm.Mov(RegOp(RDI), RegOp(R12), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(layout.FutexWake|layout.FutexPrivate), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(1), 64)   
	r.asm.Xor(RegOp(R10), RegOp(R10))
	r.asm.Xor(RegOp(R8), RegOp(R8))
	r.asm.Xor(RegOp(R9), RegOp(R9))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysFutex), 64)
	r.asm.Syscall()
	
    // Free Stack (using R13)
    r.asm.Mov(RegOp(RDI), RegOp(R13), 64)
    r.asm.Mov(RegOp(RSI), ImmOp(1024*1024), 64)
    r.asm.Mov(RegOp(RAX), ImmOp(11), 64) // SYS_munmap
    r.asm.Syscall()

	r.asm.Mov(RegOp(RDI), ImmOp(0), 64)
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysExit), 64)
	r.asm.Syscall()

	r.asm.PatchJump(doneLabel)
}

func (r *Runtime) EmitAsyncTaskAwait(handle Register) {
	loop := r.asm.NewLabel()
	r.asm.Label(loop)

	r.asm.Mov(RegOp(RAX), NewMem(handle, 0), 32)
	r.asm.Test(RAX, RAX)
	done := r.asm.JccRel(CondNe, 0)
	
	r.asm.Push(handle)
	r.asm.Mov(RegOp(RDI), RegOp(handle), 64)
	r.asm.Mov(RegOp(RSI), ImmOp(layout.FutexWait|layout.FutexPrivate), 64)
	r.asm.Mov(RegOp(RDX), ImmOp(0), 64)
	r.asm.Xor(RegOp(R10), RegOp(R10))
	r.asm.Mov(RegOp(RAX), ImmOp(layout.SysFutex), 64)
	r.asm.Syscall()
	r.asm.Pop(handle)

	r.asm.JmpRel(loop)
	r.asm.PatchJump(done)
	r.asm.Mov(RegOp(RAX), NewMem(handle, 8), 64)
}


// EmitProcessCreate forks the process. 
/// EmitProcessCreate forks the process. 
func (r *Runtime) EmitProcessCreate(fn *ir.Function, argCount int) {
	// 1. Syscall Clone (Fork)
	r.asm.Mov(RegOp(RAX), ImmOp(56), 64)
	r.asm.Mov(RegOp(RDI), ImmOp(17), 64) // SIGCHLD
	r.asm.Mov(RegOp(RSI), ImmOp(0), 64)
	r.asm.Syscall()

	// 2. Check Result (RAX)
	r.asm.Test(RAX, RAX)
	parentLabel := r.asm.JccRel(CondNe, 0) // Jump if PID != 0 (Parent)

	// --- CHILD CODE ---
	
	// Restore args
	abiRegs := []Register{RDI, RSI, RDX, RCX, R8}
	preservedRegs := []Register{R12, R13, R14, R15, RBX}
	
	for i := 0; i < argCount; i++ {
		r.asm.Mov(RegOp(abiRegs[i]), RegOp(preservedRegs[i]), 64)
	}

	// Call User Function
	r.asm.CallRelative(fn.Name())

	// Exit Child CLEANLY (Flush Buffers)
	// Move return value to RDI (Exit Code)
	r.asm.Mov(RegOp(RDI), RegOp(RAX), 64) 
	
	// FIX: Call libc 'exit' to flush printf buffers
	r.asm.CallRelative("exit")
	
	// Safety fallback (should never be reached)
	r.asm.Mov(RegOp(RAX), ImmOp(60), 64)
	r.asm.Syscall()

	// --- PARENT CODE ---
	r.asm.PatchJump(parentLabel)
}