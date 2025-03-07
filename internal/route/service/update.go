package service

import (
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
	"fmt"
)

type UpdateRoute struct {
	repo rr.RouteRepository
}

func (svc *UpdateRoute) Execute(request rm.UpdateRouteRequest) (routeId int64, err error) {
	if res, err := svc.repo.IsRouteStarted(request.Id); err != nil {
		return 0, err
	} else if res {
		return 0, fmt.Errorf("route was started and cannot be updated")
	}
	return svc.repo.Update(request.ToEntity())
}
