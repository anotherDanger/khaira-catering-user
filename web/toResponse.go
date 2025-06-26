package web

import (
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse[T any](c *fiber.Ctx, code int, status string, data T) error {
	return c.Status(code).JSON(&Response[T]{
		Code:   code,
		Status: status,
		Data:   data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, status string, message string) error {
	return c.Status(code).JSON(&Response[any]{
		Code:   code,
		Status: status,
		Data:   message,
	})
}
