package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/internal/properties"
	"smart-house-hunt/middleware"
)

func SetupPropertyRoutes(api *gin.RouterGroup, db *sql.DB, cfg *config.Config) {
	repo := properties.NewRepository(db)
	handler := properties.NewHandler(repo)

	landlordAPI := api.Group("/properties",
		middleware.RequireAuth(cfg.JWTSecret),
		middleware.RequireRole("landlord"),
	)
	{
		landlordAPI.POST("", handler.Create)
		landlordAPI.PUT("/:id", handler.Update)
		landlordAPI.DELETE("/:id", handler.Delete)
	}
}