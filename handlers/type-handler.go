package handlers

// func GetTypes(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	var type []models.ProjectType

// 	db.Find(&type)

// 	if len(type) == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"status": "error", "message": "no type found", "data": "nil",
// 		})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status": "sucess", "message": "Types Found", "data": company,
// 	})
// }