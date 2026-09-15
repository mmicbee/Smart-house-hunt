package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
)

const (
	ContextUserID = "userID"
	ContextRole   = "role"
)

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authentication required"})
			return
		}

		claims, err := auth.ParseToken(secret, strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}

	return func(c *gin.Context) {
		role := c.GetString(ContextRole)
		if !allowed[role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "you do not have permission to access this resource"})
			return
		}
		c.Next()
	}
}