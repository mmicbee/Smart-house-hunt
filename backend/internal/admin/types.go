package admin

import (
	"time"

	"smart-house-hunt/internal/properties"
)

type AdminUser struct {
	ID        int64     `json:"id"`
	RoleID    int64     `json:"role_id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AdminProperty struct {
	ID            int64                      `json:"id"`
	LandlordID    int64                      `json:"landlord_id"`
	LandlordName  string                     `json:"landlord_name"`
	LandlordEmail string                     `json:"landlord_email"`
	Title         string                     `json:"title"`
	Description   string                     `json:"description"`
	Location      string                     `json:"location"`
	Price         int64                      `json:"price"`
	Bedrooms      int                        `json:"bedrooms"`
	Bathrooms     int                        `json:"bathrooms"`
	Amenities     string                     `json:"amenities"`
	PropertyType  string                     `json:"property_type"`
	Status        string                     `json:"status"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	Images        []properties.PropertyImage `json:"images"`
}

type UpdatePropertyStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UserStats struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Suspended int `json:"suspended"`
	Tenants   int `json:"tenants"`
	Landlords int `json:"landlords"`
	Agents    int `json:"agents"`
	Admins    int `json:"admins"`
}

type PropertyStats struct {
	Total     int            `json:"total"`
	Pending   int            `json:"pending"`
	Approved  int            `json:"approved"`
	Rejected  int            `json:"rejected"`
	Suspended int            `json:"suspended"`
	ByType    map[string]int `json:"by_type"`
}

type ViewingStats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Accepted  int `json:"accepted"`
	Rejected  int `json:"rejected"`
	Cancelled int `json:"cancelled"`
}

type RecentActivityItem struct {
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

type PlatformReports struct {
	Users      UserStats            `json:"users"`
	Properties PropertyStats        `json:"properties"`
	Viewings   ViewingStats         `json:"viewings"`
	Favorites  int                  `json:"total_favorites"`
	Recent     []RecentActivityItem `json:"recent_activity"`
}
