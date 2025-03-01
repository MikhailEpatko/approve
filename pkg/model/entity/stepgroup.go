package entity

import (
	"approve/internal/common"
)

type StepType string

const (
	PARALLEL   StepType = "PARALLEL"
	SEQUENTIAL StepType = "SEQUENTIAL"
)

type StepGroupEntity struct {
	Id       int64         `db:"id"`
	RouteId  int64         `db:"route_id"`
	Name     string        `db:"name"`
	Number   int           `db:"number"`
	Status   common.Status `db:"status"`
	StepType StepType      `db:"step_type"`
	Deleted  bool          `db:"deleted"`
}

func NewStepGroupEntity(
	routeId int64,
	name string,
	number int,
	stepType StepType,
) *StepGroupEntity {
	return &StepGroupEntity{
		RouteId:  routeId,
		Name:     name,
		Number:   number,
		Status:   common.NEW,
		StepType: stepType,
		Deleted:  false,
	}
}
