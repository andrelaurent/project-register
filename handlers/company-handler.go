package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	company := new(models.Company)

	err := c.BodyParser(company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Something's wrong with your input", "data": err,
		})
	}

	err = db.Create(&company).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not create company", "data": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "Company has created", "data": company,
	})
}
