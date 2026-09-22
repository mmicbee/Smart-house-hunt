package properties

import "database/sql"

type Repository struct {
	db *sql.DB
}

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
		SELECT id, landlord_id, title, description, location, price, bedrooms, bathrooms, amenities, status, created_at, updated_at
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