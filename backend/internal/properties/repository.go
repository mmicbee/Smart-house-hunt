package properties

import "database/sql"

type Repository struct {
	db *sql.DB
}

const propertyColumns = "id, landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, status, created_at, updated_at"

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(property *Property) error {
	query := `
		INSERT INTO properties (landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, status, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		property.LandlordID,
		property.Title,
		property.Description,
		property.Location,
		property.Price,
		property.Bedrooms,
		property.Bathrooms,
		property.Amenities,
		property.Status,
	).Scan(
		&property.ID,
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
}

func (r *Repository) FindByID(id int64) (*Property, error) {
	query := `
		SELECT ` + propertyColumns + `
		FROM properties
		WHERE id = ?
	`

	property := &Property{}
	err := r.db.QueryRow(query, id).Scan(
		&property.ID,
		&property.LandlordID,
		&property.Title,
		&property.Description,
		&property.Location,
		&property.Price,
		&property.Bedrooms,
		&property.Bathrooms,
		&property.Amenities,
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return property, nil
}

func (r *Repository) ListByLandlord(landlordID int64) ([]Property, error) {
	rows, err := r.db.Query(`
		SELECT `+propertyColumns+`
		FROM properties
		WHERE landlord_id = ?
		ORDER BY id DESC
	`, landlordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Property, 0)
	for rows.Next() {
		var property Property
		if err := rows.Scan(
			&property.ID,
			&property.LandlordID,
			&property.Title,
			&property.Description,
			&property.Location,
			&property.Price,
			&property.Bedrooms,
			&property.Bathrooms,
			&property.Amenities,
			&property.Status,
			&property.CreatedAt,
			&property.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, property)
	}

	return list, rows.Err()
}

func (r *Repository) ListApproved() ([]Property, error) {
	rows, err := r.db.Query(`
		SELECT `+propertyColumns+`
		FROM properties
		WHERE status = 'approved'
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Property, 0)
	for rows.Next() {
		var property Property
		if err := rows.Scan(
			&property.ID,
			&property.LandlordID,
			&property.Title,
			&property.Description,
			&property.Location,
			&property.Price,
			&property.Bedrooms,
			&property.Bathrooms,
			&property.Amenities,
			&property.Status,
			&property.CreatedAt,
			&property.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, property)
	}

	return list, rows.Err()
}

func (r *Repository) ListImagesByProperty(propertyID int64) ([]PropertyImage, error) {
	rows, err := r.db.Query(`
		SELECT id, property_id, file_path, created_at
		FROM property_images
		WHERE property_id = ?
		ORDER BY id
	`, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]PropertyImage, 0)
	for rows.Next() {
		var image PropertyImage
		if err := rows.Scan(&image.ID, &image.PropertyID, &image.FilePath, &image.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, image)
	}

	return list, rows.Err()
}

func (r *Repository) Update(property *Property) error {
	query := `
		UPDATE properties
		SET title = ?, description = ?, location = ?, price = ?,
			bedrooms = ?, bathrooms = ?, amenities = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING status, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		property.Title,
		property.Description,
		property.Location,
		property.Price,
		property.Bedrooms,
		property.Bathrooms,
		property.Amenities,
		property.ID,
	).Scan(
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
	)
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM properties WHERE id = ?", id)
	return err
}

func (r *Repository) AddImage(image *PropertyImage) error {
	query := `
		INSERT INTO property_images (property_id, file_path)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	return r.db.QueryRow(query, image.PropertyID, image.FilePath).Scan(
		&image.ID,
		&image.CreatedAt,
	)
}