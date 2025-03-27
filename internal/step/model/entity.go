package model

import (
	am "approve/internal/approver/model"
	"approve/internal/common"
)

type StepEntity struct {
	Id            int64            `db:"id"`
	StepGroupId   int64            `db:"step_group_id"`
	Name          string           `db:"name"`
	Number        int              `db:"number"`
	Status        common.Status    `db:"status"`
	ApproverOrder common.OrderType `db:"approver_order"`
	IsApproved    bool             `db:"is_approved"`
}

func (e StepEntity) ToFullResponse(approvers []am.ApproverFullResponse) StepFullResponse {
	return StepFullResponse{
		Id:            e.Id,
		StepGroupId:   e.StepGroupId,
		Name:          e.Name,
		Number:        e.Number,
		Status:        e.Status,
		ApproverOrder: e.ApproverOrder,
		IsApproved:    e.IsApproved,
		Approvers:     approvers,
	}
}
