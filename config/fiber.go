package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/logger"
)

func NewFiber(cfg *AppConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: LoadConfig().AppName,
	})

	app.Use(logger.New())
	app.Use(fiberrecover.New())

	return app
}