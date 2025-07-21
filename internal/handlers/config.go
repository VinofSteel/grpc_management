package handlers

import (
	"context"
	"log/slog"
	"strings"

	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
	"github.com/vinofsteel/grpc-management/internal/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Queries   database.Queries
	Validator validation.ValidationProvider
}

type Handlers struct {
	Queries   database.Queries
	Validator validation.ValidationProvider

	proto_user.UnimplementedUserServiceServer
}

func New(config Config) *Handlers {
	return &Handlers{
		Queries:   config.Queries,
		Validator: config.Validator,
	}
}

// validateRequest is a generic helper function that validates any struct and returns gRPC error
func (h *Handlers) validateRequest(ctx context.Context, data any, operation string) error {
	if err := h.Validator.ValidateData(data); err != nil {
		slog.WarnContext(ctx, "Validation failed", "operation", operation, "errors", err.Errors)
		return status.Errorf(codes.InvalidArgument, "validation failed: %s", strings.Join(err.Errors, "; "))
	}
	return nil
}
