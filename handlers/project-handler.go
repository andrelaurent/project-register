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

func getProjectUniqueNumber(db *gorm.DB, projectTypeID uint, year int, companyID uint, clientID uint) (uint, error) {
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
	err = db.Unscoped().Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", projectTypeID, year, companyID, clientID).Order("unique_no ASC").Find(&projects).Error
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

func CreateProject(c *fiber.Ctx) error {
	db := database.DB.Db

	var project models.Project

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid body request",
		})
	}

	var company models.Company
	err := db.First(&company, "id = ?", project.CompanyID).Error
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
	project.Company = company

	var client models.Client
	err = db.First(&client, "id = ?", project.ClientID).Error
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
	project.Client = client

	var projectType models.ProjectType
	err = db.First(&projectType, "id = ?", project.ProjectTypeID).Error
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
	project.ProjectType = projectType

	if project.ProspectID != "" {
		var prospect models.Prospect
		err = db.First(&prospect, "prospect_id = ?", project.ProspectID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Prospect not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create project",
			})
		}
		project.Prospect = prospect
	}

	uniqueNumber, err := getProjectUniqueNumber(db, project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	numString := fmt.Sprintf("%02d", uniqueNumber)
	projectId := project.ProjectType.ProjectTypeCode + "/" + project.Company.CompanyCode + "/" + project.Client.ClientCode + "/" + numString + strconv.Itoa(project.Year)
	projectTitle := fmt.Sprintf("%s: %s", projectId, project.ProjectName)

	project.UniqueNO = int(uniqueNumber)
	project.ProjectID = projectId
	project.ProjectTitle = projectTitle
	project.IsDeleted = false

	if err := db.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func GetAllProjects(c *fiber.Ctx) error {
	db := database.DB.Db

	var projects []models.Project
	if err := db.Preload("Company").Preload("ProjectType").Preload("Client").Preload("Prospect", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Company").Preload("ProjectType").Preload("Client")
	}).Find(&projects).Error; err != nil {
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
	db.Model(&projects).Count(&totalCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Projects found",
		"size":    totalCount,
		"data":    projects,
	})
}

func GetProjectById(c *fiber.Ctx) error {
	db := database.DB.Db

	id := c.Params("id")

	var project models.Project
	if err := db.Preload("Company").Preload("ProjectType").Preload("Client").Preload("Prospect", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Company").Preload("ProjectType").Preload("Client")
	}).Find(&project, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "project not found",
			"data":    "null",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "project found",
		"data": project,
	})
}

func UpdateProject(c *fiber.Ctx) error {
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

	id := input["ID"]
	result := db.Preload("ProjectType").Preload("Company").Preload("Client").First(&project, "project_id = ?", id)
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
		if err := db.First(&projectType, "project_type_code = ?", val.(string)).Error; err != nil {
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
	if val, ok := input["project_name"]; ok && val.(string) != "" {
		project.ProjectName = val.(string)
	}
	if val, ok := input["no"]; ok && val.(int) != 0 {
		project.UniqueNO = val.(int)
	}
	if val, ok := input["year"]; ok && val.(int) != 0 {
		isPresent = true
		project.Year = val.(int)
	}
	if val, ok := input["manager"]; ok && val.(string) != "" {
		project.Pic = val.(string)
	}
	if val, ok := input["status"]; ok && val.(string) != "" {
		project.ProjectStatus = val.(string)
	}
	if val, ok := input["amount"]; ok && val.(float64) != 0 {
		project.ProjectAmount = val.(float64)
	}
	if val, ok := input["company_id"]; ok && val.(string) != "" {
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
				"message": "Failed to update project",
				"data":    nil,
			})
		}
		project.CompanyID = company.ID
		project.Company = company
	}
	if val, ok := input["client_id"]; ok && val.(string) != "" {
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
		uniqueNumber, err := getProjectUniqueNumber(db, project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create project",
			})
		}

		numString := fmt.Sprintf("%02d", uniqueNumber)
		projectId := project.ProjectType.ProjectTypeCode + "/" + project.Company.CompanyCode + "/" + project.Client.ClientCode + "/" + numString + "/" + strconv.Itoa(project.Year)
		projectTitle := fmt.Sprintf("%s: %s", projectId, project.ProjectName)

		project.UniqueNO = int(uniqueNumber)
		project.ProjectID = projectId
		project.ProjectTitle = projectTitle
	} else {
		project.ProjectTitle = fmt.Sprintf("%s: %s", project.ProjectID, project.ProjectName)
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
		"message": "Project updated",
		"data":    nil,
	})
}

func DeleteProject(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"ID"`
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
	result := db.Find(&project, "project_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Project not found",
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

func HardDeleteProject(c *fiber.Ctx) error {
	db := database.DB.Db

	type DeleteRequest struct {
		ID string `json:"ID"`
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
	result := db.Unscoped().Find(&project, "project_id = ?", id.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Project not found",
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

func RecoverProject(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		ProjectID string `json:"project_id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var project models.Project
	if err := db.Unscoped().Where("project_id = ? AND is_deleted = true", request.ProjectID).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Project not found",
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
		"message": "Project recovered",
	})
}

func SearchProjects(c *fiber.Ctx) error {
	db := database.DB.Db
	searchQuery := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var projects []models.Project

	if searchQuery != "" {
		db.Preload("Company").Preload("ProjectType").Preload("Client").Preload("Prospect", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Company").Preload("ProjectType").Preload("Client")
		}).Where("LOWER(project_name) LIKE ?", fmt.Sprintf("%%%s%%", strings.ToLower(searchQuery))).Offset((page - 1) * limit).Limit(limit).Find(&projects)
	} else {
		db.Preload("Company").Preload("ProjectType").Preload("Client").Preload("Prospect", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Company").Preload("ProjectType").Preload("Client")
		}).Offset((page - 1) * limit).Limit(limit).Find(&projects)
	}

	if len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No prospect found",
			"data":    nil,
		})
	}

	totalCount := len(projects)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Projects found",
		"size":    totalCount,
		"data":    projects,
	})

}

func FilterAllProjects(c *fiber.Ctx) error {
	db := database.DB.Db

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	companyID, _ := strconv.Atoi(c.Query("company", "0"))
	projectTypeID, _ := strconv.Atoi(c.Query("type", "0"))
	clientID, _ := strconv.Atoi(c.Query("client", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))

	query := db.Model(&models.Project{}).Preload("Company").Preload("ProjectType").Preload("Client").Preload("Prospect", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Company").Preload("ProjectType").Preload("Client")
	})

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
			"message": "No prospect found",
			"data":    nil,
		})
	}

	var totalCount int64
	query.Model(&models.Project{}).Count(&totalCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Projects found",
		"size":    totalCount,
		"data":    projects,
	})
}
