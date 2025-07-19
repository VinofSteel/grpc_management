package pkg

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// DBProvider is the abstraction for all database use in the application
type DBProvider interface {
	GetConnection(ctx context.Context) (*sqlx.DB, error)
	Close() error
}
