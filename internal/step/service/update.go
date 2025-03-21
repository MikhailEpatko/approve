package service

import (
	cm "approve/internal/common"
	sm "approve/internal/step/model"
	"errors"
)

type UpdateStepRepository interface {
	Update(step sm.StepEntity) (int64, error)
	IsRouteStarted(stepId int64) (bool, error)
}

type UpdateStep struct {
	repo UpdateStepRepository
}

func (svc *UpdateStep) Execute(request sm.UpdateStepRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = errors.New("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
