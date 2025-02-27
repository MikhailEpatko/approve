package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type RouteRepository struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) *RouteRepository {
	return &RouteRepository{db}
}

func (r *RouteRepository) FindAll() ([]RouteEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *RouteRepository) FindByIdIn(ids []int) ([]RouteEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *RouteRepository) FindById(id int) (RouteEntity, error) {
	return RouteEntity{}, NOT_IMPLEMENTED
}

func (r *RouteRepository) Save(route RouteEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *RouteRepository) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *RouteRepository) Update(route RouteEntity) error {
	return NOT_IMPLEMENTED
}
