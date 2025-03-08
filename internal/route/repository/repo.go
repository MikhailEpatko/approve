package repository

import (
	rm "approve/internal/route/model"
	"github.com/jmoiron/sqlx"
)

type RouteRepository interface {
	Save(entity rm.RouteEntity) (int64, error)
	StartRouteTx(tx *sqlx.Tx, id int64) error
	Update(route rm.RouteEntity) (int64, error)
	IsRouteStarted(routeId int64) (bool, error)
	FinishRoute(tx *sqlx.Tx, routeId int64, isRouteApproved bool) error
}

type routeRepo struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) RouteRepository {
	return &routeRepo{db}
}

func (r *routeRepo) Save(
	route rm.RouteEntity,
) (int64, error) {
	res, err := r.db.NamedExec(
		"insert into route (name, description, status) values (:name, :description, :status)",
		route,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *routeRepo) StartRouteTx(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec("update route set status = 'STARTED' where id = $1", routeId)
	return err
}

func (r *routeRepo) Update(route rm.RouteEntity) (routId int64, err error) {
	_, err = r.db.NamedExec(
		`update route 
     set name = :name,
       description = :description
     where id = :id`,
		route,
	)
	if err == nil {
		routId = route.Id
	}
	return routId, err
}

func (r *routeRepo) IsRouteStarted(routeId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (select 1 from route where id = $1 and status in ('STARTED', 'FINISHED'))`,
		routeId,
	)
	return res, err
}

func (r *routeRepo) FinishRoute(
	tx *sqlx.Tx,
	routeId int64,
	isRouteApproved bool,
) error {
	_, err := tx.Exec(
		`update route 
     set 
       status = 'FINISHED',
       is_approved = $2
     where id = $1`,
		routeId,
		isRouteApproved,
	)
	return err
}
