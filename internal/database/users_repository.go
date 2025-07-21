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

type ListUserByUsernameParams struct {
	Username    string `json:"username" db:"username"`
	ListDeleted bool   `json:"list_deleted" db:"list_deleted"`
}

type ListUserByIdParams struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ListDeleted bool      `json:"list_deleted" db:"list_deleted"`
}

type ListUsersParams struct {
	ListDeleted bool  `json:"list_deleted" db:"list_deleted"`
	Limit       int32 `json:"limit" db:"limit"`
	Offset      int32 `json:"offset" db:"offset"`
}

type InsertUserParams struct {
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type UpdateUserPasswordParams struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Password string    `json:"password" db:"password"`
}

type DeleteUserParams struct {
	ID   uuid.UUID `json:"id" db:"id"`
	Hard bool      `json:"hard" db:"hard"`
}

// Interface
type UsersRepository interface {
	ListUserByEmail(ctx context.Context, params ListUserByEmailParams) (*User, error)
	ListUserByUsername(ctx context.Context, params ListUserByUsernameParams) (*User, error)
	ListUserById(ctx context.Context, params ListUserByIdParams) (*User, error)
	ListUsers(ctx context.Context, params ListUsersParams) ([]*User, error)
	InsertUser(ctx context.Context, params InsertUserParams) (*User, error)
	UpdateUserPassword(ctx context.Context, params UpdateUserPasswordParams) (*User, error)
	DeleteUser(ctx context.Context, params DeleteUserParams) error
}
