# Enums

---

## Basic Enum

```arc
enum Direction {
    NORTH
    SOUTH
    EAST
    WEST
}

let dir = Direction.NORTH
```

A named set of discrete values. `enum` declares a type whose only valid values are the
members listed in its body. Members are accessed through the type name with dot notation
— `Direction.NORTH`, never a bare `NORTH`. By default members are assigned integer
values starting at 0 and incrementing by 1 in declaration order. The underlying type
is `int32` unless specified otherwise. Enums are value types — assigning an enum to
another variable copies the value. They are the right choice anywhere a value can only
meaningfully be one of a known fixed set of things — states, directions, modes, results.

---

## Explicit Values

```arc
enum HttpCode {
    OK           = 200
    NOT_FOUND    = 404
    SERVER_ERROR = 500
}

let code = HttpCode.OK
```

Enum members can be assigned explicit integer values. When a member has an explicit
value, the next member without an explicit value increments from it. Explicit values
do not need to be contiguous, sequential, or start from zero. They can be any integer
literal that fits in the underlying type. This is the right form when the enum values
correspond to external constants — HTTP status codes, OS error codes, protocol opcodes,
hardware register values — where the numeric value itself carries meaning and must match
an external specification exactly.

---

## Typed Enum

```arc
enum Color: uint8 {
    RED   = 0xFF
    GREEN = 0x0F
    BLUE  = 0x0A
}

enum Flags: uint32 {
    NONE       = 0x00000000
    READ       = 0x00000001
    WRITE      = 0x00000002
    EXECUTE    = 0x00000004
}
```

An enum with an explicit underlying type declared after the colon. Any fixed-width
integer type is valid — `uint8`, `uint16`, `uint32`, `uint64`, `int8`, `int16`,
`int32`, `int64`. The underlying type controls the memory size of the enum and the
range of valid member values. Typed enums are important for C interop where the ABI
requires a specific integer width, for packed structures where field size matters,
and for bitmask enums where `uint32` or `uint64` are the natural carriers. The compiler
enforces that all member values fit within the declared underlying type.

---

## Bitmask Enum

```arc
enum Permission: uint32 {
    NONE    = 0x00000000
    READ    = 0x00000001
    WRITE   = 0x00000002
    EXECUTE = 0x00000004
    ALL     = 0x00000007
}

let perms = Permission.READ | Permission.WRITE
```

An enum where members represent individual bits and combinations of members represent
sets of flags. By convention bitmask enums use an unsigned underlying type and member
values that are powers of two so each member occupies exactly one bit. Combined values
like `ALL` are expressed as the OR of the individual bits. Bitwise operators work
directly on enum values — the result is the same underlying type. Bitmask enums are the
arc equivalent of C's `#define` flag constants and are the standard pattern for
permission sets, option flags, hardware register fields, and any situation where multiple
independent binary states need to be packed into a single value.

---

## Enum in Switch

```arc
enum Status {
    OK
    PENDING
    ERROR
    TIMEOUT
}

switch status {
    case Status.OK:
        io.print("success")
    case Status.PENDING:
        io.print("waiting")
    case Status.ERROR:
        io.print("failed")
    case Status.TIMEOUT:
        io.print("timed out")
    default:
        io.print("unknown")
}
```

Enums and switch work naturally together — switch on an enum value and each case
handles one member. The compiler does not currently enforce exhaustive matching, so
a `default` branch handles any value not covered by an explicit case. When all members
are handled explicitly, `default` can be omitted. Using enum members as case values
rather than raw integers keeps the switch readable and safe — adding a new member to
the enum makes unhandled cases visible in the source rather than silently falling
through to a numeric default.

---

## Enum as Function Parameter

```arc
enum LogLevel {
    DEBUG
    INFO
    WARN
    ERROR
}

func log(level: LogLevel, msg: string) {
    if level == LogLevel.ERROR {
        io.stderr(msg)
        return
    }
    io.print(msg)
}

log(LogLevel.WARN, "disk space low")
```

Enums work as parameter types, return types, and field types in interfaces. Using an
enum as a parameter type constrains callers to only valid values — it is impossible to
pass an arbitrary integer where an enum is expected without an explicit cast. This is
the primary safety advantage of enums over raw integer constants. The function signature
communicates intent precisely: this parameter is not just any `int32`, it is specifically
a `LogLevel`. Callsites are self-documenting because the enum member name appears
explicitly rather than a bare number.

---

## Enum as Interface Field

```arc
enum State {
    IDLE
    CONNECTING
    CONNECTED
    DISCONNECTED
}

interface Connection {
    host:  string
    port:  int32
    state: State
}

var conn = Connection{host: "localhost", port: 8080, state: State.IDLE}
conn.state = State.CONNECTING
```

An enum can be a field in an interface, constraining that field to only valid members
of the enum. The field is stored as the enum's underlying integer type in memory —
there is no boxing or indirection. Initializing the field uses the same dot-notation
member access as everywhere else. Comparing the field against enum members in conditions
is the standard pattern for state machines and mode tracking inside a type.

---

## Enum Casting

```arc
enum HttpCode {
    OK           = 200
    NOT_FOUND    = 404
    SERVER_ERROR = 500
}

let raw: int32     = 200
let code: HttpCode = HttpCode(raw)

let back: int32    = int32(code)
```

An enum can be cast to and from its underlying integer type using the standard
`type(value)` cast syntax. Casting an integer to an enum produces an enum value with
that integer's bit pattern — the compiler does not check that the integer corresponds
to a declared member. Casting an enum to its underlying type gives back the raw integer.
This is necessary when interfacing with C APIs that traffic in raw integers, reading
enum values from binary data, or comparing against constants from an external header.
Outside of these cases, prefer using enum members directly rather than working with
raw integers.