package handler

import (
	cm "approve/internal/common"
	"approve/internal/resolution/model"
	svc "approve/internal/resolution/service"
	"errors"
	"github.com/gofiber/fiber"
	"go.uber.org/zap"
	"strconv"
)

func CreateResolution(c *fiber.Ctx) {
	var request model.CreateResolutionRequest
	if err := c.BodyParser(&request); err != nil {
		cm.Logger.Error("error parsing request body", zap.Error(err))
		_ = cm.ErrResponse(c, 500, err.Error())
		return
	}

	resolutionId, err := svc.CreateResolution(request)
	if err != nil {
		switch {
		case errors.Is(err, &cm.RequestValidationError{}) ||
			errors.Is(err, svc.ErrCommentShouldBeProvided) ||
			errors.Is(err, svc.ErrApproverNotFound) ||
			errors.Is(err, svc.ErrApproverIsNotStarted):
			_ = cm.ErrResponse(c, 400, err.Error())
		default:
			_ = cm.ErrResponse(c, 500, err.Error())
		}
		cm.Logger.Error("error creating resolution", zap.Error(err))
		return
	}
	cm.Logger.Info("created resolution", zap.String("resolutionId", strconv.FormatInt(resolutionId, 10)))

	if err = cm.OkResponse(c, resolutionId); err != nil {
		_ = cm.ErrResponse(c, 500, "error returning created resolution id")
		cm.Logger.Error("error returning created resolution id", zap.Error(err))
		return
	}
}
