# Functions

---

## Basic Function

```arc
func add(a: int32, b: int32) int32 {
    return a + b
}
```

The fundamental building block. `func` declares a named function with zero or more typed parameters and an optional return type. Parameters are always typed — there is no inference on function signatures. The return type comes after the parameter list with no arrow or colon separator. If the function returns nothing, the return type is omitted entirely. The body is always a braced block. Functions are not expressions in the general sense — they are declarations. Calling a function that returns a value without using that value is allowed but the compiler may warn depending on the type.

---

## No Return Value

```arc
func log(msg: string) {
    io.print(msg)
}

func reset(self s: Server) {
    s.port = 0
    s.name = ""
}
```

When a function produces no output, the return type is simply omitted from the signature. There is no `void` keyword in regular arc code — absence of a return type means the function returns nothing. This keeps signatures clean and removes a word that carries no information. The function may still use `return` as a bare statement to exit early, it just does not return a value.

---

## Return Tuple

```arc
func divide(a: int32, b: int32) (int32, bool) {
    if b == 0 {
        return (0, false)
    }
    return (a / b, true)
}

let (result, ok) = divide(10, 2)
```

A function can return multiple values by wrapping them in a tuple return type. The types are listed in parentheses separated by commas. At the callsite the tuple is immediately destructured into individual named variables using the same parenthesized syntax on the left of the assignment. This is the standard arc pattern for returning a result alongside a success flag, an error value, or any secondary output — it avoids output parameters, avoids single-field wrapper types, and keeps the intent explicit at both the declaration and the callsite. Tuples exist for return values — they are not a general-purpose first-class type you store and pass around.

---

## Pass by Value

```arc
func double(x: int32) int32 {
    return x * 2
}

let n = 10
let result = double(n)
// n is still 10
```

All function parameters are pass-by-value by default. When you pass a variable to a function, the function receives its own independent copy. Whatever the function does to that copy has no effect on the caller's variable. This is true for all types — integers, structs, strings. There are no hidden references, no implicit sharing. The caller's state is completely isolated from what happens inside the function. This is the safe default and covers the vast majority of function calls. When you need the function to modify the caller's variable, you opt in explicitly with a mutable reference.

---

## Mutable Reference Parameter

```arc
func increment(x: &mut int32) {
    x += 1
}

func reset(p: &mut Point) {
    p.x = 0
    p.y = 0
}

let n = 10
increment(&n)
// n is now 11
```

A parameter typed as `&mut T` receives a mutable reference to the caller's variable rather than a copy. The `&mut` in the parameter type is the declaration that this function will modify the caller's value. The `&` at the callsite is the acknowledgment that the caller is aware and consenting to mutation. Both sides must be explicit — a function cannot silently capture a reference, and a caller cannot accidentally pass one. There is no immutable reference type in arc — `&` exclusively means mutable. If you do not want mutation, pass by value. The asymmetry is intentional: reading is free and implicit, mutation requires ceremony.

---

## Method — Self Pattern

```arc
interface Client {
    port: int32
    name: string
}

func connect(self c: Client, host: string) bool {
    // c is the receiver
    return true
}

var c = Client{port: 8080, name: "test"}
c.connect("localhost")
```

A method is a regular function whose first parameter is `self`. The `self` keyword binds the function to the interface type of that parameter, enabling dot-notation calls on instances of that type. The name after `self` is the local variable name for the receiver inside the function body — by convention a short abbreviation of the type name. Methods are not declared inside the interface block — they are declared as top-level functions that happen to have a self receiver. This means methods and interfaces are decoupled: you can add methods to a type anywhere in the package without touching the type declaration.

---

## Method — Mutable Self

```arc
func grow(self &mut v: vector[int32], amount: int32) {
    v.capacity += amount
}

var items = vector[int32]{}
grow(&items, 64)
// or
items.grow(64)
```

When a method needs to modify the receiver itself rather than just reading from it, the self parameter is declared as `&mut`. This follows the same mutable reference rules as regular parameters — the mutation is explicit in the declaration and explicit at the callsite. Dot-notation calls handle the `&` automatically when the receiver is `var`, since the compiler knows the variable is heap-allocated and can take a reference to it. For `let` receivers the callsite `&` is required explicitly.

---

## deinit

```arc
interface Connection {
    fd:   int32
    host: string
}

deinit(self c: Connection) {
    close(c.fd)
}

var conn = Connection{fd: 5, host: "localhost"}
// when ref count hits 0, close(fd) fires automatically
```

`deinit` is a special lifecycle function that fires automatically when a `var` declaration's reference count drops to zero. It is declared with the same self-receiver syntax as a method but it is not called manually — ever. The compiler inserts the call at the exact point the last reference drops. `deinit` is the right place for any cleanup that must happen when the value is no longer used — closing file descriptors, releasing GPU handles, unregistering callbacks, flushing buffers. `deinit` only fires for `var`. It never fires for `let`, `const`, or `new`. A type can only have one `deinit`.

---

## Async Function

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

An `async` function runs on a kernel thread and may suspend at `await` points without blocking that thread. Marking a function `async` enables the use of `await` inside its body and signals to callers that the function is non-blocking. The return type is the actual value type — not a wrapped future or promise type. The compiler handles the suspension and resumption machinery invisibly. Async functions are called the same way as regular functions, with `await` at the callsite when the caller also needs to suspend until the result is ready.

---

## Await

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

`await` suspends the current async function until the awaited call completes, then resumes with the result. It can appear anywhere an expression is valid inside an async function — in assignments, conditions, return statements, or as a standalone call. `await` is only valid inside an `async` function — using it in a non-async context is a compile error. Multiple awaits in sequence execute one after another. Awaiting a function that is not async is a compile error — the type system enforces that you only await things that can actually suspend.

---

## Async Lambda

```arc
onClick(args, async (item: string) => {
    await process(item)
})

button.on_click(async () => {
    await save_state()
})

router.handle("/api/data", async (req: Request, res: Response) => {
    let data = await db.query("SELECT * FROM records")
    await res.send_json(data)
})
```

An anonymous function prefixed with `async` that can use `await` in its body. Async lambdas follow the same rules as async functions — they run on kernel threads and may suspend at await points. The `async` prefix is what enables `await` inside the body, not what determines whether the lambda runs asynchronously — all lambdas passed as callbacks run on kernel threads regardless. Omitting `async` from a callback lambda is ergonomic shorthand for when the body does not need to suspend. Both forms are valid and both run on kernel threads.

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
```

An anonymous function without the `async` prefix. Can be assigned to a function-typed field, passed as a callback parameter, or used as an event handler. Sync lambdas cannot use `await` — the compiler enforces this. They are the right choice for callbacks that do purely synchronous work: updating state, logging, simple transformations. The syntax is `(params) => { body }`. When the lambda body is assigned to an interface field the field type must match the lambda's parameter and return types exactly.

---

## Function as Interface Field

```arc
interface Button {
    onClick: func(int32, int32) void
}

interface TcpServer {
    port:      int32
    onReceive: async func([]byte) void
}
```

Function types can be stored as fields in an interface, making them first-class members of a type. The field type declares the full signature — parameter types, return type, and whether the callable is async. Assigning a lambda or named function to the field is valid as long as the signature matches. The field can be null when unset and should be checked before calling. This is the arc pattern for callbacks, event handlers, and any design where behavior needs to be swapped or injected at runtime without virtual dispatch.

---

## Calling a Function Field

```arc
func press(self b: Button, x: int32, y: int32) {
    if b.onClick != null {
        b.onClick(x, y)
    }
}

func handle_data(self s: TcpServer, data: []byte) {
    if s.onReceive != null {
        s.onReceive(data)
    }
}
```

A function stored in an interface field is called with the same dot-notation and argument syntax as a regular function call. Because function fields can be null when unset, a null check before calling is the standard pattern. The compiler does not automatically guard the call — the check is your responsibility. For async function fields the call site does not need `await` unless the caller itself is async and needs to suspend until the handler completes. If the caller does not await, the handler runs on a kernel thread and the caller continues immediately.

---

## Process — Concurrent Execution

```arc
let handle = process func(x: int32) {
    work(x)
}(1000)

process func() {
    heavy_computation()
}()
```

`process` spawns a function as an independent concurrent unit on its own kernel thread. The function runs concurrently with the caller — execution of the caller does not pause. The spawned function captures its arguments at the point of the call. `process` is the right tool for fire-and-forget background work, parallel pipelines, and anything that should run independently without the caller waiting. It is distinct from `async`/`await` — `async` is for suspension and resumption within a cooperative flow, `process` is for launching independent concurrent work.

---

## Generic Function

```arc
func swap[T](a: &mut T, b: &mut T) {
    let tmp: T = a
    a = b
    b = tmp
}

func find[T](arr: vector[T], val: T) isize {
    for let i: usize = 0; i < arr.len; i++ {
        if arr[i] == val {
            return isize(i)
        }
    }
    return -1
}

let x = 10
let y = 20
swap(&x, &y)
```

A function parameterized over one or more types, declared by listing type parameters in square brackets after the function name. The type parameter acts as a placeholder that is resolved at each callsite — the compiler generates a concrete version of the function for each distinct type it is called with. Generic functions allow a single implementation to work correctly across multiple types without runtime overhead or dynamic dispatch. Type parameters are inferred from the argument types at the callsite — you rarely need to provide them explicitly. Constraints on type parameters are not yet part of the syntax — a generic function currently accepts any type.

---

## GPU Function

```arc
gpu func kernel(data: float32, n: usize) {
    let idx   = thread_id()
    data[idx] = data[idx] * 2.0
}

async func main() {
    let result = await kernel(data, n)
}
```

A function that executes on an accelerator — GPU, NPU, or other compute target. The `gpu` prefix tells the compiler to compile this function for the build target's accelerator backend, which is set in the build configuration (CUDA, Metal, XLA, or others). All parameters are implicitly GPU-bound — the compiler handles data transfer between host and device memory. Inside a `gpu func`, `thread_id()` returns the index of the current parallel execution lane. GPU functions are called with `await` from async host code — the call dispatches the kernel to the accelerator and suspends the caller until the result is ready. GPU functions cannot call regular arc functions and cannot use `extern` blocks.