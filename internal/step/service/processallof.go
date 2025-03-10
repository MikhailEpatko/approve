package service

import (
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

type ProcessAllOffStep struct {
	finishStep   FinishStepAndStartNext
	approverRepo ar.ApproverRepository
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
	existNotFinishedApprovers, err := cm.SafeExecuteBool(nil, func() (bool, error) {
		return svc.approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	})
	if err == nil && existNotFinishedApprovers {
		if info.ApproverOrder == cm.SERIAL {
			err = svc.approverRepo.StartNextApprover(tx, info.StepId, info.ApproverId)
		}
	} else {
		err = svc.finishStep.Execute(tx, info, isResolutionApproved)
	}
	return cm.ErrorOrNil("process serial step error", err)
}
