package server

import (
	approverHandler "approve/internal/approver/handler"
	resolutionHandler "approve/internal/resolution/handler"
	routeHandler "approve/internal/route/handler"
	stepHandler "approve/internal/step/handler"
	stepGroupHandler "approve/internal/stepgroup/handler"
	"github.com/gofiber/fiber"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/route/template", routeHandler.CreateRouteTemplate)
	api.Put("/route", routeHandler.UpdateRoute)
	api.Put("/route/:routeId", routeHandler.StartRoute)
	api.Post("/route/by-filter", routeHandler.FindByFilter)
	api.Get("/route/:routeId", routeHandler.FindFullRouteById)
	api.Delete("/route/:routeId", routeHandler.DeleteRouteById)
	// Todo: api.Post("/route/create-from-template/:routeTemplateId", routeHandler.CreateNewRouteFromTemplate)

	api.Post("/group/template", stepGroupHandler.CreateStepGroupTemplate)
	api.Put("/group", stepGroupHandler.UpdateStepGroup)
	api.Delete("/group/:stepGroupId", stepGroupHandler.DeleteStepGroupById)

	api.Post("/step/template", stepHandler.CreateStepTemplate)
	api.Put("/step", stepHandler.UpdateStep)
	api.Delete("/step/:stepId", stepHandler.DeleteStepById)

	api.Post("/approver", approverHandler.CreateApproverTemplate)
	api.Put("/approver", approverHandler.UpdateApprover)
	api.Delete("/approver/:approverId", approverHandler.DeleteApproverById)

	api.Post("/resolution", resolutionHandler.CreateResolution)
}
