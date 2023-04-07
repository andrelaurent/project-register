package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/company")
	v4 := api.Group("/type")
	v1.Post("/create", handlers.CreateCompany)
	v4.Post("/create", handlers.CreateType)
	v4.Get("/read", handlers.GetTypes)
	v4.Put("/update/:id", handlers.UpdateType)
	v1.Get("/read", handlers.GetAllCompanies)
	v1.Post("/create", handlers.CreateCompany)
	v1.Put("/:id", handlers.UpdateCompany)
	v1.Delete("/:id", handlers.DeleteCompanyByID)

	v2 := api.Group("/client")
	v2.Get("/read", handlers.GetAllClients)
	v2.Get("/read/:id", handlers.GetClientByID)
	v2.Get("/search", handlers.SearchClient)
	v2.Post("/create", handlers.CreateClient)
	v2.Put("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClientByID)

	v3 := api.Group("/manager")
	v3.Get("/read", handlers.GetAllManagers)
	v3.Post("/create", handlers.CreateManager)
	v3.Put("/:id", handlers.UpdateManager)
	v3.Delete("/:id", handlers.DeleteManagerByID)

	v5 := api.Group("/projects")
	v5.Post("/create", handlers.CreateProject)
	v5.Get("/read", handlers.GetAllProjects)

	v6 := api.Group("/prospect")
	v6.Post("/create", handlers.CreateProspect)
	v6.Get("/read", handlers.GetAllProspects)
}
