package model

import (
	"approve/internal/common"
	"approve/internal/resolution/model"
)

type ApproverFullResponse struct {
	Id         int64                    `json:"id"`
	StepId     int64                    `json:"step_id"`
	Guid       string                   `json:"guid"`
	Name       string                   `json:"name"`
	Position   string                   `json:"position"`
	Email      string                   `json:"email"`
	Number     int                      `json:"number"`
	Status     common.Status            `json:"status"`
	Resolution model.ResolutionResponse `json:"resolution,omitempty"`
}
