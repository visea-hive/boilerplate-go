package helpers

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	FieldName   string `json:"field"`
	RuleMessage string `json:"message"`
}

type ErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

func GenerateErrorValidationResponse(err error) *ErrorResponse {
	var errs []ValidationError

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, v := range validationErrs {
			errs = append(errs, ValidationError{
				FieldName:   v.Field(),
				RuleMessage: v.Tag(),
			})
		}
	} else {
		errs = append(errs, ValidationError{
			FieldName:   "general",
			RuleMessage: err.Error(),
		})
	}

	return &ErrorResponse{Errors: errs}
}
