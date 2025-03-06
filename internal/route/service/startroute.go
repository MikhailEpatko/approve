package service

import (
	ar "approve/internal/approver/repository"
	"approve/internal/common"
	resr "approve/internal/resolution/repository"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type StartRoute struct {
	transaction    common.Transaction
	routeRepo      rr.RouteRepository
	stepGroupRepo  gr.StepGroupRepository
	stepRepo       sr.StepRepository
	approverRepo   ar.ApproverRepository
	resolutionRepo resr.ResolutionRepository
}

func NewStartRoute(
	transaction common.Transaction,
	routeRepo rr.RouteRepository,
	stepGroupRepo gr.StepGroupRepository,
	stepRepo sr.StepRepository,
	approverRepo ar.ApproverRepository,
	resolutionRepo resr.ResolutionRepository,
) *StartRoute {
	return &StartRoute{
		transaction:    transaction,
		routeRepo:      routeRepo,
		stepGroupRepo:  stepGroupRepo,
		stepRepo:       stepRepo,
		approverRepo:   approverRepo,
		resolutionRepo: resolutionRepo,
	}
}

func (svc *StartRoute) Execute(routeId int64) (err error) {
	tx, err := svc.transaction.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed start route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	err = svc.startRote(tx, routeId)
	return err
}

func (svc *StartRoute) startRote(
	tx *sqlx.Tx,
	routeId int64,
) error {
	err := svc.routeRepo.StartRouteTx(tx, routeId)
	if err == nil {
		err = svc.stargGroups(tx, routeId)
	}
	return err
}

func (svc *StartRoute) stargGroups(
	tx *sqlx.Tx,
	routeId int64,
) error {
	group, err := svc.stepGroupRepo.StartGroupsTx(tx, routeId)
	if err == nil && group.Id > 0 {
		err = svc.startSteps(tx, group)
	}
	return err
}

func (svc *StartRoute) startSteps(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) error {
	steps, err := svc.stepRepo.StartStepsTx(tx, group)
	if err == nil && len(steps) > 0 {
		for _, step := range steps {
			err = svc.startApprovers(tx, step)
			if err != nil {
				break
			}
		}
	}
	return err
}

func (svc *StartRoute) startApprovers(
	tx *sqlx.Tx,
	step sm.StepEntity,
) error {
	return svc.approverRepo.StartApproversTx(tx, step)
}
