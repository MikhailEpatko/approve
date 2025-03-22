package service

import (
	cm "approve/internal/common"
	sm "approve/internal/step/model"
	stepRepo "approve/internal/step/repository"
	"errors"
)

func UpdateStep(request sm.UpdateStepRequest) (routeId int64, err error) {
	isRouteStarted, err := stepRepo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = errors.New("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return stepRepo.Update(request.ToEntity()) })
}
