package service

import (
	am "approve/internal/approver/model"
	ar "approve/internal/approver/repository"
	cm "approve/internal/common"
	"approve/internal/database"
	rr "approve/internal/route/repository"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	gr "approve/internal/stepgroup/repository"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strconv"
)

var ErrSourceRouteIsNotTemplate = errors.New("source route is not template")

func CreateNewRouteFromTemplate(routeTemplateId int64) (newRouteId int64, err error) {
	tx, err := database.DB.Beginx()
	defer func() {
		if err != nil {
			txErr := tx.Rollback()
			err = fmt.Errorf("failed creating new route from template: %w, %w", err, txErr)
		} else {
			err = tx.Commit()
		}
	}()

	return cm.SafeExecuteG(err, func() (int64, error) {
		return createNewRoute(tx, routeTemplateId)
	})
}

func createNewRoute(
	tx *sqlx.Tx,
	routeTemplateId int64,
) (newRouteId int64, err error) {
	template, err := rr.FindByIdTx(tx, routeTemplateId)
	if err != nil {
		err = fmt.Errorf("failed finding route template: %w", err)
	} else if template.Status != cm.TEMPLATE {
		err = ErrSourceRouteIsNotTemplate
	}
	newRouteId, err = cm.SafeExecuteG(err, func() (int64, error) {
		return rr.SaveTx(tx, template.ToNewRoute())
	})
	err = cm.SafeExecute(err, func() error {
		return createNewStepGroups(tx, routeTemplateId, newRouteId)
	})
	return newRouteId, err
}

func createNewStepGroups(
	tx *sqlx.Tx,
	routeTemplateId int64,
	newRouteId int64,
) error {
	templateStepGroups, err := gr.FindByRouteIdTx(tx, routeTemplateId)
	if err != nil {
		return fmt.Errorf("failed finding step group templates: %w", err)
	}
	if len(templateStepGroups) == 0 {
		cm.Logger.Info(
			"step group templates not found",
			zap.String("routeTemplateId", strconv.FormatInt(routeTemplateId, 10)),
		)
		return nil
	}
	numberToTemplateId := make(map[int]int64)
	newStepGroups := make([]gm.StepGroupEntity, 0, len(templateStepGroups))
	for _, stepGroup := range templateStepGroups {
		numberToTemplateId[stepGroup.Number] = stepGroup.Id
		newStepGroups = append(newStepGroups, stepGroup.ToNewStepGroup(newRouteId))
	}
	savedNewStepGroups, err := gr.SaveAll(tx, newStepGroups)
	if err != nil {
		return fmt.Errorf("failed saving new step groups (routeTemplateId = %d): %w", routeTemplateId, err)
	}
	return createNewSteps(tx, numberToTemplateId, savedNewStepGroups)
}

func createNewSteps(
	tx *sqlx.Tx,
	groupNumberToGroupTemplateId map[int]int64,
	savedNewStepGroups []gm.StepGroupEntity,
) error {
	for _, newGroup := range savedNewStepGroups {
		groupTemplateId := groupNumberToGroupTemplateId[newGroup.Number]
		stepTemplates, err := sr.FindByGroupIdTx(tx, groupTemplateId)
		if err != nil {
			return fmt.Errorf("failed finding step templates by newGroup template id: %d %w", groupTemplateId, err)
		}
		if len(stepTemplates) == 0 {
			cm.Logger.Info(
				"step templates not found",
				zap.String("groupTemplateId", strconv.FormatInt(groupTemplateId, 10)),
			)
			continue
		}
		stepNumberToStepTemplateId := make(map[int]int64)
		newSteps := make([]sm.StepEntity, 0, len(stepTemplates))
		for _, stepTemplate := range stepTemplates {
			newSteps = append(newSteps, stepTemplate.ToNewStep(newGroup.Id))
			stepNumberToStepTemplateId[stepTemplate.Number] = stepTemplate.Id
		}
		savedNewSteps, err := sr.SaveAll(tx, newSteps)
		if err != nil {
			return fmt.Errorf("failed saving new steps for new group (groupTemplateId: %d) %w", groupTemplateId, err)
		}
		err = createNewApprovers(tx, stepNumberToStepTemplateId, savedNewSteps)
		if err != nil {
			return fmt.Errorf(
				"failed creating new approvers for step group (groupTemplateId = %d): %w",
				groupTemplateId,
				err,
			)
		}
	}
	return nil
}

func createNewApprovers(
	tx *sqlx.Tx,
	stepNumberToGroupTemplateId map[int]int64,
	savedNewSteps []sm.StepEntity,
) error {
	var newApprovers []am.ApproverEntity
	for _, newStep := range savedNewSteps {
		stepTemplateId := stepNumberToGroupTemplateId[newStep.Number]
		approverTemplates, err := ar.FindByStepIdTx(tx, stepTemplateId)
		if err != nil {
			return fmt.Errorf("failed finding approver templates by stepTemplateId: %d %w", stepTemplateId, err)
		}
		if len(approverTemplates) == 0 {
			cm.Logger.Info(
				"approver templates not found",
				zap.String("stepTemplateId", strconv.FormatInt(stepTemplateId, 10)),
			)
			continue
		}
		for _, approverTemplate := range approverTemplates {
			newApprovers = append(newApprovers, approverTemplate.ToNewApprover(newStep.Id))
		}
	}
	if len(newApprovers) != 0 {
		err := ar.SaveAll(tx, newApprovers)
		if err != nil {
			return fmt.Errorf("failed saving new approvers for new steps: %w", err)
		}
	}
	return nil
}
