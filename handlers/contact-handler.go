package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateContact(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact

	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid body request",
			"data":    err,
		})
	}

	if err := db.Create(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create contact",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contact created",
		"data":    contact,
	})

}

func GetAllContacts(c *fiber.Ctx) error {
	db := database.DB.Db

	var contacts []models.Contact

	if err := db.Order("id ASC").Find(&contacts).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contacts not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "error",
		"message": "Contacts not found",
		"data":    contacts,
	})
}

func GetContactById(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact
	id := c.Params("id")

	if err := db.Find(&contact, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Contact not found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prospects found",
		"data":    contact,
	})
}
