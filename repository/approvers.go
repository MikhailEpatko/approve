package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository interface {
	FindAll() ([]ApproverEntity, error)
	FindByIdIn(ids []int) ([]ApproverEntity, error)
	FindById(id int) (ApproverEntity, error)
	Save(approver ApproverEntity) (int, error)
	Delete(id int) error
	Update(approver ApproverEntity) error
}

type approverRepo struct {
	db *sqlx.DB
}

func NewApproverRepository(db *sqlx.DB) ApproverRepository {
	return &approverRepo{db}
}

func (r *approverRepo) FindAll() ([]ApproverEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *approverRepo) FindByIdIn(ids []int) ([]ApproverEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *approverRepo) FindById(id int) (ApproverEntity, error) {
	return ApproverEntity{}, NOT_IMPLEMENTED
}

func (r *approverRepo) Save(approver ApproverEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *approverRepo) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *approverRepo) Update(approver ApproverEntity) error {
	return NOT_IMPLEMENTED
}
