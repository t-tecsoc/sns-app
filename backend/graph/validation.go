package graph

import (
	"backend/graph/customerr"
	"backend/graph/validation"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ValidateInputModel(m any) *gqlerror.Error {
	validationErrors, err := validation.ValidateModel(m)
	if err != nil {
		return &gqlerror.Error{
			Message:    customerr.ErrorMessage(customerr.InternalServerError),
			Extensions: customerr.InternalServerErrorExtension(),
		}
	}
	if len(validationErrors) > 0 {
		return &gqlerror.Error{
			Message:    customerr.ErrorMessage(customerr.BadUserInput),
			Extensions: customerr.BadUserInputExtension(validationErrors),
		}
	}
	return nil
}
