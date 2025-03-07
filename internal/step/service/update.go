package service

import (
	sm "approve/internal/step/model"
	sr "approve/internal/step/repository"
	"fmt"
)

type UpdateStep struct {
	repo sr.StepRepository
}

func (svc *UpdateStep) Execute(request sm.UpdateStepRequest) (routeId int64, err error) {
	if res, err := svc.repo.IsRouteStarted(request.Id); err != nil {
		return 0, err
	} else if res {
		return 0, fmt.Errorf("route was started and cannot be updated")
	}
	return svc.repo.Update(request.ToEntity())
}
