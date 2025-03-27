package handler

import (
	cm "approve/internal/common"
	"approve/internal/step/model"
	svc "approve/internal/step/service"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateStepTemplate(c *fiber.Ctx) {
	var request model.CreateStepRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	stepId, err := svc.CreateStepTemplate(request)
	if err != nil {
		switch {
		case errors.Is(err, &cm.RequestValidationError{}):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error creating step template", zap.Error(err))
		return
	}
	cm.Logger.Info("created step template", zap.String("stepId", strconv.FormatInt(stepId, 10)))

	if err = cm.OkResponse(c, stepId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning created step template id")
		cm.Logger.Error("error returning created step template id", zap.Error(err))
		return
	}
}
