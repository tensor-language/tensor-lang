# arc Language Grammar (Version 2.4 - Core Syntax)

Grammar rules:
 * `interface` is the single type declaration — no struct/class split
 * memory model is decided at the callsite, not the declaration
 * `const` = immutable value, stack, never changes
 * `let` = stack allocated, copied on assignment
 * `var` = heap allocated, ref counted, shareable, nullable
 * `new` = heap allocated, manual, no ref count, you own it
 * `&mut` in parameter type signals a mutable reference
 * `&` at call site passes a mutable reference (matches &mut param)
 * passing without `&` is always pass-by-value (caller's variable unaffected)
 * `&` exclusively means mutable reference — there is no immutable reference type
 * `[]T` is a slice view — ptr + length, no allocation
 * `[N]T` is a fixed-size array
 * `vector[T]` is an owned dynamic array — heap allocated
 * `buf[0..4]` produces a slice view, `buf[0]` produces a single element
 * `memptr(&val)` gets address of a variable, `memptr(x)` casts a value to pointer
 * null valid for var declarations only, compiler enforces this
 * cast uses type(value) syntax, not cast<T>(value)
 * gpu func for accelerator kernels, target set in build config

Grammar to not add to the parser files:
    * Empty initializer

# Empty initializer
You usually disallow "empty" initializer lists as statements. An array literal {1, 2, 3} on its own line does nothing anyway.
```arc
{ 
    print("data") 
}
```

## Comments
```arc
// Single-line comment

/*
   Multi-line comment
   spanning multiple lines
*/
```

## Import, declaration
```arc
import "some/path/package"

import yourpackage "some/path/package"

import (
    // Standard Lib
    "std/io"
    "std/net"

    // Third Party (Explicit URLs)
    "github.com/user/physics"
    "github.com/user/graphics"
    
    // Other Hosts
    "gitlab.com/company/lib"
)
```

## Namespaces, declaration
```arc
namespace main
```

## Memory Model

Arc has four declaration keywords that tell you everything about ownership and lifetime.
The type declaration never carries the memory decision — the callsite does.
```arc
const MAX = 100                           // immutable, stack, never changes
let point = Point{x: 1, y: 2}            // stack allocated, copied on assignment
var client = Client{name: "x", port: 80}  // heap allocated, ref counted, shareable
let node = new Node{}                     // heap allocated, manual, you own it
```

- `const` — value never changes, compiler can inline and optimize freely
- `let` — lives on the stack, when scope ends it's gone, assignment makes a copy
- `var` — lives on the heap, ref counted, safe to share across scopes, can be null
- `new` — lives on the heap, no ref count, no deinit, you manage it manually

The difference between `var` and `new` is ownership:
```arc
var node = Node{}       // heap, ref counted — compiler manages lifetime
let node = new Node{}   // heap, manual — you manage lifetime
```

## Variables, stack with type
```arc
let x: int32 = 42
x = 100
```

Stack allocated. When this scope ends, x is gone. Assignment copies the value.

## Variables, stack with inference
```arc
let x = 42
x = 100
```

## Variables, heap ref counted with type
```arc
var client: Client = Client{name: "test", port: 8080}
```

Heap allocated, ref counted. When all references drop, deinit fires and memory is freed.
Safe to pass around, safe to store in other types. Can be assigned null.

## Variables, heap ref counted with inference
```arc
var client = Client{name: "test", port: 8080}
```

## Variables, nullable
```arc
var client: Client = null
```

Only valid for `var` declarations. The compiler will reject `let x: int32 = null` —
null only makes sense for heap allocated ref counted types where ownership is shared.

## Constants, immutable with type
```arc
const x: int32 = 42
```

## Constants, immutable with inference
```arc
const x = 42
```

## new — Manual Heap Allocation

`new` allocates on the heap with no ref counting and no deinit.
You own it. You free it. Designed for kernel drivers, hot paths, and
anywhere `var` ref counting is too heavy.
```arc
// allocate single type, zero initialized
let node = new Node{}

// allocate with field initialization
let node = new Node{value: 10, next: null}

// allocate fixed-size array
let buf = new [4096]byte
let ids = new [256]uint32

// always pair with defer delete
let node = new Node{}
defer delete(node)

let buf = new [4096]byte
defer delete(buf)
```

- `new` never triggers `deinit` — that is only for `var`
- `delete` must be called manually — use `defer` to avoid leaks
- Safe to use in interrupt handlers and spinlock-held sections
- No ref count overhead — safe in hot paths

## new, common patterns
```arc
// alloc fixed buffer then zero
let buf = new [4096]byte
mem_zero(buf, sizeof(buf))
defer delete(buf)

// alloc struct for C interop
let addr = new SockAddrIn{}
mem_zero(addr, sizeof(SockAddrIn))
defer delete(addr)
connect(fd, memptr(&addr), uint32(sizeof(SockAddrIn)))

// alloc then use then free
let node = new Node{value: 42}
defer delete(node)
process(node)
```

## Array and Slice Types

Arc has three distinct sequence types:
```arc
// []T — slice view, ptr + length, no allocation
// a window into existing memory, never owns it
let view: []byte = ...
let header: []byte = buffer[0..4]

// [N]T — fixed-size array, size known at compile time
let buf: [4096]byte = ...
let ids: [256]uint32 = ...

// vector[T] — owned dynamic array, heap allocated, growable
let owned: vector[byte] = ...
let items: vector[int32] = {1, 2, 3, 4, 5}
```

Fixed-size arrays on the stack:
```arc
// stack allocated fixed array
let buf: [1024]byte
mem_zero(buf, sizeof(buf))

// heap allocated fixed array — manual
let buf = new [4096]byte
defer delete(buf)
```

Indexing rules:
```arc
let buffer: [4096]byte = ...

let chunk: []byte  = buffer[0..4]  // range index — produces a slice view
let b:     byte    = buffer[2]     // single index — produces a single element
```

Slices work with any type:
```arc
let ints: vector[int32] = {1, 2, 3, 4, 5}
let middle: []int32 = ints[1..4]   // view of elements 1, 2, 3

func process(data: []byte) {
    // read-only view, no copy, no allocation
}
```

## Basic Types, Fixed-Width Integers
```arc
// Signed: int8, int16, int32, int64
let i: int32 = -500

// Unsigned: uint8, uint16, uint32, uint64
let u: uint64 = 10000
```

## Basic Types, Architecture Dependent
```arc
// Unsigned pointer-sized integer (x64 = uint64, x86 = uint32)
// Use for: Array indexing, Memory sizes, Loop counters
let len: usize = 100

// Signed pointer-sized integer (x64 = int64, x86 = int32)
// Use for: Pointer offsets, C functions returning ssize_t
let offset: isize = -4
```

## Basic Types, Floating Point
```arc
let f32: float32 = 3.14
let f64: float64 = 2.71828
```

## Basic Types, Aliases (Semantic)
```arc
// 'byte' is an alias for 'uint8' (Raw data)
let b: byte = 255

// 'bool' is an alias for 'uint8' (1=true, 0=false)
let flag: bool = true

// 'char' is an alias for 'uint32' (Unicode Code Point)
let r: char = 'a'
```

## Basic Types, String (Composite)
```arc
// High-level string (ptr + length)
let s: string = "hello"
```

## Basic Types, Qualified (Namespaced)
```arc
let client: net.Socket = ...
let config: json.Config = ...
```

## Literals, boolean
```arc
let flag: bool = true
let enabled: bool = false
```

## Literals, null
```arc
// null is only valid for var declarations (heap ref counted types)
// compiler error if assigned to let or const
var client: Client = null
var server: net.Server = null

if client == null {
    // handle null case
}
```

## Literals, character
```arc
let ch: char = 'a'
let digit: char = '5'
```

## Literals, character escapes
```arc
let newline: char = '\n'
let tab: char = '\t'
let backslash: char = '\\'
let quote: char = '\''
let null_char: char = '\0'
```

## Literals, string
```arc
let msg: string = "Hello, World!"
let empty: string = ""
```

## Literals, string escapes
```arc
let msg: string = "Hello\nWorld"
let path: string = "C:\\Users\\file"
let quote: string = "He said \"hello\""
let tab: string = "Column1\tColumn2"
```

## Functions, basic
```arc
func add(a: int32, b: int32) int32 {
    return a + b
}
```

## Functions, no return
```arc
func print(msg: string) {
    
}
```

## Functions, async
```arc
async func fetch_data(url: string) string {
    let response = await http.get(url)
    return response.body
}

async func process_items(items: vector[string]) {
    for item in items {
        await process(item)
    }
}
```

## Functions, await usage
```arc
async func main() {
    let data = await fetch_data("https://api.example.com")
    
    let result1 = await task1()
    let result2 = await task2()
    
    if await check_status() {
        // do something
    }
}
```

## Functions, async callback
```arc
onClick(args, async (item: string) => {
    await process(item)
})

some.fetch(args, async (url: string, timeout: int32) => {
    let resp = await http.get(url, timeout)
    return resp.body
})

button.on_click(async () => {
    await save_state()
})
```

## Functions, gpu
```arc
// All params are gpu bound, compiler maps to build target
gpu func kernel(data: float32, n: usize) {
    let idx = thread_id()
    data[idx] = data[idx] * 2.0
}

async func main() {
    let result = await kernel(data, n)
}
```

## Function Return Tuples
```arc
func divide(a: int32, b: int32) (int32, bool) {
    if b == 0 {
        return (0, false)
    }
    return (a / b, true)
}

let (result, ok) = divide(10, 2)
```

## Methods, self pattern declaration
```arc
interface Client {
    port: int32
}

// 'self' keyword binds the function to the interface as a method
// allows dot notation on the instance: c.connect("localhost")
func connect(self c: Client, host: string) bool {
    return true
}
```

## Methods, usage
```arc
var c = Client{port: 8080}
c.connect("localhost")

async func example() {
    let data = await c.fetch_data()
}
```

## deinit

`deinit` is called automatically when a `var` declaration's ref count reaches zero.
It is never called for `let` or `new` declarations.
```arc
interface Client {
    name: string
    port: int32
}

deinit(self c: Client) {
    io.close(c.port)
    // ref count hit 0 — clean up any resources here
}

// deinit fires automatically, you never call it manually
var c = Client{name: "test", port: 8080}
// ... c goes out of scope, ref count hits 0, deinit fires

// new — deinit never fires, you call delete manually
let c = new Client{name: "test", port: 8080}
defer delete(c)
```

## Type Casting
```arc
// Primitive casts - type(value) syntax
let x = int32(3.14)     // float to int
let y = float64(42)     // int to float
let z = uint8(flags)    // narrow cast
let n = usize(count)    // to pointer-sized int
```

## memptr — Two forms
```arc
// memptr(&val) — get address of variable as memory pointer
// memptr(x)    — cast integer value to memory pointer

// 1. Get address of a variable
let x = 100
let p = memptr(&x)

// 2. Cast integer value to pointer
let transient = memptr(-1)
sqlite3_bind_text(stmt, 1, val, len, transient)

// 3. Pointer arithmetic (via usize)
let addr     = usize(some_ptr)
let next_ptr = memptr(addr + 16)

// 4. Pass address to extern output params
var db: sqlite3 = null
sqlite3_open("test.db", memptr(&db))
```

## Raw Memory Operations
```arc
// mem_zero — fill memory with zeros
mem_zero(buf, sizeof(buf))
mem_zero(addr, sizeof(SockAddrIn))

// mem_copy — copy memory, regions must not overlap
mem_copy(dest, src, sizeof(Point))

// mem_move — copy memory, overlap safe
mem_move(dest, src, 1024)

// mem_compare — compare memory regions, 0 = equal
let diff = mem_compare(a, b, sizeof(Point))
```

## Control Flow, if-else
```arc
if condition {
    
} else if condition2 {
    
} else {
    
}
```

## Control Flow, for loop (C-style)
```arc
for let i = 0; i < 10; i = i + 1 {
    // loop body
}

for let i = 0; i < 10; i++ {
    // loop body
}
```

## Control Flow, for loop (while-style)
```arc
let j = 5
for j > 0 {
    j--
}
```

## Control Flow, for loop (infinite)
```arc
let counter = 0
for {
    counter++
    if counter >= 10 {
        break
    }
}
```

## Control Flow, for-in loop (iterators)
```arc
let items: vector[int32] = {1, 2, 3, 4, 5}
for item in items {
    // use item
}

let scores: map[string]int32 = {"alice": 100, "bob": 95}
for key, value in scores {
    // use key and value
}

for i in 0..10 {
    // i goes from 0 to 9
}
```

## Control Flow, break / continue / defer / return
```arc
break
continue
defer delete(ptr)
return value
```

## Control Flow, switch
```arc
let status = 2

switch status {
    case 0:
        io.print("OK")
    case 1:
        io.print("Pending")
    case 2:
        io.print("Error")
    default:
        io.print("Unknown")
}

switch status {
    case 1, 2, 3:
        io.print("Active/Pending")
    case 4:
        io.print("Completed")
    default:
        io.print("Unknown")
}
```

## Operators, arithmetic
```arc
let sum  = a + b
let diff = a - b
let prod = a * b
let quot = a / b
let rem  = a % b
```

## Operators, bitwise
```arc
let b_or  = a | b
let b_xor = a ^ b
let b_and = a & b
let shl   = a << 2
let shr   = a >> 1
let b_not = ~a
```

## Operators, compound assignment
```arc
x += 5
x -= 3
x *= 2
x /= 4
x %= 3
```

## Operators, increment/decrement
```arc
i++
++i
i--
--i
```

## Operators, comparison
```arc
let eq = a == b
let ne = a != b
let lt = a < b
let le = a <= b
let gt = a > b
let ge = a >= b
```

## Operators, logical
```arc
let and = a && b
let or  = a || b
```

## Operators, unary
```arc
let neg = -value
let not = !flag
```

## Extern, C interoperability
```arc
// Stars are only allowed inside extern blocks
// This is the C boundary — all raw pointer details live here
extern c {
    func printf(*byte, ...) int32
    func sleep "usleep" (int32) int32
    func usleep(int32) int32
}

printf("hello\n")
```

## Enums
```arc
enum Status {
    OK
    ERROR
    PENDING
}

enum HttpCode {
    OK = 200
    NOT_FOUND = 404
    SERVER_ERROR = 500
}

enum Color: uint8 {
    RED   = 0xFF0000
    GREEN = 0x00FF00
    BLUE  = 0x0000FF
}
```

## Generics, interface
```arc
interface Box<T> {
    value: T
}

func get(self b: Box<T>) T {
    return b.value
}

func set(self b: Box<T>, val: T) {
    b.value = val
}
```

## Generics, multiple type parameters
```arc
interface Pair<K, V> {
    key: K
    value: V
}

interface Result<T, E> {
    data: T
    error: E
    success: bool
}
```

## Generics, functions
```arc
func swap<T>(a: &mut T, b: &mut T) {
    let tmp: T = a
    a = b
    b = tmp
}

let x = 10
let y = 20
swap(&x, &y)

func find<T>(arr: vector[T], val: T) isize {
    for let i: usize = 0; i < arr.len; i++ {
        if arr[i] == val {
            return isize(i)
        }
    }
    return -1
}
```

## Execution Context, process
```arc
let handle = process func(x: int32) { 
    work(x) 
}(1000)

process func() {  
    work(x)
}()
```

## Execution Context, async
```arc
async func(x: int32) { 
    work(x) 
}(1000)

async func() {  
    work(x)
}()
```

## Async Event Handlers (property assignment)
```arc
handler.onEvent = (data: EventData) => { 
    process_immediate(data)
    update_state()
}

stream.onData = (chunk: []byte) => {
    buffer.append(chunk)
    validate(chunk)
}
```

## Async Event Handlers (with await capability)
```arc
handler.onEvent = async (data: EventData) => { 
    let result = await process_async(data)
    fmt.print(result.status)
    await store.save(result)
}

stream.onData = async (chunk: []byte) => {
    let validated = await validate_async(chunk)
    await buffer.write(validated)
    await notify_listeners(validated)
}
```

## Async Method Calls (callback parameter)
```arc
service.request(args, (result: Result) => {
    fmt.print("Request completed")
    handle_result(result)
})

network.fetch(url, timeout, (response: Response) => {
    fmt.printf("Status: %d\n", response.status)
})
```

## Async Method Calls (callback with await capability)
```arc
service.request(args, async (result: Result) => {
    let processed = await transform(result)
    fmt.print(processed.data)
})

router.handle("/api/data", async (req: Request, res: Response) => {
    let data = await db.query("SELECT * FROM records")
    await res.send_json(data)
})
```

**Note:** Both forms are async (run on smart threads). The `async` keyword only determines
whether `await` is allowed inside the lambda body. Omitting `async` is ergonomic shorthand
for callbacks that don't need to suspend.

## Functions, async callback (indirect invocation)
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
    fmt.printf("Processed: %s\n", processed)
})

evt.send("Hello, World!")
```

## Functions, async callback (property assignment style)
```arc
interface TcpServer {
    port: int32
    onReceive: async func([]byte) void
}

func handle_data(self s: TcpServer, data: []byte) {
    if s.onReceive != null {
        s.onReceive(data)
    }
}

var server = TcpServer{port: 8080}

server.onReceive = async (data: []byte) => {
    let decoded = await decode_packet(data)
    await store_in_db(decoded)
}

server.handle_data(received_bytes)
```

## Functions, sync callback (indirect invocation)
```arc
interface Button {
    onClick: func(int32, int32) void
}

func press(self b: Button, x: int32, y: int32) {
    if b.onClick != null {
        b.onClick(x, y)
    }
}

var button = Button{}

button.onClick = (x: int32, y: int32) => {
    fmt.printf("Clicked at (%d, %d)\n", x, y)
}

button.press(100, 200)
```

## Quick Reference
```arc
// Memory model at a glance
const MAX = 100                            // immutable, stack
let point = Point{x: 1, y: 2}             // stack, copied
var client = Client{name: "x", port: 80}  // heap, ref counted
var client: Client = null                  // heap, nullable
let node = new Node{}                      // heap, manual, you own it
let buf  = new [4096]byte                  // heap, manual fixed array

// new and delete always paired
let node = new Node{value: 42}
defer delete(node)

let buf = new [4096]byte
defer delete(buf)

// Mutation at a glance
func read(x: int32) {}        // pass by value, caller unchanged
func mutate(x: &mut int32) {} // mutable reference, caller changes

let n = 10
read(n)                        // safe, n unchanged
mutate(&n)                     // explicit, n will change

// Low level at a glance
let p    = memptr(&val)        // address of val
let p    = memptr(-1)          // integer as pointer
let p    = memptr(addr + 16)   // pointer arithmetic
let sz   = sizeof(SockAddrIn)  // compile-time size
let al   = alignof(float64)    // compile-time alignment
let bits = bitcast(uint32, f)  // reinterpret bits
let n    = len(s)              // length, type decides cost
mem_zero(buf, sizeof(buf))     // zero memory
mem_copy(dst, src, n)          // copy memory, no overlap
mem_move(dst, src, n)          // copy memory, overlap safe
mem_compare(a, b, n)           // compare memory, 0 = equal
```