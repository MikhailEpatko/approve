package repository

import (
	. "approve/common"
	. "approve/model/entity"
	"github.com/jmoiron/sqlx"
)

type ResolutionRepository interface {
	FindAll() ([]ResolutionEntity, error)
	FindByIdIn(ids []int) ([]ResolutionEntity, error)
	FindById(id int) (ResolutionEntity, error)
	Save(resolution ResolutionEntity) (int, error)
	Delete(id int) error
	Update(resolution ResolutionEntity) error
}

type resolutionRepo struct {
	db *sqlx.DB
}

func NewResolutionRepository(db *sqlx.DB) ResolutionRepository {
	return &resolutionRepo{db}
}

func (r *resolutionRepo) FindAll() ([]ResolutionEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *resolutionRepo) FindByIdIn(ids []int) ([]ResolutionEntity, error) {
	return nil, NOT_IMPLEMENTED
}

func (r *resolutionRepo) FindById(id int) (ResolutionEntity, error) {
	return ResolutionEntity{}, NOT_IMPLEMENTED
}

func (r *resolutionRepo) Save(resolution ResolutionEntity) (int, error) {
	return 0, NOT_IMPLEMENTED
}

func (r *resolutionRepo) Delete(id int) error {
	return NOT_IMPLEMENTED
}

func (r *resolutionRepo) Update(resolution ResolutionEntity) error {
	return NOT_IMPLEMENTED
}
