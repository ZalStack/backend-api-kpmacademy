package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate          = validator.New()
	reHasUpper        = regexp.MustCompile(`[A-Z]`)
	reHasLower        = regexp.MustCompile(`[a-z]`)
	reHasNumber       = regexp.MustCompile(`\d`)
	reHasSpecial      = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func init() {
	_ = validate.RegisterValidation("strong_password", validateStrongPassword) //nolint:errcheck // init-time registration always succeeds
}

func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	return reHasUpper.MatchString(password) &&
		reHasLower.MatchString(password) &&
		reHasNumber.MatchString(password) &&
		reHasSpecial.MatchString(password)
}

func ValidateStruct(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	result := make([]ValidationError, 0, len(validationErrors))
	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())
		message := fmt.Sprintf("Field '%s' failed on the '%s' tag", field, e.Tag())
		result = append(result, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return result
}
