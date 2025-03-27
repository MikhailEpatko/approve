package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/step/service"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

var errorDeletingMsg = "error deleting step by id"

func DeleteStepById(c *fiber.Ctx) {
	stepIdStr := c.Params("stepId")
	stepId, err := strconv.ParseInt(stepIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing step id: %s, %v", stepIdStr, err))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	if stepId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid step id: %s", stepIdStr))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	err = svc.DeleteStepById(stepId)
	if err != nil {
		_ = cm.ErrResponse(c, 500, err.Error())
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
	}

	cm.Logger.Info("deleted step", zap.String("stepId", stepIdStr))
	_ = cm.OkResponse(c, nil)
	return
}
