package handler

import (
	cm "approve/internal/common"
	"approve/internal/route/model"
	svc "approve/internal/route/service"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateRouteTemplate(c *fiber.Ctx) {
	var request model.CreateRouteTemplateRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	routeId, err := svc.CreateRouteTemplate(request)
	if err != nil {
		switch {
		case errors.As(err, &cm.RequestValidationError{}):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error creating route template", zap.Error(err))
		return
	}
	cm.Logger.Info("created route", zap.String("routeId", strconv.FormatInt(routeId, 10)))

	if err := cm.OkResponse(c, routeId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning created route template id")
		cm.Logger.Error("error returning created route template id", zap.Error(err))
		return
	}
}
