package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateManager(c *fiber.Ctx) error {
	db := database.DB.Db
	manager := new(models.Manager)

	err := c.BodyParser(manager)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&manager).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create manager", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Manager has created", "data": manager})
}

func GetAllManagers(c *fiber.Ctx) error {
	db := database.DB.Db
	var managers []models.Client

	db.Find(&managers)

	if len(managers) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "managers not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Manager Found", "data": managers})
}
