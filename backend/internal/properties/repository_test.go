package properties

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
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
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
		property_type TEXT NOT NULL DEFAULT 'apartment',
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
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}
	return db
}

func TestRepositoryCreateAndFindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	prop := &Property{
		LandlordID:  10,
		Title:       "Test Villa",
		Description: "A nice villa",
		Location:    "Nairobi",
		Price:       50000,
		Bedrooms:    3,
		Bathrooms:   2,
		Amenities:   "parking,pool",
		Status:      "pending",
	}

	if err := repo.Create(prop); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if prop.ID == 0 {
		t.Fatalf("expected non-zero property ID after Create")
	}

	found, err := repo.FindByID(prop.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if found.Title != "Test Villa" || found.LandlordID != 10 {
		t.Errorf("unexpected property content: %+v", found)
	}

	// Non-existent ID
	_, err = repo.FindByID(9999)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected sql.ErrNoRows for missing id, got %v", err)
	}
}

func TestRepositoryUpdateAndDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	prop := &Property{
		LandlordID:  1,
		Title:       "Initial Title",
		Description: "Initial Desc",
		Location:    "Eldoret",
		Price:       12000,
		Bedrooms:    1,
		Bathrooms:   1,
		Amenities:   "water",
		Status:      "pending",
	}
	if err := repo.Create(prop); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	prop.Title = "Updated Title"
	prop.Price = 14000
	prop.Bedrooms = 2
	if err := repo.Update(prop); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	updated, err := repo.FindByID(prop.ID)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	if updated.Title != "Updated Title" || updated.Price != 14000 || updated.Bedrooms != 2 {
		t.Errorf("property was not updated correctly: %+v", updated)
	}

	if err := repo.Delete(prop.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	_, err = repo.FindByID(prop.ID)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected ErrNoRows after Delete, got %v", err)
	}
}

func TestRepositoryListByLandlord(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	p1 := &Property{LandlordID: 1, Title: "P1", Description: "D1", Location: "L1", Price: 10000, Status: "pending"}
	p2 := &Property{LandlordID: 2, Title: "P2", Description: "D2", Location: "L2", Price: 20000, Status: "approved"}
	p3 := &Property{LandlordID: 1, Title: "P3", Description: "D3", Location: "L3", Price: 30000, Status: "approved"}

	_ = repo.Create(p1)
	_ = repo.Create(p2)
	_ = repo.Create(p3)

	list1, err := repo.ListByLandlord(1)
	if err != nil {
		t.Fatalf("ListByLandlord returned error: %v", err)
	}
	if len(list1) != 2 {
		t.Fatalf("expected 2 properties for landlord 1, got %d", len(list1))
	}
	// Ordered by id DESC
	if list1[0].Title != "P3" || list1[1].Title != "P1" {
		t.Errorf("unexpected list order or contents: %+v", list1)
	}

	list2, err := repo.ListByLandlord(2)
	if err != nil {
		t.Fatalf("ListByLandlord returned error: %v", err)
	}
	if len(list2) != 1 || list2[0].Title != "P2" {
		t.Errorf("unexpected list for landlord 2: %+v", list2)
	}

	emptyList, err := repo.ListByLandlord(99)
	if err != nil {
		t.Fatalf("ListByLandlord returned error: %v", err)
	}
	if len(emptyList) != 0 {
		t.Errorf("expected empty list for landlord 99, got %d", len(emptyList))
	}
}

func TestRepositoryListApproved(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	p1 := &Property{LandlordID: 1, Title: "P1", Description: "D1", Location: "L1", Price: 10000, Status: "pending"}
	p2 := &Property{LandlordID: 1, Title: "P2", Description: "D2", Location: "L2", Price: 20000, Status: "approved"}
	p3 := &Property{LandlordID: 2, Title: "P3", Description: "D3", Location: "L3", Price: 30000, Status: "rejected"}
	p4 := &Property{LandlordID: 2, Title: "P4", Description: "D4", Location: "L4", Price: 40000, Status: "approved"}

	_ = repo.Create(p1)
	_ = repo.Create(p2)
	_ = repo.Create(p3)
	_ = repo.Create(p4)

	approved, err := repo.ListApproved()
	if err != nil {
		t.Fatalf("ListApproved returned error: %v", err)
	}
	if len(approved) != 2 {
		t.Fatalf("expected 2 approved properties, got %d", len(approved))
	}
	// Ordered by id DESC
	if approved[0].Title != "P4" || approved[1].Title != "P2" {
		t.Errorf("unexpected approved items or order: %+v", approved)
	}
}

func TestRepositoryImages(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	prop := &Property{LandlordID: 1, Title: "P1", Description: "D1", Location: "L1", Price: 10000, Status: "approved"}
	_ = repo.Create(prop)

	img1 := &PropertyImage{PropertyID: prop.ID, FilePath: "uploads/property-images/1/img1.png"}
	img2 := &PropertyImage{PropertyID: prop.ID, FilePath: "uploads/property-images/1/img2.png"}
	if err := repo.AddImage(img1); err != nil {
		t.Fatalf("AddImage returned error: %v", err)
	}
	if err := repo.AddImage(img2); err != nil {
		t.Fatalf("AddImage returned error: %v", err)
	}

	images, err := repo.ListImagesByProperty(prop.ID)
	if err != nil {
		t.Fatalf("ListImagesByProperty returned error: %v", err)
	}
	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
	if images[0].FilePath != "uploads/property-images/1/img1.png" || images[1].FilePath != "uploads/property-images/1/img2.png" {
		t.Errorf("unexpected images returned: %+v", images)
	}

	emptyImgs, err := repo.ListImagesByProperty(999)
	if err != nil {
		t.Fatalf("ListImagesByProperty returned error: %v", err)
	}
	if len(emptyImgs) != 0 {
		t.Errorf("expected 0 images, got %d", len(emptyImgs))
	}
}

func TestRepositoryListApprovedFilters(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewRepository(db)

	p1 := &Property{LandlordID: 1, Title: "Kilimani 2BHK", Description: "D1", Location: "Kilimani, Nairobi", Price: 35000, Bedrooms: 2, Bathrooms: 2, Amenities: "wifi, parking", Status: "approved"}
	p2 := &Property{LandlordID: 1, Title: "Westlands 1BHK", Description: "D2", Location: "Westlands, Nairobi", Price: 20000, Bedrooms: 1, Bathrooms: 1, Amenities: "wifi", Status: "approved"}
	p3 := &Property{LandlordID: 1, Title: "Karen Mansion", Description: "D3", Location: "Karen, Nairobi", Price: 80000, Bedrooms: 4, Bathrooms: 3, Amenities: "swimming pool, garden, parking", Status: "approved"}
	p4 := &Property{LandlordID: 1, Title: "Pending House", Description: "D4", Location: "Kilimani, Nairobi", Price: 25000, Bedrooms: 2, Bathrooms: 1, Amenities: "wifi", Status: "pending"}

	p5 := &Property{LandlordID: 1, Title: "Roysambu Bedsitter", Description: "D5", Location: "Roysambu, Nairobi", Price: 9000, Bedrooms: 0, Bathrooms: 1, PropertyType: "bedsitter", Status: "approved"}
	p6 := &Property{LandlordID: 1, Title: "Parklands Studio", Description: "D6", Location: "Parklands, Nairobi", Price: 18000, Bedrooms: 0, Bathrooms: 1, PropertyType: "studio", Status: "approved"}

	for _, p := range []*Property{p1, p2, p3, p4, p5, p6} {
		if err := repo.Create(p); err != nil {
			t.Fatalf("failed to create property: %v", err)
		}
	}

	// 1. Filter by location
	res, err := repo.ListApproved(FilterParams{Location: "karen"})
	if err != nil || len(res) != 1 || res[0].Title != "Karen Mansion" {
		t.Errorf("location filter failed, got %+v (err: %v)", res, err)
	}

	// 2. Filter by price range
	minP := int64(20000)
	maxP := int64(40000)
	res, err = repo.ListApproved(FilterParams{MinPrice: &minP, MaxPrice: &maxP})
	if err != nil || len(res) != 2 {
		t.Errorf("price range filter failed, expected 2 got %d (err: %v)", len(res), err)
	}

	// 3. Filter by bedrooms and bathrooms
	beds := 2
	baths := 2
	res, err = repo.ListApproved(FilterParams{Bedrooms: &beds, Bathrooms: &baths})
	if err != nil || len(res) != 1 || res[0].Title != "Kilimani 2BHK" {
		t.Errorf("bedrooms/bathrooms filter failed, got %+v (err: %v)", res, err)
	}

	// 4. Filter by amenities
	res, err = repo.ListApproved(FilterParams{Amenities: "parking"})
	if err != nil || len(res) != 2 {
		t.Errorf("amenities filter failed, expected 2 got %d (err: %v)", len(res), err)
	}

	// 5. Combined filters
	res, err = repo.ListApproved(FilterParams{
		Location:  "nairobi",
		MaxPrice:  &maxP,
		Bedrooms:  &beds,
		Amenities: "wifi",
	})
	if err != nil || len(res) != 1 || res[0].Title != "Kilimani 2BHK" {
		t.Errorf("combined filters failed, got %+v (err: %v)", res, err)
	}

	// 6. Filter by property type (bedsitter)
	res, err = repo.ListApproved(FilterParams{PropertyType: "bedsitter"})
	if err != nil || len(res) != 1 || res[0].Title != "Roysambu Bedsitter" {
		t.Errorf("bedsitter filter failed, got %+v (err: %v)", res, err)
	}

	// 7. Filter by property type (studio)
	res, err = repo.ListApproved(FilterParams{PropertyType: "studio"})
	if err != nil || len(res) != 1 || res[0].Title != "Parklands Studio" {
		t.Errorf("studio filter failed, got %+v (err: %v)", res, err)
	}
}
