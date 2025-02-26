package model

import (
	"time"
)

type ApproveType string

const (
	PARALLEL_ANY_OF    ApproveType = "PARALLEL_ANY_OF"
	PARALLEL_ALL_OF    ApproveType = "PARALLEL_ALL_OF"
	SEQUENTIAL_ALL_OFF ApproveType = "SEQUENTIAL_ALL_OFF"
)

type StepEntity struct {
	Id          int
	StepGroupId int
	Name        string
	Number      int
	Status      Status
	ApproveType ApproveType
	CreatedAt   time.Time
	ModifiedAt  time.Time
	Deleted     bool
}

func NewStepEntity(
	stepGroupId int,
	name string,
	number int,
	approveType ApproveType,
) *StepEntity {
	return &StepEntity{
		StepGroupId: stepGroupId,
		Name:        name,
		Number:      number,
		Status:      NEW,
		ApproveType: approveType,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Deleted:     false,
	}
}
