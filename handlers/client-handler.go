package handlers

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateClientAuditEntry(action string, client models.Client) error {
	db := database.DB.Db

	audit := models.ClientAudit{
		ClientID:   client.ID,
		ClientCode: client.ClientCode,
		ClientName: client.ClientName,
		Action:     action,
		Date:       time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := db.Create(&audit).Error; err != nil {
		return err
	}

	return nil
}

func CreateClient(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	err := c.BodyParser(&client)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	var existingClient models.Client
	if err := db.Where("client_code = ?", client.ClientCode).First(&existingClient).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "Client code already exists", "data": nil})
	}

	err = db.Create(&client).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create client", "data": err})
	}

	if err := CreateClientAuditEntry("create", client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create client",
		})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Client has created", "data": client})
}

func GetLatestClient(c *fiber.Ctx) error {
	db := database.DB.Db

	var client models.Client

	if err := db.Order("id DESC").First(&client).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contact not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": client,
	})
}

func GetAllClients(c *fiber.Ctx) error {
	db := database.DB.Db

	var clients []models.Client

	if err := db.Order("id ASC").Find(&clients).Preload("Locations", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City").Preload("Province")
	}).Find(&clients).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contacts not found",
		})
	}

	if len(clients) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Clients not found", "data": nil})
	}

	var total int64
	db.Model(&models.Client{}).Count(&total)

	response := fiber.Map{
		"status":     "success",
		"message":    "Clients Found",
		"data":       clients,
		"totalItems": total,
	}

	return c.Status(200).JSON(response)
}

func UpdateClient(c *fiber.Ctx) error {
	db := database.DB.Db

	var client models.Client
	id := c.Params("id")

	if err := db.First(&client, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Client not found",
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve client",
			"data":    err.Error(),
		})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
			"data":    err.Error(),
		})
	}

	if err := updateClientFields(&client, updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update client",
			"data":    err.Error(),
		})
	}

	if err := db.Save(&client).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to save updated client",
			"data":    err.Error(),
		})
	}

	if err := CreateClientAuditEntry("update", client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update client",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Client updated",
		"data":    client,
	})
}

func updateClientFields(client *models.Client, updateData map[string]interface{}) error {
	for key, value := range updateData {
		switch key {
		case "client_code":
			if code, ok := value.(string); ok {
				client.ClientCode = code
			}
		case "client_name":
			if name, ok := value.(string); ok {
				client.ClientName = name
			}
		case "alias":
			if alias, ok := value.(string); ok {
				client.Alias = alias
			}
		case "website":
			if website, ok := value.(string); ok {
				client.Website = website
			}
		case "client_social_presence":
			if socialPresence, ok := value.(map[string]interface{}); ok {
				updateSocialPresenceFields(&client.SocialPresence, socialPresence)
			}
		case "subsidiary":
			if subsidiary, ok := value.(map[string]interface{}); ok {
				updateSubsidiaryFields(&client.Subsidiary, subsidiary)
			}
		case "date":
			if dateStr, ok := value.(string); ok {
				client.Date = dateStr
			}
		}
	}
	return nil
}

func updateSocialPresenceFields(socialPresence *models.SocialPresence, data map[string]interface{}) {
	for key, value := range data {
		switch key {
		case "facebook":
			if fb, ok := value.(string); ok {
				socialPresence.Facebook = fb
			}
		case "twitter":
			if twitter, ok := value.(string); ok {
				socialPresence.Twitter = twitter
			}
		case "linkedin":
			if linkedin, ok := value.(string); ok {
				socialPresence.Linkedin = linkedin
			}
		case "github":
			if github, ok := value.(string); ok {
				socialPresence.Github = github
			}
		case "other":
			if other, ok := value.([]interface{}); ok {
				socialPresence.Other = make([]string, len(other))
				for i, other := range other {
					if str, ok := other.(string); ok {
						socialPresence.Other[i] = str
					}
				}
			}
		}
	}
}

func updateSubsidiaryFields(subsidiary *models.Subsidiary, data map[string]interface{}) {
	for key, value := range data {
		switch key {
		case "subsidiaries":
			if subs, ok := value.([]interface{}); ok {
				subsidiary.Subsidiaries = make([]string, len(subs))
				for i, subs := range subs {
					if str, ok := subs.(string); ok {
						subsidiary.Subsidiaries[i] = str
					}
				}
			}
		case "immidiate_parents":
			if iparents, ok := value.([]interface{}); ok {
				subsidiary.ImmidiateParents = make([]string, len(iparents))
				for i, iparents := range iparents {
					if str, ok := iparents.(string); ok {
						subsidiary.ImmidiateParents[i] = str
					}
				}
			}
		case "ultimate_parents":
			if uparents, ok := value.([]interface{}); ok {
				subsidiary.UltimateParents = make([]string, len(uparents))
				for i, uparents := range uparents {
					if str, ok := uparents.(string); ok {
						subsidiary.UltimateParents[i] = str
					}
				}
			}
		}
	}
}

func GetClientByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var client models.Client

	id := c.Params("id")

	err := db.Preload("Locations", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City").Preload("Province")
	}).Find(&client, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client retrieved", "data": client})
}

func SearchClient(c *fiber.Ctx) error {
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

	var clients []models.Client
	var total int64

	if err := db.Model(&models.Client{}).Where("client_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search clients",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("client_name ILIKE ?", "%"+searchQuery+"%").Find(&clients).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search clients",
		})
	}

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

	return c.Status(fiber.StatusOK).JSON(response)
}

func DeleteClient(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var client models.Client
	var locations models.Locations
	result := db.Where("id = ?", id).Delete(&client)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete client", "data": result.Error})
	}

	if err := db.Model(&locations).Where("client_id", client.ID).Delete(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete associated locations",
			"data":    err.Error(),
		})
	}

	if err := CreateClientAuditEntry("soft delete", client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete client",
		})
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
	var locations models.Locations
	result := db.Unscoped().Where("id = ?", id).Delete(&client)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete client", "data": result.Error})
	}

	if err := db.Model(&locations).Where("client_id", client.ID).Delete(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete associated locations",
			"data":    err.Error(),
		})
	}

	if err := CreateClientAuditEntry("hard delete", client); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete client",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Client not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Client has been deleted from database", "data": result.RowsAffected})
}

// func RecoverClient(c *fiber.Ctx) error {
// 	db := database.DB.Db

// 	var request struct {
// 		ClientCode string `json:"client_code"`
// 	}

// 	if err := c.BodyParser(&request); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Invalid input",
// 			"data":    nil,
// 		})
// 	}

// 	var client models.Client
// 	if err := db.Unscoped().Where("client_code = ? AND deleted_at IS NOT NULL", request.ClientCode).First(&client).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Client not found",
// 			"data":    nil,
// 		})
// 	}

// 	if err := db.Unscoped().Model(&client).Update("deleted_at", nil).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to recover client",
// 			"data":    nil,
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "Client recovered",
// 	})
// }
