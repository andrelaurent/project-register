package handlers

import (
	"math"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	location := new(models.Locations)

	err := c.BodyParser(location)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	err = db.Create(&location).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create Location", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Location has created", "data": location})
}

func GetAllLocations(c *fiber.Ctx) error {
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

	var locations []models.Locations

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&locations)

	if len(locations) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Locations not found", "data": nil})
	}

	var total int64
	db.Model(&models.Locations{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Locations Found",
		"data":        locations,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetLocationByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var location models.Locations

	id := c.Params("id")

	err := db.Find(&location, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Location retrieved", "data": location})
}

func SearchLocation(c *fiber.Ctx) error {
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

	var locations []models.Locations
	var total int64

	if err := db.Model(&models.Locations{}).Where("location_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Locations",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("location_name ILIKE ?", "%"+searchQuery+"%").Find(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Locations",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Locations Found",
		"data":        locations,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	var location models.Locations

	if err := c.BodyParser(location); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	id := c.Params("id")

	existingLocation := models.Locations{}
	err := db.First(&existingLocation, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found", "data": nil})
	}

	existingLocation.Address = location.Address
	existingLocation.City = location.City
	existingLocation.Province = location.Province
	existingLocation.PostalCode = location.PostalCode
	existingLocation.Country = location.Country
	existingLocation.Geo = location.Geo

	err = db.Save(&existingLocation).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update location",
		})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Locations Found", "data": location})
}

func DeleteLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var Location models.Locations

	result := db.Where("id = ?", id).Delete(&Location)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Location", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Location has been deleted", "data": result.RowsAffected})
}

func HardDeleteLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var Location models.Locations

	result := db.Unscoped().Where("id = ?", id).Delete(&Location)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Location", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Location has been deleted from database", "data": result.RowsAffected})
}

func RecoverLocation(c *fiber.Ctx) error {
	db := database.DB.Db
	locationID := c.Query("ID")

	var location models.Locations
	if err := db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", locationID).First(&location).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Location not found",
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover Location",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&location).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover Location",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Location recovered",
	})
}
