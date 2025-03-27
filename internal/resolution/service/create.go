package service

import (
	approverRepo "approve/internal/approver/repository"
	cm "approve/internal/common"
	"approve/internal/database"
	resm "approve/internal/resolution/model"
	resolutionRepo "approve/internal/resolution/repository"
	ss "approve/internal/step/service"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrInvalidApproverOrder    = errors.New("invalid approver order")
	ErrCommentShouldBeProvided = errors.New("comment should be provided")
	ErrApproverNotFound        = errors.New("approver was not found")
	ErrApproverIsNotStarted    = errors.New("approver is not started")
)

func CreateResolution(
	request resm.CreateResolutionRequest,
) (resolutionId int64, err error) {
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("errors while create resolution: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	info, err := cm.SafeExecuteG(err, func() (resm.ApprovingInfoEntity, error) {
		return validateRequest(tx, request)
	})
	resolutionId, err = cm.SafeExecuteG(err, func() (int64, error) {
		return resolutionRepo.SaveTx(tx, request.ToEntity())
	})
	err = cm.SafeExecute(err, func() error { return approverRepo.FinishApprover(tx, request.ApproverId) })
	if err == nil {
		switch info.ApproverOrder {
		case cm.PARALLEL_ALL_OF, cm.SERIAL:
			err = ss.ProcessAllOfStep(tx, info, request.IsApproved)
		case cm.PARALLEL_ANY_OF:
			err = ss.ProcessAnyOfStep(tx, info, request.IsApproved)
		default:
			err = ErrInvalidApproverOrder
		}
	}
	return resolutionId, cm.ErrorOrNil("can't create resolution", err)
}

func validateRequest(
	tx *sqlx.Tx,
	request resm.CreateResolutionRequest,
) (info resm.ApprovingInfoEntity, err error) {
	err = cm.Validate(request)
	if err != nil {
		return info, cm.RequestValidationError{Message: err.Error()}
	}
	if !request.IsApproved && strings.TrimSpace(request.Comment) == "" {
		return info, ErrCommentShouldBeProvided
	}
	info, err = resolutionRepo.ApprovingInfoTx(tx, request.ApproverId)
	switch {
	case err != nil:
		break
	case info.ApproverId == 0:
		err = ErrApproverNotFound
	case info.ApproverStatus != cm.STARTED:
		err = ErrApproverIsNotStarted
	}
	return info, err
}
