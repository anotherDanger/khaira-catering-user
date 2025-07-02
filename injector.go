//go:build wireinject
// +build wireinject

package main

import (
	"khaira-catering-user/controller"
	"khaira-catering-user/helper"
	"khaira-catering-user/repository"
	"khaira-catering-user/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

var ServerSet = wire.NewSet(
	helper.NewDb,
	repository.NewRepositoryImpl,
	service.NewServiceImpl,
	controller.NewControllerImpl,
	NewServer,
	helper.NewElasticClient,
	helper.NewValidator,
)

func InitServer() (*fiber.App, func(), error) {
	wire.Build(ServerSet)

	return nil, nil, nil
}
