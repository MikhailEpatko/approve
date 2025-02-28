package repository

import (
	. "approve/pkg/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepGroupRepository interface {
	FindByRouteId(id int64) ([]StepGroupEntity, error)
	Save(stepGroup StepGroupEntity) (int64, error)
	Update(stepGroup StepGroupEntity) error
}

type stepGroupRepo struct {
	db *sqlx.DB
}

func NewStepGroupRepository(db *sqlx.DB) StepGroupRepository {
	return &stepGroupRepo{db}
}

func (r *stepGroupRepo) FindByRouteId(id int64) ([]StepGroupEntity, error) {
	groups := []StepGroupEntity{}
	err := r.db.Select(&groups, "select * from step_group where route_id = $1", id)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *stepGroupRepo) Save(stepGroup StepGroupEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into step_group (route_id, name, number, status, step_type)
     values (:route_id, :name, :number, :status, :step_type)`,
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
			 step_type = :step_type 
		 where id = :id`,
		&stepGroup,
	)
	if err != nil {
		return err
	}
	return nil
}
