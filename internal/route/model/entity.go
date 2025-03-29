package model

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"
	"strings"
)

type RouteEntity struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      cm.Status `db:"status"`
	IsApproved  bool      `db:"is_approved"`
}

func (e *RouteEntity) ToResponse() RouteResponse {
	return RouteResponse{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		Status:      e.Status,
	}
}

func (e *RouteEntity) ToFullResponse(stepGroups []gm.StepGroupFullResponse) FullRouteResponse {
	return FullRouteResponse{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		Status:      e.Status,
		StepGroups:  stepGroups,
	}
}

func (e *RouteEntity) ToNewRoute() RouteEntity {
	return RouteEntity{
		Name:        removeTemplatePrefix(e.Name),
		Description: e.Description,
		Status:      cm.NEW,
	}
}

func removeTemplatePrefix(s string) string {
	res, _ := strings.CutPrefix(s, "Template ")
	return res
}
