package service

import (
	cm "approve/internal/common"
	"approve/internal/database"
	rs "approve/internal/route/service"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrStepNotFound = errors.New("step not found")
)

func UpdateStep(request sm.UpdateStepRequest) (routeId int64, err error) {
	err = cm.Validate(request)
	if err != nil {
		return 0, fmt.Errorf("updating step: failed validating request: %w", err)
	}
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating step: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	err = cm.SafeExecute(err, func() error { return validateStep(tx, request) })
	return cm.SafeExecuteG(err, func() (int64, error) { return stepRepo.Update(tx, request.ToEntity()) })
}

func validateStep(
	tx *sqlx.Tx,
	request sm.UpdateStepRequest,
) error {
	IsRouteProcessing, err := stepRepo.IsRouteProcessing(tx, request.Id)
	if err != nil {
		return fmt.Errorf("failed updating step: %w", err)
	}

	if IsRouteProcessing {
		return fmt.Errorf("failed updating step: %w", rs.ErrRouteAlreadyStarted)
	}

	step, err := stepRepo.FindByIdTx(tx, request.Id)
	switch {
	case err != nil:
		return err
	case step.Id == 0:
		return ErrStepNotFound
	}
	return nil
}
