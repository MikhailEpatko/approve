package service

import (
	. "approve/internal/route/model"
	. "approve/internal/route/repository"
	"github.com/go-playground/validator/v10"
)

type FindRoutesByFilter struct {
	repo     RouteRepository
	validate *validator.Validate
}

func NewFindRoutesByApproverGuid(
	repo RouteRepository,
	validate *validator.Validate,
) *FindRoutesByFilter {
	return &FindRoutesByFilter{repo, validate}
}

func (service *FindRoutesByFilter) Execute(filter FilterRouteRequest) ([]RouteResponse, int64, error) {
	err := service.validate.Struct(filter)
	if err != nil {
		return nil, 0, err
	}
	entities, total, err := service.repo.FindByfilter(filter)
	if err != nil {
		return nil, 0, err
	}
	result := make([]RouteResponse, len(entities))
	for i, entity := range entities {
		result[i] = entity.ToResponse()
	}
	return result, total, nil
}
