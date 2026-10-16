package viewing

type CreateViewingRequest struct {
	PropertyID    int64  `json:"property_id"`
	PreferredDate string `json:"preferred_date"`
	Message       string `json:"message"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}
