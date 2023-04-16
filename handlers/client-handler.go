package handlers

import (
	"math"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateClient(c *fiber.Ctx) error {
	db := database.DB.Db
	client := new(models.Client)

	err := c.BodyParser(client)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	if client.ClientCode == "" || client.ClientName == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Client ID and name are required", "data": nil})
	}

	var existingClient models.Client
	if err := db.Where("client_code = ?", client.ClientCode).First(&existingClient).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "Client code already exists", "data": nil})
	}

	err = db.Create(&client).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create client", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Client has created", "data": client})
}

func GetAllClients(c *fiber.Ctx) error {
	db := database.DB.Db

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10

	offset := (page - 1) * limit

	var clients []models.Client

	db.Limit(limit).Offset(offset).Find(&clients)

	if len(clients) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Clients not found", "data": nil})
	}

	var total int64
	db.Model(&models.Client{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Clients Found",
		"data":        clients,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetClientByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	err := db.Find(&client, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client retrieved", "data": client})
}

func SearchClient(c *fiber.Ctx) error {
	db := database.DB.Db
	req := new(models.Client)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var clients []models.Client
	if err := db.Where("client_name LIKE ?", "%"+req.ClientName+"%").Find(&clients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search clients",
		})
	}

	return c.JSON(clients)
}

func UpdateClient(c *fiber.Ctx) error {

	type updateClient struct {
		ClientCode string `json:"client_code"`
		ClientName string `json:"client_name"`
	}

	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	db.Find(&client, "id = ?", id)

	if client == (models.Client{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "client not found", "data": nil})
	}

	var updateClientData updateClient
	err := c.BodyParser(&updateClientData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	client.ClientCode = updateClientData.ClientCode
	client.ClientName = updateClientData.ClientName

	db.Save(&client)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "clients Found", "data": client})
}

func DeleteClient(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var client models.Client
	result := db.Where("id = ?", id).Delete(&client)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete client", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client has been deleted", "data": result.RowsAffected})
}

func HardDeleteClient(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var client models.Client
	result := db.Unscoped().Where("id = ?", id).Delete(&client)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete client", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client has been deleted", "data": result.RowsAffected})
}

func RecoverClient(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	err := db.Find(&client, "id = ?", id).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	if !client.DeletedAt.Time.IsZero() {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Client is not deleted", "data": nil})
	}

	err = db.Unscoped().Model(&client).Where("id = ?", id).Update("deleted_at", nil).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to reload client", "data": err})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client recovered"})
}
