package service

import (
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	resm "approve/internal/resolution/model"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gr "approve/internal/stepgroup/repository"
	"github.com/jmoiron/sqlx"
)

type FinishStepGroupAndStartNext struct {
	routeRepo    rr.RouteRepository
	groupRepo    gr.StepGroupRepository
	stepRepo     sr.StepRepository
	approverRepo ar.ApproverRepository
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
