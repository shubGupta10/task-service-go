package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		//checking if the header exists and also if it starts with "Bearer " (space included).
		//If the header is missing or doesn’t start correctly → you return a 401 Unauthorized error.
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// You are parsing the token using the JWT library
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		//If token parsing failed, or the token is not valid → return 401 Unauthorized.
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		//If the token is valid → you call c.Next() to pass control to the next handler
		return c.Next()
	}
}
