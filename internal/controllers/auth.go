package controllers

import (
	"os"
	"testapp/internal/config"
	"testapp/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	//This creates a map (like an object in JS) to temporarily hold the request body data (name, email, password).
	var data map[string]string

	//get the data from body
	//c.BodyParser parses the incoming JSON request body and stores it in data.
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	//hash the password
	//[]byte(data["password"]) converts the string to bytes because bcrypt in Go works with byte slices.
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: string(password),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parsing the request body. Same as register.
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	//Try to find the user in the database where the email matches.
	//First(&user) tries to fetch the first match
	var user models.User
	config.DB.Where("email = ?", data["email"]).First(&user)

	//If no user is found, GORM leaves the ID as zero (default for uint).
	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	//Compares the hashed password from DB with the plain password from input.
	//If the password doesn’t match, return 401 Unauthorized.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	//Getting your JWT secret from .env.
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // fallback
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{"token": tokenString})

}
