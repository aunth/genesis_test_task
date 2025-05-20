package database

import (
	"database/sql"
	"fmt"
	"log"

	"weather-service/internal/database/migrations"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
