package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProjectTypeAuditEntry(action string, projectType models.ProjectType) error {
	db := database.DB.Db

	audit := models.ProjectTypeAudit{
		ProjectTypeID:   projectType.ID,
		ProjectTypeCode: projectType.ProjectTypeCode,
		ProjectTypeName: projectType.ProjectTypeName,
		Action:     action,
		Date:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := db.Create(&audit).Error; err != nil {
		return err
	}

	return nil
}

func GetAllProjectTypes(c *fiber.Ctx) error {
	db := database.DB.Db

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var projectTypes []models.ProjectType

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&projectTypes)

	if len(projectTypes) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project types not found", "data": nil})
	}

	var total int64
	db.Model(&models.ProjectType{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Project Types Found",
		"data":        projectTypes,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetProjectTypeByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var projectType models.ProjectType

	id := c.Params("id")

	err := db.Find(&projectType, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Project type retrieved", "data": projectType})
}

func SearchProjectType(c *fiber.Ctx) error {
	db := database.DB.Db

	searchQuery := c.Query("keyword")
	if searchQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search keyword is required",
		})
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var projectTypes []models.ProjectType
	var total int64

	if err := db.Model(&models.ProjectType{}).Where("project_type_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search project types",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("project_type_name ILIKE ?", "%"+searchQuery+"%").Find(&projectTypes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search project types",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Project Types Found",
		"data":        projectTypes,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateType(c *fiber.Ctx) error {
	db := database.DB.Db
	var projectType models.ProjectType

	err := c.BodyParser(projectType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	if projectType.ProjectTypeCode == "" || projectType.ProjectTypeName == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Project type ID and name are required", "data": nil})
	}

	var existingProjectType models.ProjectType
	if err := db.Where("project_type_code = ?", projectType.ProjectTypeCode).First(&existingProjectType).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "Project type code already exists", "data": nil})
	}

	err = db.Create(&projectType).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create project type", "data": err})
	}

	if err := CreateProjectTypeAuditEntry("create", projectType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project type",
		})
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

	if err := CreateProjectTypeAuditEntry("update", projectType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create update type",
		})
	}

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

	if err := CreateProjectTypeAuditEntry("soft delete", projectType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete project type",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Project type has been deleted from database", "data": result.RowsAffected})
}

func HardDeleteProjectType(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var projectType models.ProjectType

	result := db.Unscoped().Where("id = ?", id).Delete(&projectType)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete project type", "data": result.Error})
	}

	if err := CreateProjectTypeAuditEntry("hard delete", projectType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete project type",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Project type not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Project type has been deleted", "data": result.RowsAffected})
}

func RecoverProjectType(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		ProjectTypeCode string `json:"project_type_code"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var projectType models.ProjectType
	if err := db.Unscoped().Where("project_type_code = ? AND deleted_at IS NOT NULL", request.ProjectTypeCode).First(&projectType).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Project type not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&projectType).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover project type",
			"data":    nil,
		})
	}
	
	if err := CreateProjectTypeAuditEntry("recover", projectType); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to recover project type",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Project type recovered",
	})
}
