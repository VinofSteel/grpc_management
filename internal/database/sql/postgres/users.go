package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vinofsteel/grpc-management/internal/database"
)

func (q *PSQLQueries) ListUserByEmail(ctx context.Context, params database.ListUserByEmailParams) (*database.User, error) {
	slog.InfoContext(ctx, "Listing user by email", "email", params.Email, "layer", "repository", "driver", "psql")

	query := `SELECT
		id, created_at, updated_at, email, username, password 
			FROM users 
			WHERE email = :email`

	if !params.ListDeleted {
		query += ` AND deleted_at IS NULL`
	}

	var user database.User
	rows, err := q.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying user by email", "error", err, "email", params.Email)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			slog.ErrorContext(ctx, "Error scanning user by email", "error", err, "email", params.Email)
			return nil, err
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (q *PSQLQueries) ListUsersByIDs(ctx context.Context, params database.ListUsersByIDsParams) ([]*database.User, error) {
	slog.InfoContext(ctx, "Listing users by ids", "ids", params.IDs.Strings(), "layer", "repository", "driver", "psql")

	query := `SELECT
		id, created_at, updated_at, email, username, password 
		FROM users 
		WHERE id = ANY(:ids)`

	if !params.ListDeleted {
		query += ` AND deleted_at IS NULL`
	}

	idStrings := params.IDs.Strings()
	interfaces := make([]any, len(idStrings))
	for i, v := range idStrings {
		interfaces[i] = v
	}

	// Create a struct with the array parameter
	queryParams := struct {
		IDs interface {
			driver.Valuer
			sql.Scanner
		} `db:"ids"`
		ListDeleted bool `db:"list_deleted"`
	}{
		IDs:         pq.Array(interfaces),
		ListDeleted: params.ListDeleted,
	}

	var users []*database.User
	rows, err := q.db.NamedQueryContext(ctx, query, queryParams)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying users by IDs", "error", err, "ids", params.IDs.Strings())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user database.User
		if err := rows.StructScan(&user); err != nil {
			slog.ErrorContext(ctx, "Error scanning user from rows", "error", err)
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, "Error iterating over user rows", "error", err)
		return nil, err
	}

	return users, nil
}

func (q *PSQLQueries) ListUserByUsername(ctx context.Context, params database.ListUserByUsernameParams) (*database.User, error) {
	slog.InfoContext(ctx, "Listing user by username", "username", params.Username, "layer", "repository", "driver", "psql")

	query := `SELECT
		id, created_at, updated_at, email, username, password 
			FROM users 
			WHERE username = :username`

	if !params.ListDeleted {
		query += ` AND deleted_at IS NULL`
	}

	var user database.User
	rows, err := q.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying user by username", "error", err, "username", params.Username)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			slog.ErrorContext(ctx, "Error scanning user by username", "error", err, "username", params.Username)
			return nil, err
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (q *PSQLQueries) ListUserById(ctx context.Context, params database.ListUserByIdParams) (*database.User, error) {
	slog.InfoContext(ctx, "Listing user by id", "id", params.ID, "layer", "repository", "driver", "psql")

	query := `SELECT
		id, created_at, updated_at, email, username, password 
			FROM users 
			WHERE id = :id`

	if !params.ListDeleted {
		query += ` AND deleted_at IS NULL`
	}

	var user database.User
	rows, err := q.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		slog.ErrorContext(ctx, "Error querying user by id", "error", err, "id", params.ID)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			slog.ErrorContext(ctx, "Error scanning user by id", "error", err, "id", params.ID)
			return nil, err
		}
		return &user, nil
	}

	return nil, sql.ErrNoRows
}

func (q *PSQLQueries) InsertUser(ctx context.Context, params database.InsertUserParams) (*database.User, error) {
	slog.InfoContext(ctx, "Creating user", "email", params.Email, "username", params.Username, "layer", "repository", "driver", "psql")

	query := `INSERT INTO users 
		(email, username, password) VALUES (:email, :username, :password) 
			RETURNING id, created_at, updated_at, email, username, password`

	var user database.User
	rows, err := q.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		slog.ErrorContext(ctx, "Error inserting user", "error", err, "email", params.Email, "username", params.Username)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			slog.ErrorContext(ctx, "Error scanning inserted user", "error", err, "email", params.Email, "username", params.Username)
			return nil, err
		}
		return &user, nil
	}

	return nil, err
}

func (q *PSQLQueries) UpdateUserPassword(ctx context.Context, params database.UpdateUserPasswordParams) (*database.User, error) {
	slog.InfoContext(ctx, "Updating user password", "user_id", params.UserID, "layer", "repository", "driver", "psql")

	query := `UPDATE users SET password = :password, updated_at = CURRENT_TIMESTAMP WHERE id = :user_id
		RETURNING id, created_at, updated_at, email, username, password`

	var user database.User
	rows, err := q.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		slog.ErrorContext(ctx, "Error updating user password", "error", err, "user_id", params.UserID)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			slog.ErrorContext(ctx, "Error scanning updated user", "error", err, "user_id", params.UserID)
			return nil, err
		}
		return &user, nil
	}

	return nil, err
}

func (q *PSQLQueries) DeleteUser(ctx context.Context, params database.DeleteUserParams) error {
	slog.InfoContext(ctx, "Deleting user", "id", params.ID, "hard", params.Hard, "layer", "repository", "driver", "psql")

	tx, err := q.db.BeginTxx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "Error beginning transaction on DeleteUser", "error", err, "id", params.ID)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.ErrorContext(ctx, "Could not rollback in DeleteUser after panic", "error", rollbackErr, "id", params.ID)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.ErrorContext(ctx, "Could not rollback in DeleteUser", "error", rollbackErr, "id", params.ID)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				slog.ErrorContext(ctx, "Could not commit in DeleteUser", "error", commitErr, "id", params.ID)
				err = commitErr
			}
		}
	}()

	if !params.Hard {
		softDeleteParams := struct {
			ID        uuid.UUID `db:"id"`
			DeletedAt time.Time `db:"deleted_at"`
		}{
			ID:        params.ID,
			DeletedAt: time.Now().UTC(),
		}

		softQuery := `UPDATE users 
            SET deleted_at = :deleted_at WHERE id = :id`

		_, err = tx.NamedExecContext(ctx, softQuery, softDeleteParams)
		if err != nil {
			slog.ErrorContext(ctx, "Error executing soft delete", "error", err, "id", params.ID)
			return err
		}
	}

	if params.Hard {
		hardDeleteParams := struct {
			UserID uuid.UUID `db:"user_id"`
		}{
			UserID: params.ID,
		}

		userDeleteParams := struct {
			ID uuid.UUID `db:"id"`
		}{
			ID: params.ID,
		}

		hardQuerySessions := `DELETE FROM sessions WHERE user_id = :user_id`
		hardQueryUsers := `DELETE FROM users WHERE id = :id`

		_, err = tx.NamedExecContext(ctx, hardQuerySessions, hardDeleteParams)
		if err != nil {
			slog.ErrorContext(ctx, "Error executing hard delete on sessions", "error", err, "id", params.ID)
			return err
		}

		_, err = tx.NamedExecContext(ctx, hardQueryUsers, userDeleteParams)
		if err != nil {
			slog.ErrorContext(ctx, "Error executing hard delete on users", "error", err, "id", params.ID)
			return err
		}
	}

	return nil
}
