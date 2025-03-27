package service

import (
	am "approve/internal/approver/model"
	approverRepo "approve/internal/approver/repository"
	"approve/internal/common"
	resm "approve/internal/resolution/model"
	resolutionRepo "approve/internal/resolution/repository"
	rm "approve/internal/route/model"
	routeRepo "approve/internal/route/repository"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	gm "approve/internal/stepgroup/model"
	stepGroupRepo "approve/internal/stepgroup/repository"
	"fmt"
)

func FindFullRouteById(routeId int64) (fullRoute rm.FullRouteResponse, err error) {
	route, err := routeRepo.FindById(routeId)
	if err != nil {
		err = fmt.Errorf("error find full route by id: %v", err)
	} else if route.Id == 0 {
		err = fmt.Errorf("error find full route by id: %v", ErrRouteNotFound)
	}
	fullStepGroups, err := common.SafeExecuteG(err, func() ([]gm.StepGroupFullResponse, error) {
		return findStepGroups(routeId)
	})
	return common.SafeExecuteG(err, func() (rm.FullRouteResponse, error) {
		return route.ToFullResponse(fullStepGroups), nil
	})
}

func findStepGroups(routeId int64) ([]gm.StepGroupFullResponse, error) {
	stepGroups, err := stepGroupRepo.FindByRouteId(routeId)
	if err != nil {
		err = fmt.Errorf("error find step groups by route id: %v", err)
	}
	fullSteps, err := common.SafeExecuteG(err, func() ([]sm.StepFullResponse, error) {
		return findSteps(stepGroups)
	})
	return common.SafeExecuteG(err, func() ([]gm.StepGroupFullResponse, error) {
		groupIdToSteps := make(map[int64][]sm.StepFullResponse)
		for _, step := range fullSteps {
			groupIdToSteps[step.StepGroupId] = append(groupIdToSteps[step.StepGroupId], step)
		}
		fullStepGroups := make([]gm.StepGroupFullResponse, len(stepGroups))
		for i, group := range stepGroups {
			fullStepGroups[i] = group.ToFullResponse(groupIdToSteps[group.Id])
		}
		return fullStepGroups, nil
	})
}

func findSteps(stepGroups []gm.StepGroupEntity) (fullSteps []sm.StepFullResponse, err error) {
	stepGroupIds := make([]int64, len(stepGroups))
	for i, group := range stepGroups {
		stepGroupIds[i] = group.Id
	}
	steps, err := stepRepo.FindByStepGroupIds(stepGroupIds)
	fullApprovers, err := common.SafeExecuteG(err, func() ([]am.ApproverFullResponse, error) {
		return findApprovers(steps)
	})
	return common.SafeExecuteG(err, func() ([]sm.StepFullResponse, error) {
		stepIdToApprovers := make(map[int64][]am.ApproverFullResponse)
		for _, approver := range fullApprovers {
			stepIdToApprovers[approver.StepId] = append(stepIdToApprovers[approver.StepId], approver)
		}
		fullSteps = make([]sm.StepFullResponse, len(steps))
		for i, step := range steps {
			fullSteps[i] = step.ToFullResponse(stepIdToApprovers[step.Id])
		}
		return fullSteps, nil
	})
}

func findApprovers(steps []sm.StepEntity) (fullApprovers []am.ApproverFullResponse, err error) {
	stepIds := make([]int64, len(steps))
	for i, step := range steps {
		stepIds[i] = step.Id
	}
	approvers, err := approverRepo.FindByStepIds(stepIds)
	fullResolutions, err := common.SafeExecuteG(err, func() ([]resm.ResolutionResponse, error) {
		return findResolutions(approvers)
	})

	return common.SafeExecuteG(err, func() ([]am.ApproverFullResponse, error) {
		approverIdToResolution := make(map[int64]resm.ResolutionResponse)
		for _, resolution := range fullResolutions {
			approverIdToResolution[resolution.ApproverId] = resolution
		}
		fullApprovers = make([]am.ApproverFullResponse, len(approvers))
		for i, approver := range approvers {
			fullApprovers[i] = approver.ToFullResponse(approverIdToResolution[approver.Id])
		}
		return fullApprovers, nil
	})
}

func findResolutions(approvers []am.ApproverEntity) (fullResolutions []resm.ResolutionResponse, err error) {
	approverIds := make([]int64, len(approvers))
	for i, approver := range approvers {
		approverIds[i] = approver.Id
	}
	resolutions, err := resolutionRepo.FindByApproverIds(approverIds)
	return common.SafeExecuteG(err, func() ([]resm.ResolutionResponse, error) {
		fullResolutions = make([]resm.ResolutionResponse, len(resolutions))
		for i, resolution := range resolutions {
			fullResolutions[i] = resolution.ToResponse()
		}
		return fullResolutions, nil
	})
}
