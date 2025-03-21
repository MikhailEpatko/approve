package service

import (
	cm "approve/internal/common"
	cfg "approve/internal/config"
	sm "approve/internal/step/model"
	gm "approve/internal/stepgroup/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type StartRouteRepository interface {
	StartRoute(tx *sqlx.Tx, id int64) error
}

type StartFirstGroupRepository interface {
	StartFirstGroup(tx *sqlx.Tx, routeId int64) (gm.StepGroupEntity, error)
}

type StartStepRepository interface {
	StartSteps(tx *sqlx.Tx, group gm.StepGroupEntity) ([]sm.StepEntity, error)
}

type StartApproverRepository interface {
	StartStepApprovers(tx *sqlx.Tx, step sm.StepEntity) error
}

type StartRoute struct {
	transaction   cfg.Transaction
	routeRepo     StartRouteRepository
	stepGroupRepo StartFirstGroupRepository
	stepRepo      StartStepRepository
	approverRepo  StartApproverRepository
}

func (svc *StartRoute) Execute(routeId int64) (err error) {
	tx, err := svc.transaction.Begin()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed start route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	return cm.SafeExecute(err, func() error { return svc.startRote(tx, routeId) })
}

func (svc *StartRoute) startRote(
	tx *sqlx.Tx,
	routeId int64,
) error {
	err := svc.routeRepo.StartRoute(tx, routeId)
	return cm.SafeExecute(err, func() error { return svc.stargGroups(tx, routeId) })
}

func (svc *StartRoute) stargGroups(
	tx *sqlx.Tx,
	routeId int64,
) error {
	group, err := svc.stepGroupRepo.StartFirstGroup(tx, routeId)
	if err == nil && group.Id > 0 {
		err = svc.startSteps(tx, group)
	}
	return err
}

func (svc *StartRoute) startSteps(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) error {
	steps, err := svc.stepRepo.StartSteps(tx, group)
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
	return svc.approverRepo.StartStepApprovers(tx, step)
}
