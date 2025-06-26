package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetProducts(c *fiber.Ctx) error
}
