package model

import (
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
)

type ApproverEntity struct {
	Id       int64     `db:"id"`
	StepId   int64     `db:"step_id"`
	Guid     string    `db:"guid"`
	Name     string    `db:"name"`
	Position string    `db:"position"`
	Email    string    `db:"email"`
	Number   int       `db:"number"`
	Status   cm.Status `db:"status"`
}

func (e *ApproverEntity) ToFullResponse(resolution resm.ResolutionResponse) ApproverFullResponse {
	return ApproverFullResponse{
		Id:         e.Id,
		StepId:     e.StepId,
		Guid:       e.Guid,
		Name:       e.Name,
		Position:   e.Position,
		Email:      e.Email,
		Number:     e.Number,
		Status:     e.Status,
		Resolution: resolution,
	}
}
