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
	isStepApproved, err := common.SafeExecuteBool(err, func() (bool, error) {
		return svc.stepRepo.CalculateAndSetIsApproved(
			tx,
			info.StepId,
			info.ApproverOrder,
			isResolutionApproved,
		)
	})
	existNotFinishedSteps, err := common.SafeExecuteBool(err, func() (bool, error) {
		return svc.stepRepo.ExistsNotFinishedStepsInGroup(tx, info.StepGroupId)
	})
	if err == nil && existNotFinishedSteps {
		if info.StepOrder == common.SERIAL {
			var nextStepId int64
			nextStepId, err = svc.stepRepo.StartNextStepTx(tx, info.StepGroupId, info.StepId)
			err = common.SafeExecute(err, func() error { return svc.approverRepo.StartApproversTx(tx, nextStepId) })
		}
	} else {
		err = common.SafeExecute(err, func() error { return svc.finishGroupAndStartNext.Execute(tx, info, isStepApproved) })
	}
	if err != nil {
		return fmt.Errorf("error finish step: %w", err)
	}
	return nil
}
