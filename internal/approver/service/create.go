package service

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
)

type CreateApproverTemplate struct {
	approverRepo ar.ApproverRepository
}

func (svc *CreateApproverTemplate) Execute(
	request am.CreateApproverRequest,
) (int64, error) {
	return svc.approverRepo.Save(request.ToEntity())
}
