package database

import (
	"context"

	"github.com/google/uuid"
)

// Parameters
type ListUserByEmailParams struct {
	Email       string `json:"email" db:"email"`
	ListDeleted bool   `json:"list_deleted" db:"list_deleted"`
}

type ListUsersByIDsParams struct {
	IDs         uuid.UUIDs `json:"ids" db:"ids"`
	ListDeleted bool       `json:"list_deleted" db:"list_deleted"`
}

type ListUserByUsernameParams struct {
	Username    string `json:"username" db:"username"`
	ListDeleted bool   `json:"list_deleted" db:"list_deleted"`
}

type ListUserByIdParams struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ListDeleted bool      `json:"list_deleted" db:"list_deleted"`
}

type InsertUserParams struct {
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type UpdateUserPasswordParams struct {
	UserID   string `json:"user_id" db:"user_id"`
	Password string `json:"password" db:"password"`
}

type DeleteUserParams struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Hard bool      `json:"hard" db:"hard"`
}

// Interface
type UsersRepository interface {
	ListUserByEmail(ctx context.Context, params ListUserByEmailParams) (User, error)
	ListUserByUsername(ctx context.Context, params ListUserByUsernameParams) (User, error)
	ListUserById(ctx context.Context, params ListUserByIdParams) (User, error)
	ListUsersByIDs(ctx context.Context, params ListUsersByIDsParams) ([]User, error)
	InsertUser(ctx context.Context, params InsertUserParams) (User, error)
	UpdateUserPassword(ctx context.Context, params UpdateUserPasswordParams) (User, error)
	DeleteUser(ctx context.Context, params DeleteUserParams) error
}
