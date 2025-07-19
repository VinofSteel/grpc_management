package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vinofsteel/grpc-management/pkg"
)

type SQLQueries struct {
	db *sqlx.DB
}

func New(ctx context.Context, provider pkg.DBProvider) (*SQLQueries, error) {
	db, err := provider.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	return &SQLQueries{
		db: db,
	}, nil
}
