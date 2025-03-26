package service

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
)

func CreateApproverTemplate(
	request am.CreateApproverRequest,
) (int64, error) {
	err := cm.Validate(request)
	if err != nil {
		return 0, cm.RequestValidationError{Message: err.Error()}
	}
	return approverRepo.Save(request.ToEntity())
}
