package service

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	"approve/internal/common"
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CreateFullRoute struct {
	transaction   common.Transaction
	routeRepo     rr.RouteRepository
	stepGroupRepo gr.StepGroupRepository
	stepRepo      sr.StepRepository
	approverRepo  ar.ApproverRepository
}

func NewCreateRoute(
	transaction common.Transaction,
	routeRepo rr.RouteRepository,
	groupRepo gr.StepGroupRepository,
	stepsRepo sr.StepRepository,
	approverRepo ar.ApproverRepository,
) *CreateFullRoute {
	return &CreateFullRoute{
		transaction,
		routeRepo,
		groupRepo,
		stepsRepo,
		approverRepo,
	}
}

func (svc *CreateFullRoute) Execute(request rm.CreateRouteRequest) (int64, error) {
	tx, err := svc.transaction.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	routeId, err := svc.createRoute(tx, request)
	if err != nil {
		txErr := tx.Rollback()
		return routeId, fmt.Errorf("failed to create route: %w, %w", err, txErr)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return routeId, nil
}

func (svc *CreateFullRoute) createRoute(
	tx *sqlx.Tx,
	request rm.CreateRouteRequest,
) (int64, error) {
	routeId, err := svc.routeRepo.SaveTx(
		tx,
		request.Name,
		request.Description,
	)
	if err != nil {
		return 0, err
	}
	err = svc.createStepGroups(tx, routeId, request.StepGroups)
	if err != nil {
		return 0, err
	}
	return routeId, nil
}

func (svc *CreateFullRoute) createStepGroups(
	tx *sqlx.Tx,
	routeId int64,
	groups []gm.CreateStepGroupRequest,
) error {
	toSave := make([]gm.StepGroupEntity, len(groups))
	for i, group := range groups {
		toSave[i] = group.ToEntity(routeId)
	}
	saved, err := svc.stepGroupRepo.SaveAllTxReturning(tx, toSave)
	if err != nil {
		return err
	}
	numberToSavedGroupMap := numberToSavedGroupMap(saved)
	return svc.createSteps(tx, groups, numberToSavedGroupMap)
}

func (svc *CreateFullRoute) createSteps(
	tx *sqlx.Tx,
	groups []gm.CreateStepGroupRequest,
	numberToSavedGroupMap map[int]gm.StepGroupEntity,
) error {
	var toSave []sm.StepEntity
	var steps []sm.CreateStepRequest
	for _, group := range groups {
		steps = append(steps, group.Steps...)
		entities := make([]sm.StepEntity, len(group.Steps))
		for i, step := range group.Steps {
			entities[i] = step.ToEntity(numberToSavedGroupMap[group.Number].Id)
		}
		toSave = append(toSave, entities...)
	}
	saved, err := svc.stepRepo.SaveAllTxReturning(tx, toSave)
	if err != nil {
		return err
	}
	numberToSavedStepMap := numberToSavedStepMap(saved)
	return svc.createApprovers(tx, steps, numberToSavedStepMap)
}

func (svc *CreateFullRoute) createApprovers(
	tx *sqlx.Tx,
	steps []sm.CreateStepRequest,
	numberToSavedStepMap map[int]sm.StepEntity,
) error {
	var toSave []am.ApproverEntity
	for _, step := range steps {
		entities := make([]am.ApproverEntity, len(step.Approvers))
		for i, approver := range step.Approvers {
			entities[i] = approver.ToEntity(numberToSavedStepMap[step.Number].Id)
		}
		toSave = append(toSave, entities...)
	}
	return svc.approverRepo.SaveAllTx(tx, toSave)
}

func numberToSavedGroupMap(slice []gm.StepGroupEntity) map[int]gm.StepGroupEntity {
	m := make(map[int]gm.StepGroupEntity, len(slice))
	for _, group := range slice {
		m[group.Number] = group
	}
	return m
}

func numberToSavedStepMap(slice []sm.StepEntity) map[int]sm.StepEntity {
	m := make(map[int]sm.StepEntity, len(slice))
	for _, step := range slice {
		m[step.Number] = step
	}
	return m
}
