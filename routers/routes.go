package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/company")
	v1.Get("/read", handlers.GetAllCompanies)
	v1.Post("/create", handlers.CreateCompany)
	v1.Put("/:id", handlers.UpdateCompany)
	v1.Delete("/:id", handlers.DeleteCompanyByID)

	v2 := api.Group("/client")
	v2.Get("/read", handlers.GetAllClients)
	v2.Post("/create", handlers.CreateClient)
	v2.Put("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClientByID)

	v3 := api.Group("/manager")
	v3.Get("/read", handlers.GetAllManagers)
	v3.Post("/create", handlers.CreateManager)
	v3.Put("/:id", handlers.UpdateManager)
	v3.Delete("/:id", handlers.DeleteManagerByID)
}
