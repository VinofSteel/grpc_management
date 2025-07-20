package validation

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Errors []string
}

func (ve *ValidationError) Error() string {
	return "validation error"
}

// DatabaseProvider is the abstraction for all validator use in the application
type ValidationProvider interface {
	ValidateData(data any) *ValidationError
}

// Creates a new Validate based ValidationProvider
func NewValidateValidationrovider(ctx context.Context) ValidationProvider {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return New(ctx, validate)
}
