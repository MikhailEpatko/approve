package model

type Decision string

const (
	UNKNOWN  Decision = "UNKNOWN"
	ACCEPT   Decision = "ACCEPT"
	REJECT   Decision = "REJECT"
	REVISION Decision = "REVISION"
)

type ResolutionEntity struct {
	Id         int      `db:"id"`
	ApproverId int      `db:"approver_id"`
	Decision   Decision `db:"decision"`
	Comment    string   `db:"comment"`
	Deleted    bool     `db:"deleted"`
}

func NewResolutionEntity(
	approverId int,
) *ResolutionEntity {
	return &ResolutionEntity{
		ApproverId: approverId,
		Decision:   UNKNOWN,
		Deleted:    false,
	}
}
