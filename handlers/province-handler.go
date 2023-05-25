package handlers

import (
	"math"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateProvince(c *fiber.Ctx) error {
	db := database.DB.Db
	province := new(models.Province)

	err := c.BodyParser(province)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&province).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create Province", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Province has created", "data": province})
}

func GetAllProvinces(c *fiber.Ctx) error {
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

	var provinces []models.Province

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&provinces)

	if len(provinces) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Provinces not found", "data": nil})
	}

	var total int64
	db.Model(&models.Province{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Provinces Found",
		"data":        provinces,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetProvinceByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var province models.Province

	id := c.Params("id")

	err := db.Find(&province, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Province not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Province retrieved", "data": province})
}

func SearchProvince(c *fiber.Ctx) error {
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

	var Provinces []models.Province
	var total int64

	if err := db.Model(&models.Province{}).Where("province_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Provinces",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("province_name ILIKE ?", "%"+searchQuery+"%").Find(&Provinces).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Provinces",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Provinces Found",
		"data":        Provinces,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateProvince(c *fiber.Ctx) error {

	type updateProvince struct {
		// ProvinceCode string `json:"Province_code"`
		provinceName string `json:"province_name"`
	}

	db := database.DB.Db
	var Province models.Province

	id := c.Params("id")

	db.Find(&Province, "id = ?", id)

	if Province == (models.Province{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Province not found", "data": nil})
	}

	var updateProvinceData updateProvince
	err := c.BodyParser(&updateProvinceData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	// Province.ProvinceCode = updateProvinceData.ProvinceCode
	Province.ProvinceName = updateProvinceData.provinceName

	db.Save(&Province)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Provinces Found", "data": Province})
}

func DeleteProvince(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var Province models.Province

	result := db.Where("id = ?", id).Delete(&Province)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Province", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Province not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Province has been deleted", "data": result.RowsAffected})
}

func HardDeleteProvince(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var Province models.Province

	result := db.Unscoped().Where("id = ?", id).Delete(&Province)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Province", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Province not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Province has been deleted from database", "data": result.RowsAffected})
}

func RecoverProvince(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		ID string `json:"ID"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var Province models.Province
	if err := db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", request.ID).First(&Province).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Province not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&Province).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover Province",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Province recovered",
	})
}
