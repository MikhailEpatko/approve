package server

import (
	cfg "approve/internal/config"
	routeHandler "approve/internal/route/handler"
	"github.com/gofiber/fiber"
)

func RegisterRoutes(app *fiber.App) (err error) {
	appCfg := cfg.NewAppConfig()
	cfg.ConnectDatabase(appCfg)
	//app = fiber.New()
	api := app.Group("/api")

	api.Post("/route", routeHandler.CreateRouteTemplate)
	api.Put("/route", routeHandler.UpdateRoute)
	api.Put("/route/:routeId", routeHandler.StartRoute)
	api.Post("/route/by-filter", routeHandler.FindByFilter)

	return err
}
