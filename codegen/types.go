// codegen/types.go
package codegen

import (
	"errors"

	"github.com/arc-language/arc-lang/ast"
	"github.com/arc-language/arc-lang/builder/types"
)

// TypeGenerator maps Arc AST types to IR types.
// It also maintains a registry of named struct field names so that
// struct GEPs can be resolved by field name rather than index.
type TypeGenerator struct {
	// structFields maps struct name → ordered field names.
	// Populated by RegisterStruct when InterfaceDecls are processed.
	structFields map[string][]string
}

func NewTypeGenerator() *TypeGenerator {
	return &TypeGenerator{
		structFields: make(map[string][]string),
	}
}

// RegisterStruct records the field name order for a named struct so that
// FieldIndex can resolve member access by name.
func (tg *TypeGenerator) RegisterStruct(name string, fieldNames []string) {
	names := make([]string, len(fieldNames))
	copy(names, fieldNames)
	tg.structFields[name] = names
}

// FieldIndex returns the zero-based field index for fieldName in a struct
// named structName, or -1 if not found.
func (tg *TypeGenerator) FieldIndex(structName, fieldName string) int {
	if fields, ok := tg.structFields[structName]; ok {
		for i, n := range fields {
			if n == fieldName {
				return i
			}
		}
	}
	return -1
}

// GenType translates an Arc AST TypeRef into an IR types.Type.
func (tg *TypeGenerator) GenType(t ast.TypeRef) types.Type {
	if t == nil {
		return types.Void
	}

	switch typ := t.(type) {

	case *ast.NamedType:
		return tg.genNamedType(typ.Name)

	case *ast.PointerType:
		elem := tg.GenType(typ.Base)
		return types.NewPointer(elem)

	case *ast.MutRefType:
		// &mut T — represented as a plain pointer in the IR.
		elem := tg.GenType(typ.Base)
		return types.NewPointer(elem)

	case *ast.RefType:
		// &T — also a pointer.
		elem := tg.GenType(typ.Base)
		return types.NewPointer(elem)

	case *ast.ArrayType:
		// [N]T — fixed-size array.
		// The length expression must resolve to a constant; we handle only
		// BasicLit integers here and default to 0 otherwise.
		elem := tg.GenType(typ.Elem)
		length := int64(0)
		if bl, ok := typ.Len.(*ast.BasicLit); ok && bl.Kind == "INT" {
			if n, err := parseInt(bl.Value); err == nil {
				length = n
			}
		}
		return types.NewArray(elem, length)

	case *ast.SliceType:
		// []T — { *T, i64 } (data pointer + length)
		elem := tg.GenType(typ.Elem)
		return types.NewStruct("slice", []types.Type{
			types.NewPointer(elem),
			types.I64,
		}, false)

	case *ast.VectorType:
		// vector[T] — { *T, i64, i64 } (data, length, capacity)
		elem := tg.GenType(typ.Elem)
		return types.NewStruct("vector", []types.Type{
			types.NewPointer(elem),
			types.I64,
			types.I64,
		}, false)

	case *ast.MapType:
		// map[K]V — opaque pointer to runtime map structure.
		return types.NewPointer(types.Void)

	case *ast.TupleType:
		fieldTypes := make([]types.Type, len(typ.Elems))
		for i, e := range typ.Elems {
			fieldTypes[i] = tg.GenType(e)
		}
		return types.NewStruct("", fieldTypes, false)

	case *ast.FuncType:
		var params []types.Type
		for _, p := range typ.Params {
			params = append(params, tg.GenType(p))
		}
		var ret types.Type = types.Void
		if len(typ.Results) > 0 {
			ret = tg.GenType(typ.Results[0])
		}
		fnType := types.NewFunction(ret, params, false)
		fnType.IsAsync = typ.IsAsync
		return types.NewPointer(fnType) // function values are pointers
	}

	return types.Void
}

// GenExternType translates a TypeRef that appears inside an extern block
// (e.g. parameter and return types of ExternFunc) into an IR type.
// Extern-specific type nodes such as PointerType and RefType are legal here;
// everything else delegates to GenType.
func (tg *TypeGenerator) GenExternType(t ast.TypeRef) types.Type {
	if t == nil {
		return types.Void
	}
	return tg.GenType(t)
}

func (tg *TypeGenerator) genNamedType(name string) types.Type {
	switch name {
	// Signed integers
	case "int8":
		return types.I8
	case "int16":
		return types.I16
	case "int32", "int":
		return types.I32
	case "int64":
		return types.I64
	// Unsigned integers
	case "uint8", "byte":
		return types.U8
	case "uint16":
		return types.U16
	case "uint32":
		return types.U32
	case "uint64":
		return types.U64
	// Platform-width integers
	case "usize":
		return types.U64
	case "isize":
		return types.I64
	// Floats
	case "float32":
		return types.F32
	case "float64", "float":
		return types.F64
	// Special
	case "bool":
		return types.I1
	case "char":
		return types.U32 // Unicode scalar value
	case "void":
		return types.Void
	case "string":
		// Arc strings are fat pointers: { *u8, i64 }
		return types.NewStruct("string", []types.Type{
			types.NewPointer(types.U8),
			types.I64,
		}, false)
	case "ThreadHandle":
		return types.NewPointer(types.Void)
	default:
		// Named struct — return an opaque reference; the actual definition
		// is looked up in the module's type table during codegen.
		return types.NewStruct(name, nil, false)
	}
}

// errNotInt is returned by parseInt when the input is not a valid integer.
var errNotInt = errors.New("not an integer")

// parseInt is a small helper used for array-length literals.
func parseInt(s string) (int64, error) {
	var v int64
	if len(s) > 2 && s[:2] == "0x" {
		for _, ch := range s[2:] {
			v <<= 4
			switch {
			case ch >= '0' && ch <= '9':
				v |= int64(ch - '0')
			case ch >= 'a' && ch <= 'f':
				v |= int64(ch-'a') + 10
			case ch >= 'A' && ch <= 'F':
				v |= int64(ch-'A') + 10
			}
		}
		return v, nil
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, errNotInt
		}
		v = v*10 + int64(ch-'0')
	}
	return v, nil
}