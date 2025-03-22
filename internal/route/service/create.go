package service

import (
	"approve/internal/route/model"
	"approve/internal/route/repository"
)

func CreateRouteTemplate(request model.CreateRouteTemplateRequest) (routeId int64, err error) {
	return repository.Save(request.ToEntity())
}
