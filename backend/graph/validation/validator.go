package validation

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	hhmmRegex = regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`)
)

func init() {
	validate = validator.New()

	// register time validation
	if err := validate.RegisterValidation("HH:mm", func(fl validator.FieldLevel) bool {
		return hhmmRegex.MatchString(fl.Field().String())
	}, false); err != nil {
		log.Fatalln("failed to register validation")
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "入力は必須です"
	case "len":
		return fmt.Sprintf("%s文字で入力してください", fe.Param())
	case "gte":
		return fmt.Sprintf("%s文字以上で入力してください", fe.Param())
	case "lte":
		return fmt.Sprintf("%s文字以下で入力してください", fe.Param())
	case "timezone":
		return "IANA Time Zone databaseの形式で入力してください"
	case "HH:mm":
		return "00:00 ~ 23:59の間で入力してください"
	case "oneof":
		return fmt.Sprintf("「%s」の中のいずれかの値を入力してください", strings.Replace(fe.Param(), " ", " or ", -1))
	default:
		return fe.Error()
	}
}

func ValidateModel(model any) (map[string]string, error) {
	if err := validate.Struct(model); err != nil {
		e := err.(*validator.InvalidValidationError)

		errors := err.(validator.ValidationErrors)
		validationErrors := make(map[string]string, len(errors))
		for _, ve := range errors {
			validationErrors[ve.StructNamespace()] = msgForTag(ve)
		}

		return validationErrors, e
	}
	return nil, nil
}
