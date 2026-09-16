package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		message := fmt.Sprintf("Field '%s' failed on the '%s' tag", field, err.Tag())
		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return validationErrors
}
