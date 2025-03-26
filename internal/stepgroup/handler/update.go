package handler

import (
	cm "approve/internal/common"
	rs "approve/internal/route/service"
	"approve/internal/stepgroup/model"
	svc "approve/internal/stepgroup/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func UpdateStepGroup(c *fiber.Ctx) {
	var request model.UpdateStepGroupRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request", zap.Error(err))
		_ = cm.ErrResponse(c, 400, err.Error())
		return
	}

	stepGroupId, err := svc.UpdateStepGroup(request)
	if err != nil {
		if errors.Is(err, svc.ErrStepGroupNotFound) ||
			errors.Is(err, rs.ErrRouteAlreadyStarted) ||
			errors.Is(err, cm.RequestValidationError{}) {
			_ = cm.ErrResponse(c, 400, err.Error())
		} else {
			_ = cm.ErrResponse(c, 500, fmt.Sprintf("error updating step group: %s", err.Error()))
		}
		cm.Logger.Error("error updating step group", zap.Error(err))
		return
	}
	cm.Logger.Info("updated step group", zap.String("stepGroupId", strconv.FormatInt(stepGroupId, 10)))

	if err := cm.OkResponse(c, stepGroupId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning updated step group id")
		cm.Logger.Error("error returning updated step group id", zap.Error(err))
		return
	}
}
