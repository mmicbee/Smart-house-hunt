package properties

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/middleware"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(c *gin.Context) {
	landlordID := c.GetInt64(middleware.ContextUserID)

	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Location = strings.TrimSpace(req.Location)
	req.Amenities = strings.TrimSpace(req.Amenities)

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "title is required"})
		return
	}
	if req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "description is required"})
		return
	}
	if req.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "location is required"})
		return
	}
	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "price must be greater than zero"})
		return
	}
	if req.Bedrooms < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bedrooms cannot be negative"})
		return
	}
	if req.Bathrooms < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bathrooms cannot be negative"})
		return
	}

	property := &Property{
		LandlordID:  landlordID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Price:       req.Price,
		Bedrooms:    req.Bedrooms,
		Bathrooms:   req.Bathrooms,
		Amenities:   req.Amenities,
		Status:      "pending",
	}

	if err := h.repo.Create(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create property"})
		return
	}

	c.JSON(http.StatusCreated, toResponse(property))
}