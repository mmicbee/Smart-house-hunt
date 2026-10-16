package admin

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

func (h *Handler) ListUsers(c *gin.Context) {
	status := c.Query("status")
	role := c.Query("role")
	search := c.Query("search")

	users, err := h.repo.ListUsers(status, role, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	currentAdminID := c.GetInt64(middleware.ContextUserID)

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	if id == currentAdminID {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you cannot suspend your own admin account"})
		return
	}

	var req UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	if req.Status != "active" && req.Status != "suspended" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be either 'active' or 'suspended'"})
		return
	}

	if err := h.repo.UpdateUserStatus(id, req.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user status"})
		return
	}

	updatedUser, err := h.repo.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "user status updated",
			"user_id": id,
			"status":  req.Status,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user status updated",
		"user":    updatedUser,
	})
}

func (h *Handler) ListProperties(c *gin.Context) {
	status := c.Query("status")
	propType := c.Query("property_type")

	properties, err := h.repo.ListProperties(status, propType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load properties"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"properties": properties})
}

func (h *Handler) GetProperty(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	prop, err := h.repo.GetPropertyByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"property": prop})
}

func (h *Handler) UpdatePropertyStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	var req UpdatePropertyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Status = strings.ToLower(strings.TrimSpace(req.Status))
	switch req.Status {
	case "pending", "approved", "rejected", "suspended":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be one of: pending, approved, rejected, suspended"})
		return
	}

	prop, err := h.repo.UpdatePropertyStatus(id, req.Status)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update property status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "property status updated",
		"property": prop,
	})
}

func (h *Handler) ApproveProperty(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	prop, err := h.repo.UpdatePropertyStatus(id, "approved")
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not approve property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "property approved",
		"property": prop,
	})
}

func (h *Handler) RejectProperty(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	prop, err := h.repo.UpdatePropertyStatus(id, "rejected")
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not reject property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "property rejected",
		"property": prop,
	})
}

func (h *Handler) DeleteProperty(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid property id"})
		return
	}

	if err := h.repo.DeleteProperty(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"message": "property not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete property"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "property deleted"})
}

func (h *Handler) GetReports(c *gin.Context) {
	reports, err := h.repo.GetReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load system reports"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}
