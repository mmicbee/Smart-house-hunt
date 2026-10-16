package viewing

import (
	"database/sql"
	"errors"
)

var ErrViewingNotFound = errors.New("viewing request not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(req *ViewingRequest) error {
	query := `
		INSERT INTO viewing_requests (user_id, property_id, preferred_date, message, status)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id, status, created_at, updated_at
	`
	if req.Status == "" {
		req.Status = StatusPending
	}

	return r.db.QueryRow(
		query,
		req.UserID,
		req.PropertyID,
		req.PreferredDate,
		req.Message,
		req.Status,
	).Scan(
		&req.ID,
		&req.Status,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
}

func (r *Repository) HasPendingRequest(userID, propertyID int64) (bool, error) {
	var exists int
	err := r.db.QueryRow(`
		SELECT 1 FROM viewing_requests
		WHERE user_id = ? AND property_id = ? AND status = ?
		LIMIT 1
	`, userID, propertyID, StatusPending).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Repository) FindByID(id int64) (*ViewingRequestDetail, error) {
	query := `
		SELECT
			v.id, v.user_id, v.property_id, v.preferred_date, v.message, v.status, v.created_at, v.updated_at,
			p.title, p.location, p.price, p.landlord_id,
			landlord.name, landlord.email,
			tenant.name, tenant.email,
			COALESCE((SELECT file_path FROM property_images WHERE property_id = p.id LIMIT 1), '')
		FROM viewing_requests v
		JOIN properties p ON v.property_id = p.id
		JOIN users landlord ON p.landlord_id = landlord.id
		JOIN users tenant ON v.user_id = tenant.id
		WHERE v.id = ?
	`

	detail := &ViewingRequestDetail{}
	err := r.db.QueryRow(query, id).Scan(
		&detail.ID,
		&detail.UserID,
		&detail.PropertyID,
		&detail.PreferredDate,
		&detail.Message,
		&detail.Status,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.PropertyTitle,
		&detail.PropertyLocation,
		&detail.PropertyPrice,
		&detail.LandlordID,
		&detail.LandlordName,
		&detail.LandlordEmail,
		&detail.TenantName,
		&detail.TenantEmail,
		&detail.PropertyImage,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrViewingNotFound
	}
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (r *Repository) ListByTenant(tenantID int64) ([]ViewingRequestDetail, error) {
	query := `
		SELECT
			v.id, v.user_id, v.property_id, v.preferred_date, v.message, v.status, v.created_at, v.updated_at,
			p.title, p.location, p.price, p.landlord_id,
			landlord.name, landlord.email,
			tenant.name, tenant.email,
			COALESCE((SELECT file_path FROM property_images WHERE property_id = p.id LIMIT 1), '')
		FROM viewing_requests v
		JOIN properties p ON v.property_id = p.id
		JOIN users landlord ON p.landlord_id = landlord.id
		JOIN users tenant ON v.user_id = tenant.id
		WHERE v.user_id = ?
		ORDER BY v.created_at DESC
	`

	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]ViewingRequestDetail, 0)
	for rows.Next() {
		var detail ViewingRequestDetail
		if err := rows.Scan(
			&detail.ID,
			&detail.UserID,
			&detail.PropertyID,
			&detail.PreferredDate,
			&detail.Message,
			&detail.Status,
			&detail.CreatedAt,
			&detail.UpdatedAt,
			&detail.PropertyTitle,
			&detail.PropertyLocation,
			&detail.PropertyPrice,
			&detail.LandlordID,
			&detail.LandlordName,
			&detail.LandlordEmail,
			&detail.TenantName,
			&detail.TenantEmail,
			&detail.PropertyImage,
		); err != nil {
			return nil, err
		}
		list = append(list, detail)
	}

	return list, rows.Err()
}

func (r *Repository) ListByLandlord(landlordID int64) ([]ViewingRequestDetail, error) {
	query := `
		SELECT
			v.id, v.user_id, v.property_id, v.preferred_date, v.message, v.status, v.created_at, v.updated_at,
			p.title, p.location, p.price, p.landlord_id,
			landlord.name, landlord.email,
			tenant.name, tenant.email,
			COALESCE((SELECT file_path FROM property_images WHERE property_id = p.id LIMIT 1), '')
		FROM viewing_requests v
		JOIN properties p ON v.property_id = p.id
		JOIN users landlord ON p.landlord_id = landlord.id
		JOIN users tenant ON v.user_id = tenant.id
		WHERE p.landlord_id = ?
		ORDER BY v.created_at DESC
	`

	rows, err := r.db.Query(query, landlordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]ViewingRequestDetail, 0)
	for rows.Next() {
		var detail ViewingRequestDetail
		if err := rows.Scan(
			&detail.ID,
			&detail.UserID,
			&detail.PropertyID,
			&detail.PreferredDate,
			&detail.Message,
			&detail.Status,
			&detail.CreatedAt,
			&detail.UpdatedAt,
			&detail.PropertyTitle,
			&detail.PropertyLocation,
			&detail.PropertyPrice,
			&detail.LandlordID,
			&detail.LandlordName,
			&detail.LandlordEmail,
			&detail.TenantName,
			&detail.TenantEmail,
			&detail.PropertyImage,
		); err != nil {
			return nil, err
		}
		list = append(list, detail)
	}

	return list, rows.Err()
}

func (r *Repository) UpdateStatus(id int64, status string) error {
	res, err := r.db.Exec(`
		UPDATE viewing_requests
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrViewingNotFound
	}

	return nil
}
