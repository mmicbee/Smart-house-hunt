package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/internal/auth"
	"smart-house-hunt/internal/users"
	"smart-house-hunt/middleware"
)

func Setup(router *gin.Engine, db *sql.DB) {
	router.Use(middleware.CORS())

	usersRepo := users.NewRepository(db)
	authHandler := auth.NewAuthHandler(usersRepo)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Smart House Hunt API is running",
		})
	})

	api := router.Group("/api")
	{
		authAPI := api.Group("/auth")
		{
			authAPI.POST("/register", authHandler.Register)
		}
	}
}