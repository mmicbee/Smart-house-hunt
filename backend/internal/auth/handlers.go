package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/users"
)

type AuthHandler struct {
	users  *users.Repository
	secret string
}

func NewAuthHandler(usersRepo *users.Repository, secret string) *AuthHandler {
	return &AuthHandler{users: usersRepo, secret: secret}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "name is required"})
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "a valid email address is required"})
		return
	}
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password must be at least 6 characters"})
		return
	}
	if req.Role == "" {
		req.Role = "tenant"
	}
	switch req.Role {
	case "tenant", "landlord", "agent":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "role must be one of: tenant, landlord, agent"})
		return
	}

	if _, err := h.users.FindByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "an account with this email already exists"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not check existing account"})
		return
	}

	hashed, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not secure password"})
		return
	}

	roleID, err := h.users.RoleIDByName(req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not resolve role"})
		return
	}

	user := &users.User{
		RoleID:   roleID,
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
	}
	if err := h.users.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create account"})
		return
	}

	roleName, _ := h.users.RoleByID(user.RoleID)
	c.JSON(http.StatusCreated, UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  roleName,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
		return
	}

	user, err := h.users.FindByEmail(req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not find account"})
		return
	}

	if !CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid email or password"})
		return
	}

	roleName, _ := h.users.RoleByID(user.RoleID)

	token, err := GenerateToken(h.secret, user.ID, roleName, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create session token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  roleName,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}