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
	v1.Put("/:id", handlers.UpdateCompany)
	v1.Delete("/:id", handlers.DeleteCompany)
	v1.Delete("/hard/:id", handlers.HardDeleteCompany)

	v2 := api.Group("/client")
	v2.Post("/create", handlers.CreateClient)
	v2.Get("/read", handlers.GetAllClients)
	v2.Get("/read/:id", handlers.GetClientByID)
	v2.Get("/search", handlers.SearchClient)
	v2.Put("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClient)
	v2.Delete("/hard/:id", handlers.HardDeleteClient)

	v3 := api.Group("/type")
	v3.Post("/create", handlers.CreateType)
	v3.Get("/read", handlers.GetProjectTypes)
	v3.Get("/read/:id", handlers.GetProjectTypeByID)
	v3.Put("/update/:id", handlers.UpdateProjectType)
	v3.Delete("/:id", handlers.DeleteProjectType)
	v3.Delete("/hard/:id", handlers.HardDeleteProjectType)

	v4 := api.Group("/projects")
	v4.Post("/create", handlers.CreateProject)
	v4.Get("/read", handlers.GetAllProjects)
	v4.Patch("/update", handlers.UpdateProject)
	v4.Delete("/delete/soft", handlers.DeleteProject)
	v4.Delete("/delete/hard", handlers.DeleteProjectFromSystem)
	v4.Post("/recover", handlers.RecoverProject)
	v4.Get("/search", handlers.SearchProjects)

	v5 := api.Group("/prospect")
	v5.Post("/create", handlers.CreateProspect)
	v5.Get("/read", handlers.GetAllProspects)
	v5.Patch("/update", handlers.UpdateProspect)
	v5.Delete("/delete", handlers.DeleteProspect)
	v5.Delete("/system-delete", handlers.DeleteProspectFromSystem)
	v5.Post("/convert", handlers.ConvertToProject)
	v5.Post("/recover", handlers.RecoverProspect)
	v5.Get("/search", handlers.SearchProspects)
}
