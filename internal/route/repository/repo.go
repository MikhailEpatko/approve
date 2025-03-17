package repository

import (
	rm "approve/internal/route/model"

	"github.com/jmoiron/sqlx"
)

type RouteRepository interface {
	Save(entity rm.RouteEntity) (int64, error)
	StartRoute(tx *sqlx.Tx, id int64) error
	Update(route rm.RouteEntity) (int64, error)
	IsRouteStarted(routeId int64) (bool, error)
	FinishRoute(tx *sqlx.Tx, routeId int64, isGroupApproved bool) error
	GetById(id int64) (rm.RouteEntity, error)
	DeleteById(routeId int64) error
}

type routeRepo struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) RouteRepository {
	return &routeRepo{db}
}

func (r *routeRepo) Save(
	route rm.RouteEntity,
) (res int64, err error) {
	err = r.db.Get(
		&res,
		"insert into route (name, description, status, is_approved) values ($1, $2, $3, $4) returning id",
		route.Name,
		route.Description,
		route.Status,
		route.IsApproved,
	)
	return res, err
}

func (r *routeRepo) StartRoute(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec("update route set status = 'STARTED' where id = $1", routeId)
	return err
}

func (r *routeRepo) Update(route rm.RouteEntity) (int64, error) {
	_, err := r.db.NamedExec(
		`update route 
     set 
       name = :name,
       description = :description
     where id = :id`,
		route,
	)
	return route.Id, err
}

func (r *routeRepo) IsRouteStarted(routeId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (select 1 
                    from route
                    where id = $1 
                    and status in ('STARTED', 'FINISHED'))`,
		routeId,
	)
	return res, err
}

func (r *routeRepo) FinishRoute(
	tx *sqlx.Tx,
	routeId int64,
	isGroupApproved bool,
) error {
	_, err := tx.Exec(
		`update route 
     set 
       status = 'FINISHED',
       is_approved = $2
     where id = $1`,
		routeId,
		isGroupApproved,
	)
	return err
}

func (r *routeRepo) GetById(id int64) (res rm.RouteEntity, err error) {
	err = r.db.Get(&res, "select * from route where id = $1", id)
	return res, err
}

func (r *routeRepo) DeleteById(routeId int64) error {
	_, err := r.db.Exec("delete from route where id = $1", routeId)
	return err
}
