package users

import "time"

type User struct {
	ID        int64     `json:"id"`
	RoleID    int64     `json:"role_id"`
	RoleName  string    `json:"role,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}