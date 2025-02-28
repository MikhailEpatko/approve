package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type RouteRepository interface {
	FindAll() ([]RouteEntity, error)
	FindByIdIn(ids []int) ([]RouteEntity, error)
	FindById(id int) (RouteEntity, error)
	Save(route RouteEntity) (int, error)
	Delete(id int) error
	Update(route RouteEntity) error
}

type routeRepo struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) RouteRepository {
	return &routeRepo{db}
}

func (r *routeRepo) FindAll() ([]RouteEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *routeRepo) FindByIdIn(ids []int) ([]RouteEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *routeRepo) FindById(id int) (RouteEntity, error) {
	return RouteEntity{}, NOT_IMPLEMENTED
}

func (r *routeRepo) Save(route RouteEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *routeRepo) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *routeRepo) Update(route RouteEntity) error {
	return NOT_IMPLEMENTED
}
