package model

import (
	resm "approve/internal/resolution/model"
)

type ApproverEntity struct {
	Id       int64  `db:"id"`
	StepId   int64  `db:"step_id"`
	Guid     string `db:"guid"`
	Name     string `db:"name"`
	Position string `db:"position"`
	Email    string `db:"email"`
	Number   int    `db:"number"`
	Active   bool   `db:"active"`
	Deleted  bool   `db:"deleted"`
}

func (a ApproverEntity) ToResolutionEntity() resm.ResolutionEntity {
	return resm.ResolutionEntity{ApproverId: a.Id}
}
