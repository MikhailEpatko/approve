package service

import (
	rm "approve/internal/route/model"
	rr "approve/internal/route/repository"
)

type FindRoutesByFilter struct {
	repo rr.FindByFilterRouteRepository
}

func (service *FindRoutesByFilter) Execute(filter rm.FilterRouteRequest) ([]rm.RouteResponse, int64, error) {
	entities, total, err := service.repo.FindByfilter(filter)
	if err != nil {
		return nil, 0, err
	}
	result := make([]rm.RouteResponse, len(entities))
	for i, entity := range entities {
		result[i] = entity.ToResponse()
	}
	return result, total, nil
}
