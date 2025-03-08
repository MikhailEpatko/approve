package service

import (
	ar "approve/internal/approver/repository"
	"approve/internal/common"
	resm "approve/internal/resolution/model"
	resr "approve/internal/resolution/repository"
	service2 "approve/internal/step/service"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	ERROR_CREATE_RESOLUTION_W = "can't create resolution. Cause: %w"
)

type CreateResolution struct {
	transaction           common.Transaction
	approverRepo          ar.ApproverRepository
	resolutionRepo        resr.ResolutionRepository
	processeParallelAllOf service2.ProcessParallelAllOf
	processParallelAnyOf  service2.ProcessParallelAnyOf
	processSerial         service2.ProcessSerialStep
}

func (svc *CreateResolution) CreateResolution(
	request resm.CreateResolutionRequest,
) (resolutionId int64, err error) {
	tx, err := svc.transaction.Begin()
	if err != nil {
		return resolutionId, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("errors while create resolution: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()
	info, err := svc.validateRequest(tx, request)
	if err != nil {
		return 0, fmt.Errorf(ERROR_CREATE_RESOLUTION_W, err)
	}
	resolutionId, err = svc.resolutionRepo.SaveTx(tx, request.ToEntity())
	if err != nil {
		return 0, fmt.Errorf(ERROR_CREATE_RESOLUTION_W, err)
	}
	err = svc.approverRepo.DeativateTx(tx, request.ApproverId)
	if err != nil {
		return 0, fmt.Errorf(ERROR_CREATE_RESOLUTION_W, err)
	}
	switch info.ApproverOrder {
	case common.PARALLEL_ALL_OF:
		err = svc.processeParallelAllOf.Execute(tx, info, request.IsApproved)
	case common.PARALLEL_ANY_OF:
		err = svc.processParallelAnyOf.Execute(tx, info, request.IsApproved)
	case common.SERIAL:
		err = svc.processSerial.Execute(tx, info, request.IsApproved)
	}
	if err != nil {
		return 0, fmt.Errorf(ERROR_CREATE_RESOLUTION_W, err)
	}
	return resolutionId, nil

}

func (svc *CreateResolution) validateRequest(
	tx *sqlx.Tx,
	request resm.CreateResolutionRequest,
) (info resm.ApprovingInfoEntity, err error) {
	if !request.IsApproved && strings.TrimSpace(request.Comment) == "" {
		return info, errors.New("comment should be provided")
	}
	info, err = svc.resolutionRepo.ApprovingInfoTx(tx, request.ApproverId)
	switch {
	case err != nil:
		break
	case info.Guid == "":
		err = errors.New("approver was not found")
	case !info.Active:
		err = errors.New("approver is not active")
	case info.StepStatus != common.STARTED:
		err = errors.New("step is not started")
	}

	//TODO: check if requester is approver (has approver's guid in jwt token)

	return info, err
}
