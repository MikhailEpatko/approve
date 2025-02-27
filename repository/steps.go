package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepRepository struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) *StepRepository {
	return &StepRepository{db}
}

func (r *StepRepository) FindAll() ([]StepEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *StepRepository) FindByIdIn(ids []int) ([]StepEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *StepRepository) FindById(id int) (StepEntity, error) {
	return StepEntity{}, NOT_IMPLEMENTED
}

func (r *StepRepository) Save(step StepEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *StepRepository) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *StepRepository) Update(step StepEntity) error {
	return NOT_IMPLEMENTED
}
