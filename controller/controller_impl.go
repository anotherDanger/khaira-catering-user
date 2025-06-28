package controller

import (
	"fmt"
	"khaira-catering-user/domain"
	"khaira-catering-user/service"
	"khaira-catering-user/web"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ControllerImpl struct {
	svc service.Service
}

func NewControllerImpl(svc service.Service) Controller {
	return &ControllerImpl{
		svc: svc,
	}
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {
	result, err := ctrl.svc.GetProducts(c.Context())
	if err != nil {
		fmt.Println(err)
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Not OK", "Error")
	}

	return web.SuccessResponse[[]*domain.Products](c, fiber.StatusOK, "OK", result)
}

func (ctrl *ControllerImpl) Login(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}
	result, err := ctrl.svc.Login(c.Context(), user.Username, user.Password)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[*domain.User](c, fiber.StatusOK, "OK", result)
}

func (ctrl *ControllerImpl) Register(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	result, err := ctrl.svc.Register(c.Context(), &user)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[*domain.User](c, fiber.StatusCreated, "OK", result)
}

func (ctrl *ControllerImpl) AddToCart(c *fiber.Ctx) error {
	var product domain.Products
	username := c.Params("username")
	quantity := c.Params("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}
	err = c.BodyParser(&product)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	err = ctrl.svc.AddToCart(c.Context(), username, &product, quantityInt)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusCreated, "OK", "Success")
}

func (ctrl *ControllerImpl) GetCart(c *fiber.Ctx) error {
	username := c.Params("username")
	cart, err := ctrl.svc.GetCart(c.Context(), username)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[[]*domain.CartItem](c, fiber.StatusOK, "OK", cart)
}

func (ctrl *ControllerImpl) DeleteCartItem(c *fiber.Ctx) error {
	username := c.Params("username")
	productID := c.Params("product_id")
	err := ctrl.svc.DeleteCartItem(c.Context(), username, productID)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "OK", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "OK", "Success")
}

func (ctrl *ControllerImpl) DeleteCartItemByQuantity(c *fiber.Ctx) error {
	username := c.Params("username")
	productId := c.Params("product_id")
	quantity := c.Params("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	err = ctrl.svc.DeleteCartItemByQuantity(c.Context(), username, productId, quantityInt)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "OK", "Success")
}

func (ctrl *ControllerImpl) CreateOrder(c *fiber.Ctx) error {
	var order domain.Checkout
	err := c.BodyParser(&order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	err = ctrl.svc.CreateOrder(c.Context(), &order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "Error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "OK", "Success")
}
