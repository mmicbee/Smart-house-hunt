package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/internal/auth"
	"smart-house-hunt/internal/users"
	"smart-house-hunt/middleware"
)

func Setup(router *gin.Engine, db *sql.DB, cfg *config.Config) {
	router.Use(middleware.CORS())
	router.Static("/uploads", "storage")

	usersRepo := users.NewRepository(db)
	authHandler := auth.NewAuthHandler(usersRepo, cfg.JWTSecret)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Smart House Hunt API is running",
		})
	})

	api := router.Group("/api")
	{
		SetupPropertyRoutes(api, db, cfg)
		SetupFavoriteRoutes(api, db, cfg)
		SetupViewingRoutes(api, db, cfg)

		api.GET("/me", middleware.RequireAuth(cfg.JWTSecret), func(c *gin.Context) {
			userID := c.GetInt64(middleware.ContextUserID)
			user, err := usersRepo.FindByID(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "could not load profile"})
				return
			}
			if user.Status == "suspended" {
				c.JSON(http.StatusForbidden, gin.H{"message": "your account has been suspended"})
				return
			}

			roleName, _ := usersRepo.RoleByID(user.RoleID)
			c.JSON(http.StatusOK, gin.H{
				"id":     user.ID,
				"name":   user.Name,
				"email":  user.Email,
				"role":   roleName,
				"status": user.Status,
			})
		})

		authAPI := api.Group("/auth")
		{
			authAPI.POST("/register", authHandler.Register)
			authAPI.POST("/login", authHandler.Login)
			authAPI.POST("/logout", middleware.RequireAuth(cfg.JWTSecret), authHandler.Logout)
		}

		SetupAdminRoutes(api, db, cfg)
	}
}