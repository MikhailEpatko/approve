package repository

import (
	. "approve/internal/step/model"
	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindByGroupId(id int64) ([]StepEntity, error)
	Save(step StepEntity) (int64, error)
	Update(step StepEntity) error
	SaveAllTxReturning(tx *sqlx.Tx, steps []StepEntity) ([]StepEntity, error)
}

type stepRepo struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) StepRepository {
	return &stepRepo{db}
}

func (r *stepRepo) FindByGroupId(id int64) ([]StepEntity, error) {
	var steps []StepEntity
	err := r.db.Select(&steps, "select * from step where step_group_id = $1", id)
	if err != nil {
		return nil, err
	}
	return steps, nil
}

func (r *stepRepo) Save(step StepEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into step (step_group_id, name, number, status, approve_type)
     values (:step_group_id, :name, :number, :status, :approve_type)`,
		&step,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *stepRepo) Update(step StepEntity) error {
	_, err := r.db.NamedExec(
		`update step
     set 
       name = :name, 
       number = :number, 
       status = :status, 
       approve_type = :approve_type
     where id = :id`,
		&step,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *stepRepo) SaveAllTxReturning(
	tx *sqlx.Tx,
	steps []StepEntity,
) ([]StepEntity, error) {
	saved := make([]StepEntity, 0, len(steps))
	rows, err := tx.NamedQuery(
		`insert into step (step_group_id, name, number, status, approve_type)
     values (:step_group_id, :name, :number, :status, :approve_type)
     returning *`,
		&steps,
	)
	if err != nil {
		return nil, err
	}
	step := StepEntity{}
	for rows.Next() {
		err = rows.StructScan(&step)
		if err != nil {
			return nil, err
		}
		saved = append(saved, step)
	}
	return saved, nil
}
