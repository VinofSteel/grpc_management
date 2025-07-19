package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vinofsteel/grpc-management/pkg"
)

// Database provider for postgres connections
type PostgresProvider struct {
	connStr string
	db      *sqlx.DB
}

func (p *PostgresProvider) GetConnection(ctx context.Context) (*sqlx.DB, error) {
	if p.db != nil {
		if err := p.db.PingContext(ctx); err != nil {
			err := p.db.Close()
			return nil, err
		}

		return p.db, nil
	}

	db, err := sqlx.ConnectContext(ctx, "postgres", p.connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %w", err)
	}

	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	p.db = db
	return db, nil
}

func (p *PostgresProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}

	return nil
}

// Create a New PSQL db provider where needed
func newPostgresProvider(user, password, host, port, dbName string) *PostgresProvider {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	return &PostgresProvider{
		connStr: connStr,
	}
}

func NewPostgresDatabaseProvider() pkg.DBProvider {
	return newPostgresProvider(
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
}
