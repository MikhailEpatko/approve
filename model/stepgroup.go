package model

import "time"

type StepGroupStatus string
type StepType string

const (
	PARALLEL   StepType = "PARALLEL"
	SEQUENTIAL StepType = "SEQUENTIAL"
)

type StepGroup struct {
	Id         int
	RouteId    int
	Name       string
	Number     int
	Status     Status
	StepType   StepType
	CreatedAt  time.Time
	ModifiedAt time.Time
	Deleted    bool
}

func NewStepGroup(
	routeId int,
	name string,
	number int,
	stepType StepType,
) *StepGroup {
	return &StepGroup{
		RouteId:    routeId,
		Name:       name,
		Number:     number,
		Status:     NEW,
		StepType:   stepType,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Deleted:    false,
	}
}
