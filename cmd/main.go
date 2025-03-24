package main

import (
	cfg "approve/internal/database"
	"approve/internal/server"
	"github.com/gofiber/fiber"
)

func main() {
	cfg.Connect()

	app := fiber.New()
	server.SetupRoutes(app)

}
