package service

import stepGroupRepo "approve/internal/stepgroup/repository"

func DeleteStepGroupById(stepGroupId int64) (err error) {
	return stepGroupRepo.DeleteById(stepGroupId)
}
