package model

type CreateApproverRequest struct {
	StepId   int64  `json:"step_id"  validate:"required,min=1"`
	Guid     string `json:"guid"     validate:"required,max=36"`
	Name     string `json:"name"     validate:"required,max=150"`
	Position string `json:"position" validate:"required,max=150"`
	Email    string `json:"email"    validate:"required,max=50,email"`
	Number   int    `json:"number"   validate:"required,min=1,max=20"`
}

func (r CreateApproverRequest) ToEntity() ApproverEntity {
	return ApproverEntity{
		StepId:   r.StepId,
		Guid:     r.Guid,
		Name:     r.Name,
		Position: r.Position,
		Email:    r.Email,
		Number:   r.Number,
	}
}

type UpdateApproverRequest struct {
	Id       int64  `json:"id"       validate:"required,min=1"`
	Guid     string `json:"guid"     validate:"required,max=36"`
	Name     string `json:"name"     validate:"required,max=150"`
	Position string `json:"position" validate:"required,max=150"`
	Email    string `json:"email"    validate:"required,max=50,email"`
	Number   int    `json:"number"   validate:"required,min=1,max=20"`
}

func (r UpdateApproverRequest) ToEntity() ApproverEntity {
	return ApproverEntity{
		Id:       r.Id,
		Guid:     r.Guid,
		Name:     r.Name,
		Position: r.Position,
		Email:    r.Email,
		Number:   r.Number,
	}
}
