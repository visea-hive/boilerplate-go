package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single field validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// tagMessages maps validator tags to human-readable message templates.
// Use %s as a placeholder for the tag's parameter value (e.g. min=8 → "8").
var tagMessages = map[string]string{
	"required":  "is required",
	"email":     "must be a valid email address",
	"url":       "must be a valid URL",
	"uuid":      "must be a valid UUID",
	"uuid4":     "must be a valid UUIDv4",
	"numeric":   "must be a numeric value",
	"alpha":     "must contain only alphabetic characters",
	"alphanum":  "must contain only alphanumeric characters",
	"min":       "must be at least %s characters",
	"max":       "must be at most %s characters",
	"len":       "must be exactly %s characters",
	"gt":        "must be greater than %s",
	"gte":       "must be greater than or equal to %s",
	"lt":        "must be less than %s",
	"lte":       "must be less than or equal to %s",
	"oneof":     "must be one of: %s",
	"eq":        "must equal %s",
	"ne":        "must not equal %s",
	"startswith": "must start with %s",
	"endswith":   "must end with %s",
	"contains":  "must contain %s",
}

func tagToMessage(fe validator.FieldError) string {
	template, ok := tagMessages[fe.Tag()]
	if !ok {
		return fmt.Sprintf("failed '%s' validation", fe.Tag())
	}
	if strings.Contains(template, "%s") && fe.Param() != "" {
		return fmt.Sprintf(template, fe.Param())
	}
	return template
}

// GenerateErrorValidationResponse parses a validation error into a structured Response.
func GenerateErrorValidationResponse(err error) Response {
	var errs []ValidationError

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, v := range validationErrs {
			errs = append(errs, ValidationError{
				Field:   v.Field(),
				Message: tagToMessage(v),
			})
		}
	} else {
		errs = append(errs, ValidationError{
			Field:   "general",
			Message: err.Error(),
		})
	}

	return ErrorResponse("Validation failed", errs)
}
