# arc Standard Library Collections (Version 1.2)

> **Note**: `vector<T>`, `map<K,V>`, and `array<T,N>` are NOT compiler builtins.
> They are generic types defined in the standard library.
> See the Generics section in grammar_1.2.md for implementation patterns.

## Fixed-Size Arrays

### array - Stack-Allocated Fixed-Size Collection
```arc
// Declaration with type and size
let arr: array<int32, 5> = {1, 2, 3, 4, 5}
let coords: array<float64, 3> = {1.0, 2.0, 3.0}
```

### Array Initialization
```arc
// Initialized with values
let nums: array<int32, 5> = {1, 2, 3, 4, 5}

// Zero/default initialization
let zeros: array<int32, 100> = {}

// Type inference (size still required)
let items: array<_, 3> = {10, 20, 30}

// Using constant for size
const BUFFER_SIZE: usize = 1024
let buffer: array<byte, BUFFER_SIZE> = {}
```

### Array Access
```arc
let arr: array<int32, 5> = {1, 2, 3, 4, 5}

// Read element
let val = arr[2]

// Write element
arr[3] = 42

// Get pointer to element
let ptr: *int32 = &arr[0]

// Get pointer to entire array
let base_ptr: *int32 = &arr
```

### Array Iteration
```arc
let arr: array<int32, 5> = {1, 2, 3, 4, 5}

// Iterate over elements


for item in arr {
    io.printf("%d\n", item)
}

// Iterate with index
for let i: usize = 0; i < 5; i++ {
    arr[i] = arr[i] * 2
}
```

## Dynamic Vectors

### vector - Heap-Allocated Dynamic Array
```arc
// Declaration
let v: vector<int32> = {}

// With initial capacity hint (optimization)
let v2: vector<string> = vector<string>.with_capacity(100)
```

### Vector Initialization
```arc
// Empty vector
let empty: vector<int32> = {}

// Initialized with values
let nums: vector<int32> = {1, 2, 3, 4, 5}

// Type inference
let items = {10, 20, 30}

// Pre-allocated capacity
let buffer: vector<byte> = vector<byte>.with_capacity(1024)
```

### Vector Operations
```arc
let v: vector<int32> = {}

// Add elements
v.push(10)
v.push(20)
v.push(30)

// Get length
let len = v.len()

// Get capacity
let cap = v.capacity()

// Access elements
let first = v[0]
v[1] = 100

// Remove last element
let last = v.pop()  // Returns optional value

// Clear all elements
v.clear()

// Check if empty
if v.is_empty() {
    // vector is empty
}
```

### Vector Iteration
```arc
let items: vector<string> = {"apple", "banana", "cherry"}

// Iterate over elements
for item in items {
    io.printf("%s\n", item)
}

// Iterate with index
for let i: usize = 0; i < items.len(); i++ {
    io.printf("%d: %s\n", i, items[i])
}
```

### Vector Memory Layout
```arc
// vector<T> is typically:
struct vector<T> {
    data: *T      // Pointer to heap-allocated array
    len: usize    // Current number of elements
    cap: usize    // Allocated capacity
}

// Getting raw pointer (for syscalls, FFI)
let v: vector<byte> = {1, 2, 3, 4, 5}
let ptr: *byte = v.data()
let length = v.len()

// Pass to syscall
syscall(SYS_WRITE, STDOUT, ptr, length)
```

## Hash Maps

### map - Hash-Based Key-Value Store
```arc
// Declaration
let m: map<string, int32> = {}

// With capacity hint
let m2: map<int32, string> = map<int32, string>.with_capacity(100)
```

### Map Initialization
```arc
// Empty map
let empty: map<string, int32> = {}

// Initialized with key-value pairs
let scores: map<string, int32> = {"alice": 100, "bob": 95, "charlie": 87}

// Type inference
let config = {"host": "localhost", "port": "8080"}
```

### Map Operations
```arc
let scores: map<string, int32> = {}

// Insert/Update
scores["alice"] = 100
scores["bob"] = 95

// Get value (returns optional)
let alice_score = scores.get("alice")
if alice_score != null {
    io.printf("Alice: %d\n", *alice_score)
}

// Access with default
let score = scores.get_or_default("charlie", 0)

// Check if key exists
if scores.contains("alice") {
    // key exists
}

// Remove key
scores.remove("bob")

// Get number of entries
let count = scores.len()

// Clear all entries
scores.clear()

// Check if empty
if scores.is_empty() {
    // map is empty
}
```

### Map Iteration
```arc
let scores: map<string, int32> = {"alice": 100, "bob": 95}

// Iterate over key-value pairs
for key, value in scores {
    io.printf("%s: %d\n", key, value)
}

// Iterate over keys only
for key in scores.keys() {
    io.printf("Key: %s\n", key)
}

// Iterate over values only
for value in scores.values() {
    io.printf("Value: %d\n", value)
}
```

### Map Memory Layout
```arc
// map<K, V> is typically:
struct Entry<K, V> {
    key: K
    value: V
    hash: usize
}

struct map<K, V> {
    buckets: vector<vector<Entry<K, V>>>
    count: usize
    load_factor: float32
}
```

## Collection Performance Characteristics

**array<T, N>:**
- Fixed size known at compile time
- Stack allocated (very fast, no heap allocations)
- No bounds checking overhead in release mode
- Best for: Small, fixed-size data, performance-critical code

**vector<T>:**
- Dynamic size, grows as needed
- Heap allocated, amortized O(1) push
- Contiguous memory (cache-friendly)
- Best for: Lists, buffers, dynamic arrays

**map<K, V>:**
- Average O(1) insert/lookup/delete
- Heap allocated, hash-based
- Memory overhead for buckets and collision handling
- Best for: Key-value lookups, caches, associations

## Common Patterns

### Stack Buffer Pattern
```arc
// Fixed-size temporary buffer
let buffer: array<byte, 4096> = {}
let bytes_read = syscall(SYS_READ, fd, &buffer, 4096)
```

### Dynamic Buffer Pattern
```arc
// Growing buffer for unknown size
let buffer: vector<byte> = {}
for {
    let chunk: array<byte, 1024> = {}
    let n = read_chunk(&chunk)
    if n == 0 { break }
    
    for let i: usize = 0; i < n; i++ {
        buffer.push(chunk[i])
    }
}
```

### Cache Pattern
```arc
// Memoization/caching
let cache: map<string, int32> = {}

func expensive_computation(key: string) int32 {
    if cache.contains(key) {
        return *cache.get(key)
    }
    
    let result = compute(key)
    cache[key] = result
    return result
}
```

### Configuration Pattern
```arc
// Simple key-value config
let config: map<string, string> = {
    "host": "localhost",
    "port": "8080",
    "timeout": "30"
}

let host = config.get_or_default("host", "0.0.0.0")
let port = parse_int(config.get_or_default("port", "80"))
```