package model

import (
	"approve/internal/common"
	"approve/internal/step/model"
)

type CreateStepGroupRequest struct {
	Name      string                    `json:"name"       validate:"required,min=1,max=155"`
	Number    int                       `json:"number"     validate:"required,min=1,max=20"`
	Status    common.Status             `json:"status"     validate:"required,oneof=TEMPLATE NEW"`
	StepOrder common.OrderType          `json:"step_order" validate:"required,oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
	Steps     []model.CreateStepRequest `json:"steps"      validate:"required,min=1,max=10"`
}

func (request CreateStepGroupRequest) ToEntity(routeId int64) StepGroupEntity {
	return StepGroupEntity{
		RouteId:   routeId,
		Name:      request.Name,
		Number:    request.Number,
		Status:    request.Status,
		StepOrder: request.StepOrder,
	}
}
