package repository

import (
	cfg "approve/internal/config"
	rm "approve/internal/route/model"

	"github.com/jmoiron/sqlx"
)

func Save(
	route rm.RouteEntity,
) (res int64, err error) {
	err = cfg.DB.Get(
		&res,
		"insert into route (name, description, status, is_approved) values ($1, $2, $3, $4) returning id",
		route.Name,
		route.Description,
		route.Status,
		route.IsApproved,
	)
	return res, err
}

func StartRoute(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec("update route set status = 'STARTED' where id = $1", routeId)
	return err
}

func Update(route rm.RouteEntity) (int64, error) {
	_, err := cfg.DB.NamedExec(
		`update route 
     set 
       name = :name,
       description = :description
     where id = :id`,
		route,
	)
	return route.Id, err
}

func IsRouteStarted(routeId int64) (res bool, err error) {
	err = cfg.DB.Get(
		&res,
		`select exists (select 1 
                    from route
                    where id = $1 
                    and status in ('STARTED', 'FINISHED'))`,
		routeId,
	)
	return res, err
}

func FinishRoute(
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

func GetById(id int64) (res rm.RouteEntity, err error) {
	err = cfg.DB.Get(&res, "select * from route where id = $1", id)
	return res, err
}

func DeleteById(routeId int64) error {
	_, err := cfg.DB.Exec("delete from route where id = $1", routeId)
	return err
}
