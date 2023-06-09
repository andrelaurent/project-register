package handlers

import (
	"math"
	"strconv"
	"time"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/gofiber/fiber/v2"
)

func CreateCompanyAuditEntry(action string, company models.Company) error {
	db := database.DB.Db

	audit := models.CompanyAudit{
		CompanyID:   company.ID,
		CompanyCode: company.CompanyCode,
		CompanyName: company.CompanyName,
		Action:      action,
		Date:        time.Now().Format("2006-01-02 15:04:05"),
	}

	if err := db.Create(&audit).Error; err != nil {
		return err
	}

	return nil
}

func CreateCompany(c *fiber.Ctx) error {
	db := database.DB.Db
	var company models.Company

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

	if err := CreateCompanyAuditEntry("create", company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create company",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success", "message": "Company has created", "data": company,
	})
}

func GetAllCompanies(c *fiber.Ctx) error {
	db := database.DB.Db

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var companies []models.Company

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&companies)

	if len(companies) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Companies not found", "data": nil})
	}

	var total int64
	db.Model(&models.Company{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Companies Found",
		"data":        companies,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
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

	var companies []models.Company
	var total int64

	if err := db.Model(&models.Company{}).Where("company_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Companys",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("company_name ILIKE ?", "%"+searchQuery+"%").Find(&companies).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Companys",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Companys Found",
		"data":        companies,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
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

	if err := CreateCompanyAuditEntry("update", company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update company",
		})
	}


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

	if err := CreateCompanyAuditEntry("soft delete", company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete company",
		})
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

	if err := CreateCompanyAuditEntry("hard delete", company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete company",
		})
	}


	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Company not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Company has been deleted from database", "data": result.RowsAffected})
}

func RecoverCompany(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		CompanyCode string `json:"company_code"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var company models.Company
	if err := db.Unscoped().Where("company_code = ? AND deleted_at IS NOT NULL", request.CompanyCode).First(&company).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Company not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&company).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover company",
			"data":    nil,
		})
	}

	if err := CreateCompanyAuditEntry("recover", company); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to recover company",
		})
	}


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Company recovered",
	})
}
