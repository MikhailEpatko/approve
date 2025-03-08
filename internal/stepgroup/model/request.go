package model

import (
	"approve/internal/common"
)

type CreateStepGroupRequest struct {
	RouteId   int64            `json:"route_id"   validate:"required,min=1"`
	Name      string           `json:"name"       validate:"required,min=1,max=155"`
	Number    int              `json:"number"     validate:"required,min=1,max=10"`
	StepOrder common.OrderType `json:"step_order" validate:"required,oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SERIAL"`
}

func (r CreateStepGroupRequest) ToEntity() StepGroupEntity {
	return StepGroupEntity{
		RouteId:   r.RouteId,
		Name:      r.Name,
		Number:    r.Number,
		Status:    common.TEMPLATE,
		StepOrder: r.StepOrder,
	}
}

type UpdateStepGroupRequest struct {
	Id        int64            `json:"id"         validate:"required,min=1"`
	Name      string           `json:"name"       validate:"required,min=1,max=155"`
	Number    int              `json:"number"     validate:"required,min=1,max=10"`
	StepOrder common.OrderType `json:"step_order" validate:"required,oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SERIAL"`
}

func (r UpdateStepGroupRequest) ToEntity() StepGroupEntity {
	return StepGroupEntity{
		Id:        r.Id,
		Name:      r.Name,
		Number:    r.Number,
		StepOrder: r.StepOrder,
	}
}
