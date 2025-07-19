package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vinofsteel/grpc-management/internal/database"
	"github.com/vinofsteel/grpc-management/pkg"
)

type Queries struct {
	db *sqlx.DB
	
	database.UsersRepository
}

func New(ctx context.Context, provider pkg.DBProvider) (*Queries, error) {
	db, err := provider.GetConnection(ctx)
	if err != nil {
		return nil, err
	}

	return &Queries{
		db: db,
	}, nil
}
