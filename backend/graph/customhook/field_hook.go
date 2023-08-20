package customhook

import (
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

func FieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	for _, hook := range []func(*ast.Definition, *ast.FieldDefinition, *modelgen.Field) (*modelgen.Field, error){
		ValidationFieldHook,
		CustomTagHook,
	} {
		var err error
		f, err = hook(td, fd, f)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}
