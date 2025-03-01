package step

import (
	. "approve/pkg/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindByGroupId(id int64) ([]StepEntity, error)
	Save(step StepEntity) (int64, error)
	Update(step StepEntity) error
}

type stepRepo struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) StepRepository {
	return &stepRepo{db}
}

func (r *stepRepo) FindByGroupId(id int64) ([]StepEntity, error) {
	steps := []StepEntity{}
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
