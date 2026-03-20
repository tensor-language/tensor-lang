# Package Management System

Arc's package management uses **pure HTTPS downloads** with a **local cache**. No system package managers. No global installs. No sudo required.

All external dependencies are downloaded to `~/.arc/` and linked locally.

---

## Supported Operating Systems

Arc supports the following operating systems with **direct package CDN mappings**:

| Distribution | OS Family | Hostname Prefix | Package Example |
|--------------|-----------|-----------------|-----------------|
| **Debian** | debian | `debian.org/` | `libsqlite3-dev` |
| **Ubuntu** | debian | `ubuntu.org/` | `libsqlite3-dev` |
| **macOS** | darwin | `brew.sh/` | `sqlite` |
| **Windows** | windows | `vcpkg.io/` | `sqlite3` |
| **NixOS** (optional) | nix | `nixos.org/` | `sqlite` |

### Universal Fallback & Override

For **all other distributions** (Fedora, Arch, Alpine, FreeBSD, etc.), Arc uses **vcpkg.io** as the universal package source:

```go
default "vcpkg.io/sqlite3"  // Works on any platform (Linux, macOS, Windows, BSD)
```

**Additionally, vcpkg.io can be used as an explicit override for any OS**, including the ones listed above:

```go
require c (
    sqlite3 v3.36 (
        // You can explicitly use vcpkg on any OS
        debian   "vcpkg.io/sqlite3"     // Override Debian's default
        ubuntu   "vcpkg.io/sqlite3"     // Override Ubuntu's default
        macos    "vcpkg.io/sqlite3"     // Override macOS's default
        windows  "vcpkg.io/sqlite3"     // Standard Windows source
        default  "vcpkg.io/sqlite3"
    )
)
```

This means:
- **Fedora, CentOS Stream, Rocky Linux, AlmaLinux** → `vcpkg.io`
- **Arch Linux, Manjaro, EndeavourOS** → `vcpkg.io`
- **Alpine Linux, Void Linux, Gentoo** → `vcpkg.io`
- **openSUSE, NixOS** (if not using `nixos.org`) → `vcpkg.io`
- **FreeBSD, OpenBSD, NetBSD** → `vcpkg.io`
- **Any OS can explicitly choose vcpkg.io** for consistency

**Note:** All package sources use public CDNs with no authentication required. Paid/subscription-based distributions like Red Hat Enterprise Linux (RHEL) should use the `vcpkg.io` fallback.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│  Source Code (main.arc)                                     │
├─────────────────────────────────────────────────────────────┤
│  import c "sqlite3"                                         │
│  import c "curl"                                            │
│                                                             │
│  extern c {                                                 │
│      // Declarations only                                  │
│      opaque struct sqlite3 {}                              │
│      func sqlite3_open(*byte, **sqlite3) int32             │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Build Configuration (ax.mod)                               │
├─────────────────────────────────────────────────────────────┤
│  require c (                                                │
│      sqlite3 v3.36 (                                        │
│          debian  "debian.org/libsqlite3-dev"                │
│          ubuntu  "ubuntu.org/libsqlite3-dev"                │
│          macos   "brew.sh/sqlite"                           │
│          nixos   "nixos.org/sqlite"                         │
│          windows "vcpkg.io/sqlite3"                         │
│          default "vcpkg.io/sqlite3"                         │
│      )                                                      │
│  )                                                          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Compiler Resolution                                        │
├─────────────────────────────────────────────────────────────┤
│  1. Parse imports: ["sqlite3", "curl"]                      │
│  2. Load ax.mod                                             │
│  3. Detect OS: ubuntu 22.04                                 │
│  4. Match "ubuntu" → "ubuntu.org/libsqlite3-dev"            │
│  5. Build URL: https://ubuntu.org/packages/...              │
│  6. Download to: ~/.arc/cache/ubuntu.org/sqlite3/3.36/      │
│  7. Extract libs to: ~/.arc/libs/                           │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Linker                                                     │
├─────────────────────────────────────────────────────────────┤
│  -L~/.arc/libs -lsqlite3 -lcurl                             │
└─────────────────────────────────────────────────────────────┘
```

---

## How It Actually Works

### All Downloads Are HTTPS-Based

When you write:
```go
require c (
    sqlite3 v3.36 (
        debian  "debian.org/libsqlite3-dev"
        ubuntu  "ubuntu.org/libsqlite3-dev"
        macos   "brew.sh/sqlite"
        nixos   "nixos.org/sqlite"
        windows "vcpkg.io/sqlite3"
        default "vcpkg.io/sqlite3"
    )
)
```

The compiler:
1. **Constructs an HTTPS URL:**
   ```
   # Ubuntu
   https://ubuntu.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz
   
   # macOS
   https://brew.sh/bottles/sqlite/3.43.2/darwin-arm64.tar.gz
   
   # NixOS
   https://nixos.org/packages/sqlite/3.36.0/linux-x86_64.tar.gz
   
   # Windows (or any unlisted distro)
   https://vcpkg.io/packages/sqlite3/3.36.0/windows-x86_64.tar.gz
   
   # Arch Linux (uses default)
   https://vcpkg.io/packages/sqlite3/3.36.0/linux-x86_64.tar.gz
   ```

2. **Downloads to local cache:**
   ```
   ~/.arc/cache/ubuntu.org/libsqlite3-dev/3.36/linux-x86_64/
   ~/.arc/cache/brew.sh/sqlite/3.43.2/darwin-arm64/
   ~/.arc/cache/nixos.org/sqlite/3.36.0/linux-x86_64/
   ~/.arc/cache/vcpkg.io/sqlite3/3.36.0/windows-x86_64/
   ~/.arc/cache/vcpkg.io/sqlite3/3.36.0/linux-x86_64/      # Arch, Fedora, etc.
   ```

3. **Extracts compiled libraries:**
   ```
   ~/.arc/libs/libsqlite3.a
   ~/.arc/libs/libsqlite3.so
   ~/.arc/include/sqlite3.h
   ```

4. **Links locally:**
   ```bash
   -L~/.arc/libs -lsqlite3
   ```

**No system package manager involved. No sudo. No global installs.**

---

## Cache Structure

```
~/.arc/
├── cache/                          # Downloaded packages
│   ├── debian.org/
│   │   └── libsqlite3-dev/
│   │       └── 3.36/
│   │           └── linux-x86_64/
│   │
│   ├── ubuntu.org/
│   │   └── libsqlite3-dev/
│   │       └── 3.36/
│   │           └── linux-x86_64/
│   │
│   ├── brew.sh/
│   │   └── sqlite/
│   │       └── 3.43.2/
│   │           ├── darwin-x86_64/  # Intel Mac
│   │           └── darwin-arm64/   # Apple Silicon
│   │
│   ├── nixos.org/
│   │   └── sqlite/
│   │       └── 3.36.0/
│   │           └── linux-x86_64/
│   │
│   └── vcpkg.io/                   # Windows + all other distros + optional override
│       └── sqlite3/
│           └── 3.36.0/
│               ├── linux-x86_64/   # Fedora, Arch, Alpine, etc.
│               ├── darwin-arm64/   # macOS (if using vcpkg override)
│               └── windows-x86_64/ # Windows
│
├── libs/                           # Extracted libraries
│   ├── libsqlite3.a
│   ├── libsqlite3.so              # Linux
│   ├── libsqlite3.dylib           # macOS
│   └── sqlite3.dll                # Windows
│
├── include/                        # Extracted headers
│   ├── sqlite3.h
│   ├── curl/
│   └── boost/
│
└── bin/                            # Extracted tools (if any)
    └── sqlite3
```

---

## Source Code: Simple Imports

Your source code declares **what** you need:

```arc
namespace main

import c "sqlite3"     // Just the library name
import c "curl"

extern c {
    // Declarations only - no package manager directives
    opaque struct sqlite3 {}
    opaque struct CURL {}
    
    func sqlite3_open(*byte, **sqlite3) int32
    func curl_easy_init() *CURL
}

func main() {
    let db: *sqlite3 = null
    sqlite3_open("app.db", &db)
    
    let curl = curl_easy_init()
    // ...
}
```

**Key Points:**
- `import c "sqlite3"` — No hostnames, no URLs, just the library name
- `extern c` blocks contain **only declarations** — no `lib` directives, no linking info
- Code is portable and clean

---

## Build Configuration: `ax.mod`

### Standard Configuration

```go
module github.com/username/project
arc 1.0

require (
    github.com/arc-lang/io v1.2
)

require c (
    // Cross-platform library with OS-specific CDNs
    sqlite3 v3.36 (
        debian   "debian.org/libsqlite3-dev"
        ubuntu   "ubuntu.org/libsqlite3-dev"
        macos    "brew.sh/sqlite"
        nixos    "nixos.org/sqlite"               // Optional NixOS support
        windows  "vcpkg.io/sqlite3"
        default  "vcpkg.io/sqlite3"               // Fedora, Arch, Alpine, FreeBSD, etc.
    )
    
    curl v7.80 (
        debian   "debian.org/libcurl4-openssl-dev"
        ubuntu   "ubuntu.org/libcurl4-openssl-dev"
        macos    "brew.sh/curl"
        nixos    "nixos.org/curl"
        windows  "vcpkg.io/curl"
        default  "vcpkg.io/curl"
    )
)

require cpp (
    boost v1.82 (
        debian   "debian.org/libboost-all-dev"
        ubuntu   "ubuntu.org/libboost-all-dev"
        macos    "brew.sh/boost"
        nixos    "nixos.org/boost"
        windows  "vcpkg.io/boost"
        default  "vcpkg.io/boost"
    )
)
```

### Using vcpkg.io Everywhere (Maximum Consistency)

```go
module github.com/username/project
arc 1.0

require c (
    // Use vcpkg.io on all platforms for maximum consistency
    sqlite3 v3.36 (
        debian   "vcpkg.io/sqlite3"
        ubuntu   "vcpkg.io/sqlite3"
        macos    "vcpkg.io/sqlite3"
        windows  "vcpkg.io/sqlite3"
        default  "vcpkg.io/sqlite3"
    )
    
    curl v7.80 (
        debian   "vcpkg.io/curl"
        ubuntu   "vcpkg.io/curl"
        macos    "vcpkg.io/curl"
        windows  "vcpkg.io/curl"
        default  "vcpkg.io/curl"
    )
)
```

### Simplified Configuration (vcpkg.io only)

```go
module github.com/username/project
arc 1.0

require c (
    // When you want maximum simplicity - works everywhere
    sqlite3 v3.36 (
        default "vcpkg.io/sqlite3"
    )
    
    curl v7.80 (
        default "vcpkg.io/curl"
    )
)
```

---

## Resolution Process

### Step 1: Compiler Scans Imports

```arc
// file1.arc
import c "sqlite3"
import c "curl"

// file2.arc
import c "sqlite3"    // Duplicate - recorded once
```

**Result:** `["sqlite3", "curl"]`

### Step 2: Load `ax.mod`

Compiler reads the `require c` section.

### Step 3: OS Detection

```bash
$ arc build
Detected OS: ubuntu 22.04
OS family: debian
Architecture: linux-x86_64
```

### Step 4: Resolve URL

For `sqlite3`:

1. **Lookup in ax.mod:**
   ```go
   sqlite3 v3.36 (
       ubuntu "ubuntu.org/libsqlite3-dev"
       default "vcpkg.io/sqlite3"
   )
   ```

2. **Match OS:**
   ```
   ubuntu → "ubuntu.org/libsqlite3-dev"
   ```

3. **Construct download URL:**
   ```
   https://ubuntu.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz
   ```

4. **Check cache:**
   ```bash
   ~/.arc/cache/ubuntu.org/libsqlite3-dev/3.36/linux-x86_64/
   ```

5. **If not cached:**
   ```
   Downloading: https://ubuntu.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz
   Progress: [====================] 2.3 MB / 2.3 MB
   Extracting to: ~/.arc/cache/ubuntu.org/libsqlite3-dev/3.36/linux-x86_64/
   Copying libsqlite3.a → ~/.arc/libs/
   Copying sqlite3.h → ~/.arc/include/
   ```

6. **Link:**
   ```
   -L~/.arc/libs -lsqlite3
   ```

---

## Build Output Examples

### On Ubuntu

```bash
$ arc build
Detected OS: ubuntu 22.04 (linux-x86_64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ Found mapping: ubuntu → ubuntu.org/libsqlite3-dev
  ✓ URL: https://ubuntu.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz
  ✓ Cache: ~/.arc/cache/ubuntu.org/libsqlite3-dev/3.36/linux-x86_64/
  → Cache miss, downloading... [====================] 2.3 MB
  ✓ Extracted: libsqlite3.a, libsqlite3.so, sqlite3.h
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Found mapping: ubuntu → ubuntu.org/libcurl4-openssl-dev
  ✓ Cache hit: ~/.arc/cache/ubuntu.org/libcurl4-openssl-dev/7.80/
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

### On macOS

```bash
$ arc build
Detected OS: macOS 14.2 (darwin-arm64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ Found mapping: macos → brew.sh/sqlite
  ✓ URL: https://brew.sh/bottles/sqlite/3.43.2/darwin-arm64.tar.gz
  → Downloading... [====================] 1.8 MB
  ✓ Extracted: libsqlite3.a, libsqlite3.dylib, sqlite3.h
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Found mapping: macos → brew.sh/curl
  ✓ Cache hit: ~/.arc/cache/brew.sh/curl/8.4.0/
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

### On Windows

```powershell
> arc build
Detected OS: Windows 11 (windows-x86_64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ Found mapping: windows → vcpkg.io/sqlite3
  ✓ URL: https://vcpkg.io/packages/sqlite3/3.36.0/windows-x86_64.tar.gz
  → Downloading... [====================] 2.1 MB
  ✓ Extracted: sqlite3.lib, sqlite3.dll, sqlite3.h
  ✓ Linker: sqlite3.lib

[2/2] curl v7.80
  ✓ Found mapping: windows → vcpkg.io/curl
  ✓ Cache hit: %USERPROFILE%\.arc\cache\vcpkg.io\curl\7.80.0\
  ✓ Linker: libcurl.lib

Building project...
Linking: /LIBPATH:%USERPROFILE%\.arc\libs sqlite3.lib libcurl.lib
Build complete: .\build\myapp.exe
```

### On NixOS

```bash
$ arc build
Detected OS: NixOS 23.11 (linux-x86_64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ Found mapping: nixos → nixos.org/sqlite
  ✓ URL: https://nixos.org/packages/sqlite/3.36.0/linux-x86_64.tar.gz
  → Downloading... [====================] 2.0 MB
  ✓ Extracted: libsqlite3.a, libsqlite3.so, sqlite3.h
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Found mapping: nixos → nixos.org/curl
  → Downloading... [====================] 3.0 MB
  ✓ Extracted: libcurl.a, libcurl.so, curl.h
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

### On Arch Linux (using vcpkg.io)

```bash
$ arc build
Detected OS: Arch Linux (linux-x86_64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ No specific mapping found, using default → vcpkg.io/sqlite3
  ✓ URL: https://vcpkg.io/packages/sqlite3/3.36.0/linux-x86_64.tar.gz
  → Downloading... [====================] 2.0 MB
  ✓ Extracted: libsqlite3.a, libsqlite3.so, sqlite3.h
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Using default → vcpkg.io/curl
  ✓ URL: https://vcpkg.io/packages/curl/7.80.0/linux-x86_64.tar.gz
  → Downloading... [====================] 3.1 MB
  ✓ Extracted: libcurl.a, libcurl.so, curl.h
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

### On Fedora (using vcpkg.io)

```bash
$ arc build
Detected OS: Fedora 39 (linux-x86_64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ No specific mapping found, using default → vcpkg.io/sqlite3
  ✓ URL: https://vcpkg.io/packages/sqlite3/3.36.0/linux-x86_64.tar.gz
  ✓ Cache hit: ~/.arc/cache/vcpkg.io/sqlite3/3.36.0/linux-x86_64/
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Using default → vcpkg.io/curl
  ✓ Cache hit: ~/.arc/cache/vcpkg.io/curl/7.80.0/
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

### Using vcpkg.io Override on macOS

```bash
$ arc build
Detected OS: macOS 14.2 (darwin-arm64)
Resolving dependencies...

[1/2] sqlite3 v3.36
  ✓ Found mapping: macos → vcpkg.io/sqlite3  (explicit override)
  ✓ URL: https://vcpkg.io/packages/sqlite3/3.36.0/darwin-arm64.tar.gz
  → Downloading... [====================] 1.9 MB
  ✓ Extracted: libsqlite3.a, libsqlite3.dylib, sqlite3.h
  ✓ Linker: -lsqlite3

[2/2] curl v7.80
  ✓ Found mapping: macos → vcpkg.io/curl  (explicit override)
  ✓ Cache hit: ~/.arc/cache/vcpkg.io/curl/7.80.0/
  ✓ Linker: -lcurl

Building project...
Linking: -L~/.arc/libs -lsqlite3 -lcurl
Build complete: ./build/myapp
```

---

## Hostname Prefixes

All hostname prefixes use **public CDNs** with HTTPS downloads:

| Hostname | Purpose | Example URL |
|----------|---------|-------------|
| `debian.org/` | Debian packages | `https://debian.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz` |
| `ubuntu.org/` | Ubuntu packages | `https://ubuntu.org/packages/libsqlite3-dev/3.36/linux-x86_64.tar.gz` |
| `brew.sh/` | Homebrew bottles | `https://brew.sh/bottles/sqlite/3.43.2/darwin-arm64.tar.gz` |
| `nixos.org/` | NixOS packages | `https://nixos.org/packages/sqlite/3.36.0/linux-x86_64.tar.gz` |
| `vcpkg.io/` | Universal (Windows + all platforms) | `https://vcpkg.io/packages/sqlite3/3.36.0/windows-x86_64.tar.gz` |

**All are HTTPS downloads to `~/.arc/cache/`. No authentication required.**

---

## Platform Detection

The compiler automatically detects:

```
OS Detection:
  OS: ubuntu
  Version: 22.04
  Family: debian
  Kernel: linux
  Arch: x86_64
  Libc: glibc

Platform String: linux-x86_64
```

Platform-specific package URLs:

```
Linux x86_64:      linux-x86_64.tar.gz
Linux ARM64:       linux-aarch64.tar.gz
macOS Intel:       darwin-x86_64.tar.gz
macOS ARM:         darwin-arm64.tar.gz
Windows x64:       windows-x86_64.tar.gz
FreeBSD x64:       freebsd-x86_64.tar.gz
```

---

## Package Format

Each package is a `.tar.gz` containing:

```
sqlite3-3.36-linux-x86_64.tar.gz
├── manifest.json              # Metadata
├── lib/
│   ├── libsqlite3.a          # Static library
│   └── libsqlite3.so         # Dynamic library (Linux)
│       libsqlite3.dylib      # Dynamic library (macOS)
│       sqlite3.dll           # Dynamic library (Windows)
├── include/
│   └── sqlite3.h             # Headers
└── bin/                      # Optional tools
    └── sqlite3
```

### `manifest.json`

```json
{
  "name": "sqlite3",
  "version": "3.36.0",
  "arch": "linux-x86_64",
  "libs": ["sqlite3"],
  "headers": ["sqlite3.h"],
  "dependencies": [],
  "linker_flags": ["-lsqlite3"],
  "sha256": "a1b2c3d4..."
}
```

---

## Complete Example

### Source Code

```arc
namespace main

import c "sqlite3"
import c "curl"

extern c {
    opaque struct sqlite3 {}
    opaque struct CURL {}
    
    func sqlite3_open(*byte, **sqlite3) int32
    func curl_easy_init() *CURL
}

func main() {
    let db: *sqlite3 = null
    sqlite3_open("app.db", &db)
    
    let curl = curl_easy_init()
    // ...
}
```

### Build Config

```go
module github.com/myapp/example
arc 1.0

require c (
    sqlite3 v3.36 (
        debian   "debian.org/libsqlite3-dev"
        ubuntu   "ubuntu.org/libsqlite3-dev"
        macos    "brew.sh/sqlite"
        nixos    "nixos.org/sqlite"
        windows  "vcpkg.io/sqlite3"
        default  "vcpkg.io/sqlite3"        // All other distros
    )
    
    curl v7.80 (
        debian   "debian.org/libcurl4-openssl-dev"
        ubuntu   "ubuntu.org/libcurl4-openssl-dev"
        macos    "brew.sh/curl"
        nixos    "nixos.org/curl"
        windows  "vcpkg.io/curl"
        default  "vcpkg.io/curl"
    )
)
```

---

## Benefits

### ✅ No sudo Required
Everything in `~/.arc/`, no system pollution.

### ✅ Reproducible Builds
Same cache means same libraries, regardless of system packages.

### ✅ Version-Locked
`sqlite3 v3.36` always downloads 3.36, not "whatever system has."

### ✅ Cross-Platform Consistency
Same download mechanism on Linux, macOS, Windows, FreeBSD.

### ✅ No Dependency Conflicts
Your project's libs don't interfere with system libs.

### ✅ Portable
Entire `~/.arc/` directory can be copied or shared.

### ✅ Universal Fallback
vcpkg.io provides consistent packages for all unlisted distributions.

### ✅ Flexible Overrides
Use vcpkg.io on any OS for maximum consistency when needed.

### ✅ NixOS Support
Optional native NixOS package support for Nix users.

---

## Use Cases

### Maximum Compatibility (Recommended)
```go
require c (
    sqlite3 v3.36 (
        debian   "debian.org/libsqlite3-dev"
        ubuntu   "ubuntu.org/libsqlite3-dev"
        macos    "brew.sh/sqlite"
        nixos    "nixos.org/sqlite"
        windows  "vcpkg.io/sqlite3"
        default  "vcpkg.io/sqlite3"
    )
)
```
**Best for:** Libraries and applications targeting multiple platforms.

### Maximum Consistency
```go
require c (
    sqlite3 v3.36 (
        debian   "vcpkg.io/sqlite3"
        ubuntu   "vcpkg.io/sqlite3"
        macos    "vcpkg.io/sqlite3"
        windows  "vcpkg.io/sqlite3"
        default  "vcpkg.io/sqlite3"
    )
)
```
**Best for:** Projects requiring identical binaries across all platforms.

### Maximum Simplicity
```go
require c (
    sqlite3 v3.36 (
        default "vcpkg.io/sqlite3"
    )
)
```
**Best for:** Quick prototypes and simple projects.

---

## Summary

Arc's package management:

1. **Pure HTTPS downloads** — No system package managers
2. **Local cache** — Everything in `~/.arc/`
3. **No global installs** — No sudo, no system pollution
4. **Platform-specific URLs** — Automatic platform detection
5. **Reproducible** — Version-locked dependencies
6. **Portable** — Copy `~/.arc/` to another machine
7. **Simple mappings** — Only 5 specific OS mappings (Debian, Ubuntu, macOS, Windows, NixOS)
8. **Universal fallback** — vcpkg.io for all other distributions
9. **Flexible overrides** — vcpkg.io can be used on any OS explicitly

The hostname prefixes use **public CDNs** with HTTPS downloads:
- **`debian.org`** and **`ubuntu.org`** for Debian-based systems
- **`brew.sh`** for macOS
- **`nixos.org`** for NixOS (optional)
- **`vcpkg.io`** for Windows, all other distributions, and as an optional override for any OS

All supported platforms use publicly accessible CDNs with no authentication required.

This is simpler, cleaner, and more reliable than maintaining mappings for every Linux distribution, while still providing flexibility for projects that need either platform-specific packages or universal consistency.