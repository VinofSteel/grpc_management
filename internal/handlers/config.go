package handlers

import (
	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
)

type Config struct {
	Queries database.Queries
}

type Handlers struct {
	Queries database.Queries

	proto_user.UnimplementedUserServiceServer
}

func New(config Config) *Handlers {
	return &Handlers{
		Queries: config.Queries,
	}
}
