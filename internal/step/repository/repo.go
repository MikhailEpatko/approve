package repository

import (
	. "approve/internal/step/model"
	gm "approve/internal/stepgroup/model"
	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindByGroupId(id int64) ([]StepEntity, error)
	Save(step StepEntity) (int64, error)
	Update(step StepEntity) error
	SaveAllTx(tx *sqlx.Tx, steps []StepEntity) ([]StepEntity, error)
	StartStepsTx(tx *sqlx.Tx, group gm.StepGroupEntity) ([]StepEntity, error)
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
		`insert into step (step_group_id, name, number, status, approver_order)
     values (:step_group_id, :name, :number, :status, :approver_order)`,
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
       approver_order = :approver_order
     where id = :id`,
		&step,
	)
	return err
}

func (r *stepRepo) SaveAllTx(
	tx *sqlx.Tx,
	steps []StepEntity,
) ([]StepEntity, error) {
	saved := make([]StepEntity, 0, len(steps))
	rows, err := tx.NamedQuery(
		`insert into step (step_group_id, name, number, status, approver_order)
     values (:step_group_id, :name, :number, :status, :approver_order)
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

func (r *stepRepo) StartStepsTx(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) ([]StepEntity, error) {
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
	saved := make([]StepEntity, 0)
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
