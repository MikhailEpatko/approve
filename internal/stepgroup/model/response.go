package model

import (
	"approve/internal/common"
	"approve/internal/step/model"
)

type StepGroupFullResponse struct {
	Id         int64                    `json:"id"`
	RouteId    int64                    `json:"route_id"`
	Name       string                   `json:"name"`
	Number     int                      `json:"number"`
	Status     common.Status            `json:"status"`
	StepOrder  common.OrderType         `json:"step_order"`
	IsApproved bool                     `json:"is_approved"`
	Steps      []model.StepFullResponse `json:"steps"`
}
