package model

import (
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
