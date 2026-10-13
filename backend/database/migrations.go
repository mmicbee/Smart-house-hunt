package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB) error {
	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(content)); err != nil {
			return err
		}
	}

	return nil
}