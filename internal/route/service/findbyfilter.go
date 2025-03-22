package service

import (
	"approve/internal/route/model"
	"approve/internal/route/repository"
)

func FindByFilter(filter model.FilterRouteRequest) (result []model.RouteResponse, total int64, err error) {
	entities, total, err := repository.FindByfilter(filter)
	if err == nil {
		result = make([]model.RouteResponse, len(entities))
		for i, entity := range entities {
			result[i] = entity.ToResponse()
		}
	}
	return result, total, err
}
