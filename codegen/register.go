// codegen/register.go
package codegen

import "github.com/arc-language/arc-lang/ast"

// registerStructs walks the file and calls TypeGenerator.RegisterStruct for
// every InterfaceDecl, building the field-name index needed by genSelector and
// genCompositeLit. This must be called before genFuncBody so that field lookups
// work even when a struct is used before it is defined (Arc allows this because
// the checker resolved forward references already).
func (cg *Codegen) registerStructs(file *ast.File) {
	for _, decl := range file.Decls {
		if d, ok := decl.(*ast.InterfaceDecl); ok {
			names := make([]string, len(d.Fields))
			for i, f := range d.Fields {
				names[i] = f.Name
			}
			cg.TypeGen.RegisterStruct(d.Name, names)
		}
	}
}