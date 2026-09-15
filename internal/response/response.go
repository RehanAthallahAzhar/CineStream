package response

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func Success(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string, errs ...interface{}) error {
	var errDetails interface{}
	if len(errs) == 1 {
		errDetails = errs[0]
	} else if len(errs) > 1 {
		errDetails = errs
	}

	return c.JSON(code, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

func FormatValidationError(err error) []ValidationErrorDetail {
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		details := make([]ValidationErrorDetail, 0, len(valErrs))
		for _, e := range valErrs {
			details = append(details, ValidationErrorDetail{
				Field:   e.Field(),
				Tag:     e.Tag(),
				Message: e.Error(),
			})
		}
		return details
	}
	return nil
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	message := "Internal server error"

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		if msg, ok := he.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(code)
		}
	}

	_ = Error(c, code, message)
}
