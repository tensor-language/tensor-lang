# Control Flow

---

## if

```arc
if condition {
    // body
}
```

Executes the body block if the condition evaluates to true. The condition must be a `bool` — there is no implicit coercion from integers or pointers. Parentheses around the condition are not required. Braces are always required, even for single-statement bodies.

---

## if-else

```arc
if condition {
    // true branch
} else {
    // false branch
}
```

Extends `if` with a branch that executes when the condition is false. The `else` keyword sits on the same line as the closing brace of the if body. Both branches must be braced blocks. Exactly one branch executes. `if-else` is a statement, not an expression.

---

## if-else if-else

```arc
if status == 0 {
    io.print("ok")
} else if status == 1 {
    io.print("pending")
} else if status == 2 {
    io.print("error")
} else {
    io.print("unknown")
}
```

Chains multiple conditions in sequence. Each `else if` is evaluated only if all preceding conditions were false. The final `else` is optional and catches everything not matched by the preceding conditions. Chains can be arbitrarily long. When matching against a single value with many discrete cases, `switch` is typically cleaner.

---

## switch

```arc
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
```

Matches a single value against a list of discrete cases. Each `case` is a distinct match — there is no fallthrough between cases. The `default` branch handles anything not matched by an explicit case. `default` is optional. The matched value is evaluated once and the compiler may optimize the dispatch as a jump table depending on the density of the case values.

---

## switch, multiple values per case

```arc
switch status {
    case 1, 2, 3:
        io.print("active")
    case 4, 5:
        io.print("closing")
    default:
        io.print("unknown")
}
```

A single `case` can match multiple values separated by commas. All listed values route to the same body block. The values are OR'd together — the case matches if the switch value equals any one of them.

---

## for — C-Style

```arc
for let i = 0; i < 10; i++ {
    // i goes 0 through 9
}

for let i: usize = 0; i < len(items); i++ {
    // iterate with explicit type
}
```

A loop with an initializer, a condition, and a post statement separated by semicolons. The initializer runs once before the loop starts. The condition is checked before each iteration — when it becomes false the loop exits. The post statement runs after each iteration body completes. All three parts are optional individually but the semicolons remain when they are omitted. This form is best when you need the index explicitly, when iterating with a non-unit step, or when the loop bounds are not a simple collection.

---

## for — While Style

```arc
let j = 5
for j > 0 {
    j--
}

for conn.is_open() {
    conn.read()
}
```

A loop with only a condition. Arc does not have a separate `while` keyword. A `for` with a single boolean expression runs until that expression becomes false. The condition is evaluated before each iteration. If the condition is false on the first evaluation the body never executes.

---

## for — Infinite

```arc
for {
    process_next()
    if done {
        break
    }
}
```

A `for` with no condition runs forever until explicitly exited with `break` or `return`. This is the standard pattern for event loops, server accept loops, and any loop where the exit condition arises from inside the body.

---

## for-in — Collection Iteration

```arc
let items: vector[int32] = {1, 2, 3, 4, 5}
for item in items {
    io.print(item)
}
```

Iterates over every element in a collection, binding each element to the loop variable in turn. Works with `vector[T]`, fixed arrays `[N]T`, and slices `[]T`. The loop variable is a copy of each element — mutations to it do not affect the collection. The loop variable type is inferred from the collection's element type. For index access alongside the element, use the C-style for loop.

---

## for-in — Map Iteration

```arc
let scores: map[string]int32 = {"alice": 100, "bob": 95}
for key, value in scores {
    fmt.printf("%s: %d\n", key, value)
}
```

Iterates over every key-value pair in a map. Two loop variables are bound per iteration — the first receives the key, the second receives the value. Both are copies. Iteration order over a map is not guaranteed.

---

## for-in — Range

```arc
for i in 0..10 {
    // i goes 0 through 9
}

for i in 0..len(items) {
    // index over any length
}
```

Iterates over a half-open integer range from the start value up to but not including the end value. `0..10` produces the integers 0 through 9. The range is evaluated once at the start of the loop. The start and end can be any integer expression, not just literals.

---

## break

```arc
for item in items {
    if item == target {
        break
    }
}
```

Exits the innermost enclosing loop immediately. Execution continues at the first statement after the loop body. `break` is only valid inside a loop. When loops are nested, `break` exits only the innermost loop. There is no labeled break.

---

## continue

```arc
for item in items {
    if item < 0 {
        continue
    }
    process(item)
}
```

Skips the remainder of the current loop iteration and jumps to the next one. For C-style loops, `continue` jumps to the post statement before re-evaluating the condition. For for-in loops, `continue` advances to the next element. `continue` is only valid inside a loop and affects only the innermost enclosing loop.

---

## return

```arc
func find(items: vector[int32], val: int32) isize {
    for let i: usize = 0; i < items.len; i++ {
        if items[i] == val {
            return isize(i)
        }
    }
    return -1
}

func log(msg: string) {
    if msg == "" {
        return
    }
    io.print(msg)
}
```

Exits the current function immediately and optionally produces a value for the caller. `return` with a value exits a value-returning function. `return` with no value exits a void function early. Deferred calls scheduled with `defer` still execute when a `return` is hit.

---

## defer

```arc
let buf = new [4096]byte
defer delete(buf)

let fd = open("file.txt", O_RDONLY)
defer close(fd)

func process() {
    let conn = connect(host)
    defer conn.shutdown()

    // ... many return paths, conn.shutdown always runs
}
```

Schedules a statement to execute when the current scope exits, regardless of how it exits. Deferred statements execute in reverse order of declaration: the last `defer` registered is the first one to run. `defer` is the standard mechanism for any cleanup that must happen at scope exit: closing handles, unlocking mutexes, flushing buffers, freeing memory. Writing `defer` immediately after acquiring a resource keeps acquisition and release visually adjacent in the source.