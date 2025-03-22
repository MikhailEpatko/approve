package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

func ProcessAllOfStep(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	if !isResolutionApproved {
		err = approverRepo.FinishStepApprovers(tx, info.StepId)
		return cm.SafeExecute(err, func() error { return FinishStepAndStartNext(tx, info, isResolutionApproved) })
	}
	existNotFinishedApprovers, err := approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	if err == nil && existNotFinishedApprovers {
		if info.ApproverOrder == cm.SERIAL {
			err = approverRepo.StartNextApprover(tx, info.StepId, info.ApproverId)
		}
	} else {
		err = cm.SafeExecute(err, func() error { return FinishStepAndStartNext(tx, info, isResolutionApproved) })
	}
	return cm.ErrorOrNil("process serial or parallel_all_off step error", err)
}
