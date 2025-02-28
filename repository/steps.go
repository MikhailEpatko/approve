package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindAll() ([]StepEntity, error)
	FindByIdIn(ids []int) ([]StepEntity, error)
	FindById(id int) (StepEntity, error)
	Save(step StepEntity) (int, error)
	Delete(id int) error
	Update(step StepEntity) error
}

type stepRepo struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) StepRepository {
	return &stepRepo{db}
}

func (r *stepRepo) FindAll() ([]StepEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *stepRepo) FindByIdIn(ids []int) ([]StepEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *stepRepo) FindById(id int) (StepEntity, error) {
	return StepEntity{}, NOT_IMPLEMENTED
}

func (r *stepRepo) Save(step StepEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *stepRepo) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *stepRepo) Update(step StepEntity) error {
	return NOT_IMPLEMENTED
}
