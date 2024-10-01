package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database wraps the SQL database connection
type Database struct {
	Conn *sql.DB
}

func NewDatabase() (*Database, error) {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), port, os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func StatusCheck(ctx context.Context, db *sql.DB) error {
	var tmp bool
	return db.QueryRowContext(ctx, `SELECT true`).Scan(&tmp)
}
