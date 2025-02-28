package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type StepGroupRepository interface {
	FindAll() ([]StepGroupEntity, error)
	FindByIdIn(ids []int) ([]StepGroupEntity, error)
	FindById(id int) (StepGroupEntity, error)
	Save(StepGroup StepGroupEntity) (int, error)
	Delete(id int) error
	Update(StepGroup StepGroupEntity) error
}

type stepGroupRepo struct {
	db *sqlx.DB
}

func NewStepGroupRepository(db *sqlx.DB) StepGroupRepository {
	return &stepGroupRepo{db}
}

func (r *stepGroupRepo) FindAll() ([]StepGroupEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *stepGroupRepo) FindByIdIn(ids []int) ([]StepGroupEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *stepGroupRepo) FindById(id int) (StepGroupEntity, error) {
	return StepGroupEntity{}, NOT_IMPLEMENTED
}

func (r *stepGroupRepo) Save(StepGroup StepGroupEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *stepGroupRepo) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *stepGroupRepo) Update(StepGroup StepGroupEntity) error {
	return NOT_IMPLEMENTED
}
