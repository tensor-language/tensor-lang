// types.go - Type conversion and shape handling
package tpu

import (
	"fmt"
	"strings"

	"github.com/arc-language/arc-lang/builder/types"
)

// Shape represents tensor dimensions
type Shape struct {
	Dims     []int64
	ElemType string
}

func (s Shape) String() string {
	if len(s.Dims) == 0 {
		return s.ElemType + "[]"
	}
	dims := make([]string, len(s.Dims))
	for i, d := range s.Dims {
		dims[i] = fmt.Sprintf("%d", d)
	}
	return fmt.Sprintf("%s[%s]", s.ElemType, strings.Join(dims, ","))
}

func (g *Generator) toHloType(t types.Type) string {
	if t == nil {
		return "()" // void/unit type
	}

	switch ty := t.(type) {
	case *types.IntType:
		if ty.Signed {
			switch ty.BitWidth {
			case 1:
				return "pred[]"
			case 8:
				return "s8[]"
			case 16:
				return "s16[]"
			case 32:
				return "s32[]"
			case 64:
				return "s64[]"
			default:
				return "s32[]"
			}
		} else {
			switch ty.BitWidth {
			case 1:
				return "pred[]"
			case 8:
				return "u8[]"
			case 16:
				return "u16[]"
			case 32:
				return "u32[]"
			case 64:
				return "u64[]"
			default:
				return "u32[]"
			}
		}

	case *types.FloatType:
		switch ty.BitWidth {
		case 16:
			return "bf16[]" // Use bfloat16 by default on TPU
		case 32:
			return "f32[]"
		case 64:
			return "f64[]"
		default:
			return "f32[]"
		}

	case *types.PointerType:
		// Pointers become the underlying tensor type
		return g.toHloType(ty.ElementType)

	case *types.ArrayType:
		elemType := g.toHloElemType(ty.ElementType)
		return fmt.Sprintf("%s[%d]", elemType, ty.Length)

	case *types.StructType:
		var fields []string
		for _, f := range ty.Fields {
			fields = append(fields, g.toHloType(f))
		}
		return fmt.Sprintf("(%s)", strings.Join(fields, ", "))

	case *types.VoidType:
		return "()"

	case *types.FunctionType:
		// Function types don't directly map
		return "()"

	default:
		return "f32[]"
	}
}

func (g *Generator) toHloElemType(t types.Type) string {
	if t == nil {
		return "f32"
	}

	switch ty := t.(type) {
	case *types.IntType:
		if ty.Signed {
			return fmt.Sprintf("s%d", ty.BitWidth)
		}
		return fmt.Sprintf("u%d", ty.BitWidth)
	case *types.FloatType:
		if ty.BitWidth == 16 {
			return "bf16"
		}
		return fmt.Sprintf("f%d", ty.BitWidth)
	default:
		return "f32"
	}
}

func (g *Generator) inferShape(t types.Type) Shape {
	switch ty := t.(type) {
	case *types.ArrayType:
		elemShape := g.inferShape(ty.ElementType)
		return Shape{
			Dims:     append([]int64{ty.Length}, elemShape.Dims...),
			ElemType: elemShape.ElemType,
		}
	case *types.IntType:
		if ty.Signed {
			return Shape{ElemType: fmt.Sprintf("s%d", ty.BitWidth)}
		}
		return Shape{ElemType: fmt.Sprintf("u%d", ty.BitWidth)}
	case *types.FloatType:
		return Shape{ElemType: fmt.Sprintf("f%d", ty.BitWidth)}
	default:
		return Shape{ElemType: "f32"}
	}
}