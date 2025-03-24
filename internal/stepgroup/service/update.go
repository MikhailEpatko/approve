package service

import (
	cm "approve/internal/common"
	cfg "approve/internal/database"
	rs "approve/internal/route/service"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrStepGroupNotFound       = errors.New("route not found")
	ErrStepGroupAlreadyStarted = errors.New("route is already started")
	ErrStepGroupIsFinished     = errors.New("route is finished")
)

func UpdateStepGroup(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	err = cm.Validate(request)
	tx, err := cfg.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating route: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	err = cm.SafeExecute(err, func() error { return validateStepGroup(tx, request) })
	return cm.SafeExecuteG(err, func() (int64, error) { return stepGroupRepo.Update(tx, request.ToEntity()) })
}

func validateStepGroup(
	tx *sqlx.Tx,
	request gm.UpdateStepGroupRequest,
) error {
	isRouteStarted, err := stepGroupRepo.IsRouteProcessing(tx, request.Id)
	if err != nil {
		return fmt.Errorf("error update step group: %w", err)
	}

	if isRouteStarted {
		return fmt.Errorf("error update step group: %w", rs.ErrRouteAlreadyStarted)
	}

	group, err := stepGroupRepo.FindByIdTx(tx, request.Id)
	switch {
	case err != nil:
		return err
	case group.Id == 0:
		return ErrStepGroupNotFound
	case group.Status == cm.FINISHED:
		return ErrStepGroupIsFinished
	case group.Status == cm.STARTED:
		return ErrStepGroupAlreadyStarted
	}
	return nil
}
