package service

import (
	ar "approve/internal/approver/repository"
	"approve/internal/common"
	resm "approve/internal/resolution/model"
	rs "approve/internal/route/service"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProcessSerialStep struct {
	finishRoute            rs.FinishRoute
	approverRepo           ar.ApproverRepository
	finishStepAndStartNext FinishStepAndStartNext
}

func (svc *ProcessSerialStep) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) error {
	if !isResolutionApproved {
		return svc.finishRoute.Execute(tx, info, isResolutionApproved)
	}
	existsNotFinishedApprovers, err := svc.approverRepo.ExistNotFinishedApproversInStep(tx, info.StepId)
	if err != nil {
		return fmt.Errorf("find approver to activate error: %w", err)
	}
	if existsNotFinishedApprovers {
		if info.ApproverOrder == common.SERIAL {
			err = svc.approverRepo.StartNextApprover(tx, info.StepId, info.ApproverId)
		}
	} else {
		err = svc.finishStepAndStartNext.Execute(tx, info, isResolutionApproved)
	}
	return err
}
