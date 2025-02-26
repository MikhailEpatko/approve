package model

import "time"

type Decision string

const (
	UNKNOWN  Decision = "UNKNOWN"
	ACCEPT   Decision = "ACCEPT"
	REJECT   Decision = "REJECT"
	REVISION Decision = "REVISION"
)

type ResolutionEntity struct {
	Id         int
	ApproverId int
	Decision   Decision
	Comment    string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Deleted    bool
}

func NewResolutionEntity(
	approverId int,
) *ResolutionEntity {
	return &ResolutionEntity{
		ApproverId: approverId,
		Decision:   UNKNOWN,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Deleted:    false,
	}
}
