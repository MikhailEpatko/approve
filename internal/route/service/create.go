package service

import (
	cm "approve/internal/common"
	"approve/internal/route/model"
	"approve/internal/route/repository"
)

func CreateRouteTemplate(request model.CreateRouteTemplateRequest) (routeId int64, err error) {
	err = cm.Validate(request)
	if err != nil {
		return 0, cm.RequestValidationError{Message: err.Error()}
	}
	return repository.Save(request.ToEntity())
}
