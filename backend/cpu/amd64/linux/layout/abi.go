package layout

import (
	"bytes"
	"encoding/binary"
)

// Linux x86-64 Constants
const (
	SysMmap        = 9
	SysMprotect    = 10
	SysRtSigAction = 13
	SysClone3      = 435
	SysFutex       = 202
	SysExit        = 60

	ProtNone      = 0x0
	ProtRead      = 0x1
	ProtWrite     = 0x2
	ProtExec      = 0x4
	
	MapPrivate    = 0x02
	MapAnonymous  = 0x20
	MapStack      = 0x20000 

	CloneVM       = 0x00000100
	CloneFS       = 0x00000200
	CloneFiles    = 0x00000400
	CloneSighand  = 0x00000800
	CloneThread   = 0x00010000
	CloneSysV     = 0x02000000 
	CloneParent   = 0x00008000 

	FutexWait     = 0
	FutexWake     = 1
	FutexPrivate  = 128 
	
	SigSegv       = 11
	SaSigInfo     = 0x00000004
	SaOnStack     = 0x08000000
	SaRestart     = 0x10000000
)

// CloneArgs represents the struct required by the clone3 syscall.
// Size: 88 bytes on Linux 5.5+
type CloneArgs struct {
	Flags      uint64
	Pidfd      uint64
	ChildTid   uint64
	ParentTid  uint64
	ExitSignal uint64
	Stack      uint64
	StackSize  uint64
	Tls        uint64
	SetTid     uint64
	SetTidSize uint64
	Cgroup     uint64
}

func (c *CloneArgs) WriteTo(buf *bytes.Buffer) {
	binary.Write(buf, binary.LittleEndian, c)
}

// SigAction represents the kernel struct for signal handlers.
// Size: 152 bytes (approx)
type SigAction struct {
	Handler  uint64 
	Flags    uint64
	Restorer uint64
	Mask     uint64 
	_        [120]byte // Padding
}