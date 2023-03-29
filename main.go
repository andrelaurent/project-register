package main

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	routers.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Listen(":8080")
}
