package service

import (
	am "approve/internal/approver/model"
)

type SaveApproverRepository interface {
	Save(approver am.ApproverEntity) (int64, error)
}

type CreateApproverTemplate struct {
	approverRepo SaveApproverRepository
}

func (svc *CreateApproverTemplate) Execute(
	request am.CreateApproverRequest,
) (int64, error) {
	return svc.approverRepo.Save(request.ToEntity())
}
