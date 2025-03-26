package handler

import (
	"approve/internal/approver/model"
	svc "approve/internal/approver/service"
	cm "approve/internal/common"
	rs "approve/internal/route/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func UpdateApprover(c *fiber.Ctx) {
	var request model.UpdateApproverRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request", zap.Error(err))
		_ = cm.ErrResponse(c, 400, err.Error())
		return
	}

	approverId, err := svc.UpdateApprover(request)
	if err != nil {
		if errors.Is(err, rs.ErrRouteAlreadyStarted) ||
			errors.Is(err, svc.ErrApproverNotFound) ||
			errors.Is(err, cm.RequestValidationError{}) {
			_ = cm.ErrResponse(c, 400, err.Error())
		} else {
			_ = cm.ErrResponse(c, 500, fmt.Sprintf("error updating approver: %s", err.Error()))
		}
		cm.Logger.Error("error updating approver", zap.Error(err))
		return
	}
	cm.Logger.Info("updated approver", zap.String("approverId", strconv.FormatInt(approverId, 10)))

	if err = cm.OkResponse(c, approverId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning updated approver id")
		cm.Logger.Error("error returning updated approver id", zap.Error(err))
		return
	}
}
