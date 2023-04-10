package handlers

// import (
// 	"github.com/andrelaurent/project-register/database"
// 	"github.com/andrelaurent/project-register/models"
// 	"github.com/gofiber/fiber/v2"
// )

// func CreateManager(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	manager := new(models.Manager)

// 	err := c.BodyParser(manager)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
// 	}

// 	err = db.Create(&manager).Error
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create manager", "data": err})
// 	}

// 	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Manager has created", "data": manager})
// }

// func GetAllManagers(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	var managers []models.Manager

// 	db.Find(&managers)

// 	if len(managers) == 0 {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "managers not found", "data": nil})
// 	}

// 	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Manager Found", "data": managers})
// }

// func GetManagerByID(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	var manager models.Manager

// 	id := c.Params("id")

// 	err := db.Find(&manager, "id = ?", id).Error
// 	if err != nil {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "manager not found", "data": nil})
// 	}

// 	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "manager retrieved", "data": manager})
// }

// func SearchManager(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	req := new(models.Manager)
// 	if err := c.BodyParser(req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to parse request body",
// 		})
// 	}

// 	var managers []models.Manager
// 	if err := db.Where("manager_name LIKE ?", "%"+req.ManagerName+"%").Find(&managers).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to search managers",
// 		})
// 	}

// 	return c.JSON(managers)
// }

// func UpdateManager(c *fiber.Ctx) error {

// 	type updatemanager struct {
// 		Username string `json:"name"`
// 	}

// 	db := database.DB.Db
// 	var manager models.Manager

// 	id := c.Params("id")

// 	db.Find(&manager, "id = ?", id)

// 	if manager == (models.Manager{}) {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Manager not found", "data": nil})
// 	}

// 	var updatemanagerData updatemanager
// 	err := c.BodyParser(&updatemanagerData)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
// 	}
// 	manager.ManagerName = updatemanagerData.Username

// 	db.Save(&manager)

// 	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Managers Found", "data": manager})
// }

// func DeleteManagerByID(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	var manager models.Manager

// 	id := c.Params("id")

// 	db.Find(&manager, "id = ?", id)

// 	if manager == (models.Manager{}) {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Manager not found", "data": nil})
// 	}

// 	if manager.ManagerID == "" || manager.ManagerName == "" {
// 		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Manager ID and name are required", "data": nil})
// 	}

// 	err := db.Delete(&manager, "id = ?", id).Error
// 	if err != nil {
// 		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete manager", "data": nil})
// 	}
// 	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Manager deleted"})
// }
