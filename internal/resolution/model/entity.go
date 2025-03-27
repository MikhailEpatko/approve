package model

import "approve/internal/common"

type ResolutionEntity struct {
	Id         int64  `db:"id"`
	ApproverId int64  `db:"approver_id"`
	IsApproved bool   `db:"is_approved"`
	Comment    string `db:"comment"`
}

func (e ResolutionEntity) ToResponse() ResolutionResponse {
	return ResolutionResponse{
		Id:         e.Id,
		ApproverId: e.ApproverId,
		IsApproved: e.IsApproved,
		Comment:    e.Comment,
	}
}

type ApprovingInfoEntity struct {
	RouteId        int64            `db:"route_id"`
	StepGroupId    int64            `db:"step_group_id"`
	StepOrder      common.OrderType `db:"step_order"`
	StepId         int64            `db:"step_id"`
	ApproverOrder  common.OrderType `db:"approver_order"`
	ApproverId     int64            `db:"approver_id"`
	Guid           string           `db:"guid"`
	ApproverStatus common.Status    `db:"approver_status"`
}
