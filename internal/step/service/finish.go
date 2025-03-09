package service

import (
	ar "approve/internal/approver/repository"
	"approve/internal/common"
	resm "approve/internal/resolution/model"
	sr "approve/internal/step/repository"
	gs "approve/internal/stepgroup/service"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type FinishStepAndStartNext struct {
	stepRepo                sr.StepRepository
	approverRepo            ar.ApproverRepository
	finishGroupAndStartNext gs.FinishGroupAndStartNext
}

func (svc *FinishStepAndStartNext) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	err = svc.stepRepo.FinishStep(tx, info.StepId)
	if err != nil {
		return fmt.Errorf("can't finish step: %w", err)
	}
	isStepApproved, err := svc.stepRepo.CalculateAndSetIsApproved(
		tx,
		info.StepId,
		info.ApproverOrder,
		isResolutionApproved,
	)
	if err != nil {
		return fmt.Errorf("can't calculate and set step.is_approved: %w", err)
	}
	existNotFinishedSteps, err := svc.stepRepo.ExistsNotFinishedStepsInGroup(tx, info.StepGroupId)
	if err != nil {
		return fmt.Errorf("can't check if not finished steps in group exists: %w", err)
	}
	if existNotFinishedSteps {
		if info.StepOrder == common.SERIAL {
			nextStepId, err := svc.stepRepo.StartNextStepTx(tx, info.StepGroupId, info.StepId)
			if err != nil {
				return fmt.Errorf("can't start next step: %w", err)
			}
			err = svc.approverRepo.StartApproversTx(tx, nextStepId)
		}
	} else {
		err = svc.finishGroupAndStartNext.Execute(tx, info, isStepApproved)
	}
	return err
}
