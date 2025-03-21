package service

import (
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

type ProcessAnyOfApproverRepository interface {
	ExistNotFinishedApproversInStep(tx *sqlx.Tx, stepId int64) (bool, error)
	FinishStepApprovers(tx *sqlx.Tx, stepId int64) error
}

type ProcessAnyOfStep struct {
	approverRepo ProcessAnyOfApproverRepository
	finishStep   FinishStepAndStartNext
}

func (svc *ProcessAnyOfStep) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	if isResolutionApproved {
		err = svc.approverRepo.FinishStepApprovers(tx, info.StepId)
		return cm.SafeExecute(err, func() error { return svc.finishStep.Execute(tx, info, isResolutionApproved) })
	}
	existNotFinishedApprovers, err := svc.approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	if err == nil && !existNotFinishedApprovers {
		err = svc.finishStep.Execute(tx, info, isResolutionApproved)
	}
	return cm.ErrorOrNil("process parallel_any_of step error", err)
}
