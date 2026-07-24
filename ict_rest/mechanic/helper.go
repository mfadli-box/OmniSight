package mechanic

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log = zerolog.New(output).With().Timestamp().Logger()
}

type AppError struct {
	Code    string
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func ValidationError(msg string) *AppError {
	return &AppError{Code: "VALIDATION_ERROR", Status: http.StatusBadRequest, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: "NOT_FOUND", Status: http.StatusNotFound, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: "CONFLICT", Status: http.StatusConflict, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Status: http.StatusUnauthorized, Message: msg}
}

func Forbidden(msg string) *AppError {
	return &AppError{Code: "FORBIDDEN", Status: http.StatusForbidden, Message: msg}
}

func ExternalServiceError(msg string, err error) *AppError {
	return &AppError{Code: "EXTERNAL_SERVICE_ERROR", Status: http.StatusBadGateway, Message: msg, Err: err}
}

func InternalError(msg string, err error) *AppError {
	return &AppError{Code: "INTERNAL_ERROR", Status: http.StatusInternalServerError, Message: msg, Err: err}
}

func Error(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		event := log.Warn().
			Str("code", appErr.Code).
			Str("error", appErr.Message).
			Int("status", appErr.Status).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("request_id", c.GetString("request_id")).
			Str("user_id", c.GetString("userId"))
		if appErr.Err != nil {
			event = event.Err(appErr.Err)
		}
		event.Msg("app error")

		c.JSON(appErr.Status, gin.H{
			"code":       appErr.Code,
			"error":      appErr.Message,
			"request_id": c.GetString("request_id"),
		})
		return
	}

	log.Error().
		Err(err).
		Str("method", c.Request.Method).
		Str("path", c.Request.URL.Path).
		Str("request_id", c.GetString("request_id")).
		Str("user_id", c.GetString("userId")).
		Msg("unhandled error")

	c.JSON(http.StatusInternalServerError, gin.H{
		"code":       "INTERNAL_ERROR",
		"error":      "An unexpected error occurred",
		"request_id": c.GetString("request_id"),
	})
}
