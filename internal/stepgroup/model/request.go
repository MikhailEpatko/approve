package model

import (
	"approve/internal/common"
)

type CreateStepGroupRequest struct {
	RouteId   int64            `json:"route_id"   validate:"required,min=1"`
	Name      string           `json:"name"       validate:"required,min=1,max=155"`
	Number    int              `json:"number"     validate:"required,min=1,max=10"`
	StepOrder common.OrderType `json:"step_order" validate:"required,oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
}

func (request CreateStepGroupRequest) ToEntity() StepGroupEntity {
	return StepGroupEntity{
		RouteId:   request.RouteId,
		Name:      request.Name,
		Number:    request.Number,
		Status:    common.TEMPLATE,
		StepOrder: request.StepOrder,
	}
}
