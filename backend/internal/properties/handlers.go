package properties

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
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

	if msg := validateRequest(&req); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
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

func (h *Handler) Update(c *gin.Context) {
	landlordID := c.GetInt64(middleware.ContextUserID)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	property, err := h.repo.FindByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load property"})
		return
	}
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"message": "you can only edit your own listings"})
		return
	}

	var req PropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	if msg := validateRequest(&req); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	property.Title = req.Title
	property.Description = req.Description
	property.Location = req.Location
	property.Price = req.Price
	property.Bedrooms = req.Bedrooms
	property.Bathrooms = req.Bathrooms
	property.Amenities = req.Amenities

	if err := h.repo.Update(property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update property"})
		return
	}

	c.JSON(http.StatusOK, toResponse(property))
}

func (h *Handler) Delete(c *gin.Context) {
	landlordID := c.GetInt64(middleware.ContextUserID)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	property, err := h.repo.FindByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load property"})
		return
	}
	if property.LandlordID != landlordID {
		c.JSON(http.StatusForbidden, gin.H{"message": "you can only delete your own listings"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "property deleted"})
}

func validateRequest(req *PropertyRequest) string {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Location = strings.TrimSpace(req.Location)
	req.Amenities = strings.TrimSpace(req.Amenities)

	switch {
	case req.Title == "":
		return "title is required"
	case req.Description == "":
		return "description is required"
	case req.Location == "":
		return "location is required"
	case req.Price <= 0:
		return "price must be greater than zero"
	case req.Bedrooms < 0:
		return "bedrooms cannot be negative"
	case req.Bathrooms < 0:
		return "bathrooms cannot be negative"
	}

	return ""
}