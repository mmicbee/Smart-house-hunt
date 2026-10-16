package favorites

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
	CREATE TABLE favorites (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		property_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (user_id, property_id)
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func TestFavoritesRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	favRepo := NewRepository(db)
	propRepo := properties.NewRepository(db)

	p1 := &properties.Property{LandlordID: 1, Title: "Sunset Villa", Description: "Nice", Location: "Mombasa", Price: 30000, Status: "approved"}
	p2 := &properties.Property{LandlordID: 1, Title: "City Loft", Description: "Modern", Location: "Nairobi", Price: 40000, Status: "approved"}
	p3 := &properties.Property{LandlordID: 1, Title: "Pending Flat", Description: "Under review", Location: "Nakuru", Price: 15000, Status: "pending"}
	_ = propRepo.Create(p1)
	_ = propRepo.Create(p2)
	_ = propRepo.Create(p3)

	userID := int64(10)

	// Check initially not favorited
	isFav, err := favRepo.IsFavorite(userID, p1.ID)
	if err != nil || isFav {
		t.Fatalf("expected isFavorite false, got %v (err: %v)", isFav, err)
	}

	// Add to favorites
	fav, err := favRepo.Add(userID, p1.ID)
	if err != nil || fav == nil || fav.ID == 0 {
		t.Fatalf("failed to add favorite: %v", err)
	}

	// Verify IsFavorite is true
	isFav, err = favRepo.IsFavorite(userID, p1.ID)
	if err != nil || !isFav {
		t.Fatalf("expected isFavorite true, got %v (err: %v)", isFav, err)
	}

	// Add second property to favorites
	_, err = favRepo.Add(userID, p2.ID)
	if err != nil {
		t.Fatalf("failed to add second favorite: %v", err)
	}

	// Add pending property to favorites
	_, err = favRepo.Add(userID, p3.ID)
	if err != nil {
		t.Fatalf("failed to add third favorite: %v", err)
	}

	// ListFavoritePropertyIDs returns all 3 IDs
	ids, err := favRepo.ListFavoritePropertyIDs(userID)
	if err != nil || len(ids) != 3 {
		t.Fatalf("expected 3 favorite ids, got %d (err: %v)", len(ids), err)
	}

	// ListByUser only returns approved properties (p1 and p2, not p3)
	favProps, err := favRepo.ListByUser(userID)
	if err != nil || len(favProps) != 2 {
		t.Fatalf("expected 2 approved favorites, got %d (err: %v)", len(favProps), err)
	}

	// Remove p1
	if err := favRepo.Remove(userID, p1.ID); err != nil {
		t.Fatalf("failed to remove favorite: %v", err)
	}

	// Remove again returns ErrFavoriteNotFound
	if err := favRepo.Remove(userID, p1.ID); !errors.Is(err, ErrFavoriteNotFound) {
		t.Fatalf("expected ErrFavoriteNotFound on second remove, got %v", err)
	}

	// IsFavorite for p1 is now false
	isFav, _ = favRepo.IsFavorite(userID, p1.ID)
	if isFav {
		t.Errorf("expected isFavorite false after deletion")
	}

	// Another user has empty favorites
	otherUserProps, err := favRepo.ListByUser(999)
	if err != nil || len(otherUserProps) != 0 {
		t.Errorf("expected 0 favorites for other user, got %d", len(otherUserProps))
	}
}
