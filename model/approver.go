package model

import "time"

type ApproverEntity struct {
	Id         int
	StepId     int
	Guid       string
	Name       string
	Email      string
	Number     int
	CreatedAt  time.Time
	ModifiedAt time.Time
	Deleted    bool
}

func NewApproverEntity(
	stepId int,
	guid string,
	name string,
	email string,
	number int,
) *ApproverEntity {
	return &ApproverEntity{
		StepId:     stepId,
		Guid:       guid,
		Name:       name,
		Email:      email,
		Number:     number,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Deleted:    false,
	}
}
