package properties

import (
	"strings"
	"time"
)

type PropertyRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	Price        int64  `json:"price"`
	Bedrooms     int    `json:"bedrooms"`
	Bathrooms    int    `json:"bathrooms"`
	Amenities    string `json:"amenities"`
	PropertyType string `json:"property_type"`
}

type PropertyResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Price        int64     `json:"price"`
	Bedrooms     int       `json:"bedrooms"`
	Bathrooms    int       `json:"bathrooms"`
	Amenities    string    `json:"amenities"`
	PropertyType string    `json:"property_type"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type PropertyListItem struct {
	Property
	Images []PropertyImage `json:"images"`
}

func toResponse(property *Property) PropertyResponse {
	return PropertyResponse{
		ID:           property.ID,
		Title:        property.Title,
		Description:  property.Description,
		Location:     property.Location,
		Price:        property.Price,
		Bedrooms:     property.Bedrooms,
		Bathrooms:    property.Bathrooms,
		Amenities:    property.Amenities,
		PropertyType: property.PropertyType,
		Status:       property.Status,
		CreatedAt:    property.CreatedAt,
	}
}

type FilterParams struct {
	Location     string
	PropertyType string
	MinPrice     *int64
	MaxPrice     *int64
	Bedrooms     *int
	MinBedrooms  *int
	Bathrooms    *int
	Amenities    string
}

func detectPropertyType(title string, bedrooms int) string {
	lower := strings.ToLower(title)
	if strings.Contains(lower, "bedsitter") {
		return "bedsitter"
	}
	if strings.Contains(lower, "studio") {
		return "studio"
	}
	if strings.Contains(lower, "villa") || strings.Contains(lower, "mansionette") || strings.Contains(lower, "bungalow") || strings.Contains(lower, "house") {
		return "house"
	}
	return "apartment"
}