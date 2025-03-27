package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/route/service"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func FindFullRouteById(c *fiber.Ctx) {
	routeIdStr := c.Params("routeId")
	routeId, err := strconv.ParseInt(routeIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing route id: %s, %v", routeIdStr, err))
		cm.Logger.Error("error find route by id", zap.Error(err))
		return
	}

	if routeId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid route id: %s", routeIdStr))
		cm.Logger.Error("error finding route by id", zap.Error(err))
		return
	}

	fullRoute, err := svc.FindFullRouteById(routeId)
	if err != nil {
		switch {
		case err == svc.ErrRouteNotFound:
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error finding route by id", zap.Error(err))
	}

	if err = cm.OkResponse(c, fullRoute); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning found route")
		cm.Logger.Error("error returning found route", zap.Error(err))
		return
	}
}
