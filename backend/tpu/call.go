// call.go - Call operation emission
package tpu

import (
	"strings"

	"github.com/arc-language/arc-lang/builder/ir"
)

func (g *Generator) emitCall(i *ir.CallInst) error {
	var calleeName string
	if i.Callee != nil {
		calleeName = i.Callee.Name()
	} else if i.CalleeName != "" {
		calleeName = i.CalleeName
	} else {
		calleeName = "indirect_call"
	}

	// Check for known math intrinsics
	if hloOp, ok := g.mathIntrinsicToHlo(calleeName); ok {
		return g.emitMathIntrinsic(i, hloOp)
	}

	// Regular call - emit as custom-call or fusion
	dst := g.assignName(i)
	var args []string
	for _, op := range i.Operands() {
		args = append(args, g.getOperand(op))
	}

	retType := g.toHloType(i.Type())

	if i.Callee != nil && i.Callee.CallConv == ir.CC_TPU {
		// Call to another TPU function - use call instruction
		g.emit("%s = %s call(%s), to_apply=%s", dst, retType, strings.Join(args, ", "), sanitizeName(calleeName))
	} else {
		// External call - use custom-call
		g.emit("%s = %s custom-call(%s), custom_call_target=\"%s\"",
			dst, retType, strings.Join(args, ", "), calleeName)
	}
	return nil
}

func (g *Generator) mathIntrinsicToHlo(name string) (string, bool) {
	intrinsics := map[string]string{
		"sqrt":  "sqrt",
		"sqrtf": "sqrt",
		"exp":   "exponential",
		"expf":  "exponential",
		"log":   "log",
		"logf":  "log",
		"sin":   "sine",
		"sinf":  "sine",
		"cos":   "cosine",
		"cosf":  "cosine",
		"tanh":  "tanh",
		"tanhf": "tanh",
		"abs":   "abs",
		"absf":  "abs",
		"fabs":  "abs",
		"floor": "floor",
		"ceil":  "ceil",
		"pow":   "power",
		"powf":  "power",
		"max":   "maximum",
		"fmax":  "maximum",
		"min":   "minimum",
		"fmin":  "minimum",
	}
	op, ok := intrinsics[name]
	return op, ok
}

func (g *Generator) emitMathIntrinsic(i *ir.CallInst, hloOp string) error {
	dst := g.assignName(i)
	typ := g.toHloType(i.Type())

	var args []string
	for _, op := range i.Operands() {
		args = append(args, g.getOperand(op))
	}

	g.emit("%s = %s %s(%s)", dst, typ, hloOp, strings.Join(args, ", "))
	return nil
}