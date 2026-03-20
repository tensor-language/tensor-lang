// Package ir - Async Smart Thread instruction definitions
package ir

import (
	"fmt"
	"strings"
)

// AsyncTaskCreateInst represents spinning up a new smart thread.
// %handle = async_task_create @function(args...)
type AsyncTaskCreateInst struct {
	BaseInstruction
	Callee    *Function
	CalleeVal Value // For indirect calls or if Callee is nil
}

func (i *AsyncTaskCreateInst) String() string {
	target := ""
	if i.Callee != nil {
		target = "@" + i.Callee.Name()
	} else if i.CalleeVal != nil {
		target = formatOp(i.CalleeVal)
	} else {
		target = "<null>"
	}

	var args []string
	for _, op := range i.Ops {
		args = append(args, fmt.Sprintf("%s %s", op.Type(), formatOp(op)))
	}

	return fmt.Sprintf("%%%s = async_task_create %s(%s)", 
		i.ValName, target, strings.Join(args, ", "))
}

// AsyncTaskAwaitInst represents blocking until a task completes.
// %result = async_task_await %handle
type AsyncTaskAwaitInst struct {
	BaseInstruction
}

func (i *AsyncTaskAwaitInst) String() string {
	if len(i.Ops) == 0 {
		return fmt.Sprintf("%%%s = async_task_await <missing operand>", i.ValName)
	}
	handle := i.Ops[0]
	return fmt.Sprintf("%%%s = async_task_await %s %s", 
		i.ValName, handle.Type(), formatOp(handle))
}