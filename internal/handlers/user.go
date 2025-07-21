package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/google/uuid"
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

func (h *Handlers) CreateUser(ctx context.Context, newUser *proto_user.CreateUserRequest) (*proto_user.UserResponse, error) {
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
	return &proto_user.UserResponse{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *Handlers) GetUser(ctx context.Context, req *proto_user.GetUserRequest) (*proto_user.UserResponse, error) {
	slog.InfoContext(ctx, "Received request to get user", "id", req.Id)

	// Validate UUID format
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		slog.WarnContext(ctx, "Invalid user ID format", "id", req.Id, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	// Get user from database
	dbUser, err := h.Queries.ListUserById(ctx, database.ListUserByIdParams{
		ID:          userID,
		ListDeleted: false, // Don't include deleted users
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.WarnContext(ctx, "User not found", "id", req.Id)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		slog.ErrorContext(ctx, "Database error while fetching user", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	slog.InfoContext(ctx, "User retrieved successfully", "id", dbUser.ID, "email", dbUser.Email)
	return &proto_user.UserResponse{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *Handlers) ListUsers(ctx context.Context, req *proto_user.ListUsersRequest) (*proto_user.ListUsersResponse, error) {
	slog.InfoContext(ctx, "Received request to list users", "limit", req.Limit, "offset", req.Offset)

	// Set default pagination values if not provided
	limit := req.Limit
	offset := req.Offset

	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	// Get users from database
	dbUsers, err := h.Queries.ListUsers(ctx, database.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Database error while fetching users", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	// Convert database users to protobuf users
	users := make([]*proto_user.UserResponse, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = &proto_user.UserResponse{
			Id:        dbUser.ID.String(),
			Email:     dbUser.Email,
			Username:  dbUser.Username,
			CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	slog.InfoContext(ctx, "Users retrieved successfully", "count", len(users), "limit", limit, "offset", offset)
	return &proto_user.ListUsersResponse{
		Users: users,
	}, nil
}
