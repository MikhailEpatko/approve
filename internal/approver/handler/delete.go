package handler

import (
	svc "approve/internal/approver/service"
	cm "approve/internal/common"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

var errorDeletingMsg = "error deleting approver by id"

func DeleteApproverById(c *fiber.Ctx) {
	approverIdStr := c.Params("approverId")
	approverId, err := strconv.ParseInt(approverIdStr, 10, 64)
	if err != nil {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("error parsing approver id: %s, %v", approverIdStr, err))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	if approverId <= 0 {
		_ = cm.ErrResponse(c, 400, fmt.Sprintf("invalid approver id: %s", approverIdStr))
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
		return
	}

	err = svc.DeleteApproverById(approverId)
	if err != nil {
		_ = cm.ErrResponse(c, 500, err.Error())
		cm.Logger.Error(errorDeletingMsg, zap.Error(err))
	}

	cm.Logger.Info("deleted approver", zap.String("approverId", approverIdStr))
	_ = cm.OkResponse(c, nil)
}
