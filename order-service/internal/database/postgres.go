package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {

	connStr := "host=order-db port=5432 user=postgres password=postgres dbname=orderdb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		restaurant_id TEXT,
		total_price BIGINT,
		status TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
