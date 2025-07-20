package user

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Server struct {
	UnimplementedUserServiceServer
}

func (s *Server) CreateUser(ctx context.Context, user *NewUser) (*User, error) {
	slog.InfoContext(ctx, "Received request to create a new user", "email", user.Email, "username", user.Username)

	return &User{
		Id:        uuid.New().String(),
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: time.Now().Format("2006-01-02"),
		UpdatedAt: time.Now().Format("2006-01-02"),
	}, nil
}
