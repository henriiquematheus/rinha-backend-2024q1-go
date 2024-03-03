// connection.go

package connection

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Init(connectionString string) error {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return err
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	return nil
}

func GetPool() *pgxpool.Pool {
	return pool
}
