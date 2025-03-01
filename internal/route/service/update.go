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

func (service *UpdateRoute) Execute(request model.RouteRequest) (int64, error) {
	err := service.validate.Struct(request)
	if err == nil {
		return 0, err
	}
	return service.repo.Update(request.Name, request.Description)
}
