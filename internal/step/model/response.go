package model

import (
	"approve/internal/approver/model"
	"approve/internal/common"
)

type StepFullResponse struct {
	Id            int64                        `json:"id"`
	StepGroupId   int64                        `json:"step_group_id"`
	Name          string                       `json:"name"`
	Number        int                          `json:"number"`
	Status        common.Status                `json:"status"`
	ApproverOrder common.OrderType             `json:"approver_order"`
	IsApproved    bool                         `json:"is_approved"`
	Approvers     []model.ApproverFullResponse `json:"approvers"`
}
