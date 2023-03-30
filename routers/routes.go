package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/andrelaurent/project-register/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/company")
	// v2 := api.Group("/client")
	// v3 := api.Group("/manager")
	// v4 := api.Group("/manager")
	v1.Get("/read", handlers.GetCompanies)
	v1.Post("/create", handlers.CreateCompany)
}
