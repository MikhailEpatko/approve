package repository

import (
	. "approve/internal/stepgroup/model"
	"github.com/jmoiron/sqlx"
)

type StepGroupRepository interface {
	FindByRouteId(id int64) ([]StepGroupEntity, error)
	Save(stepGroup StepGroupEntity) (int64, error)
	Update(stepGroup StepGroupEntity) error
	SaveAllTx(tx *sqlx.Tx, groups []StepGroupEntity) ([]StepGroupEntity, error)
	StartGroupsTx(tx *sqlx.Tx, routeId int64) (StepGroupEntity, error)
}

type stepGroupRepo struct {
	db *sqlx.DB
}

func NewStepGroupRepository(db *sqlx.DB) StepGroupRepository {
	return &stepGroupRepo{db}
}

func (r *stepGroupRepo) FindByRouteId(id int64) ([]StepGroupEntity, error) {
	var groups []StepGroupEntity
	err := r.db.Select(&groups, "select * from step_group where route_id = $1", id)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *stepGroupRepo) Save(stepGroup StepGroupEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into step_group (route_id, name, number, status, step_order)
     values (:route_id, :name, :number, :status, :step_order)`,
		&stepGroup,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *stepGroupRepo) Update(stepGroup StepGroupEntity) error {
	_, err := r.db.NamedExec(
		`update step_group 
     set
		 	 name = :name, 
			 number = :number, 
			 status = :status, 
			 step_order = :step_order 
		 where id = :id`,
		&stepGroup,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *stepGroupRepo) SaveAllTx(
	tx *sqlx.Tx,
	groups []StepGroupEntity,
) ([]StepGroupEntity, error) {
	rows, err := tx.NamedQuery(
		`insert into step_group (route_id, name, number, status, step_order)
     values (:route_id, :name, :number, :status, :step_order)
     returning *`,
		&groups,
	)
	if err != nil {
		return nil, err
	}
	saved := make([]StepGroupEntity, 0, len(groups))
	group := StepGroupEntity{}
	for rows.Next() {
		err = rows.StructScan(&group)
		if err != nil {
			return nil, err
		}
		saved = append(saved, group)
	}
	return saved, nil
}

func (r *stepGroupRepo) StartGroupsTx(
	tx *sqlx.Tx,
	routeId int64,
) (group StepGroupEntity, err error) {
	rows, err := tx.Queryx(
		`update step_group 
     set status = 'STARTED'
     where route_id = $1 and number = 1
     returning *`,
		routeId,
	)
	if err == nil && rows.Next() {
		err = rows.Scan(&group)
	}
	return group, err
}
