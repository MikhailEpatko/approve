package service

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
)

func CreateApproverTemplate(
	request am.CreateApproverRequest,
) (int64, error) {
	return approverRepo.Save(request.ToEntity())
}
