package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type ResolutionRepository struct {
	db *sqlx.DB
}

func NewResolutionRepository(db *sqlx.DB) *ResolutionRepository {
	return &ResolutionRepository{db}
}

func (r *ResolutionRepository) FindAll() ([]ResolutionEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *ResolutionRepository) FindByIdIn(ids []int) ([]ResolutionEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *ResolutionRepository) FindById(id int) (ResolutionEntity, error) {
	return ResolutionEntity{}, NOT_IMPLEMENTED
}

func (r *ResolutionRepository) Save(resolution ResolutionEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *ResolutionRepository) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *ResolutionRepository) Update(resolution ResolutionEntity) error {
	return NOT_IMPLEMENTED
}
