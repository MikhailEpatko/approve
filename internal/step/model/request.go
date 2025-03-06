package model

import (
	"approve/internal/approver/model"
	"approve/internal/common"
)

type CreateStepRequest struct {
	Name          string                        `json:"name"           validate:"required,min=1,max=155"`
	Number        int                           `json:"number"         validate:"required,min=1,max=20"`
	Status        common.Status                 `json:"status"         validate:"required,oneof=TEMPLATE NEW"`
	ApproverOrder common.OrderType              `json:"approver_order" validate:"required, oneof=PARALLEL_ANY_OF PARALLEL_ALL_OF SEQUENTIAL_ALL_OFF"`
	Approvers     []model.CreateApproverRequest `json:"steps"          validate:"required,min=1,max=10"`
}

func (r CreateStepRequest) ToEntity(groupId int64) StepEntity {
	return StepEntity{
		StepGroupId:   groupId,
		Name:          r.Name,
		Number:        r.Number,
		Status:        r.Status,
		ApproverOrder: r.ApproverOrder,
	}
}
