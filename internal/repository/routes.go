package repository

import (
	. "approve/pkg/model/entity"
	"github.com/jmoiron/sqlx"
)

type RouteRepository interface {
	FindById(id int64) (RouteEntity, error)
	Save(route RouteEntity) (int64, error)
	Update(route RouteEntity) error
}

type routeRepo struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) RouteRepository {
	return &routeRepo{db}
}

func (r *routeRepo) FindById(id int64) (RouteEntity, error) {
	route := RouteEntity{}
	err := r.db.Select(&route, "select * from route where id = $1", id)
	if err != nil {
		return RouteEntity{}, err
	}
	return route, nil
}

func (r *routeRepo) Save(route RouteEntity) (int64, error) {
	res, err := r.db.NamedExec(
		"insert into route (name, description, status) values (:name, :description, :status)",
		&route,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *routeRepo) Update(route RouteEntity) error {
	_, err := r.db.NamedExec(
		`update route 
     set
			 name = :name, 
			 description = :description, 
			 status = :status 
		 where id = :id`,
		&route,
	)
	if err != nil {
		return err
	}
	return nil
}
