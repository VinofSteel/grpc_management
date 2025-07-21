package handlers

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/internal/handlers/proto_user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handlers) CreateUser(ctx context.Context, newUser *proto_user.CreateUserRequest) (*proto_user.UserResponse, error) {
	slog.InfoContext(ctx, "Received request to create a new user", "email", newUser.Email, "username", newUser.Username)

	// Inline validation struct
	if err := h.validateRequest(ctx, struct {
		Email    string `validate:"required,email"`
		Username string `validate:"required,min=3,max=50,alphanum"`
		Password string `validate:"required,password"`
	}{
		Email:    newUser.Email,
		Username: newUser.Username,
		Password: newUser.Password,
	}, "CreateUser"); err != nil {
		return nil, err
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

	slog.InfoContext(ctx, "User created successfully", "id", dbUser.ID, "email", dbUser.Email, "username", dbUser.Username)
	return &proto_user.UserResponse{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *Handlers) ListUserByID(ctx context.Context, req *proto_user.ListUserByIDRequest) (*proto_user.UserResponse, error) {
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

	slog.InfoContext(ctx, "User retrieved successfully", "id", dbUser.ID, "email", dbUser.Email, "username", dbUser.Username)
	return &proto_user.UserResponse{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *Handlers) ListUserByEmail(ctx context.Context, req *proto_user.ListUserByEmailRequest) (*proto_user.UserResponse, error) {
	slog.InfoContext(ctx, "Received request to get user", "email", req.Email)

	// Inline validation for email
	if err := h.validateRequest(ctx, struct {
		Email string `validate:"required,email"`
	}{
		Email: req.Email,
	}, "ListUserByEmail"); err != nil {
		return nil, err
	}

	// Get user from database
	dbUser, err := h.Queries.ListUserByEmail(ctx, database.ListUserByEmailParams{
		Email: req.Email,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.WarnContext(ctx, "User not found", "email", req.Email)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		slog.ErrorContext(ctx, "Database error while fetching user", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	slog.InfoContext(ctx, "User retrieved successfully", "id", dbUser.ID, "email", dbUser.Email, "username", dbUser.Username)
	return &proto_user.UserResponse{
		Id:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: dbUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (h *Handlers) ListUserByUsername(ctx context.Context, req *proto_user.ListUserByUsernameRequest) (*proto_user.UserResponse, error) {
	slog.InfoContext(ctx, "Received request to get user", "username", req.Username)

	// Inline validation for username
	if err := h.validateRequest(ctx, struct {
		Username string `validate:"required,min=3,max=50,alphanum"`
	}{
		Username: req.Username,
	}, "ListUserByUsername"); err != nil {
		return nil, err
	}

	// Get user from database
	dbUser, err := h.Queries.ListUserByUsername(ctx, database.ListUserByUsernameParams{
		Username: req.Username,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.WarnContext(ctx, "User not found", "username", req.Username)
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		slog.ErrorContext(ctx, "Database error while fetching user", "error", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	slog.InfoContext(ctx, "User retrieved successfully", "id", dbUser.ID, "email", dbUser.Email, "username", dbUser.Username)
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

	if err := h.validateRequest(ctx, struct {
		Limit  int32 `validate:"min=0,max=100"`
		Offset int32 `validate:"min=0"`
	}{
		Limit:  req.Limit,
		Offset: req.Offset,
	}, "ListUsers"); err != nil {
		return nil, err
	}

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
