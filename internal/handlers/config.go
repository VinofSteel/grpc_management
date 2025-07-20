package handlers

import (
	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
	"github.com/vinofsteel/grpc-management/internal/validation"
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
