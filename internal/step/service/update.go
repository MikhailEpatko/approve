package service

import (
	cm "approve/internal/common"
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	"fmt"
)

type UpdateStep struct {
	repo sr.StepRepository
}

func (svc *UpdateStep) Execute(request sm.UpdateStepRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
