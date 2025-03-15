package model

import cm "approve/internal/common"

type FilterRouteRequest struct {
	Guid           string    `json:"guid"   validate:"required"`
	Status         cm.Status `json:"status" validate:"required,oneof=TEMPLATE NEW STARTED FINISHED"`
	Text           string    `json:"text"`
	cm.PageRequest `json:"page"   validate:"required"`
}
