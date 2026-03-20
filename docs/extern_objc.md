# Foreign Function Interface (`extern objc`)

Arc provides direct interoperability with Objective-C frameworks. Arc handles the translation of method calls to `objc_msgSend` while exposing the memory model to the developer.

Linking is handled by `import objc`. The `extern objc` block provides declarations only.

## Quick Reference

```arc
import objc "AppKit"

extern objc {
    // 1. Map classes and methods
    class NSWindow: NSResponder {
        // Maps to [[NSWindow alloc] initWithContentRect:...]
        new "initWithContentRect:styleMask:backing:defer:" (NSRect, uint64, uint64, bool) *NSWindow
        
        // Maps to -[NSWindow setTitle:] and -[NSWindow title]
        property title: *NSString
    }
}

// 2. Subclass and implement protocols
class AppDelegate: NSApplicationDelegate {
    // Arc automatically registers this class with the ObjC runtime
    func applicationDidFinishLaunching(self d: *AppDelegate, notif: *NSNotification) {
        // Handle event
    }
}
```

---

## Memory Management

Arc does **not** use ARC (Automatic Reference Counting). You are responsible for the object lifecycle using standard Manual Reference Counting (MRC) rules.

### Retain / Release

Since `retain` and `release` are methods, you can call them like any other function. Use `defer` to ensure cleanup, just like `free()` in C.

```arc
extern objc {
    class NSObject {
        func retain(self *NSObject) *NSObject
        func release(self *NSObject) void
        func autorelease(self *NSObject) *NSObject
    }
}

func main() {
    let str = NSString.new("Hello World") // Retain count +1
    defer str.release()                   // Release on scope exit
    
    // ... use str ...
}
```

### Autorelease Pools

Objective-C APIs often return "autoreleased" objects (not owned by the caller). To prevent leaks, these must be caught by an `NSAutoreleasePool`. Use `defer` to drain the pool at the end of the scope.

```arc
extern objc {
    class NSAutoreleasePool: NSObject {
        new() *NSAutoreleasePool
        func drain(self *NSAutoreleasePool) void
    }
}

func process_data() {
    let pool = NSAutoreleasePool.new()
    defer pool.drain()
    
    // Objects created here (like factory methods) are caught by 'pool'
    let temp = NSString.stringWithUTF8String("Temporary Data")
}
```

---

## Classes and Instantiation

### The `new` Keyword

The `new` keyword in an `extern objc` block is syntactic sugar for the two-step allocation process (`alloc` + `init`).

```arc
extern objc {
    class NSWindow {
        // 1. Standard init
        new() *NSWindow
        
        // 2. Custom init with selector mapping
        new "initWithContentRect:styleMask:backing:defer:" (NSRect, uint64, uint64, bool) *NSWindow
    }
}
```

**Compiler Generation:**
```c
// new "initWith..."
id tmp = objc_msgSend(objc_getClass("NSWindow"), sel_registerName("alloc"));
result = objc_msgSend(tmp, sel_registerName("initWithContentRect:styleMask:backing:defer:"), ...);
```

### Properties

Properties generate the appropriate accessor calls.

```arc
extern objc {
    class NSWindow {
        // Read-write: generates getter and setter
        property title: *NSString
        
        // Read-only: generates getter only
        property(readonly) isVisible: bool
        
        // Custom selector mapping
        property(getter: isKeyWindow) keyWindow: bool
    }
}

// Usage
window.title = str              // [window setTitle:str]
let visible = window.isVisible  // [window isVisible]
```

---

## Subclassing and Protocols

To create a class that Objective-C can "call back" into (like a Delegate), inherit from an `extern objc` class or protocol.

**The Arc compiler automatically:**
1. Generates a new ObjC class pair at runtime (e.g., `AppDelegate`).
2. Registers methods matching the signatures.
3. Sets up the `isa` pointers.

```arc
extern objc {
    protocol NSApplicationDelegate {
        optional func applicationDidFinishLaunching "applicationDidFinishLaunching:" (
            self *Self, *NSNotification
        ) void
    }
}

// Arc Implementation
class AppDelegate: NSApplicationDelegate {
    window: *NSWindow

    // This method is exported to the ObjC runtime
    func applicationDidFinishLaunching(self d: *AppDelegate, notification: *NSNotification) {
        io.print("App Launched")
    }
}

func main() {
    let app = NSApplication.sharedApplication()
    
    // Instantiate the Arc class
    let delegate = AppDelegate{} 
    
    // Assign it - ObjC sees a valid object with an isa pointer
    app.delegate = &delegate 
    
    app.run()
}
```

---

## Method Mapping

### Selector Strings

You must specify the selector if the Arc name differs from the ObjC selector, or if the selector has colons.

```arc
extern objc {
    class NSMutableArray {
        // Simple mapping (Arc "count" -> ObjC "count")
        func count(self *NSMutableArray) NSUInteger
        
        // Complex mapping (colons are required)
        func addObject "addObject:" (self *NSMutableArray, *id) void
        func insert "insertObject:atIndex:" (self *NSMutableArray, *id, NSUInteger) void
    }
}
```

> **Safety Note:** The compiler validates that the number of colons in the selector string matches the number of arguments in the function signature.

### Static Methods

Methods without `self` are treated as Class Methods.

```arc
extern objc {
    class NSColor {
        // [+NSColor redColor]
        static func redColor() *NSColor
    }
}
```

---

## Blocks

Objective-C blocks have a specific ABI layout. Arc automatically wraps lambdas passed to `extern objc` functions into stack-allocated block structs.

```arc
extern objc {
    class NSArray {
        func enumerateObjects "enumerateObjectsUsingBlock:" (
            self *NSArray, 
            block(*id, NSUInteger, *bool) void
        ) void
    }
}

func example() {
    let arr = get_array()
    
    // Arc lambda converts to ObjC Block
    arr.enumerateObjects((obj, index, stop) => {
        io.printf("Item %d: %p\n", index, obj)
        if index > 10 {
            *stop = true
        }
    })
}
```

---

## Null Safety (`nil`)

In Objective-C, sending a message to `nil` returns 0 (or false/null) and does not crash. Arc preserves this behavior for `extern objc` types.

```arc
let win: *NSWindow = null

// This is safe. It does nothing.
win.setTitle(str) 

// This returns false/0/null.
let visible = win.isVisible 
```

---

## Type System Integration

### `id` and Casting

Use `id` for generic object pointers. Arc treats `id` as an opaque pointer that can be cast to specific class types.

```arc
extern objc {
    opaque class id {}
}

let generic: *id = array.objectAtIndex(0)
let str = cast<*NSString>(generic)
```

### Basic Types

Arc automatically bridges common primitive types when they match ABI width:

| Arc Type | ObjC Type | Note |
|----------|-----------|------|
| `bool` | `BOOL` | Mapped to `signed char` (macOS) |
| `isize` | `NSInteger` | |
| `usize` | `NSUInteger` | |
| `*byte` | `char *` | C-String |
| `Selector`| `SEL` | |

---

## Complete Example: Cocoa App

```arc
namespace main

import objc "Cocoa"

extern objc {
    type NSUInteger = uint64
    type NSInteger = int64
    type Selector = *byte
    
    const NSBackingStoreBuffered: NSUInteger = 2
    const NSWindowStyleMaskTitled: NSUInteger = 1
    const NSWindowStyleMaskClosable: NSUInteger = 2
    const NSWindowStyleMaskResizable: NSUInteger = 8
    
    struct NSRect { x: float64; y: float64; w: float64; h: float64 }
    func NSMakeRect(float64, float64, float64, float64) NSRect
    
    func sel "sel_registerName" (*byte) Selector
    
    class NSObject {
        new() *NSObject
        func release(self *NSObject) void
        func autorelease(self *NSObject) *NSObject
    }

    class NSAutoreleasePool: NSObject {
        new() *NSAutoreleasePool
        func drain(self *NSAutoreleasePool) void
    }

    class NSApplication: NSObject {
        static func sharedApplication() *NSApplication
        property delegate: *id
        func run(self *NSApplication) void
        func setActivationPolicy "setActivationPolicy:" (self *NSApplication, NSInteger) bool
    }

    class NSWindow: NSObject {
        new "initWithContentRect:styleMask:backing:defer:" (NSRect, NSUInteger, NSUInteger, bool) *NSWindow
        property title: *NSString
        func makeKeyAndOrderFront "makeKeyAndOrderFront:" (self *NSWindow, *id) void
        func center(self *NSWindow) void
    }
    
    class NSString: NSObject {
        static func from "stringWithUTF8String:" (*byte) *NSString
    }

    class NSNotification: NSObject {}
    
    protocol NSApplicationDelegate {
        optional func applicationDidFinishLaunching "applicationDidFinishLaunching:" (self *Self, *NSNotification) void
        optional func applicationShouldTerminateAfterLastWindowClosed "applicationShouldTerminateAfterLastWindowClosed:" (self *Self, *NSApplication) bool
    }
}

// Subclassing allows the ObjC runtime to call us back
class AppDelegate: NSApplicationDelegate {
    window: *NSWindow
    
    func applicationDidFinishLaunching(self d: *AppDelegate, notif: *NSNotification) {
        let style = NSWindowStyleMaskTitled | NSWindowStyleMaskClosable | NSWindowStyleMaskResizable
        let rect = NSMakeRect(0, 0, 480, 320)
        
        // Create window
        d.window = NSWindow.new(rect, style, NSBackingStoreBuffered, false)
        
        // String management
        // 'from' returns an autoreleased string, so we don't own it.
        // The property setter will copy or retain it as needed.
        d.window.title = NSString.from("Arc on MacOS")
        
        d.window.center()
        d.window.makeKeyAndOrderFront(null)
    }
    
    func applicationShouldTerminateAfterLastWindowClosed(self d: *AppDelegate, app: *NSApplication) bool {
        return true
    }
}

func main() {
    // Setup top-level pool for the main thread
    let pool = NSAutoreleasePool.new()
    defer pool.drain()

    let app = NSApplication.sharedApplication()
    app.setActivationPolicy(0) // Regular app
    
    // Create our delegate (Arc struct registered as ObjC class)
    let delegate = AppDelegate{}
    app.delegate = &delegate
    
    app.run()
}