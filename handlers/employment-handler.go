package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateEmployment(c *fiber.Ctx) error {
	db := database.DB.Db
	employment := new(models.Employment)

	err := c.BodyParser(employment)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&employment).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create employment", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Employment has created", "data": employment})
}

func GetAllEmployments(c *fiber.Ctx) error {
	db := database.DB.Db
	var employments []models.Employment

	err := db.Find(&employments).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch employments"})
	}

	return c.JSON(employments)
}


func GetEmploymentsByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var employments []models.Employment
	clientContactID := c.Params("id")

	err := db.Where("client_contact_id = ?", clientContactID).Find(&employments).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch employments"})
	}

	return c.JSON(employments)
}

func GetEmploymentsByContactID(c *fiber.Ctx) error {
	db := database.DB.Db
	var employments []models.Employment
	contactID := c.Params("id")

	err := db.Joins("JOIN client_contacts ON employments.client_contact_id = client_contacts.id").
		Joins("JOIN contacts ON client_contacts.contact_id = contacts.id").
		Where("contacts.id = ?", contactID).
		Find(&employments).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch employments"})
	}

	return c.JSON(employments)
}


func DeleteEmployment(c *fiber.Ctx) error {
	db := database.DB.Db
	employmentID := c.Params("employmentID")

	var employment models.Employment
	if err := db.First(&employment, employmentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employment not found"})
	}

	if err := db.Delete(&employment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete employment"})
	}

	return c.JSON(fiber.Map{"message": "Employment deleted"})
}
