package client

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// ConnectDB initializes and returns a PostgreSQL connection pool
func ConnectDB(URL string) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	return pool, nil
}
