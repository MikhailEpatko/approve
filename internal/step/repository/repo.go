package repository

import (
	sm "approve/internal/step/model"
	gm "approve/internal/stepgroup/model"
	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindByGroupId(id int64) ([]sm.StepEntity, error)
	Save(step sm.StepEntity) (int64, error)
	StartStepsTx(tx *sqlx.Tx, group gm.StepGroupEntity) ([]sm.StepEntity, error)
	Update(step sm.StepEntity) (int64, error)
	IsRouteStarted(stepId int64) (bool, error)
	FinishStepsByRouteId(tx *sqlx.Tx, routeId int64) error
}

type stepRepo struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) StepRepository {
	return &stepRepo{db}
}

func (r *stepRepo) FindByGroupId(id int64) ([]sm.StepEntity, error) {
	var steps []sm.StepEntity
	err := r.db.Select(&steps, "select * from step where step_group_id = $1", id)
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *stepRepo) Save(step sm.StepEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into step (step_group_id, name, number, status, approver_order)
     values (:step_group_id, :name, :number, :status, :approver_order)`,
		&step,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *stepRepo) StartStepsTx(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) ([]sm.StepEntity, error) {
	rows, err := tx.NamedQuery(
		`update step 
     set status = 'STARTED'
     where step_group_id = :id 
     and (number = 1 and 'SEQUENTIAL_ALL_OFF' = :step_order or 'SEQUENTIAL_ALL_OFF' <> :step_order)
     returning *`,
		group,
	)
	if err != nil {
		return nil, err
	}
	saved := make([]sm.StepEntity, 0)
	step := sm.StepEntity{}
	for rows.Next() {
		err = rows.StructScan(&step)
		if err != nil {
			return nil, err
		}
		saved = append(saved, step)
	}
	return saved, nil
}

func (r *stepRepo) Update(step sm.StepEntity) (stepId int64, err error) {
	_, err = r.db.NamedExec(
		`update step 
     set name = :name,
       number = :number,
       approver_order = :approver_order
     where id = :id`,
		step,
	)
	if err == nil {
		stepId = step.Id
	}
	return stepId, err
}

func (r *stepRepo) IsRouteStarted(stepId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (
       select 1 from step s 
       inner join step_group g on g.id = s.step_group_id
       inner join route r on r.id = g.route_id
       where s.id = $1 and r.status in ('STARTED', 'FINISHED'))`,
		stepId,
	)
	return res, err
}

func (r *stepRepo) FinishStepsByRouteId(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec(
		`update step 
     set status = 'FINISHED'
     where step.step_group_id in (
       select id from step_group where route_id = $1
     )`,
		routeId,
	)
	return err
}
