package model

type Decision string

const (
	UNKNOWN  Decision = "UNKNOWN"
	ACCEPT   Decision = "ACCEPT"
	REJECT   Decision = "REJECT"
	REVISION Decision = "REVISION"
)

type ResolutionEntity struct {
	Id         int64    `db:"id"`
	ApproverId int64    `db:"approver_id"`
	Decision   Decision `db:"decision"`
	Comment    string   `db:"comment"`
	Deleted    bool     `db:"deleted"`
}
