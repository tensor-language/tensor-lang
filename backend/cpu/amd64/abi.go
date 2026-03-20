package amd64

import "github.com/arc-language/arc-lang/builder/types"

// SizeOf returns the size in bytes of a type for amd64
func SizeOf(t types.Type) int {
	switch t.Kind() {
	case types.VoidKind:
		return 0
	case types.IntegerKind:
		bits := t.(*types.IntType).BitWidth
		if bits <= 8 { return 1 }
		if bits <= 16 { return 2 }
		if bits <= 32 { return 4 }
		return 8
	case types.FloatKind:
		bits := t.(*types.FloatType).BitWidth
		if bits <= 32 { return 4 }
		return 8
	case types.PointerKind, types.FunctionKind:
		return 8
	case types.ArrayKind:
		at := t.(*types.ArrayType)
		return int(at.Length) * SizeOf(at.ElementType)
	case types.StructKind:
		st := t.(*types.StructType)

		if st.Packed {
			size := 0
			for _, f := range st.Fields { size += SizeOf(f) }
			if size == 0 { return 1 }
			return size
		}

		size := 0
		for _, f := range st.Fields {
			align := AlignOf(f)
			if size%align != 0 {
				size += align - (size % align)
			}
			size += SizeOf(f)
		}

		if size == 0 { size = 1 }

		sa := AlignOf(st)
		if size%sa != 0 {
			size += sa - (size % sa)
		}
		return size
	default:
		return 8
	}
}

func AlignOf(t types.Type) int {
	switch t.Kind() {
	case types.IntegerKind:
		bits := t.(*types.IntType).BitWidth
		if bits <= 8 { return 1 }
		if bits <= 16 { return 2 }
		if bits <= 32 { return 4 }
		return 8
	case types.StructKind:
		st := t.(*types.StructType)
		if st.Packed {
			return 1
		}
		max := 1
		for _, f := range st.Fields {
			a := AlignOf(f)
			if a > max { max = a }
		}
		return max
	case types.ArrayKind:
		return AlignOf(t.(*types.ArrayType).ElementType)
	default:
		sz := SizeOf(t)
		if sz > 8 { return 8 }
		return sz
	}
}

// GetStructFieldOffset calculates the byte offset of a field at idx.
func GetStructFieldOffset(st *types.StructType, idx int) int {
	off := 0

	for i := 0; i < idx; i++ {
		f := st.Fields[i]
		if !st.Packed {
			a := AlignOf(f)
			if off%a != 0 { off += a - (off % a) }
		}
		off += SizeOf(f)
	}

	// Align the target field itself
	if !st.Packed && idx < len(st.Fields) {
		a := AlignOf(st.Fields[idx])
		if off%a != 0 { off += a - (off % a) }
	}
	return off
}