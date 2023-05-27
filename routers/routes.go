package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	v0 := api.Group("/user")
	v0.Post("/create", handlers.CreateUser)
	v0.Post("/login", handlers.UserLogin)
	v0.Get("/read", handlers.GetAllUsers)

	v1 := api.Group("/company")
	v1.Post("/create", handlers.CreateCompany)
	v1.Get("/read", handlers.GetAllCompanies)
	v1.Get("/read/:id", handlers.GetCompanyByID)
	v1.Get("/search", handlers.SearchCompany)
	v1.Put("/:id", handlers.UpdateCompany)
	v1.Delete("/:id", handlers.DeleteCompany)
	v1.Delete("/hard/:id", handlers.HardDeleteCompany)
	v1.Post("/recover", handlers.RecoverCompany)

	v2 := api.Group("/client")
	v2.Post("/create", handlers.CreateClient)
	v2.Get("/read", handlers.GetAllClients)
	v2.Get("/read/:id", handlers.GetClientByID)
	v2.Get("/search", handlers.SearchClient)
	v2.Put("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClient)
	v2.Delete("/hard/:id", handlers.HardDeleteClient)
	v2.Post("/recover", handlers.RecoverClient)

	v3 := api.Group("/type")
	v3.Post("/create", handlers.CreateType)
	v3.Get("/read", handlers.GetAllProjectTypes)
	v3.Get("/read/:id", handlers.GetProjectTypeByID)
	v3.Get("/search", handlers.SearchProjectType)
	v3.Put("/update/:id", handlers.UpdateProjectType)
	v3.Delete("/:id", handlers.DeleteProjectType)
	v3.Delete("/hard/:id", handlers.HardDeleteProjectType)
	v3.Post("/recover", handlers.RecoverProjectType)

	v4 := api.Group("/projects")
	v4.Post("/create", handlers.CreateProject)
	v4.Get("/read", handlers.GetAllProjects)
	v4.Patch("/update", handlers.UpdateProject)
	v4.Delete("/delete/soft", handlers.DeleteProject)
	v4.Delete("/delete/hard", handlers.HardDeleteProject)
	v4.Post("/recover", handlers.RecoverProject)
	v4.Get("/search", handlers.SearchProjects)

	v5 := api.Group("/prospect")
	v5.Post("/create", handlers.CreateProspect)
	v5.Get("/read", handlers.GetAllProspects)
	v5.Get("/read/:id", handlers.GetProspect)
	v5.Patch("/update", handlers.UpdateProspect)
	v5.Delete("/delete", handlers.DeleteProspect)
	v5.Delete("/system-delete", handlers.HardDeleteProspect)
	v5.Post("/convert", handlers.ConvertToProject)
	v5.Post("/recover", handlers.RecoverProspect)
	v5.Get("/search", handlers.SearchProspects)
	v5.Get("/filter", handlers.FilterAllProspects)

	v6 := api.Group("/locations")
	v6.Post("/create", handlers.CreateLocation)
	v6.Get("/read", handlers.GetAllLocations)
	v6.Get("/read/:id", handlers.GetLocationByID)
	v6.Get("/search", handlers.SearchLocation)
	v6.Put("/:id", handlers.UpdateLocation)
	v6.Delete("/:id", handlers.DeleteLocation)
	v6.Delete("/hard/:id", handlers.HardDeleteLocation)
	v6.Post("/recover", handlers.RecoverLocation)

	v7 := api.Group("/city")
	v7.Post("/create", handlers.CreateCity)
	v7.Get("/read", handlers.GetAllCities)
	v7.Get("/read/:id", handlers.GetCityByID)
	v7.Get("/search", handlers.SearchCity)
	v7.Put("/:id", handlers.UpdateCity)
	v7.Delete("/:id", handlers.DeleteCity)
	v7.Delete("/hard/:id", handlers.HardDeleteCity)
	v7.Post("/recover", handlers.RecoverCity)

	v8 := api.Group("/province")
	v8.Post("/create", handlers.CreateProvince)
	v8.Get("/read", handlers.GetAllProvinces)
	v8.Get("/read/:id", handlers.GetProvinceByID)
	v8.Get("/search", handlers.SearchProvince)
	v8.Put("/:id", handlers.UpdateProvince)
	v8.Delete("/:id", handlers.DeleteProvince)
	v8.Delete("/hard/:id", handlers.HardDeleteProvince)
	v8.Post("/recover", handlers.RecoverProvince)

	v9 := api.Group("/contact")
	v9.Post("/create", handlers.CreateContact)
	v9.Get("/read", handlers.GetAllContacts)
}
