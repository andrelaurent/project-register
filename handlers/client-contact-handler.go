package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateClientContactAuditEntry(action string, clientContact models.ClientContact) error {
	db := database.DB.Db

	audit := models.ClientContactAudit{
		ClientContactID:   clientContact.ID,
		ClientID:   clientContact.ClientID,
		ContactID: clientContact.ContactID,
		Action:      action,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := db.Create(&audit).Error; err != nil {
		return err
	}

	return nil
}

func CreateClientContact(c *fiber.Ctx) error {
	db := database.DB.Db
	var clientContact models.ClientContact

	err := c.BodyParser(&clientContact)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your inputsssss", "data": err.Error()})
	}

	client := new(models.Client)
	err = db.First(&client, clientContact.ClientID).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not retrieve client info", "data": err})
	}

	contact := new(models.Contact)
	err = db.First(&contact, clientContact.ContactID).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not retrieve contact info", "data": err})
	}

	clientContact.ClientInfo.ClientName = client.ClientName
	clientContact.ContactInfo.ContactName = contact.ContactName
	clientContact.ContactInfo.BirthDate = contact.BirthDate

	err = db.Create(&clientContact).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create client contact", "data": err})
	}

	if err := CreateClientContactAuditEntry("create", clientContact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create client_contact",
		})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Client contact has been created", "data": clientContact})
}

func GetAllClientContacts(c *fiber.Ctx) error {
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

	var clientContacts []models.ClientContact

	db.Order("id ASC").Limit(limit).Offset(offset).Preload("Employments").Find(&clientContacts)

	if len(clientContacts) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client contacts not found", "data": nil})
	}

	var total int64
	db.Model(&models.ClientContact{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Client contacts Found",
		"data":        clientContacts,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetClientContactByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var ClientContact models.ClientContact

	id := c.Params("id")

	err := db.Find(&ClientContact, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client contact not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client contact retrieved", "data": ClientContact})
}

func SearchClientContact(c *fiber.Ctx) error {
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

	var ClientContacts []models.ClientContact
	var total int64

	if err := db.Model(&models.ClientContact{}).Where("ClientContact_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search ClientContacts",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("ClientContact_name ILIKE ?", "%"+searchQuery+"%").Find(&ClientContacts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search ClientContacts",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "ClientContacts Found",
		"data":        ClientContacts,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// func UpdateClientContact(c *fiber.Ctx) error {

// 	type updateClientContact struct {
// 		ClientContactCode string `json:"ClientContact_code"`
// 		ClientContactName string `json:"ClientContact_name"`
// 	}

// 	db := database.DB.Db
// 	var ClientContact models.ClientContact

// 	id := c.Params("id")

// 	db.Find(&ClientContact, "id = ?", id)

// 	if reflect.DeepEqual(ClientContact, models.ClientContact{}) {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client contact not found", "data": nil})
// 	}

// 	var updateClientContactData updateClientContact
// 	err := c.BodyParser(&updateClientContactData)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
// 	}

// 	ClientContact.ClientContactCode = updateClientContactData.ClientContactCode
// 	ClientContact.ClientContactName = updateClientContactData.ClientContactName

// 	db.Save(&ClientContact)

// 	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client contacts Found", "data": ClientContact})
// }

func DeleteClientContact(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var ClientContact models.ClientContact
	result := db.Where("id = ?", id).Delete(&ClientContact)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete Client contact", "data": result.Error})
	}

	if err := CreateClientContactAuditEntry("soft delete", ClientContact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete client_contact",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client contact not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client contact has been deleted", "data": result.RowsAffected})
}

func HardDeleteClientContact(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var ClientContact models.ClientContact
	result := db.Unscoped().Where("id = ?", id).Delete(&ClientContact)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete client contact", "data": result.Error})
	}

	if err := CreateClientContactAuditEntry("hard delete", ClientContact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete client_contact",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client contact not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client contact has been deleted from database", "data": result.RowsAffected})
}

// func RecoverClientContact(c *fiber.Ctx) error {
// 	db := database.DB.Db

// 	var request struct {
// 		ClientContactCode string `json:"ClientContact_code"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Invalid input",
// 			"data":    nil,
// 		})
// 	}

// 	var ClientContact models.ClientContact
// 	if err := db.Unscoped().Where("ClientContact_code = ? AND deleted_at IS NOT NULL", request.ClientContactCode).First(&ClientContact).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "ClientContact not found",
// 			"data":    nil,
// 		})
// 	}

// 	if err := db.Unscoped().Model(&ClientContact).Update("deleted_at", nil).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to recover ClientContact",
// 			"data":    nil,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "ClientContact recovered",
// 	})
// }
