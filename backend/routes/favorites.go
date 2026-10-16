package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/internal/favorites"
	"smart-house-hunt/internal/properties"
	"smart-house-hunt/middleware"
)

func SetupFavoriteRoutes(api *gin.RouterGroup, db *sql.DB, cfg *config.Config) {
	favRepo := favorites.NewRepository(db)
	propRepo := properties.NewRepository(db)
	handler := favorites.NewHandler(favRepo, propRepo)

	favs := api.Group("/favorites", middleware.RequireAuth(cfg.JWTSecret))
	{
		favs.GET("", handler.List)
		favs.GET("/ids", handler.FavoriteIDs)
		favs.GET("/check/:id", handler.Check)
		favs.POST("", handler.Add)
		favs.POST("/:id", handler.Add)
		favs.DELETE("/:id", handler.Remove)
	}
}
