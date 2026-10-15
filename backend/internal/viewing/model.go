package viewing

import "time"

const (
	StatusPending   = "pending"
	StatusAccepted  = "accepted"
	StatusRejected  = "rejected"
	StatusCancelled = "cancelled"
)

type ViewingRequest struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	PropertyID    int64     `json:"property_id"`
	PreferredDate string    `json:"preferred_date"`
	Message       string    `json:"message"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ViewingRequestDetail struct {
	ViewingRequest
	PropertyTitle    string `json:"property_title"`
	PropertyLocation string `json:"property_location"`
	PropertyPrice    int    `json:"property_price"`
	LandlordID       int64  `json:"landlord_id"`
	LandlordName     string `json:"landlord_name,omitempty"`
	LandlordEmail    string `json:"landlord_email,omitempty"`
	TenantName       string `json:"tenant_name,omitempty"`
	TenantEmail      string `json:"tenant_email,omitempty"`
	PropertyImage    string `json:"property_image,omitempty"`
}
