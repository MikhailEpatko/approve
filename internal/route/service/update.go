package service

import (
	"approve/internal/route/model"
	"approve/internal/route/repository"
	"github.com/go-playground/validator/v10"
)

type UpdateRoute struct {
	repo     repository.RouteRepository
	validate *validator.Validate
}

func NewUpdateRoute(
	repo repository.RouteRepository,
	validate *validator.Validate,
) *UpdateRoute {
	return &UpdateRoute{repo, validate}
}

func (service *UpdateRoute) Execute(
	id int64,
	request model.RouteRequest,
) (int64, error) {
	entity := request.ToEntity(id)
	return service.repo.Update(entity)
}
