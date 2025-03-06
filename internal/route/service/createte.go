package service

import (
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
)

type CreateRouteTemplate struct {
	repo rr.RouteRepository
}

func (svc *CreateRouteTemplate) Execute(request rm.CreateRouteTemplateRequest) (routeId int64, err error) {
	return svc.repo.Save(request.ToEntity())
}
