package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userValidation struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=3,max=50,alphanum"`
	Password string `validate:"required,password"`
}

func toUserValidation(email, username, password string) *userValidation {
	return &userValidation{
		Email:    email,
		Username: username,
		Password: password,
	}
}

func (h *Handlers) CreateUser(ctx context.Context, newUser *proto_user.NewUser) (*proto_user.User, error) {
	slog.InfoContext(ctx, "Received request to create a new user", "email", newUser.Email, "username", newUser.Username)

	userValidation := toUserValidation(newUser.Email, newUser.Username, newUser.Password)
	if err := h.Validator.ValidateData(userValidation); err != nil {
		slog.WarnContext(ctx, "User validation failed", "errors", err.Errors)
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %s", strings.Join(err.Errors, "; "))
	}

	// Check if user already exists
	userWithExistingEmail, err := h.Queries.ListUserByEmail(ctx, database.ListUserByEmailParams{
		Email: newUser.Email,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.ErrorContext(ctx, "Database error while checking existing user", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if userWithExistingEmail != nil {
		slog.WarnContext(ctx, "User already exists", "email", newUser.Email)
		return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", newUser.Email)
	}

	userWithExistingUsername, err := h.Queries.ListUserByUsername(ctx, database.ListUserByUsernameParams{
		Username: newUser.Username,
	})
	if err != nil && err != sql.ErrNoRows {
		slog.ErrorContext(ctx, "Database error while checking existing user", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if userWithExistingUsername != nil {
		slog.WarnContext(ctx, "User already exists", "username", newUser.Username)
		return nil, status.Errorf(codes.AlreadyExists, "user with username %s already exists", newUser.Username)
	}

	// Encrypting user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		slog.ErrorContext(ctx, "Error encrypting user's password", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	newUser.Password = string(hashedPassword)

	// Create user in database
	dbUser, err := h.Queries.InsertUser(ctx, database.InsertUserParams{
		Email:    newUser.Email,
		Username: newUser.Username,
		Password: newUser.Password,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create user in database", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	slog.InfoContext(ctx, "User created successfully", "id", dbUser.ID, "email", dbUser.Email)

	return &proto_user.User{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02"),
	}, nil
}
