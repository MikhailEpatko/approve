package model

import (
	"approve/internal/common"
	"approve/internal/step/model"
)

type CreateStepGroupRequest struct {
	Name     string                    `json:"name"      validate:"required,min=1,max=155"`
	Number   int                       `json:"number"    validate:"required,min=1,max=20"`
	Status   common.Status             `json:"status"    validate:"required.oneof=TEMPLATE NEW STARTED FINISHED"`
	StepType StepType                  `json:"step_type" validate:"required,oneof=PARALLEL SEQUENTIAL"`
	Steps    []model.CreateStepRequest `json:"steps"     validate:"required,min=1,max=30"`
}

func (request CreateStepGroupRequest) ToEntity(routeId int64) StepGroupEntity {
	return StepGroupEntity{
		RouteId:  routeId,
		Name:     request.Name,
		Number:   request.Number,
		Status:   request.Status,
		StepType: request.StepType,
	}
}
