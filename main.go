package main

import (
	"khaira-catering-user/controller"
	"khaira-catering-user/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer(handler controller.Controller) *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://khatering.shop,https://khatering.netlify.app",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Get("/v1/products", handler.GetProducts)
	app.Post("/v1/login", handler.Login)
	app.Post("/v1/register", handler.Register)

	protectedRoute := app.Group("/api")
	protectedRoute.Use(middleware.MyMiddleware)

	protectedRoute.Post("/v1/cart/:username/:quantity", handler.AddToCart)
	protectedRoute.Get("/v1/cart/:username", handler.GetCart)
	protectedRoute.Delete("/v1/cart/:username/:product_id", handler.DeleteCartItem)
	protectedRoute.Delete("/v1/cart/:username/:product_id/:quantity", handler.DeleteCartItemByQuantity)
	protectedRoute.Post("/v1/checkout", handler.CreateOrder)
	protectedRoute.Get("/v1/history/:username", handler.GetOrderHistory)

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
