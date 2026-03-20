# Arc Standard Library - Core Modules (Version 1.0)

## Overview
Arc modules are imported directly by name (e.g., `import "http"`). 

**Implicitly Imported (Prelude)**
The following are available in every file without import:
*   **Types:** `int*`, `uint*`, `bool`, `char`, `byte`, `void`
*   **Collections:** `vector<T>`, `map<K,V>`, `string`, `array<T, N>`
*   **Builtins:** `print()`, `panic()`, `len()`, `sizeof()`, `cast<T>()`

---

## 1. System & I/O

### `io`
Universal interfaces for input and output streams.
```arc
// Interfaces (Traits)
interface Reader { func read(buf: vector<byte>) int32 }
interface Writer { func write(buf: vector<byte>) int32 }
interface Closer { func close() }

// Key functions
func Copy(dst: Writer, src: Reader) int64
func Pipe() (Reader, Writer)

// Types
struct Buffer { ... } // In-memory byte buffer

```

### `fs`
File system manipulation and file I/O.
```arc
struct File { ... }
struct Stat { size: int64, mode: uint32, ... }

let client = tcp.Client()

// Constructors / Actions
func Open(path: string) File
func Create(path: string) File
func Remove(path: string)
func ReadFile(path: string) vector<byte> // Utility: read all
func WriteFile(path: string, data: vector<byte>)
func Walk(path: string, cb: func(path: string, info: Stat))
```

### `path`
Cross-platform file path manipulation (handles `/` vs `\` automatically).
```arc
func Join(elem: ...string) string
func Base(path: string) string
func Dir(path: string) string
func Ext(path: string) string
func Abs(path: string) string
```

### `os`
Operating system primitives and process management.
```arc
// Environment & Args
func Env(key: string) string
func SetEnv(key: string, val: string)
let args: vector<string> // Command line arguments

// Process
func Exit(code: int32)
func Pid() int32
func Wd() string // Working directory

// Execution
func Exec(cmd: string, args: vector<string>) Result<Process>
```

---

## 2. Networking

### `net`
Shared networking primitives and IP handling.
```arc
struct IP { ... }

func ParseIP(s: string) IP
func LookupHost(host: string) vector<string>
```

### `tcp`
TCP specific implementation.
```arc
struct Connection {
    func read(buf: vector<byte>) int32
    func write(buf: vector<byte>) int32
    func close()
}

struct Listener {
    func accept() Connection
}

// Initializers
async func Client(addr: string) Connection
async func Server(port: int32) Listener
```

### `udp`
UDP specific implementation.
```arc
struct Socket {
    func read_from(buf: vector<byte>) (int32, net.IP)
    func write_to(buf: vector<byte>, addr: net.IP) int32
}

// Initializers
async func Client(addr: string) Socket // Connects to specific remote
async func Server(port: int32) Socket // Listens on port
```

### `http`
Production-ready HTTP/1.1 and HTTP/2 client and server.
```arc
struct Request { 
    method: string
    url: url.URL
    headers: map<string, string>
    body: io.Reader 
}

struct Response { 
    status: int32
    body: vector<byte> 
}

// Client
struct Client {
    // Config like timeout, proxy, etc
}

func NewClient() Client
async func Get(url: string) Response
async func Post(url: string, body: string, type: string) Response

// Server
struct Server {
    func handle(path: string, handler: async func(req, res))
    async func listen()
}

// Initializer
func NewServer(port: int32) Server
```

### `url`
URL parsing and manipulation.
```arc
struct URL {
    scheme: string
    host: string
    path: string
    query: map<string, string>
}

func Parse(raw: string) URL
func Encode(str: string) string // Percent encoding
```

---

## 3. Data & Encoding

### `json`
Fast, standard JSON serialization.
```arc
// Converts struct/map to string
func Stringify(v: any) string 

// Converts string to struct/map
func Parse<T>(json: string) T 

// Validation
func Valid(json: string) bool
```

### `encoding`
Binary text encoders.
```arc
namespace hex {
    func Encode(src: vector<byte>) string
    func Decode(src: string) vector<byte>
}

namespace base64 {
    func Encode(src: vector<byte>) string
    func Decode(src: string) vector<byte>
}

namespace csv {
    func Reader(r: io.Reader) Iterator
    func Writer(w: io.Writer) Encoder
}
```

### `regex`
Regular expression matching and replacement (PCRE style).
```arc
struct Pattern {
    func match(text: string) bool
    func find(text: string) string
    func replace(text: string, repl: string) string
}

func Compile(expr: string) Pattern
```

---

## 4. Utilities

### `time`
Date, time, and duration management.
```arc
struct Time { ... }

// Constructors
func Now() Time
func Unix(sec: int64) Time

// Duration
func Sleep(ms: int64) async
func Since(t: Time) int64

// Formatting
const RFC3339 = "2006-01-02T15:04:05Z07:00"
func Parse(layout: string, value: string) Time
```

### `math`
Standard mathematical constants and functions.
```arc
const PI: float64
const E: float64

func Abs(x: any) any
func Min(a: any, b: any) any
func Max(a: any, b: any) any
func Sqrt(x: float64) float64
func Pow(x: float64, y: float64) float64
func Sin(x: float64) float64
func Cos(x: float64) float64
```

### `rand`
Random number generation.
```arc
// Fast, non-crypto random
func Int(max: int32) int32
func Float() float64
func Shuffle<T>(v: vector<T>)

// Cryptographically secure
func CryptoBytes(n: int32) vector<byte>
```

### `log`
Structured logging with support for JSON output.
```arc
func Info(msg: string, fields: map<string, any>...)
func Warn(msg: string, fields: map<string, any>...)
func Error(msg: string, fields: map<string, any>...)
func SetLevel(lvl: int32)
func SetFormat(fmt: string) // "text" or "json"
```

### `sync`
Synchronization primitives for safe concurrency.
```arc
struct Mutex {
    func lock()
    func unlock()
}

struct WaitGroup {
    func add(delta: int32)
    async func wait()
}

struct Once {
    func do(fn: func())
}
```