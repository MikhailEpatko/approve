package service

import (
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	sm "approve/internal/step/model"
	gs "approve/internal/stepgroup/service"
	"github.com/jmoiron/sqlx"
)

type FinishStepRepository interface {
	FinishStep(tx *sqlx.Tx, stepId int64) error
	CalculateAndSetIsApproved(
		tx *sqlx.Tx,
		stepId int64,
		approverOrder cm.OrderType,
		isResolutionApproved bool,
	) (res bool, err error)
	ExistsNotFinishedStepsInGroup(x *sqlx.Tx, stepGroupId int64) (bool, error)
	StartNextStep(tx *sqlx.Tx, stepGroupId int64, stepId int64) (int64, error)
}

type StartApproverRepository interface {
	StartStepApprovers(tx *sqlx.Tx, step sm.StepEntity) error
}

type FinishStepAndStartNext struct {
	stepRepo                FinishStepRepository
	approverRepo            StartApproverRepository
	finishGroupAndStartNext gs.FinishStepGroupAndStartNext
}

func (svc *FinishStepAndStartNext) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isResolutionApproved bool,
) (err error) {
	err = svc.stepRepo.FinishStep(tx, info.StepId)
	isStepApproved, err := cm.SafeExecuteG(err, func() (bool, error) {
		return svc.stepRepo.CalculateAndSetIsApproved(
			tx,
			info.StepId,
			info.ApproverOrder,
			isResolutionApproved,
		)
	})
	existNotFinishedSteps, err := cm.SafeExecuteG(err, func() (bool, error) {
		return svc.stepRepo.ExistsNotFinishedStepsInGroup(tx, info.StepGroupId)
	})
	if err == nil && existNotFinishedSteps {
		if info.StepOrder == cm.SERIAL {
			var nextStepId int64
			nextStepId, err = svc.stepRepo.StartNextStep(tx, info.StepGroupId, info.StepId)
			err = cm.SafeExecute(err, func() error {
				var step = sm.StepEntity{
					Id:            nextStepId,
					ApproverOrder: info.ApproverOrder,
				}
				return svc.approverRepo.StartStepApprovers(tx, step)
			})
		}
	} else {
		err = cm.SafeExecute(err, func() error { return svc.finishGroupAndStartNext.Execute(tx, info, isStepApproved) })
	}
	return cm.ErrorOrNil("error finish step", err)
}
