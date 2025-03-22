package service

import (
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
)

func CreateStepTemplate(
	request sm.CreateStepRequest,
) (int64, error) {
	return stepRepo.Save(request.ToEntity())
}
