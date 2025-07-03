package controller

import (
	"khaira-catering-user/domain"
	"khaira-catering-user/service"
	"khaira-catering-user/web"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ControllerImpl struct {
	svc      service.Service
	validate *validator.Validate
}

func NewControllerImpl(svc service.Service, validate *validator.Validate) Controller {
	return &ControllerImpl{
		svc:      svc,
		validate: validate,
	}
}

func (ctrl *ControllerImpl) GetProducts(c *fiber.Ctx) error {
	result, err := ctrl.svc.GetProducts(c.Context())
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "cannot get products")
	}

	return web.SuccessResponse[[]*domain.Products](c, fiber.StatusOK, "ok", result)
}

func (ctrl *ControllerImpl) Login(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	if err := ctrl.validate.Struct(user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "validation failed: invalid input")
	}

	result, err := ctrl.svc.Login(c.Context(), user.Username, user.Password)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[*domain.User](c, fiber.StatusOK, "ok", result)
}

func (ctrl *ControllerImpl) Register(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	if err := ctrl.validate.Struct(user); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "validation failed: invalid input")
	}

	result, err := ctrl.svc.Register(c.Context(), &user)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[*domain.User](c, fiber.StatusCreated, "ok", result)
}

func (ctrl *ControllerImpl) AddToCart(c *fiber.Ctx) error {
	var product domain.Products
	username := c.Params("username")
	quantity := c.Params("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}
	err = c.BodyParser(&product)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	err = ctrl.svc.AddToCart(c.Context(), username, &product, quantityInt)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusCreated, "ok", "Success")
}

func (ctrl *ControllerImpl) GetCart(c *fiber.Ctx) error {
	username := c.Params("username")
	cart, err := ctrl.svc.GetCart(c.Context(), username)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[[]*domain.CartItem](c, fiber.StatusOK, "ok", cart)
}

func (ctrl *ControllerImpl) DeleteCartItem(c *fiber.Ctx) error {
	username := c.Params("username")
	productID := c.Params("product_id")
	err := ctrl.svc.DeleteCartItem(c.Context(), username, productID)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "ok", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "ok", "Success")
}

func (ctrl *ControllerImpl) DeleteCartItemByQuantity(c *fiber.Ctx) error {
	username := c.Params("username")
	productId := c.Params("product_id")
	quantity := c.Params("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	err = ctrl.svc.DeleteCartItemByQuantity(c.Context(), username, productId, quantityInt)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "ok", "Success")
}

func (ctrl *ControllerImpl) CreateOrder(c *fiber.Ctx) error {
	var order domain.Checkout
	err := c.BodyParser(&order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	if err := ctrl.validate.Struct(order); err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "validation failed: invalid input")
	}

	err = ctrl.svc.CreateOrder(c.Context(), &order)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", err.Error())
	}

	return web.SuccessResponse[any](c, fiber.StatusNoContent, "ok", "Success")
}

func (ctrl *ControllerImpl) GetOrderHistory(c *fiber.Ctx) error {
	username := c.Params("username")
	res, err := ctrl.svc.GetOrderHistory(c.Context(), username)
	if err != nil {
		return web.ErrorResponse(c, fiber.StatusBadRequest, "error", "cannot get orders")
	}

	return web.SuccessResponse[[]*domain.Checkout](c, fiber.StatusOK, "ok", res)
}
