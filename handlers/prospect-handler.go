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

	var project models.Project

	err := db.Unscoped().Model(&project).Select("MAX(unique_no), MIN(unique_no)").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", projectTypeID, year, companyID, clientID).Row().Scan(&maxUniqueNumber, &minUniqueNumber)
	if maxUniqueNumber == 0 && minUniqueNumber == 0 {
		return 1, nil
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	var projects []models.Project
	err = db.Unscoped().Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", projectTypeID, year, companyID, clientID).Order("unique_no ASC").Find(&project).Error
	if err != nil {
		return 0, err
	}
	for i := int(minUniqueNumber); i <= int(maxUniqueNumber); i++ {
		found := false
		for _, project := range projects {
			if int(i) == project.UniqueNO {
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

	var project models.Project

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid body request",
		})
	}

	var company models.Company
	if err := db.First(&company, "id = ?", project.CompanyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Company not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	project.Company = company

	var client models.Client
	if err := db.First(&client, "id = ?", project.ClientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Client not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	project.Client = client

	var projectType models.ProjectType
	code := "PRP"
	if err := db.First(&projectType, "project_type_code = ?", code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Project type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}
	project.ProjectType = projectType

	uniqueNumber, err := getNextUniqueNumber(db, project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	numString := fmt.Sprintf("%02d", uniqueNumber)
	projectId := project.ProjectType.ProjectTypeCode + "/" + project.Company.CompanyCode + "/" + project.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(project.Year)
	projectTitle := fmt.Sprintf("%s: %s", projectId, project.ProspectName)

	project.UniqueNO = int(uniqueNumber)
	project.ProspectID = projectId
	project.ProspectTitle = projectTitle
	project.IsDeleted = false

	if err := db.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func GetAllProspects(c *fiber.Ctx) error {
	db := database.DB.Db

	// page, _ := strconv.Atoi(c.Query("page", "1"))
	// limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// var prospects []models.Prospect
	// if err := db.Preload("Company").Preload("ProjectType").Preload("Client").Offset((page - 1) * limit).Limit(limit).Find(&prospects).Error; err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Could not find prospects",
	// 		"data":    nil,
	// 	})
	// }

	var projects []models.Project
	if err := db.Preload("Company").Preload("ProjectType").Preload("Client").Find(&projects).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not find prospects",
			"data":    nil,
		})
	}

	if len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No project found",
			"data":    projects,
		})
	}

	var totalCount int64
	db.Model(&projects).Count(&totalCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prospects found",
		"size":    totalCount,
		"data":    projects,
	})
}

func GetProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	id := c.Params("id")

	var project models.Project
	if err := db.Preload("Company").Preload("ProjectType").Preload("Client").Find(&project, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    "null",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "projects found",
		"data":    project,
	})
}

func UpdateProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	var isPresent bool
	isPresent = false

	var project models.Project
	var input map[string]interface{}
	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	// id := input["prospect_id"]
	id := c.Params("id")
	result := db.Preload("ProjectType").Preload("Company").Preload("Client").First(&project, "id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    nil,
		})
	}

	if val, ok := input["type_id"]; ok && val.(float64) != 0 {
		isPresent = true
		var projectType models.ProjectType
		if err := db.First(&projectType, "id = ?", val.(float64)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Type not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update project",
				"data":    nil,
			})
		}
		project.ProjectTypeID = projectType.ID
		project.ProjectType = projectType
	}
	if val, ok := input["prospect_name"]; ok && val.(string) != "" {
		project.ProspectName = val.(string)
	}
	// if val, ok := input["no"]; ok && val.(int) != 0 {
	// 	prospect.UniqueNO = val.(int)
	// }
	// if val, ok := input["year"]; ok && val.(float64) != 0 {
	// 	isPresent = true
	// 	prospect.Year = val.(int)
	// }
	if val, ok := input["year"]; ok {
		isPresent = true
		if year, ok := val.(float64); ok {
			project.Year = int(year)
		}
	}
	if val, ok := input["manager"]; ok && val.(string) != "" {
		project.Pic = val.(string)
	}
	if val, ok := input["status"]; ok && val.(string) != "" {
		project.ProspectStatus = val.(string)
	}
	if val, ok := input["amount"]; ok && val.(float64) != 0 {
		project.ProspectAmount = val.(float64)
	}
	if val, ok := input["company_id"]; ok && val.(float64) != 0 {
		isPresent = true
		var company models.Company
		if err := db.First(&company, "id = ?", val.(float64)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Company not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update project",
				"data":    nil,
			})
		}
		project.CompanyID = company.ID
		project.Company = company
	}
	if val, ok := input["client_id"]; ok && val.(float64) != 0 {
		isPresent = true
		var client models.Client
		if err := db.First(&client, "id = ?", val.(float64)).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Client not found",
					"data":    nil,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update project",
				"data":    nil,
			})
		}
		project.ClientID = client.ID
		project.Client = client
	}
	if val, ok := input["jira"]; ok {
		project.Jira = val.(bool)
	}
	if val, ok := input["clockify"]; ok {
		project.Clockify = val.(bool)
	}
	if val, ok := input["pcs"]; ok {
		project.Pcs = val.(bool)
	}
	if val, ok := input["pms"]; ok {
		project.Pms = val.(bool)
	}
	if isPresent {
		uniqueNumber, err := getNextUniqueNumber(db, project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create project",
			})
		}

		numString := fmt.Sprintf("%02d", uniqueNumber)
		projectId := project.ProjectType.ProjectTypeCode + "/" + project.Company.CompanyCode + "/" + project.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(project.Year)
		projectTitle := fmt.Sprintf("%s: %s", projectId, project.ProspectName)

		project.UniqueNO = int(uniqueNumber)
		project.ProspectID = projectId
		project.ProspectTitle = projectTitle
	} else {
		project.ProspectTitle = fmt.Sprintf("%s: %s", project.ProspectID, project.ProspectName)
	}

	if err := db.Save(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update project",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "project updated",
		"data":    project,
	})
}

func DeleteProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"prospect_id"`
	}
	var project models.Project
	var id DeleteRequest

	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}
	result := db.Find(&project, "prospect_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    nil,
		})
	}

	if err := db.Model(&project).Updates(map[string]interface{}{"deleted_at": time.Now(), "is_deleted": true}).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete project",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "project deleted",
	})
}

func HardDeleteProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"prospect_id"`
	}
	var project models.Project
	var id DeleteRequest

	if err := c.BodyParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}
	result := db.Unscoped().Find(&project, "prospect_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Delete(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete project",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "project deleted",
	})
}

// func ConvertToProject(c *fiber.Ctx) error {
// 	db := database.DB.Db

// 	type RequestId struct {
// 		ProspectID string `json:"prospect_id"`
// 		TypeID     uint   `json:"type_id"`
// 	}

// 	var id RequestId
// 	if err := c.BodyParser(&id); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Invalid input",
// 			"data":    nil,
// 		})
// 	}

// 	var prospect models.Prospect
// 	if err := db.Preload("ProjectType").Preload("Company").Preload("Client").Where("prospect_id = ?", id.ProspectID).First(&prospect).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Prospect not found",
// 			"data":    nil,
// 		})
// 	}

// 	// projectId := prospect.ProspectID
// 	// projectId = strings.ReplaceAll(projectId, "PROSPECT/", "")
// 	// projectTitle := prospect.ProspectTitle
// 	// projectTitle = strings.ReplaceAll(projectTitle, "PROSPECT/", "")

// 	var projectType models.ProjectType
// 	if err := db.First(&projectType, "id = ?", id.TypeID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"status":  "error",
// 				"message": "Type not found",
// 				"data":    nil,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to update prospect",
// 			"data":    nil,
// 		})
// 	}

// 	project := models.Project{
// 		// ProjectID:     projectId,
// 		ProjectTypeID: id.TypeID,
// 		ProjectType:   projectType,
// 		ProjectName:   prospect.ProspectName,
// 		// UniqueNO:      prospect.UniqueNO,
// 		Year:          prospect.Year,
// 		Pic:           prospect.Pic,
// 		ProjectStatus: prospect.ProspectStatus,
// 		// ProjectTitle:  projectTitle,
// 		ProjectAmount: prospect.ProspectAmount,
// 		CompanyID:     prospect.CompanyID,
// 		Company:       prospect.Company,
// 		ClientID:      prospect.ClientID,
// 		Client:        prospect.Client,
// 		ProspectID:    prospect.ProspectID,
// 		Prospect:      prospect,
// 		// IsDeleted:     false,
// 		Jira:     prospect.Jira,
// 		Clockify: prospect.Clockify,
// 		Pcs:      prospect.Pcs,
// 		Pms:      prospect.Pms,
// 	}

// 	uniqueNumber, err := getProjectUniqueNumber(db, project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to create project",
// 		})
// 	}

// 	numString := fmt.Sprintf("%02d", uniqueNumber)
// 	projectId := project.ProjectType.ProjectTypeCode + "/" + project.Company.CompanyCode + "/" + project.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(project.Year)
// 	projectTitle := fmt.Sprintf("%s: %s", projectId, project.ProjectName)

// 	project.UniqueNO = int(uniqueNumber)
// 	project.ProjectID = projectId
// 	project.ProjectTitle = projectTitle
// 	project.IsDeleted = false

// 	if err := db.First(&project, "project_id = ?", projectId).Error; err == nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Project already existed",
// 		})
// 	} else if errors.Is(err, gorm.ErrRecordNotFound) {
// 		if err := db.Create(&project).Error; err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"status":  "error",
// 				"message": "Failed to convert prospect",
// 				"data":    nil,
// 			})
// 		}
// 	} else {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":  "error",
// 			"message": "Failed to convert prospect",
// 			"data":    nil,
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"status":  "success",
// 		"message": "Prospect converted",
// 		"data":    project,
// 	})
// }

func RecoverProspect(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		ProspectID string `json:"prospect_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var project models.Project
	if err := db.Unscoped().Where("prospect_id = ? AND is_deleted = true", request.ProspectID).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&project).Updates(map[string]interface{}{"deleted_at": nil, "is_deleted": false}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to revocer project",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "project recovered",
	})
}

func SearchProspects(c *fiber.Ctx) error {
	db := database.DB.Db
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var project []models.Project

	if searchQuery != "" {
		db.Preload("ProjectType").Preload("Company").Preload("Client").Where("LOWER(prospect_name) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(searchQuery))).Offset((page - 1) * limit).Limit(limit).Find(&project)
	} else {
		db.Preload("ProjectType").Preload("Company").Preload("Client").Offset((page - 1) * limit).Limit(limit).Find(&project)
	}

	if len(project) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No project found",
			"data":    nil,
		})
	}

	totalCount := len(project)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Projects found",
		"size":    totalCount,
		"data":    project,
	})

}

func FilterAllProspects(c *fiber.Ctx) error {
	db := database.DB.Db

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	companyID, _ := strconv.Atoi(c.Query("company", "0"))
	projectTypeID, _ := strconv.Atoi(c.Query("type", "0"))
	clientID, _ := strconv.Atoi(c.Query("client", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))

	query := db.Model(&models.Project{}).Preload("Company").Preload("ProjectType").Preload("Client")

	if companyID != 0 {
		query = query.Where("company_id = ?", companyID)
	}

	if projectTypeID != 0 {
		query = query.Where("project_type_id = ?", projectTypeID)
	}

	if clientID != 0 {
		query = query.Where("client_id = ?", clientID)
	}

	if year != 0 {
		query = query.Where("year = ?", year)
	}

	var projects []models.Project
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&projects).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not find projects",
			"data":    nil,
		})
	}

	if len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No project found",
			"data":    nil,
		})
	}

	var totalCount int64
	query.Model(&models.Project{}).Count(&totalCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Prospects found",
		"size":    totalCount,
		"data":    projects,
	})
}
