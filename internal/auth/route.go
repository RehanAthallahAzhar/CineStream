package auth

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, handler *Handler) {
	authGroup := e.Group("")

	authGroup.POST("/register", handler.Register)
	authGroup.POST("/login", handler.Login)
}
