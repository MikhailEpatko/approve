package service

import (
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProcessSerialStep struct {
	finishStep             FinishStepAndStartNext
	approverRepo           ar.ApproverRepository
	finishStepAndStartNext FinishStepAndStartNext
}

func (svc *ProcessSerialStep) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) error {
	if !isResolutionApproved {
		return svc.finishStep.Execute(tx, info, isResolutionApproved)
	}
	existsNotFinishedApprovers, err := svc.approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	if err != nil {
		return fmt.Errorf("find approver to activate error: %w", err)
	}
	if existsNotFinishedApprovers {
		if info.ApproverOrder == cm.SERIAL {
			err = svc.approverRepo.StartNextApprover(tx, info.StepId, info.ApproverId)
		}
	} else {
		err = svc.finishStepAndStartNext.Execute(tx, info, isResolutionApproved)
	}
	if err != nil {
		return fmt.Errorf("process serial step error: %w", err)
	}
	return nil
}
