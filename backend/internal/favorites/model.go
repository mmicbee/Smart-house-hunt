package favorites

import "time"

type Favorite struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	PropertyID int64     `json:"property_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FavoriteRequest struct {
	PropertyID int64 `json:"property_id"`
}
