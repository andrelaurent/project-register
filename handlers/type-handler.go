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
	project_type := new(models.ProjectType)

	err := c.BodyParser(project_type)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	if project_type.ProjectTypeCode == "" || project_type.ProjectTypeName == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "project type ID and name are required", "data": nil})
	}

	var existingProjectType models.ProjectType
	if err := db.Where("project_type_code = ?", project_type.ProjectTypeCode).First(&existingProjectType).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "project type code already exists", "data": nil})
	}

	err = db.Create(&project_type).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create project type", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "project type has created", "data": project_type})
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
