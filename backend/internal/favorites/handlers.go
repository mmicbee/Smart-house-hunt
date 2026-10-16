package favorites

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/properties"
	"smart-house-hunt/middleware"
)

type Handler struct {
	repo           *Repository
	propertiesRepo *properties.Repository
}

func NewHandler(repo *Repository, propertiesRepo *properties.Repository) *Handler {
	return &Handler{
		repo:           repo,
		propertiesRepo: propertiesRepo,
	}
}

func (h *Handler) Add(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	var propertyID int64
	var err error

	if paramID := c.Param("id"); paramID != "" {
		propertyID, err = strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
			return
		}
	} else {
		var req FavoriteRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.PropertyID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "valid property_id is required"})
			return
		}
		propertyID = req.PropertyID
	}

	prop, err := h.propertiesRepo.FindByID(propertyID)
	if errors.Is(err, sql.ErrNoRows) || (prop != nil && prop.Status != "approved") {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not verify property"})
		return
	}

	isFav, err := h.repo.IsFavorite(userID, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not check favorite status"})
		return
	}
	if isFav {
		c.JSON(http.StatusOK, gin.H{"message": "property is already in favorites"})
		return
	}

	fav, err := h.repo.Add(userID, propertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save favorite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "property saved to favorites",
		"favorite": fav,
	})
}

func (h *Handler) Remove(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	if err := h.repo.Remove(userID, id); err != nil {
		if errors.Is(err, ErrFavoriteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "favorite not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not remove favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "property removed from favorites"})
}

func (h *Handler) List(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	props, err := h.repo.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load favorites"})
		return
	}

	items := make([]properties.PropertyListItem, 0, len(props))
	for _, prop := range props {
		images, err := h.propertiesRepo.ListImagesByProperty(prop.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load property images"})
			return
		}
		items = append(items, properties.PropertyListItem{
			Property: prop,
			Images:   images,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"properties": items,
		"favorites":  items,
	})
}

func (h *Handler) Check(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	isFav, err := h.repo.IsFavorite(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not check favorite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFav})
}

func (h *Handler) FavoriteIDs(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	ids, err := h.repo.ListFavoritePropertyIDs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load favorite ids"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"favorite_ids": ids})
}
