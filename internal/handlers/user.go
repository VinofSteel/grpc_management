package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
)

func (s *Handlers) CreateUser(ctx context.Context, newUser *proto_user.NewUser) (*proto_user.User, error) {
	slog.InfoContext(ctx, "Received request to create a new user", "email", newUser.Email, "username", newUser.Username)

	return &proto_user.User{
		Id:        uuid.New().String(),
		Email:     newUser.Email,
		Username:  newUser.Username,
		CreatedAt: time.Now().Format("2006-01-02"),
		UpdatedAt: time.Now().Format("2006-01-02"),
	}, nil
}
