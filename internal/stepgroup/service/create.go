package service

import (
	gm "approve/internal/stepgroup/model"
)

type StepGroupRepository interface {
	Save(stepGroup gm.StepGroupEntity) (int64, error)
}

type CreateStepGroupTemplate struct {
	repo StepGroupRepository
}

func (svc *CreateStepGroupTemplate) Execute(
	request gm.CreateStepGroupRequest,
) (int64, error) {
	return svc.repo.Save(request.ToEntity())
}
