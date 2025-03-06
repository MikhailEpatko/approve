package repository

import (
	rm "approve/internal/route/model"
	"github.com/jmoiron/sqlx"
)

type RouteRepository interface {
	Save(entity rm.RouteEntity) (int64, error)
	StartRouteTx(tx *sqlx.Tx, id int64) error
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
