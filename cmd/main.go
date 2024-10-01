package main

import (
	"database/sql"
	"fmt"
	"os"
	"rest-skeleton/pkg/config"
	"rest-skeleton/pkg/database"
	"rest-skeleton/pkg/migration"

	_ "github.com/lib/pq"
)

func main() {
	if _, ok := os.LookupEnv("APP_NAME"); !ok {
		if err := config.Setup(".env"); err != nil {
			fmt.Println("failed to setup config", err)
			os.Exit(1)
		}
	}

	db, err := database.NewDatabase()
	if err != nil {
		fmt.Println("Could not connect to database", err)
		os.Exit(1)
	}
	defer db.Conn.Close()

	if len(os.Args) < 2 {
		fmt.Println("No command requested. try with: go run cmd/main.go migrate")
		return
	}

	switch os.Args[1] {
	case "migrate":
		migrate(db.Conn)
	default:
		fmt.Println("Unknown command. Available commands: migrate")
	}
}

func migrate(db *sql.DB) {
	fmt.Println("Starting migration...")
	if err := migration.Migrate(db); err != nil {
		fmt.Println("Could not migrate database: ", err)
	} else {
		fmt.Println("Migrated database successfully")
	}
	fmt.Println("Finish migration...")
}
