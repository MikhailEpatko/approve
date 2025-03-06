package model

import (
	"approve/internal/approver/model"
	"approve/internal/common"
)

type CreateStepRequest struct {
	StepGroupId   int64                         `json:"step_group_id"  validate:"required,min=1,max=10"`
	Name          string                        `json:"name"           validate:"required,min=1,max=155"`
	Number        int                           `json:"number"         validate:"required,min=1,max=20"`
	ApproverOrder common.OrderType              `json:"approver_order" validate:"required, oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
	Approvers     []model.CreateApproverRequest `json:"steps"          validate:"required,min=1,max=10"`
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
