# arc Kernel Drivers & Modules (kernel_drivers_modules.md)

## Philosophy

Kernel modules give you **total system control**. The complexity isn't the programming - it's organizing thousands of APIs across different subsystems.

arc solves this with **namespace imports**:
- Clear module hierarchy
- Versioned dependencies  
- Cross-platform abstractions
- Raw APIs when you need them

**Remember:** It's all just machine code calling kernel functions. The imports just organize it cleanly.


# May need to organize and make wrappers for some callbacks 

### API Growth Over Time:
```
Windows NT 3.1 (1993):   ~800 APIs
Windows NT 4.0 (1996):   ~1,200 APIs
Windows 2000:            ~1,800 APIs
Windows XP:              ~2,200 APIs
Windows Vista:           ~3,000 APIs
Windows 7:               ~3,500 APIs
Windows 8:               ~4,500 APIs
Windows 10:              ~5,500 APIs
Windows 11:              ~6,000+ APIs

github.com/arc-lang/windows/kernel/
├── raw/                          # Direct extern bindings
│   ├── ntoskrnl/                # Core kernel (ntoskrnl.exe)
│   │   ├── io.arc              # Io* functions
│   │   ├── mm.arc              # Mm* functions
│   │   ├── ex.arc              # Ex* functions
│   │   ├── ke.arc              # Ke* functions
│   │   ├── rtl.arc             # Rtl* functions
│   │   ├── ps.arc              # Ps* functions
│   │   └── zw.arc              # Zw* functions
│   ├── wdf/                     # Windows Driver Framework
│   │   ├── device.arc
│   │   ├── io.arc
│   │   └── memory.arc
│   └── ndis/                    # Network drivers
│       ├── core.arc
│       └── packet.arc
├── types/                       # Type definitions
│   ├── core.arc                # NTSTATUS, UNICODE_STRING, etc.
│   ├── irp.arc                 # IRP, IO_STACK_LOCATION, etc.
│   └── device.arc              # DEVICE_OBJECT, DRIVER_OBJECT, etc.
└── abstractions/                # High-level wrappers
    ├── driver/                  # Driver helpers
    │   ├── core.arc
    │   └── callbacks.arc
    ├── io/                      # I/O abstractions
    │   └── request.arc
    └── memory/                  # Memory helpers
        └── pool.arc

extern "ntoskrnl.exe" {
    func IoCreateDevice(
        DriverObject: *types.DRIVER_OBJECT,
        DeviceExtensionSize: uint32,
        DeviceName: *types.UNICODE_STRING,
        DeviceType: uint32,
        DeviceCharacteristics: uint32,
        Exclusive: bool,
        DeviceObject: **types.DEVICE_OBJECT
    ) types.NTSTATUS
}
```

## Complete Import Reference

### Network Stack 🌐

```arc
// Core networking - TCP/IP stack, sockets, routing
import "kernel/net/core"
import "kernel/net/ipv4"
import "kernel/net/ipv6"
import "kernel/net/tcp"
import "kernel/net/udp"
import "kernel/net/icmp"

// Packet filtering & firewall (iptables/nftables backend)
import "kernel/net/netfilter"
import "kernel/net/netfilter/ipv4"
import "kernel/net/netfilter/ipv6"
import "kernel/net/netfilter/bridge"

// Raw packet access (tcpdump, wireshark use this)
import "kernel/net/packet"
import "kernel/net/raw"

// Network device drivers (ethernet, wifi)
import "kernel/net/device"
import "kernel/net/ethernet"

// Wireless subsystem
import "kernel/net/wireless"
import "kernel/net/wifi"
import "kernel/net/bluetooth"
import "kernel/net/bluetooth/hci"
import "kernel/net/bluetooth/l2cap"

// Network namespaces & virtualization
import "kernel/net/namespace"
import "kernel/net/veth"        // Virtual ethernet pairs
import "kernel/net/bridge"      // Network bridge
import "kernel/net/tun"         // TUN/TAP devices

// Socket management
import "kernel/net/sock"
import "kernel/net/inet"
```

### Filesystem & Block I/O 💾

```arc
// Virtual File System - the core abstraction
import "kernel/fs/vfs"
import "kernel/fs/dcache"       // Dentry cache
import "kernel/fs/inode"
import "kernel/fs/file"

// Specific filesystems
import "kernel/fs/ext4"
import "kernel/fs/btrfs"
import "kernel/fs/xfs"
import "kernel/fs/ntfs"
import "kernel/fs/fat"
import "kernel/fs/proc"         // /proc filesystem
import "kernel/fs/sysfs"        // /sys filesystem
import "kernel/fs/tmpfs"        // RAM-based filesystem

// Block device layer (disks, SSDs, RAID)
import "kernel/block/device"
import "kernel/block/queue"
import "kernel/block/bio"       // Block I/O operations
import "kernel/block/genhd"     // Generic hard disk

// Page cache & buffer management
import "kernel/fs/pagecache"
import "kernel/fs/buffer"

// File monitoring (inotify backend)
import "kernel/fs/notify"
import "kernel/fs/inotify"
import "kernel/fs/fanotify"

// FUSE - userspace filesystems
import "kernel/fs/fuse"

// Device mapper (LVM, encryption, etc)
import "kernel/block/dm"
import "kernel/block/dm/crypt"
import "kernel/block/dm/linear"
```

### USB 🔌

```arc
// USB core infrastructure
import "kernel/usb/core"
import "kernel/usb/urb"         // USB Request Block

// USB host controller drivers
import "kernel/usb/host"
import "kernel/usb/host/ehci"   // USB 2.0
import "kernel/usb/host/xhci"   // USB 3.0+
import "kernel/usb/host/ohci"   // USB 1.1
import "kernel/usb/host/uhci"   // USB 1.1 (Intel)

// USB device classes
import "kernel/usb/storage"     // Flash drives, external HDDs
import "kernel/usb/serial"      // USB-to-serial adapters
import "kernel/usb/hid"         // Keyboards, mice, game controllers
import "kernel/usb/audio"       // USB headsets, microphones
import "kernel/usb/video"       // Webcams (UVC)
import "kernel/usb/cdc"         // Modems, network adapters
import "kernel/usb/printer"     // USB printers
import "kernel/usb/hub"         // USB hubs

// USB gadget (device mode - your computer acts as USB device)
import "kernel/usb/gadget"
import "kernel/usb/gadget/function"
```

### PCI/PCIe 🎛️

```arc
// PCI bus enumeration & configuration
import "kernel/pci/core"
import "kernel/pci/bus"
import "kernel/pci/device"

// PCI Express features
import "kernel/pci/pcie"
import "kernel/pci/pcie/aer"    // Advanced Error Reporting
import "kernel/pci/pcie/aspm"   // Active State Power Management

// MSI/MSI-X interrupts (modern interrupt handling)
import "kernel/pci/msi"

// Device power management
import "kernel/pci/power"

// PCI hotplug
import "kernel/pci/hotplug"

// PCI configuration space access
import "kernel/pci/config"
```

### DMA & Memory 🧠

```arc
// Direct Memory Access engine
import "kernel/dma/engine"
import "kernel/dma/channel"

// DMA buffer management
import "kernel/dma/buffer"
import "kernel/dma/pool"

// DMA mapping for devices
import "kernel/dma/mapping"
import "kernel/dma/coherent"

// IOMMU (virtualization, security)
import "kernel/iommu/core"
import "kernel/iommu/intel"     // VT-d
import "kernel/iommu/amd"       // AMD-Vi

// Memory management
import "kernel/mm/alloc"        // kmalloc, vmalloc
import "kernel/mm/slab"         // Slab allocator
import "kernel/mm/page"         // Page allocation
import "kernel/mm/mmap"         // Memory mapping
import "kernel/mm/vmalloc"      // Virtual memory allocation
import "kernel/mm/mempool"      // Memory pools
import "kernel/mm/highmem"      // High memory support
```

### Graphics & Display 🎨

```arc
// Direct Rendering Manager (GPU drivers)
import "kernel/gpu/drm/core"
import "kernel/gpu/drm/device"

// GPU memory management
import "kernel/gpu/drm/gem"     // Graphics Execution Manager
import "kernel/gpu/drm/ttm"     // Translation Table Manager

// Display output (monitors, HDMI, DisplayPort)
import "kernel/gpu/drm/kms"     // Kernel Mode Setting
import "kernel/gpu/drm/crtc"    // Display controller
import "kernel/gpu/drm/encoder"
import "kernel/gpu/drm/connector"

// Specific GPU drivers
import "kernel/gpu/drm/i915"    // Intel
import "kernel/gpu/drm/amdgpu"  // AMD
import "kernel/gpu/drm/nouveau" // NVIDIA (open)
import "kernel/gpu/drm/radeon"  // AMD (legacy)

// Framebuffer (simple graphics)
import "kernel/video/fbdev"
import "kernel/video/fbmem"
```

### Input Devices ⌨️

```arc
// Input subsystem core
import "kernel/input/core"
import "kernel/input/event"

// Device types
import "kernel/input/keyboard"
import "kernel/input/mouse"
import "kernel/input/touchscreen"
import "kernel/input/gamepad"
import "kernel/input/joystick"
import "kernel/input/tablet"

// Input event injection (for automation)
import "kernel/input/uinput"

// HID (Human Interface Device) protocol
import "kernel/hid/core"
import "kernel/hid/generic"
```

### Security & Crypto 🔐

```arc
// Cryptographic API
import "kernel/crypto/core"
import "kernel/crypto/aes"
import "kernel/crypto/des"
import "kernel/crypto/rsa"
import "kernel/crypto/hash"     // SHA, MD5, etc
import "kernel/crypto/hmac"
import "kernel/crypto/random"   // Secure RNG
import "kernel/crypto/cipher"

// Security modules
import "kernel/security/selinux"
import "kernel/security/apparmor"
import "kernel/security/smack"
import "kernel/security/tomoyo"
import "kernel/security/seccomp" // Syscall filtering

// Key management
import "kernel/keys/keyring"
import "kernel/keys/encrypted"

// TPM (Trusted Platform Module)
import "kernel/security/tpm"
import "kernel/security/tpm/tpm2"

// Integrity measurement
import "kernel/security/ima"    // Integrity Measurement Architecture
import "kernel/security/evm"    // Extended Verification Module
```

### Audio 🔊

```arc
// Advanced Linux Sound Architecture
import "kernel/sound/core"
import "kernel/sound/pcm"       // Digital audio streams
import "kernel/sound/control"   // Mixer controls
import "kernel/sound/rawmidi"   // Raw MIDI

// Hardware interfaces
import "kernel/sound/pci"       // PCI sound cards
import "kernel/sound/usb"       // USB audio
import "kernel/sound/hdmi"      // HDMI audio

// Sequencer (MIDI)
import "kernel/sound/seq"
import "kernel/sound/seq/oss"   // OSS compatibility

// Timers
import "kernel/sound/timer"

// Compress offload (hardware audio processing)
import "kernel/sound/compress"
```

### Power Management ⚡

```arc
// CPU frequency scaling
import "kernel/power/cpufreq"
import "kernel/power/cpufreq/governor"

// Device power management
import "kernel/power/pm/device"
import "kernel/power/pm/runtime"  // Runtime PM
import "kernel/power/pm/qos"      // Quality of Service

// System sleep states
import "kernel/power/suspend"
import "kernel/power/hibernate"

// Thermal management
import "kernel/thermal/core"
import "kernel/thermal/governor"
import "kernel/thermal/cooling"

// Battery & AC adapter
import "kernel/power/supply"

// Idle management
import "kernel/power/cpuidle"
```

### Interrupt Handling ⚡

```arc
// Core interrupt management
import "kernel/irq/core"
import "kernel/irq/chip"
import "kernel/irq/domain"

// Threaded interrupts
import "kernel/irq/threaded"

// IRQ affinity (CPU pinning)
import "kernel/irq/affinity"

// GPIO interrupts
import "kernel/irq/gpio"

// MSI/MSI-X
import "kernel/irq/msi"
```

### Timers & Scheduling ⏱️

```arc
// High-resolution timers
import "kernel/time/hrtimer"
import "kernel/time/timer"

// Clock sources
import "kernel/time/clocksource"
import "kernel/time/clockevent"

// Tick management
import "kernel/time/tick"
import "kernel/time/nohz"       // Tickless kernel

// Process scheduler
import "kernel/sched/core"
import "kernel/sched/fair"      // CFS scheduler
import "kernel/sched/rt"        // Real-time scheduling
import "kernel/sched/deadline"  // Deadline scheduling

// CPU isolation
import "kernel/sched/isolation"
import "kernel/sched/topology"

// Wait queues
import "kernel/sched/wait"
```

### Virtualization 🖥️

```arc
// Kernel Virtual Machine
import "kernel/virt/kvm/core"
import "kernel/virt/kvm/x86"
import "kernel/virt/kvm/arm"
import "kernel/virt/kvm/mmu"    // Memory management
import "kernel/virt/kvm/irqchip"

// Virtual devices
import "kernel/virt/virtio/core"
import "kernel/virt/virtio/pci"
import "kernel/virt/virtio/net"
import "kernel/virt/virtio/block"
import "kernel/virt/virtio/scsi"
import "kernel/virt/virtio/balloon"
import "kernel/virt/virtio/console"
import "kernel/virt/virtio/gpu"

// VFIO (device passthrough)
import "kernel/virt/vfio/core"
import "kernel/virt/vfio/pci"

// Hypervisor interfaces
import "kernel/virt/hyperv"
import "kernel/virt/xen"
```

### Platform & Firmware 🔧

```arc
// ACPI (power, device enumeration)
import "kernel/acpi/core"
import "kernel/acpi/bus"
import "kernel/acpi/scan"
import "kernel/acpi/power"

// Device tree (ARM, embedded)
import "kernel/of/core"         // Open Firmware
import "kernel/of/device"
import "kernel/of/platform"

// UEFI runtime services
import "kernel/efi/runtime"
import "kernel/efi/vars"

// BIOS interfaces
import "kernel/firmware/dmi"    // Desktop Management Interface
```

### IPC & Synchronization 🔄

```arc
// Futex (fast userspace mutex)
import "kernel/futex/core"

// Semaphores, mutexes, spinlocks
import "kernel/sync/mutex"
import "kernel/sync/semaphore"
import "kernel/sync/spinlock"
import "kernel/sync/rwlock"     // Read-write locks
import "kernel/sync/rwsem"      // Read-write semaphores

// Wait queues
import "kernel/sync/wait"
import "kernel/sync/completion"

// Read-Copy-Update (lockless data structures)
import "kernel/sync/rcu"

// Atomic operations
import "kernel/sync/atomic"
```

### Tracing & Debugging 🔍

```arc
// Function tracing
import "kernel/trace/ftrace"
import "kernel/trace/function"

// Event tracing
import "kernel/trace/events"
import "kernel/trace/ring_buffer"

// Performance counters
import "kernel/perf/core"
import "kernel/perf/event"

// Kernel probes (dynamic instrumentation)
import "kernel/trace/kprobes"
import "kernel/trace/uprobes"   // Userspace probes
import "kernel/trace/tracepoint"

// BPF (programmable packet filtering & tracing)
import "kernel/bpf/core"
import "kernel/bpf/verifier"
import "kernel/bpf/syscall"

// Kernel debugging
import "kernel/debug/kgdb"      // Kernel debugger
import "kernel/debug/kdb"       // Kernel debugger
```

### I2C & SPI 📡

```arc
// I2C bus
import "kernel/i2c/core"
import "kernel/i2c/adapter"
import "kernel/i2c/device"

// SPI bus
import "kernel/spi/core"
import "kernel/spi/master"
import "kernel/spi/device"
```

### GPIO & Hardware Control 🔌

```arc
// GPIO (General Purpose I/O)
import "kernel/gpio/core"
import "kernel/gpio/chip"
import "kernel/gpio/consumer"

// PWM (Pulse Width Modulation)
import "kernel/pwm/core"

// Watchdog timers
import "kernel/watchdog/core"
```

### Serial & TTY 📟

```arc
// TTY layer
import "kernel/tty/core"
import "kernel/tty/serial"
import "kernel/tty/pty"         // Pseudo-terminals

// Serial drivers
import "kernel/serial/8250"     // Classic UART
import "kernel/serial/core"
```

### NVME & Storage 💽

```arc
// NVMe (Non-Volatile Memory Express)
import "kernel/nvme/core"
import "kernel/nvme/pci"
import "kernel/nvme/fabrics"    // NVMe over Fabrics

// SCSI
import "kernel/scsi/core"
import "kernel/scsi/disk"
import "kernel/scsi/cdrom"

// SATA
import "kernel/ata/core"
import "kernel/ata/libata"

// MD (Multiple Devices - RAID)
import "kernel/md/core"
import "kernel/md/raid"
```

---

## Complete Examples

### Example 1: Network Packet Filter (Firewall)

```arc
namespace firewall

import "kernel/net/netfilter"
import "kernel/net/ipv4"
import "kernel/net/tcp"
import "kernel/trace/events"
import "std/io"

// Blocked IP addresses
let blocked_ips: map<uint32, bool> = {}

func init() int32 {
    io.log("Loading firewall module...")
    
    // Add some IPs to blocklist
    blocked_ips[ip_to_u32("192.168.1.100")] = true
    blocked_ips[ip_to_u32("10.0.0.50")] = true
    
    // Register netfilter hook at PREROUTING (before routing decision)
    let hook = netfilter.register_hook(
        pf: netfilter.NFPROTO_IPV4,
        hooknum: netfilter.NF_INET_PRE_ROUTING,
        priority: netfilter.NF_IP_PRI_FIRST,
        callback: &packet_filter
    )
    
    if hook == null {
        io.log("Failed to register netfilter hook")
        return -1
    }
    
    io.log("Firewall loaded - blocking ${blocked_ips.len()} IPs")
    return 0
}

func packet_filter(priv: *void, skb: *netfilter.sk_buff, state: *netfilter.nf_hook_state) uint32 {
    // Parse IP header
    let iph = ipv4.ip_hdr(skb)
    let src_ip = iph.saddr
    let dst_ip = iph.daddr
    let protocol = iph.protocol
    
    // Check if source IP is blocked
    if blocked_ips.contains(src_ip) {
        trace.log("Blocked packet from ${u32_to_ip(src_ip)}")
        return netfilter.NF_DROP  // Drop the packet
    }
    
    // Block all SSH attempts from outside local network
    if protocol == ipv4.IPPROTO_TCP {
        let tcph = tcp.tcp_hdr(skb)
        let dst_port = ntohs(tcph.dest)
        
        if dst_port == 22 && !is_local_ip(src_ip) {
            trace.log("Blocked SSH from ${u32_to_ip(src_ip)}")
            return netfilter.NF_DROP
        }
    }
    
    // Allow packet to continue
    return netfilter.NF_ACCEPT
}

func cleanup() {
    io.log("Unloading firewall module")
    netfilter.unregister_hook(&packet_filter)
}

// Helpers
func ip_to_u32(ip: string) uint32 {
    // Parse "192.168.1.1" to uint32
    // Implementation omitted for brevity
    return 0
}

func u32_to_ip(ip: uint32) string {
    let a = (ip >> 24) & 0xFF
    let b = (ip >> 16) & 0xFF
    let c = (ip >> 8) & 0xFF
    let d = ip & 0xFF
    return "${a}.${b}.${c}.${d}"
}

func is_local_ip(ip: uint32) bool {
    // Check if 192.168.x.x or 10.x.x.x
    let a = (ip >> 24) & 0xFF
    let b = (ip >> 16) & 0xFF
    return (a == 192 && b == 168) || (a == 10)
}

func ntohs(val: uint16) uint16 {
    // Network to host byte order (big endian to little endian)
    return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}
```

### Example 2: Custom Block Device (RAM Disk)

```arc
namespace ramdisk

import "kernel/block/device"
import "kernel/block/queue"
import "kernel/block/bio"
import "kernel/mm/page"
import "kernel/mm/alloc"
import "std/io"

const RAMDISK_SIZE: uint64 = 100 * 1024 * 1024  // 100MB
const SECTOR_SIZE: uint32 = 512

let ramdisk_data: *byte = null
let gendisk: *block.gendisk = null

func init() int32 {
    io.log("Creating RAM disk (${RAMDISK_SIZE / 1024 / 1024}MB)...")
    
    // Allocate memory for the disk
    ramdisk_data = alloc.vmalloc(RAMDISK_SIZE)
    if ramdisk_data == null {
        io.log("Failed to allocate memory")
        return -1
    }
    
    // Zero the memory
    for let i: uint64 = 0; i < RAMDISK_SIZE; i++ {
        ramdisk_data[i] = 0
    }
    
    // Register block device
    let major = block.register_blkdev(0, "ramdisk")
    if major < 0 {
        alloc.vfree(ramdisk_data)
        return major
    }
    
    // Create request queue
    let queue = block.blk_alloc_queue(null)
    queue.set_make_request_fn(&ramdisk_make_request)
    queue.set_logical_block_size(SECTOR_SIZE)
    
    // Allocate gendisk
    gendisk = block.alloc_disk(1)  // 1 minor number
    gendisk.major = major
    gendisk.first_minor = 0
    gendisk.fops = &ramdisk_fops
    gendisk.queue = queue
    gendisk.set_capacity(RAMDISK_SIZE / SECTOR_SIZE)
    gendisk.set_disk_name("ramdisk0")
    
    // Add disk to system
    block.add_disk(gendisk)
    
    io.log("RAM disk registered as /dev/ramdisk0")
    return 0
}

func ramdisk_make_request(queue: *block.request_queue, bio: *block.bio) {
    let sector = bio.bi_iter.bi_sector
    let offset = sector * SECTOR_SIZE
    
    // Iterate over bio segments
    let iter = bio.bi_iter
    for segment in bio.bi_io_vec {
        let page_data = page.kmap_atomic(segment.bv_page)
        let page_offset = segment.bv_offset
        let len = segment.bv_len
        
        if bio.is_write() {
            // Write: copy from page to ramdisk
            for let i: uint32 = 0; i < len; i++ {
                ramdisk_data[offset + cast<uint64>(i)] = page_data[page_offset + i]
            }
        } else {
            // Read: copy from ramdisk to page
            for let i: uint32 = 0; i < len; i++ {
                page_data[page_offset + i] = ramdisk_data[offset + cast<uint64>(i)]
            }
        }
        
        page.kunmap_atomic(page_data)
        offset += cast<uint64>(len)
    }
    
    bio.endio(bio, 0)  // Complete successfully
}

let ramdisk_fops = block.block_device_operations{
    owner: null,
    open: &ramdisk_open,
    release: &ramdisk_release
}

func ramdisk_open(bdev: *block.block_device, mode: uint32) int32 {
    return 0
}

func ramdisk_release(disk: *block.gendisk, mode: uint32) {
    // Nothing to do
}

func cleanup() {
    if gendisk != null {
        block.del_gendisk(gendisk)
        block.put_disk(gendisk)
    }
    
    if ramdisk_data != null {
        alloc.vfree(ramdisk_data)
    }
    
    block.unregister_blkdev(gendisk.major, "ramdisk")
    io.log("RAM disk unloaded")
}
```

### Example 3: USB Device Driver (Custom Hardware)

```arc
namespace usb_sensor

import "kernel/usb/core"
import "kernel/usb/urb"
import "kernel/input/core"
import "std/io"

const VENDOR_ID: uint16 = 0x1234
const PRODUCT_ID: uint16 = 0x5678

struct SensorDevice {
    usb_dev: *usb.device
    input_dev: *input.input_dev
    endpoint_in: byte
    endpoint_out: byte
    urb: *usb.urb
}

let sensor_devices: vector<*SensorDevice> = {}

func init() int32 {
    io.log("Registering USB sensor driver...")
    
    let driver = usb.usb_driver{
        name: "sensor_driver",
        probe: &sensor_probe,
        disconnect: &sensor_disconnect,
        id_table: &sensor_id_table
    }
    
    return usb.usb_register(&driver)
}

let sensor_id_table = array<usb.usb_device_id, 2>{
    usb.usb_device_id{
        match_flags: usb.USB_DEVICE_ID_MATCH_DEVICE,
        idVendor: VENDOR_ID,
        idProduct: PRODUCT_ID
    },
    usb.usb_device_id{} // Terminator
}

func sensor_probe(intf: *usb.usb_interface, id: *usb.usb_device_id) int32 {
    let usb_dev = intf.to_usb_device()
    
    io.log("USB sensor connected: ${usb_dev.product_name()}")
    
    // Allocate device structure
    let sensor = SensorDevice{
        usb_dev: usb_dev
    }
    
    // Find bulk endpoints
    let iface_desc = intf.cur_altsetting
    for let i: byte = 0; i < iface_desc.desc.bNumEndpoints; i++ {
        let endpoint = &iface_desc.endpoint[i].desc
        
        if usb.usb_endpoint_is_bulk_in(endpoint) {
            sensor.endpoint_in = endpoint.bEndpointAddress
        } else if usb.usb_endpoint_is_bulk_out(endpoint) {
            sensor.endpoint_out = endpoint.bEndpointAddress
        }
    }
    
    if sensor.endpoint_in == 0 || sensor.endpoint_out == 0 {
        io.log("Failed to find required endpoints")
        return -1
    }
    
    // Create input device for sensor data
    sensor.input_dev = input.allocate_device()
    sensor.input_dev.name = "USB Sensor"
    sensor.input_dev.phys = usb_dev.devpath
    
    // Set supported events (absolute coordinates)
    input.set_capability(sensor.input_dev, input.EV_ABS, input.ABS_X)
    input.set_capability(sensor.input_dev, input.EV_ABS, input.ABS_Y)
    input.set_capability(sensor.input_dev, input.EV_ABS, input.ABS_Z)
    
    // Register input device
    if input.register_device(sensor.input_dev) != 0 {
        input.free_device(sensor.input_dev)
        return -1
    }
    
    // Allocate and submit URB for reading
    sensor.urb = usb.usb_alloc_urb(0)
    let buffer: *byte = alloc.kmalloc(64, alloc.GFP_KERNEL)
    
    usb.usb_fill_bulk_urb(
        sensor.urb,
        usb_dev,
        usb.usb_rcvbulkpipe(usb_dev, sensor.endpoint_in),
        buffer,
        64,
        &sensor_read_callback,
        &sensor
    )
    
    usb.usb_submit_urb(sensor.urb, alloc.GFP_KERNEL)
    
    sensor_devices.push(&sensor)
    usb.usb_set_intfdata(intf, &sensor)
    
    io.log("Sensor initialized successfully")
    return 0
}

func sensor_read_callback(urb: *usb.urb) {
    let sensor = cast<*SensorDevice>(urb.context)
    let data = cast<*byte>(urb.transfer_buffer)
    
    if urb.status == 0 {
        // Parse sensor data (example: 3 int16 values for X, Y, Z)
        let x = cast<int16>((cast<uint16>(data[0]) << 8) | data[1])
        let y = cast<int16>((cast<uint16>(data[2]) << 8) | data[3])
        let z = cast<int16>((cast<uint16>(data[4]) << 8) | data[5])
        
        // Report to input subsystem
        input.report_abs(sensor.input_dev, input.ABS_X, cast<int32>(x))
        input.report_abs(sensor.input_dev, input.ABS_Y, cast<int32>(y))
        input.report_abs(sensor.input_dev, input.ABS_Z, cast<int32>(z))
        input.sync(sensor.input_dev)
    }
    
    // Resubmit URB to continue reading
    usb.usb_submit_urb(urb, alloc.GFP_ATOMIC)
}

func sensor_disconnect(intf: *usb.usb_interface) {
    let sensor = cast<*SensorDevice>(usb.usb_get_intfdata(intf))
    
    io.log("USB sensor disconnected")
    
    // Cancel URB
    usb.usb_kill_urb(sensor.urb)
    usb.usb_free_urb(sensor.urb)
    
    // Unregister input device
    input.unregister_device(sensor.input_dev)
    
    // Remove from list
    for let i: usize = 0; i < sensor_devices.len(); i++ {
        if sensor_devices[i] == sensor {
            sensor_devices.remove(i)
            break
        }
    }
}

func cleanup() {
    usb.usb_deregister(&driver)
    io.log("USB sensor driver unloaded")
}
```

### Example 4: Simple Hypervisor (KVM)

```arc
namespace mini_vm

import "kernel/virt/kvm/core"
import "kernel/virt/kvm/x86"
import "kernel/mm/page"
import "std/io"

const GUEST_MEM_SIZE: uint64 = 512 * 1024 * 1024  // 512MB

func create_and_run_vm() {
    io.log("Creating VM...")
    
    // Create VM
    let vm = kvm.create_vm()
    
    // Allocate guest memory
    let guest_mem = page.alloc_pages(GUEST_MEM_SIZE / page.PAGE_SIZE)
    
    // Set memory region
    vm.set_user_memory_region(
        slot: 0,
        guest_phys_addr: 0,
        memory_size: GUEST_MEM_SIZE,
        userspace_addr: cast<uint64>(guest_mem)
    )
    
    // Create vCPU
    let vcpu = vm.create_vcpu(0)
    
    // Setup CPU state for 16-bit real mode (boot)
    setup_boot_state(vcpu)
    
    // Load guest code at 0x1000
    let guest_code: array<byte, 32> = {
        // Simple loop that outputs 'H' to port 0xE9 (QEMU debug port)
        0xB0, 0x48,                 // mov al, 'H'
        0xE6, 0xE9,                 // out 0xE9, al
        0xB0, 0x69,                 // mov al, 'i'
        0xE6, 0xE9,                 // out 0xE9, al
        0xB0, 0x0A,                 // mov al, '\n'
        0xE6, 0xE9,                 // out 0xE9, al
        0xF4                        // hlt
    }
    
    // Copy code to guest memory
    for let i: usize = 0; i < guest_code.len(); i++ {
        guest_mem[0x1000 + i] = guest_code[i]
    }
    
    io.log("Running VM...")
    
    // Main VM loop
    for {
        let run = vcpu.run()
        
        switch run.exit_reason {
            case kvm.KVM_EXIT_HLT:
                io.log("VM halted")
                break
                
            case kvm.KVM_EXIT_IO:
                handle_io(run.io)
                
            case kvm.KVM_EXIT_MMIO:
                io.log("MMIO access at 0x${run.mmio.phys_addr:X}")
                
            case kvm.KVM_EXIT_SHUTDOWN:
                io.log("VM shutdown")
                break
                
            default:
                io.log("Unhandled exit: ${run.exit_reason}")
                break
        }
    }
    
    io.log("VM stopped")
}

func setup_boot_state(vcpu: *kvm.vcpu) {
    // Setup registers
    let regs = kvm.x86.kvm_regs{}
    regs.rip = 0x1000       // Start at 0x1000
    regs.rflags = 0x2       // Reserved bit
    vcpu.set_regs(&regs)
    
    // Setup segment registers for real mode
    let sregs = kvm.x86.kvm_sregs{}
    
    // CS segment
    sregs.cs.base = 0
    sregs.cs.limit = 0xFFFF
    sregs.cs.selector = 0
    sregs.cs.type = 0xB     // Code segment
    sregs.cs.present = 1
    sregs.cs.s = 1          // Code/data segment
    
    // DS, ES, SS segments (similar setup)
    sregs.ds.base = 0
    sregs.ds.limit = 0xFFFF
    sregs.ds.selector = 0
    sregs.ds.type = 3
    sregs.ds.present = 1
    sregs.ds.s = 1
    
    sregs.es = sregs.ds
    sregs.ss = sregs.ds
    
    vcpu.set_sregs(&sregs)
}

func handle_io(io: kvm.x86.kvm_run_io) {
    if io.direction == kvm.KVM_EXIT_IO_OUT && io.port == 0xE9 {
        // QEMU debug port - print character
        let data = cast<*byte>(&io.data_offset)
        io.write_stdout(data, io.size)
    }
}
```

### Example 5: Real-Time Packet Analyzer (BPF + Netfilter)

```arc
namespace packet_analyzer

import "kernel/net/netfilter"
import "kernel/net/ipv4"
import "kernel/net/tcp"
import "kernel/net/udp"
import "kernel/bpf/core"
import "kernel/trace/events"
import "std/io"

struct PacketStats {
    total_packets: uint64
    tcp_packets: uint64
    udp_packets: uint64
    icmp_packets: uint64
    total_bytes: uint64
}

let stats: PacketStats = {}
let start_time: uint64 = 0

func init() int32 {
    io.log("Starting packet analyzer...")
    
    start_time = ktime_get_ns()
    
    // Hook at PREROUTING to see all incoming packets
    let hook = netfilter.register_hook(
        pf: netfilter.NFPROTO_IPV4,
        hooknum: netfilter.NF_INET_PRE_ROUTING,
        priority: netfilter.NF_IP_PRI_FIRST,
        callback: &analyze_packet
    )
    
    if hook == null {
        return -1
    }
    
    io.log("Packet analyzer active")
    return 0
}

func analyze_packet(priv: *void, skb: *netfilter.sk_buff, state: *netfilter.nf_hook_state) uint32 {
    stats.total_packets++
    stats.total_bytes += cast<uint64>(skb.len)
    
    // Parse IP header
    let iph = ipv4.ip_hdr(skb)
    let protocol = iph.protocol
    let src_ip = ntohl(iph.saddr)
    let dst_ip = ntohl(iph.daddr)
    
    switch protocol {
        case ipv4.IPPROTO_TCP:
            stats.tcp_packets++
            analyze_tcp(skb, src_ip, dst_ip)
            
        case ipv4.IPPROTO_UDP:
            stats.udp_packets++
            analyze_udp(skb, src_ip, dst_ip)
            
        case ipv4.IPPROTO_ICMP:
            stats.icmp_packets++
    }
    
    // Log every 10000 packets
    if stats.total_packets % 10000 == 0 {
        print_stats()
    }
    
    return netfilter.NF_ACCEPT
}

func analyze_tcp(skb: *netfilter.sk_buff, src_ip: uint32, dst_ip: uint32) {
    let tcph = tcp.tcp_hdr(skb)
    let src_port = ntohs(tcph.source)
    let dst_port = ntohs(tcph.dest)
    let flags = tcph.flags
    
    // Detect SYN flood attack
    if flags & tcp.TH_SYN != 0 && flags & tcp.TH_ACK == 0 {
        trace.log("SYN: ${ip_to_str(src_ip)}:${src_port} -> ${ip_to_str(dst_ip)}:${dst_port}")
    }
    
    // Detect port scanning
    if dst_port < 1024 && flags & tcp.TH_SYN != 0 {
        trace.log("Port scan detected: ${ip_to_str(src_ip)} -> ${dst_port}")
    }
}

func analyze_udp(skb: *netfilter.sk_buff, src_ip: uint32, dst_ip: uint32) {
    let udph = udp.udp_hdr(skb)
    let src_port = ntohs(udph.source)
    let dst_port = ntohs(udph.dest)
    
    // Log DNS queries
    if dst_port == 53 {
        trace.log("DNS query from ${ip_to_str(src_ip)}")
    }
}

func print_stats() {
    let elapsed = (ktime_get_ns() - start_time) / 1000000000  // Convert to seconds
    let pps = stats.total_packets / elapsed  // Packets per second
    let mbps = (stats.total_bytes * 8) / (elapsed * 1000000)  // Megabits per second
    
    io.log("=== Packet Statistics ===")
    io.log("Total packets: ${stats.total_packets}")
    io.log("TCP: ${stats.tcp_packets}, UDP: ${stats.udp_packets}, ICMP: ${stats.icmp_packets}")
    io.log("Total bytes: ${stats.total_bytes}")
    io.log("Throughput: ${pps} pps, ${mbps} Mbps")
}

func cleanup() {
    print_stats()
    netfilter.unregister_hook(&analyze_packet)
    io.log("Packet analyzer stopped")
}

// Helpers
func ntohl(val: uint32) uint32 {
    return ((val & 0xFF) << 24) | ((val & 0xFF00) << 8) | 
           ((val & 0xFF0000) >> 8) | ((val >> 24) & 0xFF)
}

func ntohs(val: uint16) uint16 {
    return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}

func ip_to_str(ip: uint32) string {
    let a = (ip >> 24) & 0xFF
    let b = (ip >> 16) & 0xFF
    let c = (ip >> 8) & 0xFF
    let d = ip & 0xFF
    return "${a}.${b}.${c}.${d}"
}

func ktime_get_ns() uint64 {
    // Get current time in nanoseconds
    return 0  // Implementation from kernel time subsystem
}
```

---

## Cross-Platform Abstraction Layer

```arc
// Create your own cross-platform modules
namespace "github.com/yourname/kernel/device"

#if target_os(linux)
    import "kernel/block/device" as backend
#elif target_os(windows)
    import "github.com/arc-lang/windows/kernel/io" as backend
#elif target_os(darwin)
    import "github.com/arc-lang/darwin/kernel/iokit" as backend
#endif

// Unified interface that works everywhere
struct Device {
    #if target_os(linux)
        linux_dev: *backend.block_device
    #elif target_os(windows)
        windows_dev: *backend.DEVICE_OBJECT
    #elif target_os(darwin)
        darwin_dev: *backend.IOService
    #endif
}

func create(name: string) *Device {
    let dev = Device{}
    
    #if target_os(linux)
        dev.linux_dev = backend.register_blkdev(0, name)
    #elif target_os(windows)
        backend.IoCreateDevice(...)
    #elif target_os(darwin)
        // IOKit registration
    #endif
    
    return &dev
}
```

---

## The Power of Organization

**Traditional kernel development:**
```c
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/fs.h>
#include <linux/cdev.h>
#include <linux/device.h>
#include <linux/slab.h>
// ... 50 more includes
```

**arc kernel development:**
```arc
import "kernel/fs/vfs"
import "kernel/block/device"
import "kernel/mm/alloc"
```

**Same machine code. Better organization.** 🚀

That's the entire point - making kernel programming **accessible** without sacrificing power or performance.