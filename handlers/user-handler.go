package handlers

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"time"

	"math"
	"strconv"

	"github.com/andrelaurent/project-register/database"
	"github.com/andrelaurent/project-register/models"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secretkeyforbackendprojectregistrationtest"

func UserLogin(c *fiber.Ctx) error {
	db := database.DB.Db
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	var result models.User
	if err := db.Where("email = ?", user.Email).First(&result).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	if !PasswordsMatch(user.Password, result.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, errToken := claims.SignedString([]byte(SecretKey))

	if errToken != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	if err := db.Model(&result).Update("token", token).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "project-registration-backend",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result.ID,
		"token":  token,
	})
}

func ValidateToken(c *fiber.Ctx, userId uint, token string) bool {
	var user models.User
	db := database.DB.Db

	if err := db.First(&user, userId).Error; err != nil {
		return false
	}

	return user.Token == token
}

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("project-registration-backend")

		claims := &jwt.StandardClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		if !ValidateToken(c, uint(userID), token) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token",
			})
		}

		return c.Next()
	}
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(models.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Username, email, and password are required", "data": nil})
	}

	if !govalidator.IsEmail(user.Email) {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid email address", "data": nil})
	}

	if len(user.Password) < 8 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Password must be at least 8 characters long", "data": nil})
	}

	var existingUser models.User
	if err := db.Where("id = ?", user.ID).First(&existingUser).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"status": "error", "message": "User ID already exists", "data": nil})
	}

	user.Password = HashPassword(user.Password)
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not create user", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "User created", "data": user})
}

func PasswordsMatch(password1 string, password2 string) bool {
	hash1 := sha512.Sum512([]byte(password1))
	hash2 := sha512.Sum512([]byte(password2))

	return subtle.ConstantTimeCompare(hash1[:], hash2[:]) == 1
}

func HashPassword(password string) string {
	hash := sha512.Sum512([]byte(password))
	encodedHash := base64.StdEncoding.EncodeToString(hash[:])

	return encodedHash
}

func GetAllUsers(c *fiber.Ctx) error {
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

	var Users []models.User

	db.Order("id ASC").Limit(limit).Offset(offset).Find(&Users)

	if len(Users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Users not found", "data": nil})
	}

	var total int64
	db.Model(&models.User{}).Count(&total)

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Users Found",
		"data":        Users,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(200).JSON(response)
}

func GetUserByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var User models.User

	id := c.Params("id")

	err := db.Find(&User, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User retrieved", "data": User})
}

func SearchUser(c *fiber.Ctx) error {
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

	var Users []models.User
	var total int64

	if err := db.Model(&models.User{}).Where("User_name ILIKE ?", "%"+searchQuery+"%").Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Users",
		})
	}

	if err := db.Limit(limit).Offset(offset).Where("User_name ILIKE ?", "%"+searchQuery+"%").Find(&Users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search Users",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := fiber.Map{
		"status":      "success",
		"message":     "Users Found",
		"data":        Users,
		"currentPage": page,
		"perPage":     limit,
		"totalPages":  totalPages,
		"totalItems":  total,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func UpdateUser(c *fiber.Ctx) error {

	type updateUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	db := database.DB.Db
	var User models.User

	id := c.Params("id")

	db.Find(&User, "id = ?", id)

	if User == (models.User{}) {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}

	User.Username = updateUserData.Username
	User.Email = updateUserData.Email
	User.Password = updateUserData.Password
	db.Save(&User)

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Users Found", "data": User})
}

func DeleteUser(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var User models.User
	result := db.Where("id = ?", id).Delete(&User)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete user", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User has been deleted", "data": result.RowsAffected})
}

func HardDeleteUser(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")

	var User models.User
	result := db.Unscoped().Where("id = ?", id).Delete(&User)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not delete user", "data": result.Error})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User has been deleted from database", "data": result.RowsAffected})
}

func RecoverUser(c *fiber.Ctx) error {
	db := database.DB.Db

	var request struct {
		UserCode string `json:"User_code"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    nil,
		})
	}

	var User models.User
	if err := db.Unscoped().Where("User_code = ? AND deleted_at IS NOT NULL", request.UserCode).First(&User).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
	}

	if err := db.Unscoped().Model(&User).Update("deleted_at", nil).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to recover user",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User recovered",
	})
}
