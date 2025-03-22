package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	gs "approve/internal/stepgroup/service"
	"github.com/jmoiron/sqlx"
)

func FinishStepAndStartNext(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	err = stepRepo.FinishStep(tx, info.StepId)
	isStepApproved, err := cm.SafeExecuteG(err, func() (bool, error) {
		return stepRepo.CalculateAndSetIsApproved(
			tx,
			info.StepId,
			info.ApproverOrder,
			isResolutionApproved,
		)
	})
	existNotFinishedSteps, err := cm.SafeExecuteG(err, func() (bool, error) {
		return stepRepo.ExistsNotFinishedStepsInGroup(tx, info.StepGroupId)
	})
	if err == nil && existNotFinishedSteps {
		if info.StepOrder == cm.SERIAL {
			var nextStepId int64
			nextStepId, err = stepRepo.StartNextStep(tx, info.StepGroupId, info.StepId)
			err = cm.SafeExecute(err, func() error {
				var step = sm.StepEntity{
					Id:            nextStepId,
					ApproverOrder: info.ApproverOrder,
				}
				return approverRepo.StartStepApprovers(tx, step)
			})
		}
	} else {
		err = cm.SafeExecute(err, func() error { return gs.FinishStepGroupAndStartNext(tx, info, isStepApproved) })
	}
	return cm.ErrorOrNil("error finish step", err)
}
