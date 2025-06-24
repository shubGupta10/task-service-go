package routes

import (
	"testapp/internal/controllers"
	"testapp/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(app *fiber.App) {
	todo := app.Group("/todo")
	todo.Post("/create-todo", middleware.Protected(), controllers.CreateTodo)
	todo.Get("/get-todos", middleware.Protected(), controllers.GetTodos)
}
