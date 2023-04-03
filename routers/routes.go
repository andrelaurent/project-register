package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/company")
	v1.Get("/read", handlers.GetCompanies)
	v1.Post("/create", handlers.CreateCompany)

	v2 := api.Group("/client")
	v2.Get("/read", handlers.GetAllClients)
	v2.Post("/create", handlers.CreateClient)

	v3 := api.Group("/manager")
	v3.Get("/read", handlers.GetAllManagers)
	v2.Post("/create", handlers.CreateManager)
}
