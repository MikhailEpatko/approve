package server

import (
	routeHandler "approve/internal/route/handler"
	stepHandler "approve/internal/step/handler"
	stepGroupHandler "approve/internal/stepgroup/handler"
	"github.com/gofiber/fiber"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/route", routeHandler.CreateRouteTemplate)
	api.Put("/route", routeHandler.UpdateRoute)
	api.Put("/route/:routeId", routeHandler.StartRoute)
	api.Post("/route/by-filter", routeHandler.FindByFilter)

	api.Post("/group", stepGroupHandler.CreateStepGroupTemplate)
	api.Put("/group", stepGroupHandler.UpdateStepGroup)

	api.Post("/step", stepHandler.CreateStepTemplate)
	api.Put("/step", stepHandler.UpdateStep)

}
