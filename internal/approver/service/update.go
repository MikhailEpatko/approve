package service

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	"approve/internal/database"
	rs "approve/internal/route/service"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	ErrApproverNotFound = errors.New("approver not found")
)

func UpdateApprover(request am.UpdateApproverRequest) (routeId int64, err error) {
	err = cm.Validate(request)
	if err != nil {
		return 0, fmt.Errorf("updating step: failed validating request: %w", err)
	}
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed updating approver: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	err = cm.SafeExecute(err, func() error { return validateStep(tx, request) })
	return cm.SafeExecuteG(err, func() (int64, error) { return approverRepo.Update(tx, request.ToEntity()) })
}

func validateStep(
	tx *sqlx.Tx,
	request am.UpdateApproverRequest,
) error {
	IsRouteProcessing, err := approverRepo.IsRouteProcessing(tx, request.Id)
	if err != nil {
		return fmt.Errorf("failed updating approver: %w", err)
	}

	if IsRouteProcessing {
		return fmt.Errorf("failed updating approver: %w", rs.ErrRouteAlreadyStarted)
	}

	step, err := approverRepo.FindByIdTx(tx, request.Id)
	switch {
	case err != nil:
		return err
	case step.Id == 0:
		return ErrApproverNotFound
	}
	return nil
}
