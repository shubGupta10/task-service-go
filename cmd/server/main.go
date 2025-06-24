package main

import (
	"log"
	"os"
	"testapp/internal/config"
	"testapp/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.ConnectDatabase()

	// basic health route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Server is running")
	})

	//routes
	routes.AuthRoutes(app)
	routes.TodoRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	//start the server
	log.Fatal(app.Listen(":" + port))
}
