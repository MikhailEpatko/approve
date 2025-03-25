package service

import (
	cm "approve/internal/common"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
)

func CreateStepTemplate(
	request sm.CreateStepRequest,
) (int64, error) {
	err := cm.Validate(request)
	if err != nil {
		return 0, cm.RequestValidationError{Message: err.Error()}
	}
	return stepRepo.Save(request.ToEntity())
}
