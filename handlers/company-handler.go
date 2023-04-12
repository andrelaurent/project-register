package handlers

import (
	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	company := new(models.Company)

	err := c.BodyParser(company)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Something's wrong with your input", "data": err,
		})
	}

	if company.CompanyCode == "" || company.CompanyName == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Company ID and name are required", "data": nil})
	}

	var existingCompany models.Company
	if err := db.Where("company_code = ?", company.CompanyCode).First(&existingCompany).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "Company code already exists", "data": nil})
	}

	err = db.Create(&company).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error", "message": "Could not create company", "data": err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "Company has created", "data": company,
	})
}

func GetAllCompanies(c *fiber.Ctx) error {
	db := database.DB.Db
	var company []models.Company

	db.Find(&company)

	if len(company) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": "error", "message": "no company found", "data": "nil",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "sucess", "message": "Companies Found", "data": company,
	})
}

func GetCompanyByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var company models.Company

	id := c.Params("id")

	err := db.Find(&company, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "company not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "company retrieved", "data": company})
}

func SearchCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	req := new(models.Company)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var companys []models.Company
	if err := db.Where("company_name LIKE ?", "%"+req.CompanyName+"%").Find(&companys).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search companys",
		})
	}

	return c.JSON(companys)
}

func UpdateCompany(c *fiber.Ctx) error {

	type updatecompany struct {
		CompanyCode string `json:"company_code"`
		CompanyName string `json:"company_name"`
	}

	db := database.DB.Db
	var company models.Company

	id := c.Params("id")

	db.Find(&company, "id = ?", id)

	if company == (models.Company{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "company not found", "data": nil})
	}

	var updatecompanyData updatecompany
	err := c.BodyParser(&updatecompanyData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	company.CompanyCode = updatecompanyData.CompanyCode
	company.CompanyName = updatecompanyData.CompanyName

	db.Save(&company)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "companys Found", "data": company})
}

func DeleteCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var company models.Company

	result := db.Where("id = ?", id).Delete(&company)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete company", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company has been deleted", "data": result.RowsAffected})
}

func HardDeleteCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var company models.Company

	result := db.Unscoped().Where("id = ?", id).Delete(&company)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete company", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company has been deleted", "data": result.RowsAffected})
}
