package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	routeRepo "approve/internal/route/repository"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"github.com/jmoiron/sqlx"
)

func FinishStepGroupAndStartNext(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isStepApproved bool,
) (err error) {
	err = stepGroupRepo.FinishGroup(tx, info.StepGroupId)
	isGroupApproved, err := cm.SafeExecuteG(err, func() (bool, error) {
		return stepGroupRepo.CalculateAndSetIsApproved(
			tx,
			info.StepGroupId,
			info.StepOrder,
			isStepApproved,
		)
	})
	nextGroupId, err := cm.SafeExecuteG(err, func() (int64, error) {
		return stepGroupRepo.StartNextGroup(tx, info.RouteId, info.StepGroupId)
	})
	if err == nil && nextGroupId != 0 {
		var nextStepId int64
		nextStepId, err = stepRepo.StartNextStep(tx, info.StepGroupId, info.StepId)
		err = cm.SafeExecute(err, func() error {
			var step = sm.StepEntity{
				Id:            nextStepId,
				ApproverOrder: info.ApproverOrder,
			}
			return approverRepo.StartStepApprovers(tx, step)
		})
	} else {
		err = cm.SafeExecute(err, func() error { return routeRepo.FinishRoute(tx, info.RouteId, isGroupApproved) })
	}
	return cm.ErrorOrNil("error finish step group", err)
}
