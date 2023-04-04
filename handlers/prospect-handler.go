package handlers

import (
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	var prospect models.Prospect

	if err := c.BodyParser(&prospect); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid body request",
		})
	}

	var company models.Company
	var client models.Client
	var projectType models.ProjectType
	var manager models.Manager

	err := db.First(&company, "id = '"+prospect.CompanyID+"'").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.Company = company

	err = db.First(&client, "id = '"+prospect.ClientID+"'").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Client not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.Client = client

	err = db.First(&projectType, "id = '"+prospect.ProjectTypeID+"'").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Project type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.ProjectType = projectType

	err = db.First(&manager, "id = '"+prospect.ManagerID+"'").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Manager not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.Manager = manager

	var uniqueNum int

	err = db.Order("created_at DESC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID).First(&prospect).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uniqueNum = 1
		}
	}
	uniqueNum = prospect.UniqueNO + 1
	prospectId := "PROSPECT/" + prospect.ProjectTypeID + "/" + prospect.CompanyID + "/" + prospect.ClientID + "/" + strconv.Itoa(uniqueNum) + "/" + strconv.Itoa(prospect.Year)
	prospectTitle := prospectId + ": " + prospect.ProspectName

	prospect.UniqueNO = uniqueNum
	prospect.ProspectID = prospectId
	prospect.ProspectTitle = prospectTitle

	if err := db.Create(&prospect).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(prospect)
}

func GetAllProspects(c *fiber.Ctx) error {
	db := database.DB.Db

	var prospects []models.Prospect

	db.Find(&prospects)

	if len(prospects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "No project found", "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "Projects found", "data": prospects,
	})
}
