package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/route/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateNewRouteFromTemplate(c *fiber.Ctx) {
	routeIdStr := c.Params("routeId")
	routeId, err := strconv.ParseInt(routeIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing route id: %s, %v", routeIdStr, err))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	if routeId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid route id: %s", routeIdStr))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	newRouteId, err := svc.CreateNewRouteFromTemplate(routeId)
	if err != nil {
		if errors.Is(err, svc.ErrRouteNotFound) ||
			errors.Is(err, svc.ErrSourceRouteIsNotTemplate) {
			_ = cm.ErrResponse(c, 400, err.Error())
		} else {
			_ = cm.ErrResponse(c, 500, fmt.Sprintf("error creating new route from template: %s", err.Error()))
		}
		cm.Logger.Error("error creating new route from template", zap.Error(err))
		return
	}
	cm.Logger.Info("created new route from template", zap.String("newRouteId", strconv.FormatInt(newRouteId, 10)))

	if err = cm.OkResponse(c, newRouteId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning new route id")
		cm.Logger.Error(
			"error returning new route id",
			zap.String("routeId", routeIdStr),
			zap.String("newRouteId", strconv.FormatInt(newRouteId, 10)),
			zap.Error(err),
		)
		return
	}
}
