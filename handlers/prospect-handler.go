package handlers

import (
	"errors"
	"fmt"
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
			"status":  "error",
			"message": "Invalid body request",
		})
	}

	var company models.Company
	if err := db.First(&company, "id = '"+prospect.CompanyID+"'").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.Company = company

	var client models.Client
	if err := db.First(&client, "id = ?", prospect.ClientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Client not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.Client = client

	var projectType models.ProjectType
	if err := db.First(&projectType, "id = ?", prospect.ProjectTypeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Project type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	prospect.ProjectType = projectType

	var uniqueNum int
	var findProspect models.Prospect
	if err := db.Order("unique_no DESC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID).First(&findProspect).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			uniqueNum = 1
		}
	}

	uniqueNum = findProspect.UniqueNO + 1
	numString := fmt.Sprintf("%02d", uniqueNum)
	prospectId := "PROSPECT/" + prospect.ProjectTypeID + "/" + prospect.CompanyID + "/" + prospect.ClientID + "/" + numString + "/" + strconv.Itoa(prospect.Year)
	prospectTitle := fmt.Sprintf("%s: %s", prospectId, prospect.ProspectName)

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
	db.Preload("Company").Preload("ProspectManager").Preload("ProjectType").Preload("Client").Find(&prospects)

	if len(prospects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No prospect found",
			"data":    nil,
		})
	}

	var totalCount int64
	db.Model(&prospects).Count(&totalCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prospects found",
		"size":    totalCount,
		"data":    prospects,
	})
}

func UpdateProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	var prospect models.Prospect
	var input map[string]interface{}
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	id := input["ID"]
	result := db.Preload("ProjectType").Preload("Company").Preload("Client").First(&prospect, "prospect_id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Prospect not found",
			"data":    nil,
		})
	}

	if val, ok := input["type_id"]; ok && val.(string) != "" {
		var projectType models.ProjectType
		if err := db.First(&projectType, "id = ?", val.(string)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Type not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update prospect",
				"data":    nil,
			})
		}
		prospect.ProjectTypeID = projectType.ProjectTypeID
		prospect.ProjectType = projectType
	}
	if val, ok := input["name"]; ok && val.(string) != "" {
		prospect.ProspectName = val.(string)
	}
	if val, ok := input["no"]; ok && val.(int) != 0 {
		prospect.UniqueNO = val.(int)
	}
	if val, ok := input["year"]; ok && val.(int) != 0 {
		prospect.Year = val.(int)
	}
	if val, ok := input["manager"]; ok && val.(string) != "" {
		prospect.Pic = val.(string)
	}
	if val, ok := input["status"]; ok && val.(string) != "" {
		prospect.ProspectStatus = val.(string)
	}
	if val, ok := input["amount"]; ok && val.(float64) != 0 {
		prospect.ProspectAmount = val.(float64)
	}
	if val, ok := input["company_id"]; ok && val.(string) != "" {
		var company models.Company
		if err := db.First(&company, "id = ?", val.(string)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Company not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update prospect",
				"data":    nil,
			})
		}
		prospect.CompanyID = company.CompanyID
		prospect.Company = company
	}
	if val, ok := input["client_id"]; ok && val.(string) != "" {
		var client models.Client
		if err := db.First(&client, "id = ?", val.(string)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Client not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update prospect",
				"data":    nil,
			})
		}
		prospect.ClientID = client.ClientID
		prospect.Client = client
	}
	if val, ok := input["jira"]; ok {
		prospect.Jira = val.(bool)
	}
	if val, ok := input["clockify"]; ok {
		prospect.Clockify = val.(bool)
	}
	if val, ok := input["pcs"]; ok {
		prospect.Pcs = val.(bool)
	}
	if val, ok := input["pms"]; ok {
		prospect.Pms = val.(bool)
	}

	var uniqueNum int
	var findProspect models.Prospect

	if err := db.Order("unique_no DESC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID).First(&findProspect).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			uniqueNum = 1
		}
	}

	uniqueNum = findProspect.UniqueNO + 1
	prospect.UniqueNO = uniqueNum
	numString := fmt.Sprintf("%02d", uniqueNum)
	prospect.ProspectID = "PROSPECT/" + prospect.ProjectTypeID + "/" + prospect.CompanyID + "/" + prospect.ClientID + "/" + numString + "/" + strconv.Itoa(prospect.Year)
	prospect.ProspectTitle = fmt.Sprintf("%s: %s", prospect.ProspectID, prospect.ProspectName)

	if err := db.Save(&prospect).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update prospect",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(prospect)
}

// UPDATE PROSPECT LOOP DATA

// if err := db.Save(&prospect).Error; err != nil {
// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 		"status":  "error",
// 		"message": "Failed to update prospect",
// 		"data":    nil,
// 	})
// }

// var uniqueNum int
// var findProspects []models.Prospect

// if err := db.Order("unique_no ASC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID).Find(&findProspects).Error; err != nil {
// 	if db.RowsAffected == 0 {
// 		uniqueNum = 1
// 	} else {
// 		for i := range findProspects {
// 			if findProspects[i].UniqueNO != i+2 {
// 				uniqueNum = i + 2
// 				break
// 			}
// 		}
// 	}
// }

// prospect.UniqueNO = uniqueNum
// numString := fmt.Sprintf("%02d", uniqueNum)
// prospect.ProspectID = "PROSPECT/" + prospect.ProjectTypeID + "/" + prospect.CompanyID + "/" + prospect.ClientID + "/" + numString + "/" + strconv.Itoa(prospect.Year)
// prospect.ProspectTitle = fmt.Sprintf("%s: %s", prospect.ProspectID, prospect.ProspectName)
