package service

import (
	rm "approve/internal/route/model"
)

type SaveRouteRepository interface {
	Save(entity rm.RouteEntity) (int64, error)
}

type CreateRouteTemplate struct {
	repo SaveRouteRepository
}

func (svc *CreateRouteTemplate) Execute(request rm.CreateRouteTemplateRequest) (routeId int64, err error) {
	return svc.repo.Save(request.ToEntity())
}
