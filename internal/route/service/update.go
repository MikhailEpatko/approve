package service

import (
	cm "approve/internal/common"
	rm "approve/internal/route/model"
	"approve/internal/route/repository"
	"fmt"
)

func UpdateRoute(request rm.UpdateRouteRequest) (routeId int64, err error) {
	isRouteStarted, err := repository.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return repository.Update(request.ToEntity()) })
}
