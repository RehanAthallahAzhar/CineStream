package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"CineStream/internal/config"
	"CineStream/internal/response"
)

const (
	ContextUserIDKey = "user_id"
	ContextEmailKey  = "email"
	ContextRoleKey   = "role"
)

func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, http.StatusUnauthorized, "Authorization header is missing")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return response.Error(c, http.StatusUnauthorized, "Invalid authorization header format. Format must be Bearer <token>")
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				return response.Error(c, http.StatusUnauthorized, "Invalid or expired JWT token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return response.Error(c, http.StatusUnauthorized, "Failed to parse token claims")
			}

			userID, _ := claims["user_id"].(string)
			email, _ := claims["email"].(string)
			role, _ := claims["role"].(string)

			c.Set(ContextUserIDKey, userID)
			c.Set(ContextEmailKey, email)
			c.Set(ContextRoleKey, role)

			return next(c)
		}
	}
}

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRoleVal := c.Get(ContextRoleKey)
			userRole, ok := userRoleVal.(string)
			if !ok || userRole == "" {
				return response.Error(c, http.StatusForbidden, "User role is missing from context")
			}

			for _, role := range allowedRoles {
				if strings.EqualFold(userRole, role) {
					return next(c)
				}
			}

			return response.Error(c, http.StatusForbidden, "You do not have permission to access this resource")
		}
	}
}
