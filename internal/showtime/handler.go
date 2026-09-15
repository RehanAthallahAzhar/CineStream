package showtime

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

func (h *Handler) Create(c echo.Context) error {
	var req CreateShowtimeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid JSON payload")
	}

	if err := h.validate.Struct(req); err != nil {
		valErrors := response.FormatValidationError(err)

		return response.Error(c, http.StatusBadRequest, "Validation error", valErrors)
	}

	res, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidTimeRange) || errors.Is(err, ErrInvalidUUIDFormat) {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, ErrScheduleConflict) {
			return response.Error(c, http.StatusConflict, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusCreated, "Showtime created successfully", res)
}

func (h *Handler) GetByID(c echo.Context) error {
	id := c.Param("id")

	res, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvalidUUIDFormat) {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, ErrShowtimeNotFound) {
			return response.Error(c, http.StatusNotFound, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusOK, "Showtime retrieved successfully", res)
}

func (h *Handler) List(c echo.Context) error {
	var query ListShowtimesQuery
	if err := c.Bind(&query); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid query parameters")
	}

	res, err := h.service.List(c.Request().Context(), query)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusOK, "Showtimes retrieved successfully", res)
}

func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")

	var req UpdateShowtimeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid JSON payload")
	}

	if err := h.validate.Struct(req); err != nil {
		valErrors := response.FormatValidationError(err)
		return response.Error(c, http.StatusBadRequest, "Validation error", valErrors)
	}

	res, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrInvalidUUIDFormat) || errors.Is(err, ErrInvalidTimeRange) {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, ErrShowtimeNotFound) {
			return response.Error(c, http.StatusNotFound, err.Error())
		}

		if errors.Is(err, ErrScheduleConflict) {
			return response.Error(c, http.StatusConflict, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusOK, "Showtime updated successfully", res)
}

func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")

	err := h.service.Delete(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, ErrInvalidUUIDFormat) {
			return response.Error(c, http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, ErrShowtimeNotFound) {
			return response.Error(c, http.StatusNotFound, err.Error())
		}

		return response.Error(c, http.StatusInternalServerError, "Internal server error")
	}

	return response.Success(c, http.StatusOK, "Showtime deleted successfully", nil)
}
