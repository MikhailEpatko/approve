package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/route/service"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func StartRoute(c *fiber.Ctx) {
	routeIdStr := c.Params("routeId")
	routeId, err := strconv.ParseInt(routeIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing route id: %s, %v", routeIdStr, err))
		cm.Logger.Error("error creating route template", zap.Error(err))
		return
	}

	if routeId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid route id: %s", routeIdStr))
		cm.Logger.Error("error creating route template", zap.Error(err))
		return
	}

	err = svc.StartRoute(routeId)
	if err != nil {
		switch err {
		case svc.ErrRouteNotFound, svc.ErrRouteAlreadyStarted, svc.ErrRouteIsFinished, svc.ErrRouteIsTemplate:
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error starting route", zap.Error(err))
	}
	cm.Logger.Info("started route", zap.String("routeId", routeIdStr))

	if err = cm.OkResponse(c, routeId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning started route id")
		cm.Logger.Error("error returning started route id", zap.Error(err))
		return
	}
}
