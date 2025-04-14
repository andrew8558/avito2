package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDb(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}
