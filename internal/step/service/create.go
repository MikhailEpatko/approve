package service

import (
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
)

type CreateStepTemplate struct {
	stepRepo sr.StepRepository
}

func (svc *CreateStepTemplate) Execute(
	request sm.CreateStepRequest,
) (int64, error) {
	return svc.stepRepo.Save(request.ToEntity())
}
