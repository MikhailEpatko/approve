package handler

import (
	cm "approve/internal/common"
	"approve/internal/step/model"
	svc "approve/internal/step/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func UpdateStep(c *fiber.Ctx) {
	var request model.UpdateStepRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request", zap.Error(err))
		_ = cm.ErrResponse(c, 400, err.Error())
		return
	}

	stepId, err := svc.UpdateStep(request)
	if err != nil {
		if errors.Is(err, svc.ErrStepNotFound) ||
			errors.Is(err, svc.ErrStepAlreadyStarted) ||
			errors.Is(err, svc.ErrStepIsFinished) ||
			errors.Is(err, cm.RequestValidationError{}) {
			_ = cm.ErrResponse(c, 400, err.Error())
		} else {
			_ = cm.ErrResponse(c, 500, fmt.Sprintf("error updating step: %s", err.Error()))
		}
		cm.Logger.Error("error updating step", zap.Error(err))
		return
	}
	cm.Logger.Info("updated step", zap.String("stepId", strconv.FormatInt(stepId, 10)))

	if err := cm.OkResponse(c, stepId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning updated step id")
		cm.Logger.Error("error returning updated step id", zap.Error(err))
		return
	}
}
