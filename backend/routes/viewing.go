package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/internal/properties"
	"smart-house-hunt/internal/viewing"
	"smart-house-hunt/middleware"
)

func SetupViewingRoutes(api *gin.RouterGroup, db *sql.DB, cfg *config.Config) {
	viewingRepo := viewing.NewRepository(db)
	propRepo := properties.NewRepository(db)
	handler := viewing.NewHandler(viewingRepo, propRepo)

	viewings := api.Group("/viewings", middleware.RequireAuth(cfg.JWTSecret))
	{
		viewings.POST("", handler.Create)
		viewings.GET("/mine", handler.ListMine)
		viewings.GET("/landlord", handler.ListLandlord)
		viewings.GET("/:id", handler.GetByID)
		viewings.PATCH("/:id/status", handler.UpdateStatus)
		viewings.POST("/:id/accept", handler.Accept)
		viewings.POST("/:id/reject", handler.Reject)
		viewings.POST("/:id/cancel", handler.Cancel)
	}
}
