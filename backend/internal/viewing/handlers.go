package viewing

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) Create(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	var req CreateViewingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	if req.PropertyID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "property_id is required"})
		return
	}

	date := strings.TrimSpace(req.PreferredDate)
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "preferred_date is required"})
		return
	}

	prop, err := h.propertiesRepo.FindByID(req.PropertyID)
	if errors.Is(err, sql.ErrNoRows) || (prop != nil && prop.Status != "approved") {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not verify property"})
		return
	}

	if prop.LandlordID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you cannot request a viewing for your own property"})
		return
	}

	hasPending, err := h.repo.HasPendingRequest(userID, req.PropertyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not check existing requests"})
		return
	}
	if hasPending {
		c.JSON(http.StatusConflict, gin.H{"message": "a pending viewing request already exists for this property"})
		return
	}

	viewingReq := &ViewingRequest{
		UserID:        userID,
		PropertyID:    req.PropertyID,
		PreferredDate: date,
		Message:       strings.TrimSpace(req.Message),
		Status:        StatusPending,
	}

	if err := h.repo.Create(viewingReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create viewing request"})
		return
	}

	c.JSON(http.StatusCreated, viewingReq)
}

func (h *Handler) ListMine(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	list, err := h.repo.ListByTenant(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load viewing requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"viewing_requests": list})
}

func (h *Handler) ListLandlord(c *gin.Context) {
	userID := c.GetInt64(middleware.ContextUserID)

	list, err := h.repo.ListByLandlord(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load viewing requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"viewing_requests": list})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid viewing request id"})
		return
	}

	detail, err := h.repo.FindByID(id)
	if errors.Is(err, ErrViewingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "viewing request not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load viewing request"})
		return
	}

	userID := c.GetInt64(middleware.ContextUserID)
	role := c.GetString(middleware.ContextRole)

	if userID != detail.UserID && userID != detail.LandlordID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"message": "you do not have permission to view this request"})
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid viewing request id"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != StatusAccepted && status != StatusRejected && status != StatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid status. Allowed: accepted, rejected, cancelled"})
		return
	}

	detail, err := h.repo.FindByID(id)
	if errors.Is(err, ErrViewingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "viewing request not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load viewing request"})
		return
	}

	userID := c.GetInt64(middleware.ContextUserID)
	role := c.GetString(middleware.ContextRole)

	if status == StatusAccepted || status == StatusRejected {
		if userID != detail.LandlordID && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "only the landlord can accept or reject viewing requests"})
			return
		}
	} else if status == StatusCancelled {
		if userID != detail.UserID && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "only the tenant can cancel this request"})
			return
		}
	}

	if err := h.repo.UpdateStatus(id, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update viewing request"})
		return
	}

	detail.Status = status
	c.JSON(http.StatusOK, detail)
}

func (h *Handler) Accept(c *gin.Context) {
	h.transitionStatus(c, StatusAccepted)
}

func (h *Handler) Reject(c *gin.Context) {
	h.transitionStatus(c, StatusRejected)
}

func (h *Handler) Cancel(c *gin.Context) {
	h.transitionStatus(c, StatusCancelled)
}

func (h *Handler) transitionStatus(c *gin.Context, newStatus string) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid viewing request id"})
		return
	}

	detail, err := h.repo.FindByID(id)
	if errors.Is(err, ErrViewingNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "viewing request not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load viewing request"})
		return
	}

	userID := c.GetInt64(middleware.ContextUserID)
	role := c.GetString(middleware.ContextRole)

	if newStatus == StatusAccepted || newStatus == StatusRejected {
		if userID != detail.LandlordID && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "only the landlord can accept or reject viewing requests"})
			return
		}
	} else if newStatus == StatusCancelled {
		if userID != detail.UserID && role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "only the tenant can cancel this request"})
			return
		}
	}

	if err := h.repo.UpdateStatus(id, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update viewing request"})
		return
	}

	detail.Status = newStatus
	c.JSON(http.StatusOK, detail)
}
