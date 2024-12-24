package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/pkg"
	"gorm.io/gorm"
)

func CreateApp(db *gorm.DB, emailer *pkg.Emailer, env *configs.EnvConfig) {
	app := fiber.New()

	app.Use(logger.New())

	app.Static("/images", "./images")

	InitializeRoute(app, db, emailer, env)

	app.Listen(":8080")
}
