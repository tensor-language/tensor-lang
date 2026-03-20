// Package ir - intrinsic instruction definitions
package ir

import (
	"fmt"

	"github.com/arc-language/arc-lang/builder/types"
)

// SizeOfInst represents the sizeof intrinsic (compile-time constant)
type SizeOfInst struct {
	BaseInstruction
	QueryType types.Type
}

func (i *SizeOfInst) String() string {
	return fmt.Sprintf("%%%s = sizeof %s", i.ValName, i.QueryType)
}

// AlignOfInst represents the alignof intrinsic (compile-time constant)
type AlignOfInst struct {
	BaseInstruction
	QueryType types.Type
}

func (i *AlignOfInst) String() string {
	return fmt.Sprintf("%%%s = alignof %s", i.ValName, i.QueryType)
}

// MemSetInst represents memset(dest, val, count)
type MemSetInst struct {
	BaseInstruction
}

func (i *MemSetInst) String() string {
	dest := i.Ops[0]
	val := i.Ops[1]
	count := i.Ops[2]
	return fmt.Sprintf("memset %s %s, %s %s, %s %s",
		dest.Type(), formatOp(dest),
		val.Type(), formatOp(val),
		count.Type(), formatOp(count))
}

// MemCpyInst represents memcpy(dest, src, count)
type MemCpyInst struct {
	BaseInstruction
}

func (i *MemCpyInst) String() string {
	dest := i.Ops[0]
	src := i.Ops[1]
	count := i.Ops[2]
	return fmt.Sprintf("memcpy %s %s, %s %s, %s %s",
		dest.Type(), formatOp(dest),
		src.Type(), formatOp(src),
		count.Type(), formatOp(count))
}

// MemMoveInst represents memmove(dest, src, count)
type MemMoveInst struct {
	BaseInstruction
}

func (i *MemMoveInst) String() string {
	dest := i.Ops[0]
	src := i.Ops[1]
	count := i.Ops[2]
	return fmt.Sprintf("memmove %s %s, %s %s, %s %s",
		dest.Type(), formatOp(dest),
		src.Type(), formatOp(src),
		count.Type(), formatOp(count))
}

// StrLenInst represents strlen(str) -> usize
type StrLenInst struct {
	BaseInstruction
}

func (i *StrLenInst) String() string {
	str := i.Ops[0]
	return fmt.Sprintf("%%%s = strlen %s %s",
		i.ValName, str.Type(), formatOp(str))
}

// MemChrInst represents memchr(ptr, val, count) -> *void
type MemChrInst struct {
	BaseInstruction
}

func (i *MemChrInst) String() string {
	ptr := i.Ops[0]
	val := i.Ops[1]
	count := i.Ops[2]
	return fmt.Sprintf("%%%s = memchr %s %s, %s %s, %s %s",
		i.ValName,
		ptr.Type(), formatOp(ptr),
		val.Type(), formatOp(val),
		count.Type(), formatOp(count))
}

// MemCmpInst represents memcmp(ptr1, ptr2, count) -> i32
type MemCmpInst struct {
	BaseInstruction
}

func (i *MemCmpInst) String() string {
	ptr1 := i.Ops[0]
	ptr2 := i.Ops[1]
	count := i.Ops[2]
	return fmt.Sprintf("%%%s = memcmp %s %s, %s %s, %s %s",
		i.ValName,
		ptr1.Type(), formatOp(ptr1),
		ptr2.Type(), formatOp(ptr2),
		count.Type(), formatOp(count))
}

// RaiseInst represents raise(message) - aborts execution
type RaiseInst struct {
	BaseInstruction
}

func (i *RaiseInst) String() string {
	msg := i.Ops[0]
	return fmt.Sprintf("raise %s %s", msg.Type(), formatOp(msg))
}