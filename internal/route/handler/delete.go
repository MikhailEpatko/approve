package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/route/service"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

var errorDeletingMsg = "error deleting route by id"

func DeleteRouteById(c *fiber.Ctx) {
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

	err = svc.DeleteRouteById(routeId)
	if err != nil {
		_ = cm.ErrResponse(c, 500, err.Error())
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
	}

	cm.Logger.Info("deleted route", zap.String("routeId", routeIdStr))
	_ = cm.OkResponse(c, nil)
}
