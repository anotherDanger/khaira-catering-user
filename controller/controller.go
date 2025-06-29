package controller

import "github.com/gofiber/fiber/v2"

type Controller interface {
	GetProducts(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	AddToCart(c *fiber.Ctx) error
	GetCart(c *fiber.Ctx) error
	DeleteCartItem(c *fiber.Ctx) error
	DeleteCartItemByQuantity(c *fiber.Ctx) error
	CreateOrder(c *fiber.Ctx) error
	GetOrderHistory(c *fiber.Ctx) error
}
