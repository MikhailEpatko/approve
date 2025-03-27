package service

import stepRepo "approve/internal/step/repository"

func DeleteStepById(stepId int64) (err error) {
	return stepRepo.DeleteById(stepId)
}
