package service

import (
	cm "approve/internal/common"
	cfg "approve/internal/config"
	resm "approve/internal/resolution/model"
	ss "approve/internal/step/service"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type FinishApproverRepository interface {
	FinishApprover(tx *sqlx.Tx, approverId int64) error
}

type CreateResolutionRepository interface {
	SaveTx(tx *sqlx.Tx, resolution resm.ResolutionEntity) (int64, error)
	ApprovingInfoTx(tx *sqlx.Tx, approverId int64) (resm.ApprovingInfoEntity, error)
}

type CreateResolution struct {
	transaction    cfg.Transaction
	approverRepo   FinishApproverRepository
	resolutionRepo CreateResolutionRepository
	processAnyOf   ss.ProcessAnyOfStep
	processAllOf   ss.ProcessAllOffStep
}

var (
	ErrInvalidApproverOrder    = errors.New("invalid approver order")
	ErrCommentShouldBeProvided = errors.New("comment should be provided")
	ErrApproverNotFound        = errors.New("approver was not found")
	ErrApproverIsNotStarted    = errors.New("approver is not started")
	ErrStepIsNotStarted        = errors.New("step is not started")
)

func (svc *CreateResolution) CreateResolution(
	request resm.CreateResolutionRequest,
) (resolutionId int64, err error) {
	tx, err := svc.transaction.Begin()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("errors while create resolution: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	info, err := cm.SafeExecuteG(err, func() (resm.ApprovingInfoEntity, error) {
		return svc.validateRequest(tx, request)
	})
	resolutionId, err = cm.SafeExecuteG(err, func() (int64, error) {
		return svc.resolutionRepo.SaveTx(tx, request.ToEntity())
	})
	err = cm.SafeExecute(err, func() error { return svc.approverRepo.FinishApprover(tx, request.ApproverId) })
	if err == nil {
		switch info.ApproverOrder {
		case cm.PARALLEL_ALL_OF, cm.SERIAL:
			err = svc.processAllOf.Execute(tx, info, request.IsApproved)
		case cm.PARALLEL_ANY_OF:
			err = svc.processAnyOf.Execute(tx, info, request.IsApproved)
		default:
			err = ErrInvalidApproverOrder
		}
	}
	return resolutionId, cm.ErrorOrNil("can't create resolution", err)
}

func (svc *CreateResolution) validateRequest(
	tx *sqlx.Tx,
	request resm.CreateResolutionRequest,
) (info resm.ApprovingInfoEntity, err error) {
	// TODO: check if requester is approver (has approver's guid in jwt token)
	if !request.IsApproved && strings.TrimSpace(request.Comment) == "" {
		return info, ErrCommentShouldBeProvided
	}
	info, err = svc.resolutionRepo.ApprovingInfoTx(tx, request.ApproverId)
	switch {
	case err != nil:
		break
	case info.Guid == "":
		err = ErrApproverNotFound
	case info.ApproverStatus != cm.STARTED:
		err = ErrApproverIsNotStarted
	case info.StepStatus != cm.STARTED:
		err = ErrStepIsNotStarted
	}
	return info, err
}
