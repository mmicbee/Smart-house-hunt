package properties

import "time"

type PropertyRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Price       int64  `json:"price"`
	Bedrooms    int    `json:"bedrooms"`
	Bathrooms   int    `json:"bathrooms"`
	Amenities   string `json:"amenities"`
}

type PropertyResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Price       int64     `json:"price"`
	Bedrooms    int       `json:"bedrooms"`
	Bathrooms   int       `json:"bathrooms"`
	Amenities   string    `json:"amenities"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type PropertyListItem struct {
	Property
	Images []PropertyImage `json:"images"`
}

func toResponse(property *Property) PropertyResponse {
	return PropertyResponse{
		ID:          property.ID,
		Title:       property.Title,
		Description: property.Description,
		Location:    property.Location,
		Price:       property.Price,
		Bedrooms:    property.Bedrooms,
		Bathrooms:   property.Bathrooms,
		Amenities:   property.Amenities,
		Status:      property.Status,
		CreatedAt:   property.CreatedAt,
	}
}