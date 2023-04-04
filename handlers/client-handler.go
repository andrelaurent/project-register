package handlers

import (
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
