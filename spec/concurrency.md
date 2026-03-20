# Concurrency

---

## Kernel Threads

Arc's concurrency model is built directly on kernel threads — not green threads, not a
userspace scheduler, not a runtime that multiplexes virtual threads onto OS threads the
way Go's goroutine scheduler does. Every concurrent unit in arc is a real kernel thread.

These are not traditional heavyweight OS threads. They are extremely lightweight kernel
threads that start small and grow their stack automatically on demand. You can run
thousands of them concurrently without the memory overhead associated with conventional
threads. The pool expands as demand increases and the OS handles all scheduling — you
never manage threads manually.

This design is driven by arc's goals beyond general-purpose concurrency. Multi-GPU
workloads, hardware accelerator pipelines, and direct device coordination require that
threads are real kernel entities. A userspace scheduler sitting between arc and the
kernel would either block hal-level calls, produce incorrect device synchronization, or
crash under concurrent device access patterns that are normal in GPU compute. By going
directly to kernel threads arc gets correct behavior under hardware concurrency and the
ability to coordinate multiple GPUs and accelerators simultaneously without fighting a
runtime that was never designed for it.

---

## async func

```arc
async func fetch(url: string) string {
    let response = await http.get(url)
    return response.body
}

async func save(data: string) bool {
    let ok = await db.write(data)
    return ok
}
```

A function that runs on a kernel thread and may suspend at `await` points without
blocking that thread. The `async` prefix enables `await` inside the body and signals
to callers that this function is non-blocking. The return type is the actual value type
— not a wrapped future, promise, or handle. The compiler generates the suspension and
resumption machinery invisibly. When an async function suspends at an `await`, the
kernel thread is returned to the pool and picks up other work. When the awaited result
is ready the function resumes on whatever thread is available. Because these are real
kernel threads, suspension and resumption involve no userspace scheduler — the OS
handles it directly. Async functions are the standard form for any work that involves
waiting on I/O, devices, or other async results.

---

## await

```arc
async func main() {
    let data    = await fetch("https://api.example.com")
    let result1 = await task1()
    let result2 = await task2()

    if await check_status() {
        // proceed
    }
}
```

Suspends the current async function until the awaited call completes, then resumes with
the result. `await` can appear anywhere an expression is valid inside an async function
— assignments, conditions, return statements, function arguments. Multiple awaits in
sequence execute one after another, each suspending until its call resolves before the
next begins. `await` is only valid inside an `async` function — using it in a non-async
context is a compile error. The kernel thread is not blocked during suspension — it is
returned to the pool immediately and can pick up other work. This is what makes arc's
concurrency safe for hardware-bound workloads: the kernel sees real thread suspension
and resumption rather than a userspace scheduler trying to fake it. When the awaited
call completes the function is resumed, potentially on a different kernel thread from
the one that suspended it.

---

## process

```arc
let handle = process func(x: int32) {
    work(x)
}(1000)

process func() {
    heavy_computation()
}()
```

Spawns a function as an independent concurrent unit on its own kernel thread. The
spawned function runs concurrently with the caller — the caller does not pause and does
not wait for the spawned function to complete. Arguments are captured and copied at the
point of the call. `process` is fire-and-forget: the work happens independently on a
real kernel thread and the caller moves on immediately. Because the spawned unit is a
real kernel thread, it can safely drive hardware, coordinate with devices, and run
alongside GPU kernels without any userspace scheduler interfering. Spawning thousands
of process units is practical given how lightweight these threads are. This is the right
tool for background work, parallel pipelines, independent computations, and anything
that should run concurrently without the caller caring about the result.

---

## process vs async

```arc
// async — caller suspends, waits for result
async func main() {
    let result = await fetch(url)   // main suspends until fetch completes
    process_result(result)
}

// process — caller continues immediately, work runs independently
func main() {
    process func() {
        heavy_work()                // runs on its own kernel thread
    }()
    do_other_things()               // caller continues immediately
}
```

`async`/`await` and `process` serve different concurrency needs. `async`/`await` is
cooperative suspension — you pause at a known point, wait for a specific result, and
resume. The flow reads sequentially even though the kernel thread is not blocked.
`process` is parallel dispatch — you launch work onto a kernel thread and immediately
continue. There is no waiting, no result handed back, no suspension. Use `async`/`await`
when the caller needs the result before it can proceed. Use `process` when the work is
independent and the caller has no reason to coordinate with it. Both models run on real
kernel threads — the difference is purely about whether the caller waits.

---

## Async Lambda

```arc
button.on_click(async (item: string) => {
    let result = await process(item)
    await store.save(result)
})

router.handle("/api/data", async (req: Request, res: Response) => {
    let data = await db.query("SELECT * FROM records")
    await res.send_json(data)
})

stream.onData = async (chunk: []byte) => {
    let validated = await validate(chunk)
    await buffer.write(validated)
}
```

An anonymous function prefixed with `async` that can use `await` in its body. Async
lambdas run on kernel threads and may suspend at await points, returning the thread to
the pool during suspension exactly as a named async function does. The `async` prefix
is what enables `await` inside the body — it is not what determines whether the lambda
runs on a kernel thread. All lambdas passed as callbacks run on kernel threads regardless
of whether they are marked async. Omitting `async` is ergonomic shorthand for callbacks
whose body does not need to suspend. Both forms are valid and both are real kernel thread
work. Async lambdas are the standard pattern for I/O callbacks, route handlers, event
processors, and any reactive code where the handler itself needs to perform async work.

---

## Sync Lambda

```arc
button.onClick = (x: int32, y: int32) => {
    fmt.printf("Clicked at (%d, %d)\n", x, y)
}

service.request(args, (result: Result) => {
    fmt.print("done")
    handle_result(result)
})

items.each((item: int32) => {
    process_sync(item)
})
```

An anonymous function without the `async` prefix. Cannot use `await` in its body — the
compiler enforces this. Sync lambdas run on kernel threads like all callbacks but they
run to completion without suspending — the thread stays with the lambda until the body
finishes. They are the right choice for callbacks that do purely synchronous work:
transformations, state updates, logging, simple computations. For callbacks that need
to call async functions or wait on I/O or devices, use an async lambda instead. The
distinction in the source is a single word — `async` — which makes the presence or
absence of suspension immediately visible at the callsite.

---

## Event Handler — Property Assignment

```arc
interface Server {
    port:      int32
    onConnect: async func(Connection) void
    onData:    async func([]byte) void
    onClose:   func() void
}

var server = Server{port: 8080}

server.onConnect = async (conn: Connection) => {
    let data = await conn.read()
    await process(data)
}

server.onData = async (chunk: []byte) => {
    await buffer.append(chunk)
}

server.onClose = () => {
    log("server closed")
}
```

Event handlers are function fields on an interface assigned at the callsite. The field
type declares the expected signature including whether the handler is async. Assignment
uses a lambda that matches that signature exactly. This is the standard arc pattern for
event-driven code — the type declares what events it produces, the consumer assigns
handlers for the ones it cares about. Unassigned handlers remain null and should be
checked before calling. There is no registration system, no observer list, no virtual
dispatch — just a function stored in a field, called when the event fires, running on
its own kernel thread.

---

## Callback Parameter

```arc
service.request(args, async (result: Result) => {
    let processed = await transform(result)
    fmt.print(processed.data)
})

network.fetch(url, timeout, (response: Response) => {
    fmt.printf("status: %d\n", response.status)
})

onClick(args, async (item: string) => {
    await process(item)
})
```

A lambda passed directly as an argument to a function call. The function receiving the
callback declares the expected function type as a parameter — the compiler checks that
the lambda's signature matches. Callback parameters are the standard pattern for one-off
handlers where storing the function in a field would be unnecessarily stateful. The
lambda captures its surrounding scope at the point of the call. Async callback parameters
can use `await` in their body. Sync callback parameters cannot. The distinction is always
visible at the callsite from the presence or absence of `async` before the lambda. Each
callback runs on its own kernel thread — passing thousands of callbacks is practical
given the lightweight nature of arc's threading model.

---

## Indirect Async Callback

```arc
interface Event {
    onTrigger: async func(string) void
}

func register(self evt: Event, handler: async func(string) void) {
    evt.onTrigger = handler
}

func send(self evt: Event, message: string) {
    evt.onTrigger(message)
}

var evt = Event{}

evt.register(async (msg: string) => {
    let processed = await process_message(msg)
    fmt.printf("processed: %s\n", processed)
})

evt.send("hello")
```

A function that accepts an async function type as a parameter and stores it for later
invocation. The parameter type `async func(string) void` is the full signature of the
expected handler including the async qualifier. The stored handler is called later
through the interface field using normal dot-notation call syntax. The caller of `send`
does not await the handler — it fires and the handler runs on its own kernel thread
independently. This pattern decouples the producer of events from the consumer without
any framework machinery — just a field, an assignment, and a call.

---

## gpu func

```arc
gpu func kernel(data: float32, n: usize) {
    let idx   = thread_id()
    data[idx] = data[idx] * 2.0
}

gpu func matrix_multiply(a: float32, b: float32, out: float32, n: usize) {
    let row = thread_id()
    let sum = 0.0
    for let i = 0; i < n; i++ {
        sum += a[row * n + i] * b[i]
    }
    out[row] = sum
}

async func main() {
    let result = await kernel(data, n)
}
```

A function that executes on an accelerator — GPU, NPU, or other compute target set in
the build configuration. The `gpu` prefix tells the compiler to compile this function
for the target accelerator backend: CUDA, Metal, XLA, or others. All parameters are
implicitly device-bound — the compiler handles data transfer between host and device
memory automatically. Inside a `gpu func`, `thread_id()` returns the index of the
current parallel execution lane. Each lane runs the same function body independently
over its own index. GPU functions are called with `await` from async host code — the
call dispatches the kernel to the accelerator and suspends the calling kernel thread
until all lanes complete and results are ready. Because arc uses real kernel threads,
multiple GPU kernels can be coordinated simultaneously across multiple devices without
a userspace scheduler creating contention or deadlock at the device boundary. GPU
functions cannot call regular arc functions and cannot use `extern` blocks.