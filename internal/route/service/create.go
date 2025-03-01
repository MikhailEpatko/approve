package service

import (
	"approve/internal/route/model"
	"approve/internal/route/repository"
	"github.com/go-playground/validator/v10"
)

type CreateRoute struct {
	repo     repository.RouteRepository
	validate *validator.Validate
}

func NewCreateRoute(
	repo repository.RouteRepository,
	validate *validator.Validate,
) *CreateRoute {
	return &CreateRoute{repo, validate}
}

func (service *CreateRoute) Execute(request model.RouteRequest) (int64, error) {
	err := service.validate.Struct(request)
	if err != nil {
		return 0, err
	}
	return service.repo.Save(request.Name, request.Description)
}
