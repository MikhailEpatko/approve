package model

import (
	"approve/internal/common"
	"approve/internal/stepgroup/model"
)

type RouteResponse struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      common.Status `json:"status"`
}

type FullRouteResponse struct {
	Id          int64                         `json:"id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	Status      common.Status                 `json:"status"`
	StepGroups  []model.StepGroupFullResponse `json:"groups"`
}
