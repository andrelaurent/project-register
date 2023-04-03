package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProject(c *fiber.Ctx) error {
	db := database.DB.Db
	project := new(models.Project)

	err := c.BodyParser(project)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&project).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create project", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Project has created", "data": project})
}

func GetAllProjects(c *fiber.Ctx) error {
	db := database.DB.Db
	var projects []models.Client

	db.Find(&projects)

	if len(projects) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "projects not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Client Found", "data": projects})
}
