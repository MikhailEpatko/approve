package service

import (
	cm "approve/internal/common"
	rm "approve/internal/route/model"
	"fmt"
)

type UpdateRouteRepository interface {
	Update(route rm.RouteEntity) (int64, error)
	IsRouteStarted(routeId int64) (bool, error)
}

type UpdateRoute struct {
	repo UpdateRouteRepository
}

func (svc *UpdateRoute) Execute(request rm.UpdateRouteRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteG(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
