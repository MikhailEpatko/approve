package model

import (
	"approve/internal/common"
	"approve/internal/stepgroup/model"
)

type RouteRequest struct {
	Name        string        `json:"name"         validate:"required,min=3,max=155"`
	Description string        `json:"description"  validate:"max=255"`
	Status      common.Status `json:"status"       validate:"required,oneof=TEMPLATE NEW"`
}

func (r RouteRequest) ToEntity(routeId int64) RouteEntity {
	return RouteEntity{
		Id:          routeId,
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
	}
}

type CreateRouteRequest struct {
	RouteRequest
	StepGroups []model.CreateStepGroupRequest `json:"step_groups" validate:"required,min=1,max=10"`
}

func (r CreateRouteRequest) ToEntity() RouteEntity {
	return RouteEntity{
		Name:        r.Name,
		Description: r.Description,
		Status:      r.Status,
	}
}
