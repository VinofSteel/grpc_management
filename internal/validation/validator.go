package validation

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidateProvider interface {
	RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error
	Struct(s any) error
	Var(field any, tag string) error
}

type Validator struct {
	validate ValidateProvider
}

type errorResponse struct {
	Error        bool
	FailedField  string
	Tag          string
	ErrorMessage string
}

func New(ctx context.Context, validate ValidateProvider) ValidationProvider {
	if err := validate.RegisterValidation("password", passwordTagValidation); err != nil {
		slog.ErrorContext(ctx, "Error registering passwordTagValidation", "error", err)
		os.Exit(1)
	}

	return &Validator{
		validate,
	}
}

func (v *Validator) ValidateData(data any) *ValidationError {
	if errors := structValidation(v.validate, data); len(errors) > 0 && errors[0].Error {
		var errorMessages []string

		for _, err := range errors {
			errorMessages = append(errorMessages, err.ErrorMessage)
		}

		return &ValidationError{
			Errors: errorMessages,
		}
	}

	return nil
}

// Utilities
func structValidation(validate ValidateProvider, data any) []errorResponse {
	var validationErrors []errorResponse

	errors := validate.Struct(data)
	if errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			var errResp errorResponse

			errResp.Tag = err.Tag()
			errResp.Error = true

			errResp.FailedField = strings.ToLower(err.Field())

			switch err.Tag() {
			case "required":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' is required", errResp.FailedField)
			case "email":
				errResp.ErrorMessage = "field 'email' must be a valid email address"
			case "min":
				if err.Kind() == reflect.Int || err.Kind() == reflect.Float64 {
					errResp.ErrorMessage = fmt.Sprintf("field '%s' must be at least %s", errResp.FailedField, err.Param())
				} else {
					errResp.ErrorMessage = fmt.Sprintf("field '%s' must be at least %s characters long", errResp.FailedField, err.Param())
				}
			case "max":
				if err.Kind() == reflect.Int || err.Kind() == reflect.Float64 {
					errResp.ErrorMessage = fmt.Sprintf("field '%s' must be at most %s", errResp.FailedField, err.Param())
				} else {
					errResp.ErrorMessage = fmt.Sprintf("field '%s' must be at most %s characters long", errResp.FailedField, err.Param())
				}
			case "password":
				errResp.ErrorMessage = "field 'password' must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character"
			case "alphanum":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must contain only alphanumeric characters", errResp.FailedField)
			case "alpha":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must contain only alphabetic characters", errResp.FailedField)
			case "numeric":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must be a valid number", errResp.FailedField)
			case "len":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must be exactly %s characters long", errResp.FailedField, err.Param())
			case "oneof":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must be one of [%s]", errResp.FailedField, err.Param())
			case "datetime":
				errResp.ErrorMessage = fmt.Sprintf("field '%s' must be a valid datetime in YYYY-MM-DD format", errResp.FailedField)
			default:
				errResp.ErrorMessage = fmt.Sprintf("field '%s' failed validation for tag '%s'", errResp.FailedField, err.Tag())
			}

			validationErrors = append(validationErrors, errResp)
		}
	}

	return validationErrors
}
