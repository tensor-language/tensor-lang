# External Imports (`import c`)

Arc unifies external dependency management. There are no makefiles, no `pkg-config` setups, and no complex configuration files. The compiler talks directly to the **Universal Package Manager (upkg)** to resolve, download, and link dependencies automatically.

## Syntax Overview

```arc
// 1. System/Registry Imports (Standard libraries)
import c "libc"             // Links the C runtime (-lc)
import c "sqlite3"          // Auto-installs & links libsqlite3

// 2. Decentralized Imports (Git repositories)
import c "github.com/user/lib-wrapper"  // Clones repo, reads index.toml, links
```

---

## How It Works

The Arc compiler integrates `upkg` directly into its frontend.

```mermaid
graph TD
    A[Source Code] -->|import c "name"| B(Compiler)
    B -->|Query| C{Is it a URL?}
    
    C -->|No: System Pkg| D[upkg Registry]
    D -->|Lookup deps/name| E[Read index.toml]
    
    C -->|Yes: Git URL| F[Git Clone]
    F -->|~/.upkg/src/...| G[Read index.toml]
    
    E --> H[Linker Flags]
    G --> H
    
    H -->| -lsqlite3 -lc | I[Internal ELF Linker]
    I --> J[Final Executable]
```

---

## Import Types

### 1. Registry Imports (System Libraries)
These resolve against the central `upkg` registry (the `deps/` folder). This is used for standard, high-frequency libraries.

```arc
import c "libc"
import c "sqlite3"
import c "openssl"
```

**What happens:**
1. Compiler checks `~/.upkg/deps/<name>/index.toml`.
2. It detects your OS (Ubuntu, Arch, macOS, etc.).
3. It installs the system package if missing (e.g., `apt install libsqlite3-dev`).
4. It reads the `libs` array (e.g., `["sqlite3"]`) and auto-links.

### 2. Decentralized Imports (Wrappers)
You don't need to wait for the central registry to update. You can import any C library wrapper hosted on Git.

```arc
import c "github.com/username/my-lib-wrapper"
import c "github.com/username/monorepo/packages/graphics"
```

**What happens:**
1. Compiler clones the repo to `~/.upkg/src/github.com/username/...`.
2. It looks for an `index.toml` in that directory.
3. It applies the linker flags defined in that file.

---

## The `index.toml` Configuration

To make a C library importable in Arc, you only need one file: `index.toml`. This tells the compiler what to link and `upkg` what to install.

**Example: `deps/sqlite3/index.toml`**

```toml
name = "sqlite3"

# Used by the Arc Compiler for linking
# Result: -lsqlite3
libs = ["sqlite3"] 

# Used by upkg to install the library on the host OS
[backends]
apt     = "libsqlite3-dev"
dnf     = "sqlite-devel"
brew    = "sqlite"
apk     = "sqlite-dev"
pacman  = "sqlite"
winget  = "SQLite.SQLite"
```

**Example: `deps/libc/index.toml`**

```toml
name = "libc"
libs = ["c"] # Result: -lc

[backends]
apt = "libc6-dev"
# ... other distros ...
```

---

## Complete Example

This is a real, working example. No build scripts required.

### `main.ax`

```go
namespace main

// 1. Import the C Runtime and SQLite
import c "libc"
import c "sqlite3"

// Constants
const SQLITE_OK = 0

// Struct
struct sqlite3 {}

// 2. Define C Interface
extern c {    
    // Functions
    func sqlite3_libversion() *byte
    func sqlite3_open(*byte, **sqlite3) int32
    func sqlite3_close(*sqlite3) int32
    func sqlite3_errmsg(*sqlite3) *byte
    
    // Varargs from libc
    func printf(*byte, ...) int32
}

func main() {
    let db: *sqlite3 = null

    // Call C functions directly
    printf("SQLite Version: %s\n", sqlite3_libversion())

    if sqlite3_open("test.db", &db) != SQLITE_OK {
        printf("Error: %s\n", sqlite3_errmsg(db))
    } else {
        printf("Database opened successfully!\n")
    }

    sqlite3_close(db)
}
```

### Build & Run

```bash
$ arc build main.ax -o app
[Driver] [INFO] Compiling main.ax -> app
[Pkg] Resolving system dependency: libc
[Driver] [INFO] Auto-detected libraries for 'libc': [c]
[Pkg] Resolving system dependency: sqlite3
[Driver] [INFO] Auto-detected libraries for 'sqlite3': [sqlite3]
[Driver] [INFO] Linking libraries: [c sqlite3]
[Driver] [INFO] Success! Executable created: app

$ ./app
SQLite Version: 3.45.1
Database opened successfully!
```

---

## Cache Directory Structure

All dependencies are managed in the user's home directory, keeping your project folder clean.

```text
~/.upkg/
├── deps/                   # The Central Registry (Synced from GitHub)
│   ├── libc/
│   │   └── index.toml
│   ├── sqlite3/
│   │   └── index.toml
│   └── openssl/
│       └── index.toml
│
└── src/                    # Decentralized Git Imports
    └── github.com/
        └── johndoe/
            └── my-wrapper/
                └── index.toml
```