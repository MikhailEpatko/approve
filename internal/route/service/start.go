package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	cfg "approve/internal/config"
	routeRepo "approve/internal/route/repository"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func StartRoute(routeId int64) (err error) {
	tx, err := cfg.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed start route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	return cm.SafeExecute(err, func() error { return startRote(tx, routeId) })
}

func startRote(
	tx *sqlx.Tx,
	routeId int64,
) error {
	err := routeRepo.StartRoute(tx, routeId)
	return cm.SafeExecute(err, func() error { return stargGroups(tx, routeId) })
}

func stargGroups(
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
