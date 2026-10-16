package viewing

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"smart-house-hunt/internal/properties"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	schema := `
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role_id INTEGER NOT NULL
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
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE property_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		property_id INTEGER NOT NULL,
		file_path TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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
		t.Fatalf("failed to create schema: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO users (id, name, email, password, role_id)
		VALUES
			(1, 'Landlord Joe', 'joe@example.com', 'hash', 2),
			(2, 'Tenant Alice', 'alice@example.com', 'hash', 1),
			(3, 'Tenant Bob', 'bob@example.com', 'hash', 1);
	`)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	return db
}

func TestViewingRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	propRepo := properties.NewRepository(db)

	prop := &properties.Property{
		LandlordID:  1,
		Title:       "Kilimani Heights",
		Description: "Stunning 2BR",
		Location:    "Kilimani",
		Price:       60000,
		Bedrooms:    2,
		Bathrooms:   2,
		Amenities:   "wifi,gym",
		Status:      "approved",
	}
	if err := propRepo.Create(prop); err != nil {
		t.Fatalf("failed to create property: %v", err)
	}

	img := &properties.PropertyImage{PropertyID: prop.ID, FilePath: "storage/property-images/img1.jpg"}
	if err := propRepo.AddImage(img); err != nil {
		t.Fatalf("failed to add property image: %v", err)
	}

	// 1. Create request
	req := &ViewingRequest{
		UserID:        2,
		PropertyID:    prop.ID,
		PreferredDate: "2026-10-25T14:00",
		Message:       "Would love to see the balcony view",
		Status:        StatusPending,
	}
	if err := repo.Create(req); err != nil {
		t.Fatalf("failed to create viewing request: %v", err)
	}
	if req.ID == 0 {
		t.Fatalf("expected non-zero id for created viewing request")
	}

	// 2. HasPendingRequest
	hasPending, err := repo.HasPendingRequest(2, prop.ID)
	if err != nil {
		t.Fatalf("failed to check pending request: %v", err)
	}
	if !hasPending {
		t.Fatalf("expected pending request to exist for user 2")
	}

	hasPendingBob, err := repo.HasPendingRequest(3, prop.ID)
	if err != nil {
		t.Fatalf("failed to check pending request for bob: %v", err)
	}
	if hasPendingBob {
		t.Fatalf("expected no pending request for user 3")
	}

	// 3. FindByID
	detail, err := repo.FindByID(req.ID)
	if err != nil {
		t.Fatalf("failed to find viewing by id: %v", err)
	}
	if detail.PropertyTitle != "Kilimani Heights" {
		t.Errorf("expected property title 'Kilimani Heights', got '%s'", detail.PropertyTitle)
	}
	if detail.TenantName != "Tenant Alice" {
		t.Errorf("expected tenant name 'Tenant Alice', got '%s'", detail.TenantName)
	}
	if detail.LandlordName != "Landlord Joe" {
		t.Errorf("expected landlord name 'Landlord Joe', got '%s'", detail.LandlordName)
	}
	if detail.PropertyImage != "storage/property-images/img1.jpg" {
		t.Errorf("expected property image 'storage/property-images/img1.jpg', got '%s'", detail.PropertyImage)
	}

	// 4. ListByTenant
	tenantList, err := repo.ListByTenant(2)
	if err != nil {
		t.Fatalf("failed to list by tenant: %v", err)
	}
	if len(tenantList) != 1 {
		t.Fatalf("expected 1 viewing request for tenant 2, got %d", len(tenantList))
	}

	// 5. ListByLandlord
	landlordList, err := repo.ListByLandlord(1)
	if err != nil {
		t.Fatalf("failed to list by landlord: %v", err)
	}
	if len(landlordList) != 1 {
		t.Fatalf("expected 1 viewing request for landlord 1, got %d", len(landlordList))
	}

	// 6. UpdateStatus (Accepted)
	if err := repo.UpdateStatus(req.ID, StatusAccepted); err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	detailAfterUpdate, err := repo.FindByID(req.ID)
	if err != nil {
		t.Fatalf("failed to find updated viewing: %v", err)
	}
	if detailAfterUpdate.Status != StatusAccepted {
		t.Errorf("expected status 'accepted', got '%s'", detailAfterUpdate.Status)
	}

	// Now that it's accepted, HasPendingRequest should return false
	hasPendingAfterAccept, err := repo.HasPendingRequest(2, prop.ID)
	if err != nil {
		t.Fatalf("failed to check pending after accept: %v", err)
	}
	if hasPendingAfterAccept {
		t.Errorf("expected no pending request after acceptance")
	}

	// 7. Update non-existent ID
	err = repo.UpdateStatus(99999, StatusRejected)
	if !errors.Is(err, ErrViewingNotFound) {
		t.Errorf("expected ErrViewingNotFound, got %v", err)
	}
}
