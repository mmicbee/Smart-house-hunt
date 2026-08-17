package database

import (
	"datadase/sql"
	"os"
)

func RunMigrations(db *sql.DB) error {
	sqlFile, err := os.ReadFile("database/migrations/001_initial.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))
	return err
}