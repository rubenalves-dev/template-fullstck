package platform

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDatabase(connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}
