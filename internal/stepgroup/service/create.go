package service

import (
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
)

type CreateStepGroupTemplate struct {
	repo gr.StepGroupRepository
}

func (svc *CreateStepGroupTemplate) Execute(
	request gm.CreateStepGroupRequest,
) (int64, error) {
	return svc.repo.Save(request.ToEntity())
}
