package model

import (
	"approve/internal/common"
	sm "approve/internal/step/model"
)

type StepGroupEntity struct {
	Id         int64            `db:"id"`
	RouteId    int64            `db:"route_id"`
	Name       string           `db:"name"`
	Number     int              `db:"number"`
	Status     common.Status    `db:"status"`
	StepOrder  common.OrderType `db:"step_order"`
	IsApproved bool             `db:"is_approved"`
}

func (e *StepGroupEntity) ToFullResponse(steps []sm.StepFullResponse) StepGroupFullResponse {
	return StepGroupFullResponse{
		Id:         e.Id,
		RouteId:    e.RouteId,
		Name:       e.Name,
		Number:     e.Number,
		Status:     e.Status,
		StepOrder:  e.StepOrder,
		IsApproved: e.IsApproved,
		Steps:      steps,
	}
}

func (e *StepGroupEntity) ToNewStepGroup(routeId int64) StepGroupEntity {
	return StepGroupEntity{
		RouteId:   routeId,
		Name:      e.Name,
		Number:    e.Number,
		Status:    common.NEW,
		StepOrder: e.StepOrder,
	}
}
