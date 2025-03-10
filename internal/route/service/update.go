package service

import (
	cm "approve/internal/common"
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	"fmt"
)

type UpdateRoute struct {
	repo rr.RouteRepository
}

func (svc *UpdateRoute) Execute(request rm.UpdateRouteRequest) (routeId int64, err error) {
	isRouteStarted, err := svc.repo.IsRouteStarted(request.Id)
	if err == nil && isRouteStarted {
		err = fmt.Errorf("route was started and cannot be updated")
	}
	return cm.SafeExecuteInt64(err, func() (int64, error) { return svc.repo.Update(request.ToEntity()) })
}
