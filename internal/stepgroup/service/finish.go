package service

import (
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	sm "approve/internal/step/model"
	"github.com/jmoiron/sqlx"
)

type FinishRouteRepository interface {
	FinishRoute(tx *sqlx.Tx, routeId int64, isGroupApproved bool) error
}

type FinishStepGroupRepository interface {
	FinishGroup(tx *sqlx.Tx, stepGroupId int64) error
	CalculateAndSetIsApproved(
		tx *sqlx.Tx,
		stepGroupId int64,
		stepOrder cm.OrderType,
		isStepApproved bool,
	) (bool, error)
	StartNextGroup(tx *sqlx.Tx, routeId int64, stepGroupId int64) (int64, error)
}

type FinishStepRepository interface {
	StartNextStep(tx *sqlx.Tx, stepGroupId int64, stepId int64) (int64, error)
}

type FinishApproverRepository interface {
	StartStepApprovers(tx *sqlx.Tx, step sm.StepEntity) error
}

type FinishStepGroupAndStartNext struct {
	routeRepo    FinishRouteRepository
	groupRepo    FinishStepGroupRepository
	stepRepo     FinishStepRepository
	approverRepo FinishApproverRepository
}

func (svc *FinishStepGroupAndStartNext) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isStepApproved bool,
) (err error) {
	err = svc.groupRepo.FinishGroup(tx, info.StepGroupId)
	isGroupApproved, err := cm.SafeExecuteG(err, func() (bool, error) {
		return svc.groupRepo.CalculateAndSetIsApproved(
			tx,
			info.StepGroupId,
			info.StepOrder,
			isStepApproved,
		)
	})
	nextGroupId, err := cm.SafeExecuteG(err, func() (int64, error) {
		return svc.groupRepo.StartNextGroup(tx, info.RouteId, info.StepGroupId)
	})
	if err == nil && nextGroupId != 0 {
		var nextStepId int64
		nextStepId, err = svc.stepRepo.StartNextStep(tx, info.StepGroupId, info.StepId)
		err = cm.SafeExecute(err, func() error {
			var step = sm.StepEntity{
				Id:            nextStepId,
				ApproverOrder: info.ApproverOrder,
			}
			return svc.approverRepo.StartStepApprovers(tx, step)
		})
	} else {
		err = cm.SafeExecute(err, func() error { return svc.routeRepo.FinishRoute(tx, info.RouteId, isGroupApproved) })
	}
	return cm.ErrorOrNil("error finish step group", err)
}
