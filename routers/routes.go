package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v1 := api.Group("/company")
	v1.Post("/create", handlers.CreateCompany)
	v1.Get("/read", handlers.GetAllCompanies)
	v1.Get("/read/:id", handlers.GetCompanyByID)
	v1.Get("/search", handlers.SearchCompany)
	v1.Put("/:id", handlers.UpdateCompany)
	v1.Delete("/:id", handlers.DeleteCompany)
	v1.Delete("/hard/:id", handlers.HardDeleteCompany)
	v1.Patch("/recover/:id", handlers.RecoverCompany)

	v2 := api.Group("/client")
	v2.Post("/create", handlers.CreateClient)
	v2.Get("/read", handlers.GetAllClients)
	v2.Get("/read/:id", handlers.GetClientByID)
	v2.Get("/search", handlers.SearchClient)
	v2.Put("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClient)
	v2.Delete("/hard/:id", handlers.HardDeleteClient)
	v2.Patch("/recover/:id", handlers.RecoverClient)

	v4 := api.Group("/type")
	v4.Post("/create", handlers.CreateType)
	v4.Get("/read", handlers.GetProjectTypes)
	v4.Get("/read/:id", handlers.GetProjectTypeByID)
	v4.Get("/search", handlers.SearchProjectType)
	v4.Put("/update/:id", handlers.UpdateProjectType)
	v4.Delete("/:id", handlers.DeleteProjectType)
	v4.Delete("/hard/:id", handlers.HardDeleteProjectType)
	v4.Patch("/recover/:id", handlers.RecoverProjectType)

	v5 := api.Group("/projects")
	v5.Post("/create", handlers.CreateProject)
	v5.Get("/read", handlers.GetAllProjects)

	v6 := api.Group("/prospect")
	v6.Post("/create", handlers.CreateProspect)
	v6.Get("/read", handlers.GetAllProspects)
	v6.Patch("/update", handlers.UpdateProspect)
	v6.Delete("/delete", handlers.DeleteProspect)
	v6.Delete("/system-delete", handlers.DeleteProspectFromSystem)
	v6.Post("/convert", handlers.ConvertToProject)
	v6.Post("/recover", handlers.RecoverProspect)
	v6.Get("/search", handlers.SearchProspects)
}
