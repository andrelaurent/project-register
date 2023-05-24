package handlers

import (
	"math"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateCity(c *fiber.Ctx) error {
	db := database.DB.Db
	city := new(models.City)

	err := c.BodyParser(city)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	city.ID = uuid.New()

	err = db.Create(&city).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create city", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "city has created", "data": city})
}

func GetAllCities(c *fiber.Ctx) error {
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

	var cities []models.City

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&cities)

	if len(cities) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Cities not found", "data": nil})
	}

	var total int64
	db.Model(&models.City{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Cities Found",
		"data":        cities,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetCityByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var city models.City

	id := c.Params("id")

	err := db.Find(&city, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "City not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "City retrieved", "data": city})
}

func SearchCity(c *fiber.Ctx) error {
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

	var cities []models.City
	var total int64

	if err := db.Model(&models.City{}).Where("city_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Citys",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("city_name ILIKE ?", "%"+searchQuery+"%").Find(&cities).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Citys",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Citys Found",
		"data":        cities,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateCity(c *fiber.Ctx) error {

	type updateCity struct {
		// CityCode string `json:"City_code"`
		CityName string `json:"City_name"`
	}

	db := database.DB.Db
	var city models.City

	id := c.Params("id")

	db.Find(&city, "id = ?", id)

	if city == (models.City{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "City not found", "data": nil})
	}

	var updateCityData updateCity
	err := c.BodyParser(&updateCityData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	// City.CityCode = updateCityData.CityCode
	city.CityName = updateCityData.CityName

	db.Save(&city)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Citys Found", "data": city})
}

func DeleteCity(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var city models.City

	result := db.Where("id = ?", id).Delete(&city)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete City", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "City not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "City has been deleted", "data": result.RowsAffected})
}

func HardDeleteCity(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var city models.City

	result := db.Unscoped().Where("id = ?", id).Delete(&city)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete city", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "City not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "City has been deleted from database", "data": result.RowsAffected})
}

func RecoverCity(c *fiber.Ctx) error {
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

	var City models.City
	if err := db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", request.ID).First(&City).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "City not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&City).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover city",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "City recovered",
	})
}
