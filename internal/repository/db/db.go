package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var db *pgxpool.Pool

func GetConnection() (*pgxpool.Pool, error) {
	if db == nil {
		connStr := "host=db port=5432 user=admin password=123 dbname=rinha sslmode=disable"
		conn, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			return nil, err
		}

		db = conn
	}

	return db, nil
}

func CloseConnection() {
	db.Close()
}
