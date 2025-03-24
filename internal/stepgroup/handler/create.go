package handler

import (
	cm "approve/internal/common"
	"approve/internal/stepgroup/model"
	svc "approve/internal/stepgroup/service"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateStepGroupTemplate(c *fiber.Ctx) {
	var request model.CreateStepGroupRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	groupId, err := svc.CreateStepGroupTemplate(request)
	if err != nil {
		switch {
		case errors.As(err, &cm.RequestValidationError{}):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error creating step group template", zap.Error(err))
		return
	}
	cm.Logger.Info("created step group", zap.String("groupId", strconv.FormatInt(groupId, 10)))

	if err := cm.OkResponse(c, groupId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning created step group template id")
		cm.Logger.Error("error returning created stepgroup template id", zap.Error(err))
		return
	}
}
