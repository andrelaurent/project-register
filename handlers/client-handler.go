package handlers

import (
	"time"

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
	var clients []models.Client

	db.Find(&clients)

	if len(clients) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Clients not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client Found", "data": clients})
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
		Username string `json:"name"`
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
	client.ClientName = updateClientData.Username

	db.Save(&client)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "clients Found", "data": client})
}

func DeleteClientByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	db.Find(&client, "id = ?", id)

	if client == (models.Client{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	err := db.Delete(&client, "id = ?", id).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete client", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client deleted"})
}

func SoftDeleteClientByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	db.Find(&client, "id = ?", id)

	if client == (models.Client{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	err := db.Model(&client).Update("deleted_at", time.Now()).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete client", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client soft deleted"})
}
