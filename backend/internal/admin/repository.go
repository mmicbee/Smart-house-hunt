package admin

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"smart-house-hunt/internal/properties"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListUsers(status, role, search string) ([]AdminUser, error) {
	query := `
		SELECT u.id, u.role_id, COALESCE(r.name, ''), u.name, u.email, COALESCE(u.status, 'active'), u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE 1=1
	`
	var args []any

	if s := strings.TrimSpace(status); s != "" {
		query += " AND LOWER(u.status) = ?"
		args = append(args, strings.ToLower(s))
	}
	if ro := strings.TrimSpace(role); ro != "" {
		query += " AND LOWER(r.name) = ?"
		args = append(args, strings.ToLower(ro))
	}
	if q := strings.TrimSpace(search); q != "" {
		query += " AND (LOWER(u.name) LIKE ? OR LOWER(u.email) LIKE ?)"
		likePattern := "%" + strings.ToLower(q) + "%"
		args = append(args, likePattern, likePattern)
	}

	query += " ORDER BY u.id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(
			&u.ID,
			&u.RoleID,
			&u.Role,
			&u.Name,
			&u.Email,
			&u.Status,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, u)
	}

	if list == nil {
		list = make([]AdminUser, 0)
	}
	return list, rows.Err()
}

func (r *Repository) GetUserByID(id int64) (*AdminUser, error) {
	query := `
		SELECT u.id, u.role_id, COALESCE(r.name, ''), u.name, u.email, COALESCE(u.status, 'active'), u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`
	var u AdminUser
	err := r.db.QueryRow(query, id).Scan(
		&u.ID,
		&u.RoleID,
		&u.Role,
		&u.Name,
		&u.Email,
		&u.Status,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) UpdateUserStatus(id int64, status string) error {
	query := `
		UPDATE users
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) ListProperties(status, propertyType string) ([]AdminProperty, error) {
	query := `
		SELECT p.id, p.landlord_id, COALESCE(u.name, ''), COALESCE(u.email, ''),
		       p.title, p.description, p.location, p.price, p.bedrooms, p.bathrooms,
		       p.amenities, p.property_type, p.status, p.created_at, p.updated_at
		FROM properties p
		LEFT JOIN users u ON p.landlord_id = u.id
		WHERE 1=1
	`
	var args []any

	if s := strings.TrimSpace(status); s != "" && s != "all" {
		query += " AND LOWER(p.status) = ?"
		args = append(args, strings.ToLower(s))
	}
	if pt := strings.TrimSpace(propertyType); pt != "" && pt != "all" {
		query += " AND LOWER(p.property_type) = ?"
		args = append(args, strings.ToLower(pt))
	}

	query += " ORDER BY p.id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AdminProperty
	for rows.Next() {
		var p AdminProperty
		if err := rows.Scan(
			&p.ID,
			&p.LandlordID,
			&p.LandlordName,
			&p.LandlordEmail,
			&p.Title,
			&p.Description,
			&p.Location,
			&p.Price,
			&p.Bedrooms,
			&p.Bathrooms,
			&p.Amenities,
			&p.PropertyType,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.Images = r.getImagesForProperty(p.ID)
		list = append(list, p)
	}

	if list == nil {
		list = make([]AdminProperty, 0)
	}
	return list, rows.Err()
}

func (r *Repository) GetPropertyByID(id int64) (*AdminProperty, error) {
	query := `
		SELECT p.id, p.landlord_id, COALESCE(u.name, ''), COALESCE(u.email, ''),
		       p.title, p.description, p.location, p.price, p.bedrooms, p.bathrooms,
		       p.amenities, p.property_type, p.status, p.created_at, p.updated_at
		FROM properties p
		LEFT JOIN users u ON p.landlord_id = u.id
		WHERE p.id = ?
	`
	var p AdminProperty
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.LandlordID,
		&p.LandlordName,
		&p.LandlordEmail,
		&p.Title,
		&p.Description,
		&p.Location,
		&p.Price,
		&p.Bedrooms,
		&p.Bathrooms,
		&p.Amenities,
		&p.PropertyType,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Images = r.getImagesForProperty(p.ID)
	return &p, nil
}

func (r *Repository) UpdatePropertyStatus(id int64, status string) (*AdminProperty, error) {
	query := `
		UPDATE properties
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := r.db.Exec(query, status, id)
	if err != nil {
		return nil, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}
	return r.GetPropertyByID(id)
}

func (r *Repository) DeleteProperty(id int64) error {
	res, err := r.db.Exec("DELETE FROM properties WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) getImagesForProperty(propertyID int64) []properties.PropertyImage {
	rows, err := r.db.Query(`
		SELECT id, property_id, file_path, created_at
		FROM property_images
		WHERE property_id = ?
		ORDER BY id
	`, propertyID)
	if err != nil {
		return make([]properties.PropertyImage, 0)
	}
	defer rows.Close()

	var images []properties.PropertyImage
	for rows.Next() {
		var img properties.PropertyImage
		if err := rows.Scan(&img.ID, &img.PropertyID, &img.FilePath, &img.CreatedAt); err == nil {
			images = append(images, img)
		}
	}
	if images == nil {
		images = make([]properties.PropertyImage, 0)
	}
	return images
}

func (r *Repository) GetReports() (*PlatformReports, error) {
	reports := &PlatformReports{
		Properties: PropertyStats{
			ByType: make(map[string]int),
		},
		Recent: make([]RecentActivityItem, 0),
	}

	// 1. User stats
	row := r.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'suspended' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN role_id = 1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN role_id = 2 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN role_id = 3 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN role_id = 4 THEN 1 ELSE 0 END), 0)
		FROM users
	`)
	_ = row.Scan(
		&reports.Users.Total,
		&reports.Users.Active,
		&reports.Users.Suspended,
		&reports.Users.Tenants,
		&reports.Users.Landlords,
		&reports.Users.Agents,
		&reports.Users.Admins,
	)

	// 2. Property stats
	pRow := r.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'suspended' THEN 1 ELSE 0 END), 0)
		FROM properties
	`)
	_ = pRow.Scan(
		&reports.Properties.Total,
		&reports.Properties.Pending,
		&reports.Properties.Approved,
		&reports.Properties.Rejected,
		&reports.Properties.Suspended,
	)

	// Property types breakdown
	ptRows, err := r.db.Query(`
		SELECT property_type, COUNT(*)
		FROM properties
		GROUP BY property_type
	`)
	if err == nil {
		defer ptRows.Close()
		for ptRows.Next() {
			var pType string
			var count int
			if err := ptRows.Scan(&pType, &count); err == nil {
				reports.Properties.ByType[pType] = count
			}
		}
	}

	// 3. Viewing stats (if viewing_requests table exists)
	vRow := r.db.QueryRow(`
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END), 0)
		FROM viewing_requests
	`)
	_ = vRow.Scan(
		&reports.Viewings.Total,
		&reports.Viewings.Pending,
		&reports.Viewings.Accepted,
		&reports.Viewings.Rejected,
		&reports.Viewings.Cancelled,
	)

	// 4. Favorites count
	_ = r.db.QueryRow(`SELECT COUNT(*) FROM favorites`).Scan(&reports.Favorites)

	// 5. Recent activity log
	recentItems := make([]RecentActivityItem, 0)

	// Recent properties
	pRecentRows, err := r.db.Query(`
		SELECT title, location, status, created_at
		FROM properties
		ORDER BY id DESC
		LIMIT 4
	`)
	if err == nil {
		defer pRecentRows.Close()
		for pRecentRows.Next() {
			var title, loc, status string
			var createdAt time.Time
			if err := pRecentRows.Scan(&title, &loc, &status, &createdAt); err == nil {
				recentItems = append(recentItems, RecentActivityItem{
					Type:      "property",
					Title:     fmt.Sprintf("New Property: %s", title),
					Details:   fmt.Sprintf("Location: %s (Status: %s)", loc, status),
					Timestamp: createdAt,
				})
			}
		}
	}

	// Recent users
	uRecentRows, err := r.db.Query(`
		SELECT u.name, COALESCE(r.name, 'user'), u.created_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.id DESC
		LIMIT 4
	`)
	if err == nil {
		defer uRecentRows.Close()
		for uRecentRows.Next() {
			var name, role string
			var createdAt time.Time
			if err := uRecentRows.Scan(&name, &role, &createdAt); err == nil {
				recentItems = append(recentItems, RecentActivityItem{
					Type:      "user",
					Title:     fmt.Sprintf("User Registered: %s", name),
					Details:   fmt.Sprintf("Role: %s", role),
					Timestamp: createdAt,
				})
			}
		}
	}

	// Recent viewing requests
	vRecentRows, err := r.db.Query(`
		SELECT preferred_date, status, created_at
		FROM viewing_requests
		ORDER BY id DESC
		LIMIT 4
	`)
	if err == nil {
		defer vRecentRows.Close()
		for vRecentRows.Next() {
			var date, status string
			var createdAt time.Time
			if err := vRecentRows.Scan(&date, &status, &createdAt); err == nil {
				recentItems = append(recentItems, RecentActivityItem{
					Type:      "viewing",
					Title:     fmt.Sprintf("Viewing Scheduled for %s", date),
					Details:   fmt.Sprintf("Status: %s", status),
					Timestamp: createdAt,
				})
			}
		}
	}

	// Sort recent items descending by timestamp
	sort.Slice(recentItems, func(i, j int) bool {
		return recentItems[i].Timestamp.After(recentItems[j].Timestamp)
	})

	if len(recentItems) > 8 {
		recentItems = recentItems[:8]
	}
	reports.Recent = recentItems

	return reports, nil
}
