package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rubenkristian/backend/commons"
)

func CreateApp(appConfig commons.AppConfig) {
	app := fiber.New()

	app.Use(logger.New())

	app.Static("/images", "./images")

	InitializeRoute(app, &appConfig)

	app.Listen(":8080")
}
