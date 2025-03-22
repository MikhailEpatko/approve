package service

import (
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
)

func CreateStepGroupTemplate(
	request gm.CreateStepGroupRequest,
) (int64, error) {
	return stepGroupRepo.Save(request.ToEntity())
}
