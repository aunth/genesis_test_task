package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func getConnectionString() (string, error) {
	requiredVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return "", fmt.Errorf("required environment variable %s is not set", v)
		}
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname), nil
}

func Connect() (*sql.DB, error) {
	connStr, err := getConnectionString()
	if err != nil {
		return nil, fmt.Errorf("error getting connection string: %v", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error opening database:", err)
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Println("Error connecting to the database:", err)
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
