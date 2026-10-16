package favorites

import (
	"database/sql"
	"errors"

	"smart-house-hunt/internal/properties"
)

var ErrFavoriteNotFound = errors.New("favorite not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Add(userID, propertyID int64) (*Favorite, error) {
	fav := &Favorite{
		UserID:     userID,
		PropertyID: propertyID,
	}

	query := `
		INSERT INTO favorites (user_id, property_id)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, userID, propertyID).Scan(&fav.ID, &fav.CreatedAt)
	if err != nil {
		return nil, err
	}

	return fav, nil
}

func (r *Repository) Remove(userID, propertyID int64) error {
	res, err := r.db.Exec(`
		DELETE FROM favorites
		WHERE user_id = ? AND property_id = ?
	`, userID, propertyID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrFavoriteNotFound
	}

	return nil
}

func (r *Repository) IsFavorite(userID, propertyID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow(`
		SELECT 1 FROM favorites
		WHERE user_id = ? AND property_id = ?
		LIMIT 1
	`, userID, propertyID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Repository) ListFavoritePropertyIDs(userID int64) ([]int64, error) {
	rows, err := r.db.Query(`
		SELECT property_id
		FROM favorites
		WHERE user_id = ?
		ORDER BY id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (r *Repository) ListByUser(userID int64) ([]properties.Property, error) {
	query := `
		SELECT p.id, p.landlord_id, p.title, p.description, p.location, p.price,
		       p.bedrooms, p.bathrooms, p.amenities, p.status, p.created_at, p.updated_at
		FROM favorites f
		JOIN properties p ON f.property_id = p.id
		WHERE f.user_id = ? AND p.status = 'approved'
		ORDER BY f.id DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]properties.Property, 0)
	for rows.Next() {
		var property properties.Property
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
