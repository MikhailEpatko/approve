package service

import (
	sm "approve/internal/step/model"
)

type CreateStepRepository interface {
	Save(step sm.StepEntity) (int64, error)
}

type CreateStepTemplate struct {
	stepRepo CreateStepRepository
}

func (svc *CreateStepTemplate) Execute(
	request sm.CreateStepRequest,
) (int64, error) {
	return svc.stepRepo.Save(request.ToEntity())
}
