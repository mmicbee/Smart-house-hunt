package users 

import (
	"database/sql"
) 

type Repository struct {
	db *sql.DB 
}

func Newrepository(db *sql.DB) * Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *User) error {
	query := `
	INSERT INTO users (role_id, name, email, password)
	VALUES (?, ?, ?, ?)
	RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		query,
		user.RoleID,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UPdatedAt,
	)
}

func (r *Repository) FindByEmail(email string) (*user, error) {
	query := `
	SELECT id, role_id, name, email, password, created_at, updated_at
	FROM users
	WHERE EMAIL = ?
	`

	user := &User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.Password,
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
		SELECT id, role_id, name, email, password, created_at, updated_at
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
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
