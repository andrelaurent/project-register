package routers

import (
	"github.com/andrelaurent/project-register/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Group("/auth", handlers.Authenticate)

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
	v2.Get("/latest", handlers.GetLatestClient)
	v2.Get("/read/:id", handlers.GetClientByID)
	v2.Get("/search", handlers.SearchClient)
	v2.Patch("/:id", handlers.UpdateClient)
	v2.Delete("/:id", handlers.DeleteClient)
	v2.Delete("/hard/:id", handlers.HardDeleteClient)

	v3 := api.Group("/type")
	v3.Post("/create", handlers.CreateType)
	v3.Get("/read", handlers.GetAllProjectTypes)
	v3.Get("/read/:id", handlers.GetProjectTypeByID)
	v3.Get("/search", handlers.SearchProjectType)
	v3.Put("/update/:id", handlers.UpdateProjectType)
	v3.Delete("/:id", handlers.DeleteProjectType)
	v3.Delete("/hard/:id", handlers.HardDeleteProjectType)
	v3.Post("/recover", handlers.RecoverProjectType)

	// v4 := api.Group("/projects")
	// v4.Post("/create", handlers.CreateProject)
	// v4.Get("/read", handlers.GetAllProjects)
	// v4.Get("/read/:id", handlers.GetProjectById)
	// v4.Patch("/update/:id", handlers.UpdateProject)
	// v4.Delete("/delete/soft", handlers.DeleteProject)
	// v4.Delete("/delete/hard", handlers.HardDeleteProject)
	// v4.Post("/recover", handlers.RecoverProject)
	// v4.Get("/search", handlers.SearchProjects)

	v5 := api.Group("/model")
	v5.Post("/prospect/create", handlers.CreateProject)
	v5.Get("/:model/read", handlers.GetAllProjects)
	v5.Get("/:model/read/:id", handlers.GetProject)
	v5.Patch("/:model/update/:id", handlers.UpdateProject)
	v5.Delete("/:model/delete/:id", handlers.DeleteProject)
	v5.Delete("/:model/hard/:id", handlers.HardDeleteProject)
	v5.Post("/prospect/convert/:id", handlers.ConvertToProject)
	v5.Post("/:model/recover/:id", handlers.RecoverProject)
	v5.Get("/project/search", handlers.SearchProjects)
	v5.Get("/project/filter", handlers.FilterAllProjects)

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
	v7.Get("/filter/:id", handlers.GetCityFitered)

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
	v9.Get("/latest", handlers.GetLatestContact)
	v9.Get("/read", handlers.GetAllContacts)
	v9.Get("/read/:id", handlers.GetContactById)
	v9.Patch("/update/:id", handlers.UpdateContact)
	v9.Delete("/delete/soft/:id", handlers.SoftDeleteContact)
	v9.Get("/locations/:id", handlers.GetLocationByContactID)

	v10 := api.Group("/clientcontact")
	v10.Post("/create", handlers.CreateClientContact)
	v10.Get("/read", handlers.GetAllClientContacts)

	v11 := api.Group("/employments")
	v11.Post("/create", handlers.CreateEmployment)
	v11.Get("/read", handlers.GetAllEmployments)
	v11.Get("/read/:id", handlers.GetEmploymentsByContactID)
	v11.Delete("/delete", handlers.DeleteEmployment)
}
