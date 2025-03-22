package handler

import (
	cm "approve/internal/common"
	"approve/internal/route/model"
	svc "approve/internal/route/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func UpdateRoute(c *fiber.Ctx) {
	var request model.UpdateRouteRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request", zap.Error(err))
		_ = cm.ErrResponse(c, 400, err.Error())
		return
	}

	routeId, err := svc.UpdateRoute(request)
	if err != nil {
		if errors.Is(err, svc.ErrRouteNotFound) ||
			errors.Is(err, svc.ErrRouteAlreadyStarted) ||
			errors.Is(err, svc.ErrRouteIsFinished) {
			_ = cm.ErrResponse(c, 400, err.Error())
		} else {
			_ = cm.ErrResponse(c, 500, fmt.Sprintf("error updating route: %s", err.Error()))
		}
		cm.Logger.Error("error updating route", zap.Error(err))
		return
	}
	cm.Logger.Info("updated route", zap.String("routeId", strconv.FormatInt(routeId, 10)))

	if err := cm.OkResponse(c, routeId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning updated route id")
		cm.Logger.Error("error returning updated route id", zap.Error(err))
		return
	}
}
