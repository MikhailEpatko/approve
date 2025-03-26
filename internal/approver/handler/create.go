package handler

import (
	"approve/internal/approver/model"
	svc "approve/internal/approver/service"
	cm "approve/internal/common"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateApproverTemplate(c *fiber.Ctx) {
	var request model.CreateApproverRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	approverId, err := svc.CreateApproverTemplate(request)
	if err != nil {
		switch {
		case errors.As(err, &cm.RequestValidationError{}):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error creating approver template", zap.Error(err))
		return
	}
	cm.Logger.Info("created approver template", zap.String("approverId", strconv.FormatInt(approverId, 10)))

	if err = cm.OkResponse(c, approverId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning created approver template id")
		cm.Logger.Error("error returning created approver template id", zap.Error(err))
		return
	}
}
