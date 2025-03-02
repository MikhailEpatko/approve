package model

import "approve/internal/common"

type FilterRouteRequest struct {
	Guid   string        `json:"guid"   validate:"required"`
	Status common.Status `json:"status" validate:"required.oneof=TEMPLATE NEW STARTED FINISHED"`
	Text   string        `json:"text"`
	common.PageRequest
}
