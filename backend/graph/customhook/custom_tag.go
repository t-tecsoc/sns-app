package customhook

import (
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
)

func CustomTagHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	c := fd.Directives.ForName("customTag")
	if c != nil {
		key := c.Arguments.ForName("key")
		value := c.Arguments.ForName("value")
		if key != nil && value != nil {
			f.Tag += fmt.Sprintf(" %s:\"%s\"", key.Value.Raw, strings.ReplaceAll(value.Value.Raw, "-", ":"))
		}
	}
	return f, nil
}
