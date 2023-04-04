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

func CreateType(c *fiber.Ctx) error {
	db := database.DB.Db
	projectType := new(models.ProjectType)

	err := c.BodyParser(projectType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Something's wrong with your input", "data": err,
		})
	}

	if projectType.ProjectTypeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Id cannot be empty",
		})
	}

	err = db.Create(&projectType).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not create company", "data": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "Company has created", "data": projectType,
	})
}

func UpdateType(c *fiber.Ctx) error {
	db := database.DB.Db

	type updateType struct {
		ProjectTypeName string `json:"name"`
	}

	id := c.Params("id")

	var projectType models.ProjectType
	db.Find(&projectType, "id = ?", id)

	if projectType == (models.ProjectType{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Type not found",
		})
	}

	var updatedTypeData updateType

	err := c.BodyParser(&updatedTypeData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Something's wrong with the input", "data": err,
		})
	}
	projectType.ProjectTypeName = updatedTypeData.ProjectTypeName

	db.Save(&projectType)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "type updated", "data": projectType,
	})
}

func DeleteType(c *fiber.Ctx) error {
	db := database.DB.Db

	var projectType models.ProjectType

	id := c.Params("id")

	db.Find(&projectType, "id = ?", id)

	if projectType == (models.ProjectType{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Type not found",
		})
	}

	err := db.Delete(&projectType, "id = ?", id).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "Failed to delete user", "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "error", "message": "Type deleted",
	})
}
