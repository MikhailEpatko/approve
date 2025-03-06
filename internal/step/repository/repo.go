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
