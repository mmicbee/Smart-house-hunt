package users 

import (
	"database/sql"
) 

type Repository struct {
	db *sql.DB 
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) RoleIDByName(name string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM roles WHERE name = ?", name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) RoleByID(id int64) (string, error) {
	var name string
	err := r.db.QueryRow("SELECT name FROM roles WHERE id = ?", id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *Repository) List() ([]User, error) {
	rows, err := r.db.Query(`
		SELECT u.id, u.role_id, COALESCE(r.name, ''), u.name, u.email, COALESCE(u.status, 'active'), u.created_at, u.updated_at
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.RoleName,
			&user.Name,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		all = append(all, user)
	}

	return all, rows.Err()
}

func (r *Repository) Create(user *User) error {
	if user.Status == "" {
		user.Status = "active"
	}
	query := `
	INSERT INTO users (role_id, name, email, password, status)
	VALUES (?, ?, ?, ?, ?)
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.RoleID,
		user.Name,
		user.Email,
		user.Password,
		user.Status,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	query := `
	SELECT id, role_id, name, email, password, COALESCE(status, 'active'), created_at, updated_at
	FROM users
	WHERE email = ?
	`

	user := &User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) FindByID(id int64) (*User, error) {
	query := `
		SELECT id, role_id, name, email, password, COALESCE(status, 'active'), created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &User{}

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) UpdateStatus(id int64, status string) error {
	query := `
		UPDATE users
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	res, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

