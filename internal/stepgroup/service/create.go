package service

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
)

func CreateStepGroupTemplate(
	request gm.CreateStepGroupRequest,
) (int64, error) {
	err := cm.Validate(request)
	if err != nil {
		return 0, cm.RequestValidationError{Message: err.Error()}
	}
	return stepGroupRepo.Save(request.ToEntity())
}
