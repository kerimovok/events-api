package utils

import (
	"events-api/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, appErr *errors.AppError) error {
	return c.Status(appErr.Code).JSON(Response{
		Success: false,
		Message: appErr.Message,
		Error:   appErr.Error(),
	})
}
