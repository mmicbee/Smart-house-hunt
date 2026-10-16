package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRunMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer db.Close()

	if err := RunMigrations(db); err != nil {
		t.Fatalf("RunMigrations failed: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query schema_migrations: %v", err)
	}
	t.Logf("Total applied migrations: %d", count)

	var userCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = 'admin@smarthousehunt.com'").Scan(&userCount)
	if err != nil {
		t.Fatalf("failed to query admin user: %v", err)
	}
	if userCount != 1 {
		t.Errorf("expected 1 admin user, got %d", userCount)
	}
}
