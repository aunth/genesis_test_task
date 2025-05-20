package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migration struct {
	Version int
	SQL     string
}

func RunMigrations(db *sql.DB) error {
	migrationsDir := filepath.Join("internal", "database", "migrations")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %v", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			var version int
			_, err := fmt.Sscanf(file.Name(), "%d_", &version)
			if err != nil {
				return fmt.Errorf("invalid migration filename format: %s", file.Name())
			}

			content, err := os.ReadFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
			}

			migrations = append(migrations, Migration{
				Version: version,
				SQL:     string(content),
			})
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, migration := range migrations {

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("error beginning transaction: %v", err)
		}

		_, err = tx.Exec(migration.SQL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error executing migration %d: %v", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing migration %d: %v", migration.Version, err)
		}

		fmt.Printf("Applied migration %d\n", migration.Version)
	}

	return nil
}
