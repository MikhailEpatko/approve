package model

import (
	"approve/internal/common"
)

type ApproveType string

const (
	PARALLEL_ANY_OF    ApproveType = "PARALLEL_ANY_OF"
	PARALLEL_ALL_OF    ApproveType = "PARALLEL_ALL_OF"
	SEQUENTIAL_ALL_OFF ApproveType = "SEQUENTIAL_ALL_OFF"
)

type StepEntity struct {
	Id          int64         `db:"id"`
	StepGroupId int64         `db:"step_group_id"`
	Name        string        `db:"name"`
	Number      int           `db:"number"`
	Status      common.Status `db:"status"`
	ApproveType ApproveType   `db:"approve_type"`
	Deleted     bool          `db:"deleted"`
}
