package showtime

import (
	"github.com/labstack/echo/v4"

	"CineStream/internal/auth"
	"CineStream/internal/config"
	appMiddleware "CineStream/internal/middleware"
)

func RegisterRoutes(e *echo.Echo, handler *Handler, cfg *config.Config) {
	// Public
	e.GET("/showtimes", handler.List)
	e.GET("/showtimes/:id", handler.GetByID)

	// Protected
	adminGroup := e.Group("/showtimes")
	adminGroup.Use(appMiddleware.JWTMiddleware(cfg))
	adminGroup.Use(appMiddleware.RequireRole(auth.RoleAdmin))

	adminGroup.POST("", handler.Create)
	adminGroup.PUT("/:id", handler.Update)
	adminGroup.DELETE("/:id", handler.Delete)
}
