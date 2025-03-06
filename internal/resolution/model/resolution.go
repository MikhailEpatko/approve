package model

import "approve/internal/common"

type ResolutionEntity struct {
	Id         int64           `db:"id"`
	ApproverId int64           `db:"approver_id"`
	Decision   common.Decision `db:"decision"`
	Comment    string          `db:"comment"`
	Deleted    bool            `db:"deleted"`
}
