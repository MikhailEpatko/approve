package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepGroupRepository struct {
	db *sqlx.DB
}

func NewStepGroupRepository(db *sqlx.DB) *StepGroupRepository {
	return &StepGroupRepository{db}
}

func (r *StepGroupRepository) FindAll() ([]StepGroupEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *StepGroupRepository) FindByIdIn(ids []int) ([]StepGroupEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *StepGroupRepository) FindById(id int) (StepGroupEntity, error) {
	return StepGroupEntity{}, NOT_IMPLEMENTED
}

func (r *StepGroupRepository) Save(StepGroup StepGroupEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *StepGroupRepository) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *StepGroupRepository) Update(StepGroup StepGroupEntity) error {
	return NOT_IMPLEMENTED
}
