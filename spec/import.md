# Import

---

## Package Declaration

```arc
namespace main

namespace http

namespace fmt
```

Every source file begins with a `namespace` declaration that names the package it
belongs to. The namespace is the unit of organization — all files sharing the same
namespace declaration are part of the same package. By convention, namespaces are
short, lowercase, single words with no underscores or mixed caps. The name appears
as a prefix at every callsite where the package is used — `fmt.printf`, `http.get`
— so brevity matters. A longer name taxes every caller. Avoid underscores and mixed
caps: use `httputil` not `http_util` or `httpUtil`. The namespace declaration must
appear before any imports or declarations in the file.

---

## Package Declaration, Compound Names

```arc
namespace httputil

namespace strconv

namespace bufio
```

When a package name is compound, smash the words together into a single lowercase
word rather than separating them with underscores or capitalizing word boundaries.
This is the pattern used throughout the standard library — `httputil`, `strconv`,
`bufio`, `filepath`. The result is a name that is still short, still clearly
lowercase, and still unambiguous as a callsite prefix. Underscores in namespace
names are technically permitted but strongly discouraged — they add length without
adding clarity and are not idiomatic.

---

## Single Import

```arc
import "fmt"
import "net"
import "os"
```

Imports a single package by its path string. After the import, the package is
accessible through its final path component as the local name. `import "fmt"` makes
`fmt` available — `fmt.printf`, `fmt.println`, and so on. `import "net/http"` makes
`http` available. The import path is a string literal. Importing a package and not
using it is a compile error — every import must be referenced somewhere in the file.

---

## Grouped Import

```arc
import (
    "fmt"
    "net"
    "os"
    "strings"
)
```

Multiple packages imported together in a single block delimited by parentheses. Each
path is on its own line with no separator between them. Grouped imports are the
standard form when a file needs more than one package — they keep all dependencies
visible in one place at the top of the file rather than scattered through multiple
single imports. The same rules apply: every imported package must be used.

---

## Named Import

```arc
import f "fmt"
import network "net/http"
import io "os"
```

Assigns a custom local name to the imported package. The name before the path string
is what the package is called in the current file — `f.printf`, `network.get`,
`io.open`. The actual namespace name declared inside the imported package is
irrelevant. Named imports are the right tool when two packages share the same final
path component and would otherwise collide, when the default name is ambiguous in
context, or when a shorter alias makes dense code more readable. The local name must
follow the same conventions as a namespace name: short, lowercase, no underscores,
no mixed caps.

---

## Named Import, Collision Resolution

```arc
import (
    pgdb  "database/postgres"
    mysql "database/mysql"
)

func main() {
    let pg = pgdb.connect("postgres://localhost/app")
    let my = mysql.connect("mysql://localhost/app")
}
```

The most common reason for named imports. Two packages with the same last path
component — both would resolve to `database` without an alias. Assigning distinct
local names resolves the collision cleanly at the import declaration so there is no
ambiguity at the callsite. Any names that clearly distinguish the two packages are
valid — short, lowercase aliases are the convention. Avoid underscores in aliases
just as in namespace declarations.

---

## Blank Import

```arc
import _ "drivers/usb"
import _ "init/platform"
```

Imports a package for its initialization side effects only. The `_` name explicitly
discards the package identifier — no symbols from the package are accessible in the
current file. The package's top-level initialization code runs when the program
starts. Blank import is the correct form for device drivers, platform initializers,
codec registrations, and any package whose only contribution is running setup code
at startup. Using `_` makes the intent explicit: this import is here for what it
does, not for what it exports.

---

## Dot Import

```arc
import . "math"
import . "constants"

func main() {
    let x = sqrt(2.0)    // math.sqrt, called without prefix
    let y = PI           // constants.PI, referenced without prefix
}
```

Imports all exported identifiers from a package directly into the current file's
namespace. After a dot import, `sqrt` means `math.sqrt` and `PI` means
`constants.PI` — no package qualifier needed. Dot imports reduce repetition when a
single package dominates a file and the qualifier adds no clarity. Use sparingly:
dot imports make it harder to tell where a name comes from when reading unfamiliar
code. They are most appropriate for test files and domain-specific packages whose
names are unambiguous in context.

---

## Mixed Import Block

```arc
import (
    "fmt"
    "os"

    f       "fmt/color"
    network "net/http"

    _ "drivers/gpu"
    _ "init/signals"
)
```

A single import block can contain plain imports, named imports, and blank imports
together. Blank lines within the block are conventional separators for grouping by
category — standard library, external packages, side-effect-only imports. The
compiler does not enforce any grouping convention; the blank lines are purely for
readability. All entries in the block follow the same rules regardless of which
style they use.

---

## Package Visibility

```arc
// lowercase — package-private, not accessible outside this namespace
func helper(x: int32) int32 {
    return x * 2
}

interface config {
    host: string
    port: int32
}

// uppercase — exported, accessible by any file that imports this package
func Connect(host: string, port: int32) bool {
    return true
}

interface Server {
    host: string
    port: int32
}
```

Visibility is determined entirely by the case of the first letter of the identifier.
Names beginning with an uppercase letter are exported — accessible to any file that
imports the package. Names beginning with a lowercase letter are package-private —
accessible only within files sharing the same namespace declaration. There is no
`public`, `private`, or `export` keyword. The case of the name is the declaration
of its visibility. This applies to functions, interfaces, constants, variables, and
enum members alike.

---

## Init Function

```arc
namespace devices

import "drivers/usb"
import "hal"

func init() {
    hal.register("usb", usb.driver)
    hal.set_log_level(hal.WARN)
}

func Connect(device: string) bool {
    return hal.open(device)
}
```

A function named `init` with no parameters and no return value runs automatically
when the package is first loaded, before any other code in the package executes.
`init` is called by the runtime — never called manually. A package can define
multiple `init` functions across multiple files; they all run in the order the
compiler processes the files. `init` is the right place for one-time setup that
must happen before the package is used: registering drivers, setting defaults,
validating configuration, initializing global state. Packages imported with
`import _` run their `init` functions even though their symbols are not accessible.