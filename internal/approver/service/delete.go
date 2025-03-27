package service

import stepRepo "approve/internal/approver/repository"

func DeleteApproverById(approverId int64) (err error) {
	return stepRepo.DeleteById(approverId)
}
