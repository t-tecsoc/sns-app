package customhook

import (
	"fmt"

	"regexp"

	"strings"

	"github.com/99designs/gqlgen/plugin/modelgen"
	"golang.org/x/exp/slices"
)

func ForeignKeyFieldHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	defaultTypeList := []string{"int", "string", "bool"}

	for _, model := range b.Models {
		fields := model.Fields
		for _, field := range fields {
			fieldName := field.Name
			// ポインタを外す
			r := regexp.MustCompile("^[*](.+)$")
			fieldType := r.ReplaceAllString(field.Type.String(), "$1")
			// fieldがユーザー定義型だったら
			if !slices.Contains(defaultTypeList, fieldType) {
				for _, fN := range fields {
					foreignKeyName := fmt.Sprintf("%s%sId", strings.ToUpper(string(fieldName[0])), fieldName[1:])
					// 外部キーのフィールドが存在するとき、foreignkeyのタグを追加する
					if fN.Name == foreignKeyName {
						field.Tag += fmt.Sprintf(" gorm:\"foreignKey:%s\"", foreignKeyName)
						break
					}
				}
			}
		}
	}

	return b
}
