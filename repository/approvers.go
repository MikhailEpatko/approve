package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository struct {
	db *sqlx.DB
}

func NewApproverRepository(db *sqlx.DB) *ApproverRepository {
	return &ApproverRepository{db}
}

func (r *ApproverRepository) FindAll() ([]ApproverEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *ApproverRepository) FindByIdIn(ids []int) ([]ApproverEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *ApproverRepository) FindById(id int) (ApproverEntity, error) {
	return ApproverEntity{}, NOT_IMPLEMENTED
}

func (r *ApproverRepository) Save(approver ApproverEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *ApproverRepository) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *ApproverRepository) Update(approver ApproverEntity) error {
	return NOT_IMPLEMENTED
}
