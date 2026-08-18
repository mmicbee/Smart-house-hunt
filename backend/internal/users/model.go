package users

import "time"

type User struct {
	ID          int64            `json:"id"`
	RoleID      int64          `json:"role_id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Password    string        `json:"-"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}