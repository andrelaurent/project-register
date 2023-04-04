package handlers

import (
	"log"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateProject(c *fiber.Ctx) error {
	db := database.DB.Db

	var project models.Project

	if err := c.BodyParser(&project); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", "message": "Invalid body request",
		})
	}

	var company models.Company
	var client models.Client
	var projectType models.ProjectType
	var manager models.Manager
	var prospect models.Prospect

	log.Println(project.CompanyID)

	err := db.First(&company, "id = '"+project.CompanyID+"'").Error
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

	err = db.First(&client, "id = '"+project.ClientID+"'").Error
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

	err = db.First(&projectType, "id = '"+project.ProjectTypeID+"'").Error
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

	err = db.First(&manager, "id = '"+project.ManagerID+"'").Error
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
	project.Manager = manager

	err = db.First(&prospect, "id = "+project.ProspectID+"'").Error
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

	var uniqueNum int

	err = db.Order("created_at DESC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uniqueNum = 1
		}
	}
	uniqueNum = project.UniqueNO + 1

	projectId := project.ProjectTypeID + "/" + project.CompanyID + "/" + project.ClientID + "/" + strconv.Itoa(uniqueNum) + "/" + strconv.Itoa(project.Year)
	projectTitle := projectId + ": " + project.ProjectName

	project.UniqueNO = uniqueNum
	project.ProjectID = projectId
	project.ProjectTitle = projectTitle

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

	db.Find(&projects)

	if len(projects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "No project found", "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", "message": "Projects found", "data": projects,
	})
}
