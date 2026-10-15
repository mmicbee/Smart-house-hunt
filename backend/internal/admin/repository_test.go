package admin

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

	CREATE TABLE properties (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		landlord_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		price INTEGER NOT NULL,
		bedrooms INTEGER NOT NULL DEFAULT 1,
		bathrooms INTEGER NOT NULL DEFAULT 1,
		amenities TEXT NOT NULL DEFAULT '',
		property_type TEXT NOT NULL DEFAULT 'apartment',
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (landlord_id) REFERENCES users(id)
	);

	CREATE TABLE property_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		property_id INTEGER NOT NULL,
		file_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
	);

	CREATE TABLE favorites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		property_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id, property_id)
	);

	CREATE TABLE viewing_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		property_id INTEGER NOT NULL,
		preferred_date TEXT NOT NULL,
		message TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create test schema: %v", err)
	}

	return db
}

func TestListUsers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	_, err := db.Exec(`
		INSERT INTO users (id, role_id, name, email, password, status) VALUES
		(1, 1, 'John Tenant', 'john@test.com', 'pass', 'active'),
		(2, 2, 'Jane Landlord', 'jane@test.com', 'pass', 'active'),
		(3, 4, 'Super Admin', 'admin@test.com', 'pass', 'active'),
		(4, 1, 'Suspended User', 'bad@test.com', 'pass', 'suspended');
	`)
	if err != nil {
		t.Fatalf("failed to insert test users: %v", err)
	}

	// 1. All users
	all, err := repo.ListUsers("", "", "")
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(all) != 4 {
		t.Fatalf("expected 4 users, got %d", len(all))
	}

	// 2. Filter by status
	suspended, err := repo.ListUsers("suspended", "", "")
	if err != nil {
		t.Fatalf("ListUsers suspended failed: %v", err)
	}
	if len(suspended) != 1 || suspended[0].Email != "bad@test.com" {
		t.Fatalf("expected 1 suspended user, got %d", len(suspended))
	}

	// 3. Filter by role
	landlords, err := repo.ListUsers("", "landlord", "")
	if err != nil {
		t.Fatalf("ListUsers landlord failed: %v", err)
	}
	if len(landlords) != 1 || landlords[0].Role != "landlord" {
		t.Fatalf("expected 1 landlord, got %d", len(landlords))
	}

	// 4. Search query
	found, err := repo.ListUsers("", "", "Admin")
	if err != nil {
		t.Fatalf("ListUsers search failed: %v", err)
	}
	if len(found) != 1 || found[0].Name != "Super Admin" {
		t.Fatalf("expected 1 search result, got %d", len(found))
	}
}

func TestUpdateUserStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	_, err := db.Exec(`
		INSERT INTO users (id, role_id, name, email, password, status) VALUES
		(10, 1, 'Mark User', 'mark@test.com', 'pass', 'active');
	`)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	if err := repo.UpdateUserStatus(10, "suspended"); err != nil {
		t.Fatalf("UpdateUserStatus to suspended failed: %v", err)
	}

	user, err := repo.GetUserByID(10)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if user.Status != "suspended" {
		t.Fatalf("expected status 'suspended', got '%s'", user.Status)
	}

	// Non-existent user
	if err := repo.UpdateUserStatus(999, "active"); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestListAndUpdateProperties(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	_, err := db.Exec(`
		INSERT INTO users (id, role_id, name, email, password, status) VALUES
		(1, 2, 'Landlord Bob', 'bob@test.com', 'pass', 'active');

		INSERT INTO properties (id, landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, property_type, status) VALUES
		(1, 1, 'Pending Studio', 'Nice studio', 'Kilimani', 20000, 0, 1, 'WiFi', 'studio', 'pending'),
		(2, 1, 'Approved Flat', 'Modern flat', 'Westlands', 45000, 2, 2, 'Parking, Gym', 'apartment', 'approved'),
		(3, 1, 'Rejected Bedsitter', 'Small room', 'Roysambu', 8000, 0, 1, '', 'bedsitter', 'rejected');

		INSERT INTO property_images (property_id, file_path) VALUES
		(1, 'uploads/property-images/1/img1.jpg');
	`)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// 1. List all
	all, err := repo.ListProperties("", "")
	if err != nil {
		t.Fatalf("ListProperties failed: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 properties, got %d", len(all))
	}

	// Check landlord info & images loaded
	p1, err := repo.GetPropertyByID(1)
	if err != nil {
		t.Fatalf("GetPropertyByID failed: %v", err)
	}
	if p1.LandlordName != "Landlord Bob" || p1.LandlordEmail != "bob@test.com" {
		t.Errorf("unexpected landlord info: %s, %s", p1.LandlordName, p1.LandlordEmail)
	}
	if len(p1.Images) != 1 {
		t.Errorf("expected 1 image for property 1, got %d", len(p1.Images))
	}

	// 2. Filter by status
	pending, err := repo.ListProperties("pending", "")
	if err != nil {
		t.Fatalf("ListProperties pending failed: %v", err)
	}
	if len(pending) != 1 || pending[0].Title != "Pending Studio" {
		t.Fatalf("expected 1 pending property, got %d", len(pending))
	}

	// 3. Update status (Approve)
	updated, err := repo.UpdatePropertyStatus(1, "approved")
	if err != nil {
		t.Fatalf("UpdatePropertyStatus failed: %v", err)
	}
	if updated.Status != "approved" {
		t.Fatalf("expected status 'approved', got '%s'", updated.Status)
	}

	// 4. Delete property
	if err := repo.DeleteProperty(3); err != nil {
		t.Fatalf("DeleteProperty failed: %v", err)
	}
	_, err = repo.GetPropertyByID(3)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows after delete, got %v", err)
	}
}

func TestGetReports(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	_, err := db.Exec(`
		INSERT INTO users (id, role_id, name, email, password, status) VALUES
		(1, 1, 'Tenant Tom', 'tom@test.com', 'pass', 'active'),
		(2, 2, 'Landlord Lisa', 'lisa@test.com', 'pass', 'active'),
		(3, 4, 'Admin Alex', 'alex@test.com', 'pass', 'active'),
		(4, 1, 'Bad Tenant', 'bad@test.com', 'pass', 'suspended');

		INSERT INTO properties (id, landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, property_type, status) VALUES
		(1, 2, 'Pending House', 'Quiet house', 'Karen', 80000, 3, 2, 'Garden', 'house', 'pending'),
		(2, 2, 'Approved Apartment', 'City view', 'Kilimani', 35000, 1, 1, 'WiFi', 'apartment', 'approved');

		INSERT INTO favorites (user_id, property_id) VALUES (1, 2);

		INSERT INTO viewing_requests (user_id, property_id, preferred_date, message, status) VALUES
		(1, 2, '2026-10-20 14:00', 'Excited to view', 'pending');
	`)
	if err != nil {
		t.Fatalf("failed to seed test data: %v", err)
	}

	reports, err := repo.GetReports()
	if err != nil {
		t.Fatalf("GetReports failed: %v", err)
	}

	// Check Users
	if reports.Users.Total != 4 {
		t.Errorf("expected 4 total users, got %d", reports.Users.Total)
	}
	if reports.Users.Active != 3 || reports.Users.Suspended != 1 {
		t.Errorf("expected 3 active, 1 suspended; got %d active, %d suspended", reports.Users.Active, reports.Users.Suspended)
	}
	if reports.Users.Tenants != 2 || reports.Users.Landlords != 1 || reports.Users.Admins != 1 {
		t.Errorf("unexpected role breakdown: %+v", reports.Users)
	}

	// Check Properties
	if reports.Properties.Total != 2 {
		t.Errorf("expected 2 total properties, got %d", reports.Properties.Total)
	}
	if reports.Properties.Pending != 1 || reports.Properties.Approved != 1 {
		t.Errorf("expected 1 pending, 1 approved; got pending=%d, approved=%d", reports.Properties.Pending, reports.Properties.Approved)
	}
	if reports.Properties.ByType["house"] != 1 || reports.Properties.ByType["apartment"] != 1 {
		t.Errorf("unexpected by_type breakdown: %+v", reports.Properties.ByType)
	}

	// Check Viewings & Favorites
	if reports.Viewings.Total != 1 || reports.Viewings.Pending != 1 {
		t.Errorf("expected 1 pending viewing, got %+v", reports.Viewings)
	}
	if reports.Favorites != 1 {
		t.Errorf("expected 1 favorite, got %d", reports.Favorites)
	}

	// Check Recent Activity
	if len(reports.Recent) == 0 {
		t.Error("expected non-empty recent activity list")
	}
}
