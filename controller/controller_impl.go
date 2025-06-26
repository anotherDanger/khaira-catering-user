package controller

import (
	"fmt"
	"khaira-catering-user/domain"
	"khaira-catering-user/service"
	"khaira-catering-user/web"

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
