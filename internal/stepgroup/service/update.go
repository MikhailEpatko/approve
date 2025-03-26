package service

import (
	cm "approve/internal/common"
	"approve/internal/database"
	rs "approve/internal/route/service"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrStepGroupNotFound = errors.New("step group not found")
)

func UpdateStepGroup(request gm.UpdateStepGroupRequest) (groupId int64, err error) {
	err = cm.Validate(request)
	if err != nil {
		return 0, fmt.Errorf("updating step group: failed validating request: %w", err)
	}
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating step group: %w, %w", err, txErr)
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
		return fmt.Errorf("failed updating step group: %w", err)
	}

	if isRouteStarted {
		return fmt.Errorf("failed updating step group: %w", rs.ErrRouteAlreadyStarted)
	}

	group, err := stepGroupRepo.FindByIdTx(tx, request.Id)
	switch {
	case err != nil:
		return err
	case group.Id == 0:
		return ErrStepGroupNotFound
	}
	return nil
}
