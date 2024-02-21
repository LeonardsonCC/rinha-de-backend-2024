package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func GetConnection() (*sql.DB, error) {
	if db == nil {
		connStr := "host=db port=5432 user=admin password=123 dbname=rinha sslmode=disable"
		conn, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}

		conn.SetMaxOpenConns(300)
		conn.SetMaxIdleConns(300)

		db = conn
	}

	return db, nil
}
