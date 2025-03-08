package model

type CreateResolutionRequest struct {
	ApproverId int64  `json:"approver_id" validate:"required,min=1"`
	IsApproved bool   `json:"is_approved" validate:"required"`
	Comment    string `json:"comment"`
}

func (r *CreateResolutionRequest) ToEntity() ResolutionEntity {
	return ResolutionEntity{
		ApproverId: r.ApproverId,
		IsApproved: r.IsApproved,
		Comment:    r.Comment,
	}
}
