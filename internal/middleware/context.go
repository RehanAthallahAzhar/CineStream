package middleware

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	ErrUnauthenticated = errors.New("user is unauthenticated")
	ErrInvalidUserID   = errors.New("user_id in context is invalid uuid")
)

func GetUserIDFromContext(c echo.Context) (uuid.UUID, error) {
	val := c.Get(ContextUserIDKey)
	userIDStr, ok := val.(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, ErrUnauthenticated
	}

	id, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, ErrInvalidUserID
	}

	return id, nil
}

func GetUserRoleFromContext(c echo.Context) string {
	val := c.Get(ContextRoleKey)
	role, ok := val.(string)
	if !ok {
		return ""
	}
	return role
}

func GetUserEmailFromContext(c echo.Context) string {
	val := c.Get(ContextEmailKey)
	email, ok := val.(string)
	if !ok {
		return ""
	}
	return email
}
