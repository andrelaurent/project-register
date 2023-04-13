package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func GetProjectTypes(c *fiber.Ctx) error {
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

func GetProjectTypeByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var projectType models.ProjectType

	id := c.Params("id")

	err := db.Find(&projectType, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "projectType not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "projectType retrieved", "data": projectType})
}

func SearchProjectType(c *fiber.Ctx) error {
	db := database.DB.Db
	req := new(models.ProjectType)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var projectTypes []models.ProjectType
	if err := db.Where("projectType_name LIKE ?", "%"+req.ProjectTypeName+"%").Find(&projectTypes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search projectTypes",
		})
	}

	return c.JSON(projectTypes)
}

func CreateType(c *fiber.Ctx) error {
	db := database.DB.Db
	projectType := new(models.ProjectType)

	err := c.BodyParser(projectType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	if projectType.ProjectTypeCode == "" || projectType.ProjectTypeName == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "project type ID and name are required", "data": nil})
	}

	var existingProjectType models.ProjectType
	if err := db.Where("project_type_code = ?", projectType.ProjectTypeCode).First(&existingProjectType).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "project type code already exists", "data": nil})
	}

	err = db.Create(&projectType).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create project type", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "project type has created", "data": projectType})
}

func UpdateProjectType(c *fiber.Ctx) error {

	type updateProjectType struct {
		ProjectTypeCode string `json:"project_type_code"`
		ProjectTypeName string `json:"project_type_name"`
	}

	db := database.DB.Db
	var projectType models.ProjectType

	id := c.Params("id")

	db.Find(&projectType, "id = ?", id)

	if projectType == (models.ProjectType{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	var updateProjectTypeData updateProjectType
	err := c.BodyParser(&updateProjectTypeData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	projectType.ProjectTypeCode = updateProjectTypeData.ProjectTypeCode
	projectType.ProjectTypeName = updateProjectTypeData.ProjectTypeName

	db.Save(&projectType)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "project type Found", "data": projectType})
}

func DeleteProjectType(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var projectType models.ProjectType

	result := db.Where("id = ?", id).Delete(&projectType)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete project type", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Project type has been deleted", "data": result.RowsAffected})
}

func HardDeleteProjectType(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var projectType models.ProjectType

	result := db.Unscoped().Where("id = ?", id).Delete(&projectType)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete project type", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Project type has been deleted", "data": result.RowsAffected})
}

func RecoverProjectType(c *fiber.Ctx) error {
	db := database.DB.Db
	var projectType models.ProjectType

	id := c.Params("id")

	err := db.Find(&projectType, "id = ?", id).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "ProjectType not found", "data": nil})
	}

	if !projectType.DeletedAt.Time.IsZero() {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "ProjectType is not deleted", "data": nil})
	}

	err = db.Unscoped().Model(&projectType).Where("id = ?", id).Update("deleted_at", nil).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to reload ProjectType", "data": err})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "ProjectType recovered"})
}
