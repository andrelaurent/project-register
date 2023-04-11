package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func getNextUniqueNumber(db *gorm.DB, projectTypeID uint, year int, companyID uint, clientID uint) (uint, error) {
	var maxUniqueNumber, minUniqueNumber uint
	var vacantNumbers []uint

	var prospect models.Prospect

	err := db.Unscoped().Model(&prospect).Select("MAX(unique_no), MIN(unique_no)").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", projectTypeID, year, companyID, clientID).Row().Scan(&maxUniqueNumber, &minUniqueNumber)
	if maxUniqueNumber == 0 && minUniqueNumber == 0 {
		return 1, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	var prospects []models.Prospect
	err = db.Unscoped().Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", projectTypeID, year, companyID, clientID).Order("unique_no ASC").Find(&prospects).Error
	if err != nil {
		return 0, err
	}
	for i := int(minUniqueNumber); i <= int(maxUniqueNumber); i++ {
		found := false
		for _, prospect := range prospects {
			if int(i) == prospect.UniqueNO {
				found = true
				break
			}
		}
		if !found {
			vacantNumbers = append(vacantNumbers, uint(i))
		}
	}

	if len(vacantNumbers) > 0 {
		return vacantNumbers[0], nil
	}
	return maxUniqueNumber + 1, nil
}

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
	if err := db.First(&company, "id = ?", prospect.ClientID).Error; err != nil {
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
			"message": "Failed to create prospect",
		})
	}
	prospect.ProjectType = projectType

	uniqueNumber, err := getNextUniqueNumber(db, prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create prospect",
		})
	}

	numString := fmt.Sprintf("%02d", uniqueNumber)
	prospectId := "PROSPECT/" + prospect.ProjectType.ProjectTypeID + "/" + prospect.Company.CompanyCode + "/" + prospect.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(prospect.Year)
	prospectTitle := fmt.Sprintf("%s: %s", prospectId, prospect.ProspectName)

	prospect.UniqueNO = int(uniqueNumber)
	prospect.ProspectID = prospectId
	prospect.ProspectTitle = prospectTitle
	prospect.IsDeleted = false

	if err := db.Create(&prospect).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create prospect",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(prospect)
}

func GetAllProspects(c *fiber.Ctx) error {
	db := database.DB.Db

	var prospects []models.Prospect
	db.Preload("Company").Preload("ProjectType").Preload("Client").Find(&prospects)

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

	var isPresent bool
	isPresent = false

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
		isPresent = true
		var projectType models.ProjectType
		if err := db.First(&projectType, "project_type_id = ?", val.(string)).Error; err != nil {
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
		prospect.ProjectTypeID = projectType.ID
		prospect.ProjectType = projectType
	}
	if val, ok := input["prospect_name"]; ok && val.(string) != "" {
		prospect.ProspectName = val.(string)
	}
	if val, ok := input["no"]; ok && val.(int) != 0 {
		prospect.UniqueNO = val.(int)
	}
	if val, ok := input["year"]; ok && val.(int) != 0 {
		isPresent = true
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
	if val, ok := input["company_code"]; ok && val.(string) != "" {
		isPresent = true
		var company models.Company
		if err := db.First(&company, "company_code = ?", val.(string)).Error; err != nil {
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
		prospect.CompanyID = company.ID
		prospect.Company = company
	}
	if val, ok := input["client_code"]; ok && val.(string) != "" {
		isPresent = true
		var client models.Client
		if err := db.First(&client, "client_code = ?", val.(string)).Error; err != nil {
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
		prospect.ClientID = client.ID
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
	if isPresent {
		uniqueNumber, err := getNextUniqueNumber(db, prospect.ProjectTypeID, prospect.Year, prospect.CompanyID, prospect.ClientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create prospect",
			})
		}

		numString := fmt.Sprintf("%02d", uniqueNumber)
		prospectId := "PROSPECT/" + prospect.ProjectType.ProjectTypeID + "/" + prospect.Company.CompanyCode + "/" + prospect.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(prospect.Year)
		prospectTitle := fmt.Sprintf("%s: %s", prospectId, prospect.ProspectName)

		prospect.UniqueNO = int(uniqueNumber)
		prospect.ProspectID = prospectId
		prospect.ProspectTitle = prospectTitle
	} else {
		prospect.ProspectTitle = fmt.Sprintf("%s: %s", prospect.ProspectID, prospect.ProspectName)
	}

	if err := db.Save(&prospect).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update prospect",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(prospect)
}

func DeleteProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"ID"`
	}
	var prospect models.Prospect
	var id DeleteRequest

	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}
	result := db.Find(&prospect, "prospect_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Prospect not found",
			"data":    nil,
		})
	}

	if err := db.Model(&prospect).Updates(map[string]interface{}{"deleted_at": time.Now(), "is_deleted": true}).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete prospect",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "prospect deleted",
	})
}

func DeleteProspectFromSystem(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"ID"`
	}
	var prospect models.Prospect
	var id DeleteRequest

	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}
	result := db.Unscoped().Find(&prospect, "prospect_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Prospect not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Delete(&prospect).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete prospect",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "prospect deleted",
	})
}

func ConvertToProject(c *fiber.Ctx) error {
	db := database.DB.Db

	type RequestId struct {
		ProspectID string `json:"prospect_id"`
	}

	var id RequestId
	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var prospect models.Prospect
	if err := db.Preload("ProjectType").Preload("Company").Preload("Client").Where("prospect_id = ?", id.ProspectID).First(&prospect).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Prospect not found",
			"data":    nil,
		})
	}

	projectId := prospect.ProspectID
	projectId = strings.ReplaceAll(projectId, "PROSPECT/", "")

	project := models.Project{
		ProjectID:     projectId,
		ProjectTypeID: prospect.ProjectTypeID,
		ProjectType:   prospect.ProjectType,
		ProjectName:   prospect.ProspectName,
		UniqueNO:      prospect.UniqueNO,
		Year:          prospect.Year,
		Pic:           prospect.Pic,
		ProjectStatus: prospect.ProspectStatus,
		ProjectTitle:  prospect.ProspectTitle,
		ProjectAmount: prospect.ProspectAmount,
		CompanyID:     prospect.CompanyID,
		Company:       prospect.Company,
		ClientID:      prospect.ClientID,
		Client:        prospect.Client,
		ProspectID:    prospect.ProspectID,
		Prospect:      prospect,
		IsDeleted:     false,
		Jira:          prospect.Jira,
		Clockify:      prospect.Clockify,
		Pcs:           prospect.Pcs,
		Pms:           prospect.Pms,
	}

	if err := db.First(&project, "project_id = ?", projectId).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project already existed",
		})
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&project).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to convert prospect",
				"data":    nil,
			})
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to convert prospect",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}
