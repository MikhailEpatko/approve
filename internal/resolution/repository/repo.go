package repository

import (
	. "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

type ResolutionRepository interface {
	FindByApproverId(id int64) ([]ResolutionEntity, error)
	Save(resolution ResolutionEntity) (int64, error)
	Update(resolution ResolutionEntity) error
	SaveAllTx(tx *sqlx.Tx, resolutions []ResolutionEntity) error
}

type resolutionRepo struct {
	db *sqlx.DB
}

func NewResolutionRepository(db *sqlx.DB) ResolutionRepository {
	return &resolutionRepo{db}
}

func (r *resolutionRepo) FindByApproverId(id int64) ([]ResolutionEntity, error) {
	var resolutions []ResolutionEntity
	err := r.db.Select(&resolutions, "select * from resolution where approver_id = $1", id)
	if err != nil {
		return nil, err
	}
	return resolutions, nil
}

func (r *resolutionRepo) Save(resolution ResolutionEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into resolution (approver_id, decision, comment)
     values (:approver_id, :decision, :comment)`,
		&resolution,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *resolutionRepo) Update(resolution ResolutionEntity) error {
	_, err := r.db.NamedExec(
		`update resolution
     set
       decision = :decision,
       comment = :comment
     where id = :id`,
		&resolution,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *resolutionRepo) SaveAllTx(
	tx *sqlx.Tx,
	resolutions []ResolutionEntity,
) error {
	_, err := tx.NamedExec(`insert into resolution (approver_id) values (:approver_id)`, resolutions)
	return err
}
