package service

import (
	cm "approve/internal/common"
	"approve/internal/route/model"
	"approve/internal/route/repository"
)

func FindByFilter(filter model.FilterRouteRequest) (page cm.PageResponse, err error) {
	err = cm.Validate(filter)
	if err != nil {
		return page, cm.RequestValidationError{Message: err.Error()}
	}

	entities, total, err := repository.FindByfilter(filter)
	if err == nil {
		result := make([]model.RouteResponse, len(entities))
		for i, entity := range entities {
			result[i] = entity.ToResponse()
		}
		page = cm.PageResponse{
			Result:      result,
			Total:       total,
			PageRequest: filter.PageRequest,
		}
	}
	return page, err
}
