## Control Flow, try-except
```arc
// Basic try-except block
try {
    let result = divide(10, 0)
    io.printf("Result: %d\n", result)
} except err {
    io.printf("Error: %s\n", err)
}

// Try-except with typed errors
enum FileError {
    NotFound
    PermissionDenied
    IOError
}

func read_file(path: string) string throws FileError {
    if !exists(path) {
        throw FileError.NotFound
    }
    return contents
}

try {
    let data = read_file("/tmp/config.txt")
    process(data)
} except FileError.NotFound {
    io.printf("File not found\n")
} except FileError.PermissionDenied {
    io.printf("Permission denied\n")
} except err {
    // Catch-all for other errors
    io.printf("Unexpected error\n")
}

// Try-except with finally (always executes)
try {
    let file = open("data.txt")
    process(file)
} except err {
    io.printf("Error: %s\n", err)
} finally {
    // Cleanup code always runs
    cleanup()
}

// Try-except with ref-counted classes (automatic cleanup)
try {
    let client = HttpClient{}  // Ref count = 1
    let data = client.fetch("https://api.example.com")
    process(data)
} except err {
    // If exception thrown, client is automatically released
    io.printf("Failed to fetch: %s\n", err)
}
// client ref count decremented here (or during exception unwinding)
```

**Machine Code Implementation:**

When the parser encounters `try { ... } except { ... }`, codegen generates:

1. **Setup Phase (try keyword):**
   - Save current exception handler address to exception stack
   - Generate a `.catch_label` for the except block
   - Push exception context: `{ handler: .catch_label, frame: rbp, locals: [] }`
   - Store frame pointer (rbp) for stack restoration

2. **Try Block Body:**
   - Generate normal code for try block statements
   - Track all ref-counted objects allocated (classes): register their stack locations
   - Each `let obj = MyClass{}` adds `[rbp-offset]` to the exception context's locals list

3. **Throw Statement (throw keyword):**
   - Load error value into rax
   - Call `__arc_throw(error)` runtime function which:
     - Walks exception stack to find nearest handler
     - Decrements ref counts for all registered locals in current frame
     - Restores stack pointer (rsp) and frame pointer (rbp) from exception context
     - Jumps to the `.catch_label` address with error in rax

4. **Success Path (end of try block):**
   - Pop exception context from stack
   - Clear registered locals list
   - Jump over except block to `.finally` or function continuation

5. **Except Block (except keyword):**
   - `.catch_label:` receives control with error value in rax
   - Bind error value to the except variable (e.g., `err`)
   - Generate code for except block body
   - Type-specific catches generate comparison checks before binding

6. **Finally Block (finally keyword, optional):**
   - Generate cleanup code at `.finally` label
   - Executed in both success and exception paths
   - Alternative to `defer` for exception-aware cleanup

**Stack Unwinding with Ref Counting:**
- During unwinding, the runtime walks the registered locals list
- For each ref-counted object at `[rbp-offset]`, calls `__arc_release(obj)`
- This ensures no memory leaks even when exceptions propagate through multiple frames
- Nested try blocks use an exception context stack for proper unwinding order

**Example Generated Assembly:**
```asm
; try block setup
push rbp
mov rbp, rsp
lea rax, [.catch]                ; Address of except block
call __arc_push_exception_ctx    ; Save to exception stack

; try block body
call divide
; ... code ...

; Success path
call __arc_pop_exception_ctx
jmp .done

.catch:
; Exception landed here (rax = error)
mov [rbp-8], rax                 ; Save error
call __arc_cleanup_locals        ; Release ref-counted objects
mov rdi, [rbp-8]                ; Load error
; ... except block code ...

.done:
pop rbp
ret
```
