package main

import (
	"approve/internal/database"
	"approve/internal/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func main() {
	database.Connect()
	app := fiber.New()
	app.Use(middleware.Logger())
	server.SetupRoutes(app)
	_ = app.Listen(8080)
}
