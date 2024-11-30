package utils

import (
	"events-api/pkg/validator"
	"fmt"
	"strings"

	validatorv10 "github.com/go-playground/validator/v10"
)

// ValidateRequest validates a struct using validator tags
func ValidateRequest(req interface{}) error {
	// Validate struct
	if err := validator.Validate.Struct(req); err != nil {
		// Check if it's a validation error
		if validationErrors, ok := err.(validatorv10.ValidationErrors); ok {
			var errMsgs []string
			for _, validationErr := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed validation on the '%s' tag", validationErr.Field(), validationErr.Tag()))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, "; "))
		}
		// Handle other types of errors
		return err
	}
	return nil
}
