package auth

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"CineStream/internal/response"
)

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid JSON request payload")
	}

	if err := h.validate.Struct(req); err != nil {
		valErrors := response.FormatValidationError(err)
		return response.Error(c, http.StatusBadRequest, "Validation error", valErrors)
	}

	res, err := h.service.Register(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return response.Error(c, http.StatusConflict, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusCreated, "User registered successfully", res)
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid JSON request payload")
	}

	if err := h.validate.Struct(req); err != nil {
		valErrors := response.FormatValidationError(err)

		return response.Error(c, http.StatusBadRequest, "Validation error", valErrors)
	}

	res, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return response.Error(c, http.StatusUnauthorized, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusOK, "Login successful", res)
}
