package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"smart-house-hunt/config"
	"smart-house-hunt/internal/admin"
	"smart-house-hunt/middleware"
)

func SetupAdminRoutes(api *gin.RouterGroup, db *sql.DB, cfg *config.Config) {
	repo := admin.NewRepository(db)
	handler := admin.NewHandler(repo)

	adminAPI := api.Group("/admin",
		middleware.RequireAuth(cfg.JWTSecret),
		middleware.RequireRole("admin"),
	)
	{
		adminAPI.GET("/reports", handler.GetReports)
		adminAPI.GET("/users", handler.ListUsers)
		adminAPI.PUT("/users/:id/status", handler.UpdateUserStatus)
		adminAPI.PATCH("/users/:id/status", handler.UpdateUserStatus)

		adminAPI.GET("/properties", handler.ListProperties)
		adminAPI.GET("/properties/:id", handler.GetProperty)
		adminAPI.PUT("/properties/:id/status", handler.UpdatePropertyStatus)
		adminAPI.PATCH("/properties/:id/status", handler.UpdatePropertyStatus)
		adminAPI.PUT("/properties/:id/approve", handler.ApproveProperty)
		adminAPI.POST("/properties/:id/approve", handler.ApproveProperty)
		adminAPI.PUT("/properties/:id/reject", handler.RejectProperty)
		adminAPI.POST("/properties/:id/reject", handler.RejectProperty)
		adminAPI.DELETE("/properties/:id", handler.DeleteProperty)
	}
}
