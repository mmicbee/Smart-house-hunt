package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`); err != nil {
		return err
	}

	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return err
	}
	if len(files) == 0 {
		files, err = filepath.Glob("migrations/*.sql")
		if err != nil {
			return err
		}
	}
	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)

		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&exists)
		if err != nil {
			return err
		}
		if exists > 0 {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				// If column already exists in SQLite, ignore duplicate column error
				if !strings.Contains(err.Error(), "duplicate column name") {
					return err
				}
			}
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			return err
		}
	}

	return nil
}