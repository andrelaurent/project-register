package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func GetTypes(c *fiber.Ctx) error {
	db := database.DB.Db
	var projectType []models.ProjectType

	db.Find(&projectType)

	if len(projectType) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "no type found", "data": "nil",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "sucess", "message": "Types Found", "data": projectType,
	})
}
