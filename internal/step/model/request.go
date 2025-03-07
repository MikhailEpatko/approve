package model

import (
	"approve/internal/common"
)

type CreateStepRequest struct {
	StepGroupId   int64            `json:"step_group_id"  validate:"required,min=1,max=10"`
	Name          string           `json:"name"           validate:"required,min=1,max=155"`
	Number        int              `json:"number"         validate:"required,min=1,max=20"`
	ApproverOrder common.OrderType `json:"approver_order" validate:"required, oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
}

func (r CreateStepRequest) ToEntity() StepEntity {
	return StepEntity{
		StepGroupId:   r.StepGroupId,
		Name:          r.Name,
		Number:        r.Number,
		Status:        common.TEMPLATE,
		ApproverOrder: r.ApproverOrder,
	}
}

type UpdateStepRequest struct {
	Id            int64            `json:"route_id"       validate:"required,min=1"`
	Name          string           `json:"name"           validate:"required,min=1,max=155"`
	Number        int              `json:"number"         validate:"required,min=1,max=20"`
	ApproverOrder common.OrderType `json:"approver_order" validate:"required, oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
}

func (r UpdateStepRequest) ToEntity() StepEntity {
	return StepEntity{
		Id:            r.Id,
		Name:          r.Name,
		Number:        r.Number,
		ApproverOrder: r.ApproverOrder,
	}
}
