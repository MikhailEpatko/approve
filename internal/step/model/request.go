package model

import (
	"approve/internal/approver/model"
	"approve/internal/common"
)

type CreateStepRequest struct {
	Name        string        `db:"name"`
	Number      int           `db:"number"`
	Status      common.Status `db:"status"`
	ApproveType ApproveType   `db:"approve_type"`
	Approvers   []model.CreateApproverRequest
}

func (r CreateStepRequest) ToEntity(groupId int64) StepEntity {
	return StepEntity{
		StepGroupId: groupId,
		Name:        r.Name,
		Number:      r.Number,
		Status:      r.Status,
		ApproveType: r.ApproveType,
	}
}
