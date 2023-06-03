package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateContact(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact

	if err := c.BodyParser(&contact); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid body request",
			"data":    err.Error(),
		})
	}

	// if contact.Gender != "P" || contact.Gender != "L" || contact.Gender != "Other" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Gender must be 'P' or 'L'",
	// 	})
	// }

	if err := db.Create(&contact).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create contact",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contact created",
		"data":    contact,
	})

}

func GetLatestContact(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact

	if err := db.Order("id DESC").First(&contact).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contact not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": contact,
	})
}

func GetAllContacts(c *fiber.Ctx) error {
	db := database.DB.Db

	var contacts []models.Contact

	if err := db.Order("id ASC").Preload("Locations", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City").Preload("Province")
	}).Find(&contacts).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contacts not found",
		})
	}

	size := len(contacts)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "error",
		"message": "Contacts not found",
		"size":    size,
		"data":    contacts,
	})
}

func GetContactById(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact
	id := c.Params("id")

	if err := db.Preload("Locations", func(db *gorm.DB) *gorm.DB {
		return db.Preload("City").Preload("Province")
	}).Find(&contact, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Contact not found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contacts found",
		"data":    contact,
	})
}

func UpdateContact(c *fiber.Ctx) error {
	db := database.DB.Db

	contactID := c.Params("id")

	var contact models.Contact
	if err := db.First(&contact, contactID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Contact not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve contact",
			"error":   err.Error(),
		})
	}

	var patchData map[string]interface{}
	if err := c.BodyParser(&patchData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse patch data",
			"error":   err.Error(),
		})
	}

	if err := applyPatchData(&contact, patchData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to apply patch data",
			"error":   err.Error(),
		})
	}

	if err := db.Save(&contact).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update contact",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Contact updated successfully",
		"data":    contact,
	})
}

func applyPatchData(contact *models.Contact, patchData map[string]interface{}) error {
	for key, value := range patchData {
		switch key {
		case "contact_name":
			if str, ok := value.(string); ok {
				contact.ContactName = str
			}
		case "contact_alias":
			if str, ok := value.(string); ok {
				contact.ContactAlias = str
			}
		case "gender":
			if str, ok := value.(string); ok {
				contact.Gender = str
			}
		case "contact_emails":
			if emails, ok := value.([]interface{}); ok {
				contact.Emails = make([]string, len(emails))
				for i, email := range emails {
					if str, ok := email.(string); ok {
						contact.Emails[i] = str
					}
				}
			}
		case "contact_phones":
			if phones, ok := value.([]interface{}); ok {
				contact.Phones = make([]string, len(phones))
				for i, phone := range phones {
					if str, ok := phone.(string); ok {
						contact.Phones[i] = str
					}
				}
			}
		case "contact_social_presence":
			if socialPresence, ok := value.(map[string]interface{}); ok {
				if linkedin, ok := socialPresence["linkedin"].(string); ok {
					contact.ContactSocialPresence.Linkedin = linkedin
				}
				if facebook, ok := socialPresence["facebook"].(string); ok {
					contact.ContactSocialPresence.Facebook = facebook
				}
				if twitter, ok := socialPresence["twitter"].(string); ok {
					contact.ContactSocialPresence.Twitter = twitter
				}
				if github, ok := socialPresence["github"].(string); ok {
					contact.ContactSocialPresence.Github = github
				}
				if other, ok := socialPresence["other"].([]interface{}); ok {
					contact.ContactSocialPresence.Other = make([]string, len(other))
					for i, url := range other {
						if str, ok := url.(string); ok {
							contact.ContactSocialPresence.Other[i] = str
						}
					}
				}
			}
		case "birth_date":
			if str, ok := value.(string); ok {
				contact.BirthDate = str
			}
		case "religion":
			if str, ok := value.(string); ok {
				contact.Religion = str
			}
		case "interests":
			if interests, ok := value.([]interface{}); ok {
				contact.Interests = make([]string, len(interests))
				for i, interest := range interests {
					if str, ok := interest.(string); ok {
						contact.Interests[i] = str
					}
				}
			}
		case "skills":
			if skills, ok := value.([]interface{}); ok {
				contact.Skills = make([]string, len(skills))
				for i, skill := range skills {
					if str, ok := skill.(string); ok {
						contact.Skills[i] = str
					}
				}
			}
		case "educations":
			if educations, ok := value.([]interface{}); ok {
				contact.Educations = make([]string, len(educations))
				for i, education := range educations {
					if str, ok := education.(string); ok {
						contact.Educations[i] = str
					}
				}
			}
		case "notes":
			if str, ok := value.(string); ok {
				contact.Notes = str
			}
		default:
			return fmt.Errorf("unknown field: %s", key)
		}
	}
	return nil
}

func SoftDeleteContact(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact
	var locations models.Locations
	id := c.Params("id")

	if err := db.Find(&contact, "id = ?", id).Update("deleted_at", time.Now()).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete contact",
			"data":    err.Error(),
		})
	}

	if err := db.Model(&locations).Where("contact_id = ?", contact.ID).Update("deleted_at", time.Now()).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete associated locations",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contact deleted",
		"data":    nil,
	})
}

func HardDeleteContact(c *fiber.Ctx) error {
	db := database.DB.Db

	var contact models.Contact
	var locations models.Locations
	id := c.Params("id")

	if err := db.Find(&contact, "id = ?", id).Delete(&contact).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Contact not found",
			"data":    err.Error(),
		})
	}

	if err := db.Model(&locations).Where("contact_id = ?", contact.ID).Delete(&locations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete associated locations",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Contact deleted",
		"data":    nil,
	})
}
