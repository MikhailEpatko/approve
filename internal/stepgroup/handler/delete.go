package handler

import (
	cm "approve/internal/common"
	svc "approve/internal/stepgroup/service"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

var errorDeletingMsg = "error deleting step group by id"

func DeleteStepGroupById(c *fiber.Ctx) {
	stepGroupIdStr := c.Params("stepGroupId")
	stepGroupId, err := strconv.ParseInt(stepGroupIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing step group id: %s, %v", stepGroupIdStr, err))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	if stepGroupId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid step group id: %s", stepGroupIdStr))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	err = svc.DeleteStepGroupById(stepGroupId)
	if err != nil {
		_ = cm.ErrResponse(c, 500, err.Error())
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
	}

	cm.Logger.Info("deleted step group", zap.String("stepGroupId", stepGroupIdStr))
	_ = cm.OkResponse(c, nil)
	return
}
