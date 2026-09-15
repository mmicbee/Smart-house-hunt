package properties

import "time"

type Property struct {
	ID          int64     `json:"id"`
	LandlordID  int64     `json:"landlord_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Price       int64     `json:"price"`
	Bedrooms    int       `json:"bedrooms"`
	Bathrooms   int       `json:"bathrooms"`
	Amenities   string    `json:"amenities"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}