package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
}

func SuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessResponseWithPagination(c *fiber.Ctx, message string, data interface{}, pagination *Pagination) error {
	return c.Status(http.StatusOK).JSON(APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, errors interface{}) error {
	return c.Status(statusCode).JSON(APIResponse{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, http.StatusNotFound, message, nil)
}

func BadRequestResponse(c *fiber.Ctx, message string, errors interface{}) error {
	return ErrorResponse(c, http.StatusBadRequest, message, errors)
}

func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, http.StatusForbidden, message, nil)
}

func InternalErrorResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, http.StatusInternalServerError, message, nil)
}

func ConflictResponse(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, http.StatusConflict, message, nil)
}
