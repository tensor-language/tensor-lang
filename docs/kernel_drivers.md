# arc Kernel Driver Development (kernel_driver.md)

## Philosophy

At the end of the day, kernel drivers are just **machine code calling kernel APIs**. The complexity in traditional driver development comes from:
- Poor package management
- Inconsistent APIs across OS versions
- No cross-platform abstractions
- Manual dependency tracking

arc solves this through **module organization via import paths**. The same code structure works for Windows, Linux, macOS - you just import different modules.

## The Reality: It's All Machine Code

```arc
// This arc code:
import "github.com/arc-lang/windows/kernel/io"

let device = io.create_device("MyDevice")

// Compiles to x64 machine code:
// mov rcx, [device_name_ptr]
// call IoCreateDevice
// mov [device], rax

// No magic. Just function calls to kernel exports.
```

## Raw Kernel APIs (Windows)

### Direct ntoskrnl.exe Bindings

```arc
namespace driver

// Raw Windows kernel bindings
import "github.com/arc-lang/windows/kernel/raw/ntoskrnl"

// Driver entry point (required by Windows)
extern func DriverEntry(driver_obj: *ntoskrnl.DRIVER_OBJECT, registry_path: *ntoskrnl.UNICODE_STRING) ntoskrnl.NTSTATUS {
    
    // Register dispatch routines (callbacks)
    driver_obj.MajorFunction[ntoskrnl.IRP_MJ_CREATE] = &on_create
    driver_obj.MajorFunction[ntoskrnl.IRP_MJ_CLOSE] = &on_close
    driver_obj.MajorFunction[ntoskrnl.IRP_MJ_READ] = &on_read
    driver_obj.MajorFunction[ntoskrnl.IRP_MJ_WRITE] = &on_write
    driver_obj.DriverUnload = &on_unload
    
    // Create device object
    let device: *ntoskrnl.DEVICE_OBJECT = null
    let status = ntoskrnl.IoCreateDevice(
        driver_obj,
        0,  // No device extension
        null,  // No device name
        ntoskrnl.FILE_DEVICE_UNKNOWN,
        0,
        false,
        &device
    )
    
    if status != ntoskrnl.STATUS_SUCCESS {
        return status
    }
    
    return ntoskrnl.STATUS_SUCCESS
}

// IRP handler callback
func on_read(device_obj: *ntoskrnl.DEVICE_OBJECT, irp: *ntoskrnl.IRP) ntoskrnl.NTSTATUS {
    let stack = ntoskrnl.IoGetCurrentIrpStackLocation(irp)
    let buffer: *byte = cast<*byte>(irp.AssociatedIrp.SystemBuffer)
    let length = stack.Parameters.Read.Length
    
    // Read operation logic
    // ...
    
    irp.IoStatus.Status = ntoskrnl.STATUS_SUCCESS
    irp.IoStatus.Information = length
    ntoskrnl.IoCompleteRequest(irp, ntoskrnl.IO_NO_INCREMENT)
    
    return ntoskrnl.STATUS_SUCCESS
}

func on_create(device_obj: *ntoskrnl.DEVICE_OBJECT, irp: *ntoskrnl.IRP) ntoskrnl.NTSTATUS {
    irp.IoStatus.Status = ntoskrnl.STATUS_SUCCESS
    ntoskrnl.IoCompleteRequest(irp, ntoskrnl.IO_NO_INCREMENT)
    return ntoskrnl.STATUS_SUCCESS
}

func on_close(device_obj: *ntoskrnl.DEVICE_OBJECT, irp: *ntoskrnl.IRP) ntoskrnl.NTSTATUS {
    irp.IoStatus.Status = ntoskrnl.STATUS_SUCCESS
    ntoskrnl.IoCompleteRequest(irp, ntoskrnl.IO_NO_INCREMENT)
    return ntoskrnl.STATUS_SUCCESS
}

func on_write(device_obj: *ntoskrnl.DEVICE_OBJECT, irp: *ntoskrnl.IRP) ntoskrnl.NTSTATUS {
    irp.IoStatus.Status = ntoskrnl.STATUS_SUCCESS
    ntoskrnl.IoCompleteRequest(irp, ntoskrnl.IO_NO_INCREMENT)
    return ntoskrnl.STATUS_SUCCESS
}

func on_unload(driver_obj: *ntoskrnl.DRIVER_OBJECT) {
    // Cleanup
}
```

## Raw Kernel APIs (Linux)

### Direct Linux Kernel Bindings

```arc
namespace driver

// Raw Linux kernel bindings
import "github.com/arc-lang/linux/kernel/raw/module"
import "github.com/arc-lang/linux/kernel/raw/fs"
import "github.com/arc-lang/linux/kernel/raw/printk"

// Module metadata
const MODULE_LICENSE = "GPL"
const MODULE_AUTHOR = "Your Name"
const MODULE_DESCRIPTION = "Example driver"

// Module init (called on insmod)
func init_module() int32 {
    printk.printk(printk.KERN_INFO, "Driver loaded\n")
    
    // Allocate memory
    let buffer = module.kmalloc(1024, module.GFP_KERNEL)
    if buffer == null {
        return -module.ENOMEM
    }
    
    // Register character device
    let major = fs.register_chrdev(0, "mydriver", &file_ops)
    if major < 0 {
        printk.printk(printk.KERN_ERR, "Failed to register device\n")
        return major
    }
    
    return 0
}

// Module exit (called on rmmod)
func cleanup_module() {
    printk.printk(printk.KERN_INFO, "Driver unloaded\n")
    
    // Cleanup
    fs.unregister_chrdev(major, "mydriver")
    module.kfree(buffer)
}

// File operations structure
let file_ops = fs.file_operations{
    read: &device_read,
    write: &device_write,
    open: &device_open,
    release: &device_release
}

func device_open(inode: *fs.inode, file: *fs.file) int32 {
    return 0
}

func device_read(file: *fs.file, buffer: *byte, length: uint64, offset: *int64) int64 {
    // Read logic
    return cast<int64>(length)
}

func device_write(file: *fs.file, buffer: *byte, length: uint64, offset: *int64) int64 {
    // Write logic
    return cast<int64>(length)
}

func device_release(inode: *fs.inode, file: *fs.file) int32 {
    return 0
}
```

## Abstracted Helper Libraries

### Windows Driver Helper (github.com/arc-lang/windows/kernel/driver)

```arc
namespace driver

// High-level Windows driver helpers
import "github.com/arc-lang/windows/kernel/driver"
import "github.com/arc-lang/windows/kernel/io"

func main() driver.NTSTATUS {
    // Create driver with helpers
    let dev = driver.Device.new("MyDevice", driver.DeviceType.UNKNOWN)
    
    // Register handlers with closures
    dev.on_create(func(req: *driver.Request) driver.NTSTATUS {
        io.log("Device opened")
        return driver.STATUS_SUCCESS
    })
    
    dev.on_read(func(req: *driver.Request) driver.NTSTATUS {
        let data = "Hello from driver"
        req.write_bytes(data.as_bytes())
        return driver.STATUS_SUCCESS
    })
    
    dev.on_write(func(req: *driver.Request) driver.NTSTATUS {
        let bytes = req.read_bytes()
        io.log("Received: ${bytes.len()} bytes")
        return driver.STATUS_SUCCESS
    })
    
    return driver.STATUS_SUCCESS
}
```

### Linux Driver Helper (github.com/arc-lang/linux/kernel/driver)

```arc
namespace driver

// High-level Linux driver helpers
import "github.com/arc-lang/linux/kernel/driver"
import "github.com/arc-lang/linux/kernel/log"

func init() int32 {
    // Create character device with helpers
    let dev = driver.CharDevice.new("mydriver", 0)
    
    // Register handlers
    dev.on_open(func(file: *driver.File) int32 {
        log.info("Device opened")
        return 0
    })
    
    dev.on_read(func(file: *driver.File, buffer: *byte, size: uint64) int64 {
        let data = "Hello from driver"
        buffer.copy_from(data.as_bytes(), data.len())
        return cast<int64>(data.len())
    })
    
    dev.on_write(func(file: *driver.File, buffer: *byte, size: uint64) int64 {
        log.info("Received: ${size} bytes")
        return cast<int64>(size)
    })
    
    return 0
}
```

## Cross-Platform Core Module

### The Key: Organization Through Namespaces

```arc
namespace driver

// Cross-platform device abstraction
import "github.com/arc-lang/kernel/device"  // Core abstraction
import "github.com/arc-lang/kernel/memory"
import "github.com/arc-lang/kernel/log"

// This code works on Windows, Linux, macOS
func init() device.Status {
    let dev = device.create("MyDevice")
    
    dev.on_read(func(req: *device.Request) device.Status {
        let msg = "Cross-platform driver!"
        req.write(msg.as_bytes())
        return device.OK
    })
    
    dev.on_write(func(req: *device.Request) device.Status {
        let data = req.read()
        log.info("Got ${data.len()} bytes")
        return device.OK
    })
    
    return device.OK
}

// Under the hood, github.com/arc-lang/kernel/device uses:
// - "github.com/arc-lang/windows/kernel/io" on Windows
// - "github.com/arc-lang/linux/kernel/fs" on Linux
// - "github.com/arc-lang/darwin/kernel/iokit" on macOS
```

## Filesystem Driver Example

### Windows Minifilter (Abstracted)

```arc
namespace fsfilter

import "github.com/arc-lang/windows/kernel/minifilter"
import "github.com/arc-lang/windows/kernel/io"

func init() minifilter.Status {
    let filter = minifilter.create("MyFSFilter")
    
    // Intercept file operations
    filter.on_pre_create(func(data: *minifilter.CallbackData) minifilter.PreopStatus {
        let filename = data.file_name()
        io.log("Opening: ${filename}")
        
        // Allow operation to continue
        return minifilter.PREOP_SUCCESS_NO_CALLBACK
    })
    
    filter.on_post_create(func(data: *minifilter.CallbackData) minifilter.PostopStatus {
        let filename = data.file_name()
        let status = data.io_status()
        io.log("Opened: ${filename}, status: ${status}")
        
        return minifilter.POSTOP_FINISHED_PROCESSING
    })
    
    // Block writes to protected files
    filter.on_pre_write(func(data: *minifilter.CallbackData) minifilter.PreopStatus {
        let filename = data.file_name()
        
        if filename.contains("protected") {
            io.log("Blocked write to: ${filename}")
            return minifilter.PREOP_COMPLETE  // Block the write
        }
        
        return minifilter.PREOP_SUCCESS_NO_CALLBACK
    })
    
    return filter.start()
}
```

### Linux VFS Hook (Abstracted)

```arc
namespace fsfilter

import "github.com/arc-lang/linux/kernel/vfs"
import "github.com/arc-lang/linux/kernel/log"

func init() int32 {
    let filter = vfs.create_hook()
    
    // Hook file open
    filter.on_open(func(file: *vfs.File) int32 {
        let path = file.path()
        log.info("Opening: ${path}")
        return 0  // Allow
    })
    
    // Hook file write
    filter.on_write(func(file: *vfs.File, buffer: *byte, size: uint64) int64 {
        let path = file.path()
        
        if path.contains("protected") {
            log.info("Blocked write to: ${path}")
            return -vfs.EACCES  // Block
        }
        
        return cast<int64>(size)  // Allow
    })
    
    return filter.register()
}
```

## Network Driver Example

### Raw Packet Driver (Windows NDIS)

```arc
namespace netdriver

import "github.com/arc-lang/windows/kernel/ndis"
import "github.com/arc-lang/kernel/memory"

func init() ndis.Status {
    let driver = ndis.Driver.new()
    
    driver.on_initialize(func(adapter: *ndis.Adapter) ndis.Status {
        // Setup adapter
        return ndis.STATUS_SUCCESS
    })
    
    driver.on_receive(func(adapter: *ndis.Adapter, packet: *ndis.Packet) {
        let data = packet.data()
        let length = packet.length()
        
        // Process packet
        process_packet(data, length)
        
        // Complete packet
        ndis.return_packet(packet)
    })
    
    driver.on_send(func(adapter: *ndis.Adapter, packet: *ndis.Packet) ndis.Status {
        // Send packet to hardware
        return ndis.STATUS_SUCCESS
    })
    
    return driver.register()
}

func process_packet(data: *byte, length: uint32) {
    // Parse ethernet header
    let eth_header = cast<*EthernetHeader>(data)
    
    // Check if IP packet
    if eth_header.ethertype == 0x0800 {
        let ip_header = cast<*IPHeader>(data + 14)
        // Process IP packet
    }
}

struct EthernetHeader {
    dest_mac: array<byte, 6>
    src_mac: array<byte, 6>
    ethertype: uint16
}

struct IPHeader {
    version_ihl: byte
    tos: byte
    total_length: uint16
    // ... rest of IP header
}
```

## USB Driver Example

### Cross-Platform USB Driver

```arc
namespace usbdriver

import "github.com/arc-lang/kernel/usb"  // Cross-platform USB
import "github.com/arc-lang/kernel/log"

func init() usb.Status {
    let driver = usb.Driver.new()
    
    // Match USB devices (Vendor ID: 0x1234, Product ID: 0x5678)
    driver.match_device(0x1234, 0x5678)
    
    driver.on_probe(func(device: *usb.Device) usb.Status {
        log.info("USB device connected: ${device.product_name()}")
        
        // Get interface
        let intf = device.interface(0)
        
        // Find bulk IN endpoint
        let ep_in = intf.find_endpoint(usb.ENDPOINT_IN | usb.ENDPOINT_BULK)
        
        // Find bulk OUT endpoint
        let ep_out = intf.find_endpoint(usb.ENDPOINT_OUT | usb.ENDPOINT_BULK)
        
        // Start reading
        device.start_read(ep_in, &on_data_received)
        
        return usb.OK
    })
    
    driver.on_disconnect(func(device: *usb.Device) {
        log.info("USB device disconnected")
    })
    
    return driver.register()
}

func on_data_received(device: *usb.Device, data: *byte, length: uint32) {
    log.info("Received ${length} bytes from USB device")
    
    // Echo data back
    device.write_bulk(data, length)
}
```

## Memory Management

### Safe Memory Allocation

```arc
namespace driver

import "github.com/arc-lang/kernel/memory"

func handle_request() {
    // Stack-allocated buffer (fast, auto-cleanup)
    let small_buffer: array<byte, 4096> = {}
    
    // Heap-allocated buffer (for large/dynamic sizes)
    let large_buffer = memory.allocate<byte>(1024 * 1024)  // 1MB
    defer memory.free(large_buffer)  // Auto-cleanup on scope exit
    
    // Pool allocation (kernel-specific)
    let pool = memory.allocate_pool(memory.NonPagedPool, 4096)
    defer memory.free_pool(pool)
    
    // Use buffers...
}
```

## The Power of Module Organization

```
github.com/arc-lang/
├── windows/
│   └── kernel/
│       ├── raw/
│       │   ├── ntoskrnl/      # Direct ntoskrnl.exe bindings
│       │   ├── wdf/           # Direct WDF bindings
│       │   └── ndis/          # Direct NDIS bindings
│       ├── driver/            # High-level driver helpers
│       ├── minifilter/        # Filesystem filter helpers
│       ├── io/                # I/O utilities
│       └── memory/            # Memory management
├── linux/
│   └── kernel/
│       ├── raw/
│       │   ├── module/        # Direct kernel module APIs
│       │   ├── fs/            # Direct VFS APIs
│       │   └── net/           # Direct network APIs
│       ├── driver/            # High-level driver helpers
│       ├── vfs/               # Filesystem helpers
│       └── memory/            # Memory management
├── darwin/
│   └── kernel/
│       ├── raw/
│       │   └── iokit/         # Direct IOKit bindings
│       ├── driver/            # High-level driver helpers
│       └── usb/               # USB helpers
└── kernel/                     # Cross-platform abstractions
    ├── device/                # Device abstraction
    ├── memory/                # Memory abstraction
    ├── usb/                   # USB abstraction
    ├── network/               # Network abstraction
    └── log/                   # Logging abstraction
```

## Build Configuration

### Target-Specific Compilation

```arc
// config.arc

// Conditional compilation based on target OS
#if target_os(windows)
    import "github.com/arc-lang/windows/kernel/driver"
#elif target_os(linux)
    import "github.com/arc-lang/linux/kernel/driver"
#elif target_os(darwin)
    import "github.com/arc-lang/darwin/kernel/driver"
#endif

// Build command:
// arc build --target=windows-kernel-x64 --output=mydriver.sys
// arc build --target=linux-kernel-x64 --output=mydriver.ko
// arc build --target=darwin-kernel-x64 --output=mydriver.kext
```

## The Bottom Line

**All of this compiles to the same thing: machine code calling kernel APIs.**

```
Your arc code
    ↓
Compiler (monomorphization, optimization)
    ↓
x64/ARM64 machine code
    ↓
call IoCreateDevice      (Windows)
call kmalloc             (Linux)  
call IOService::start    (macOS)
```

The **real innovation** isn't magic - it's:
1. **Module organization** - Clear import paths, versioned dependencies
2. **Cross-platform abstractions** - Write once, compile for multiple OSes
3. **Package management** - `import "github.com/user/driver-utils"`
4. **Safety** - Type checking, automatic cleanup with `defer`
5. **Modern syntax** - No C preprocessor hell, clean error handling

You're not creating a new runtime or VM. You're just providing **better organization and tooling** around the same kernel APIs that have always existed.

**That's all driver development needed in the first place.**