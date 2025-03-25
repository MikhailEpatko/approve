package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	"approve/internal/database"
	routeRepo "approve/internal/route/repository"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRouteNotFound       = errors.New("route not found")
	ErrRouteAlreadyStarted = errors.New("route is already started")
	ErrRouteIsFinished     = errors.New("route is finished")
)

func StartRoute(routeId int64) (err error) {
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed start route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	err = cm.SafeExecute(err, func() error {
		route, innerErr := routeRepo.FindByIdTx(tx, routeId)
		switch {
		case innerErr != nil:
			return innerErr
		case route.Id == 0:
			return ErrRouteNotFound
		case route.Status == cm.FINISHED:
			return ErrRouteIsFinished
		case route.Status == cm.STARTED:
			return ErrRouteAlreadyStarted
		}
		return nil
	})
	return cm.SafeExecute(err, func() error { return startRote(tx, routeId) })
}

func startRote(
	tx *sqlx.Tx,
	routeId int64,
) (err error) {
	err = routeRepo.StartRoute(tx, routeId)
	return cm.SafeExecute(err, func() error { return startGroups(tx, routeId) })
}

func startGroups(
	tx *sqlx.Tx,
	routeId int64,
) error {
	group, err := stepGroupRepo.StartFirstGroup(tx, routeId)
	if err == nil && group.Id > 0 {
		err = startSteps(tx, group)
	}
	return err
}

func startSteps(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) error {
	steps, err := stepRepo.StartSteps(tx, group)
	if err == nil && len(steps) > 0 {
		for _, step := range steps {
			err = startApprovers(tx, step)
			if err != nil {
				break
			}
		}
	}
	return err
}

func startApprovers(
	tx *sqlx.Tx,
	step sm.StepEntity,
) error {
	return approverRepo.StartStepApprovers(tx, step)
}
