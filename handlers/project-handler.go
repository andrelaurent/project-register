package handlers

import (
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
			"message": "Invalid body request",
		})
	}

	var company models.Company
	var client models.Client
	var projectType models.ProjectType
	var manager models.Manager

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

	var uniqueNum int

	err = db.Order("created_at DESC").Where("project_type_id = ? AND year = ? AND company_id = ? AND client_id = ?", project.ProjectTypeID, project.Year, project.CompanyID, project.ClientID).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			uniqueNum = 1
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	} else {
		uniqueNum = project.UniqueNO + 1
	}

	projectId := project.ProjectTypeID + "/" + project.CompanyID + "/" + project.ClientID + "/" + strconv.Itoa(uniqueNum) + "/" + strconv.Itoa(project.Year)
	projectTitle := projectId + ": " + project.ProjectName

	project.ProjectID = projectId
	project.ProjectTitle = projectTitle

	if err := db.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}
