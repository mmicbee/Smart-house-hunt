package users

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	schema := `
	CREATE TABLE roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);
	INSERT INTO roles (id, name) VALUES
	(1, 'tenant'),
	(2, 'landlord'),
	(3, 'agent'),
	(4, 'admin');

	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		role_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (role_id) REFERENCES roles(id)
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestCreateAndFindUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	user := &User{
		RoleID:   1,
		Name:     "Alice Smith",
		Email:    "alice@example.com",
		Password: "hashedpassword123",
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected non-zero user ID")
	}
	if user.Status != "active" {
		t.Errorf("expected default status 'active', got '%s'", user.Status)
	}

	byEmail, err := repo.FindByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("failed to find by email: %v", err)
	}
	if byEmail.Name != "Alice Smith" || byEmail.Status != "active" {
		t.Errorf("unexpected user from email: %+v", byEmail)
	}

	byID, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("failed to find by ID: %v", err)
	}
	if byID.Email != "alice@example.com" {
		t.Errorf("unexpected user from ID: %+v", byID)
	}
}

func TestListUsersWithRoles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	u1 := &User{RoleID: 1, Name: "Tenant Alice", Email: "alice@test.com", Password: "pass"}
	u2 := &User{RoleID: 2, Name: "Landlord Bob", Email: "bob@test.com", Password: "pass"}
	if err := repo.Create(u1); err != nil {
		t.Fatalf("failed to create u1: %v", err)
	}
	if err := repo.Create(u2); err != nil {
		t.Fatalf("failed to create u2: %v", err)
	}

	list, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 users, got %d", len(list))
	}
	if list[0].RoleName != "tenant" || list[1].RoleName != "landlord" {
		t.Errorf("unexpected role names: %s, %s", list[0].RoleName, list[1].RoleName)
	}
}

func TestUpdateUserStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	u := &User{RoleID: 1, Name: "User One", Email: "user1@test.com", Password: "pass"}
	if err := repo.Create(u); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if err := repo.UpdateStatus(u.ID, "suspended"); err != nil {
		t.Fatalf("failed to update status to suspended: %v", err)
	}

	updated, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("failed to find user: %v", err)
	}
	if updated.Status != "suspended" {
		t.Fatalf("expected status 'suspended', got '%s'", updated.Status)
	}

	if err := repo.UpdateStatus(u.ID, "active"); err != nil {
		t.Fatalf("failed to update status to active: %v", err)
	}

	reactivated, _ := repo.FindByID(u.ID)
	if reactivated.Status != "active" {
		t.Fatalf("expected status 'active', got '%s'", reactivated.Status)
	}

	// Non-existent user
	if err := repo.UpdateStatus(9999, "suspended"); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows for non-existent user, got %v", err)
	}
}

func TestRolesLookup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	id, err := repo.RoleIDByName("landlord")
	if err != nil || id != 2 {
		t.Errorf("expected landlord role_id 2, got %d, err: %v", id, err)
	}

	name, err := repo.RoleByID(4)
	if err != nil || name != "admin" {
		t.Errorf("expected admin role name, got %s, err: %v", name, err)
	}
}
