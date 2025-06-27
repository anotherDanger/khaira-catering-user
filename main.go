package main

import (
	"khaira-catering-user/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(handler controller.Controller) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://khatering.shop, http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE",
	}))

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/v1/products", handler.GetProducts)
	app.Post("/v1/login", handler.Login)
	app.Post("/v1/register", handler.Register)
	app.Post("/v1/cart/:username", handler.AddToCart)
	app.Get("/v1/cart/:username", handler.GetCart)
	app.Delete("/v1/cart/:username/:product_id", handler.DeleteCartItem)

	return app
}

func main() {
	app, cleanup, err := InitServer()
	if err != nil {
		panic(err)
	}

	defer cleanup()

	app.Listen(":8083")
}
