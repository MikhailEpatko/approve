package handler

import (
	cm "approve/internal/common"
	"approve/internal/route/model"
	svc "approve/internal/route/service"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
)

func FindByFilter(c *fiber.Ctx) {
	var filter model.FilterRouteRequest
	if err := c.BodyParser(&filter); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	result, err := svc.FindByFilter(filter)
	if err != nil {
		switch {
		case errors.As(err, &cm.RequestValidationError{}):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error finding routes by filter", zap.Error(err))
		return
	}
	cm.Logger.Info("success finding routes by filter")

	if err = cm.OkResponse(c, result); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning found routes")
		cm.Logger.Error("error returning found routes", zap.Error(err))
		return
	}
}
