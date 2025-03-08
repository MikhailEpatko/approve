package service

import (
	ar "approve/internal/approver/repository"
	resm "approve/internal/resolution/model"
	rr "approve/internal/route/repository"
	sr "approve/internal/step/repository"
	gr "approve/internal/stepgroup/repository"
	"github.com/jmoiron/sqlx"
)

type FinishRoute struct {
	routeRepo     rr.RouteRepository
	stepGroupRepo gr.StepGroupRepository
	stepRepo      sr.StepRepository
	approverRepo  ar.ApproverRepository
}

func (svc FinishRoute) Execute(
	tx *sqlx.Tx,
	info resm.ApprovingInfoEntity,
	isRouteApproved bool,
) (err error) {
	err = svc.approverRepo.DeactivateApproversByRouteId(tx, info.RouteId)
	if err != nil {
		return err
	}
	err = svc.stepRepo.FinishStepsByRouteId(tx, info.RouteId)
	if err != nil {
		return err
	}
	err = svc.stepGroupRepo.FinishGroupsByRouteId(tx, info.RouteId)
	if err != nil {
		return err
	}
	return svc.routeRepo.FinishRoute(tx, info.RouteId, isRouteApproved)
}
