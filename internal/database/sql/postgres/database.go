package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/pkg"
)

type PSQLQueries struct {
	db *sqlx.DB

	database.UsersRepository
}

func NewPSQLQueries(ctx context.Context, provider pkg.DBProvider) (*PSQLQueries, error) {
	db, err := provider.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	return &PSQLQueries{
		db: db,
	}, nil
}
