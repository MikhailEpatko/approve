package service

import (
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

type ProcessAllOfApproverRepository interface {
	ExistNotFinishedApproversInStep(tx *sqlx.Tx, stepId int64) (bool, error)
	StartNextApprover(tx *sqlx.Tx, stepId int64, approverId int64) error
	FinishStepApprovers(tx *sqlx.Tx, stepId int64) error
}

type ProcessAllOffStep struct {
	finishStep   FinishStepAndStartNext
	approverRepo ProcessAllOfApproverRepository
}

func (svc *ProcessAllOffStep) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	if !isResolutionApproved {
		err = svc.approverRepo.FinishStepApprovers(tx, info.StepId)
		return cm.SafeExecute(err, func() error { return svc.finishStep.Execute(tx, info, isResolutionApproved) })
	}
	existNotFinishedApprovers, err := svc.approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	if err == nil && existNotFinishedApprovers {
		if info.ApproverOrder == cm.SERIAL {
			err = svc.approverRepo.StartNextApprover(tx, info.StepId, info.ApproverId)
		}
	} else {
		err = cm.SafeExecute(err, func() error { return svc.finishStep.Execute(tx, info, isResolutionApproved) })
	}
	return cm.ErrorOrNil("process serial or parallel_all_off step error", err)
}
