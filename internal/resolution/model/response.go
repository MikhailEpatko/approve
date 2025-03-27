package model

type ResolutionResponse struct {
	Id         int64  `json:"id"`
	ApproverId int64  `json:"approver_id"`
	IsApproved bool   `json:"is_approved"`
	Comment    string `json:"comment"`
}
