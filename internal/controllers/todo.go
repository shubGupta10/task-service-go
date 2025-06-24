package controllers

import (
	"strconv"
	"testapp/internal/config"
	"testapp/internal/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTodo(c *fiber.Ctx) error {
	var todo map[string]string

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Convert user_id from string to uint
	userIdStr := todo["user_id"]
	var userId uint
	if parsedId, err := strconv.ParseUint(userIdStr, 10, 32); err == nil {
		userId = uint(parsedId)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	todoItems := models.Todo{
		ID:          0,
		Title:       todo["title"],
		IsCompleted: todo["is_completed"] == "true", // Convert string to bool
		UserId:      userId,
	}

	if err := config.DB.Create(&todoItems).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create todo"})
	}

	return c.JSON(todoItems)
}

func GetTodos(c *fiber.Ctx) error {
	var todos []models.Todo

	if err := config.DB.Find(&todos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve todos"})
	}

	return c.JSON(todos)
}
