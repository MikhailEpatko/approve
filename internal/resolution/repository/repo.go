package repository

import (
	. "approve/internal/resolution/model"
	"github.com/jmoiron/sqlx"
)

type ResolutionRepository interface {
	FindByApproverId(id int64) ([]ResolutionEntity, error)
	Save(resolution ResolutionEntity) (int64, error)
	Update(resolution ResolutionEntity) error
	SaveTx(tx *sqlx.Tx, resolution ResolutionEntity) (int64, error)
	ApprovingInfoTx(tx *sqlx.Tx, approverId int64) (ApprovingInfoEntity, error)
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
		`insert into resolution (approver_id, is_approved, comment)
     values (:approver_id, :is_approved, :comment)`,
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
       is_approved = :is_approved,
       comment = :comment
     where id = :id`,
		&resolution,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *resolutionRepo) SaveTx(
	tx *sqlx.Tx,
	resolution ResolutionEntity,
) (int64, error) {
	res, err := tx.NamedExec(
		`insert into resolution (approver_id, is_approved, comment)
     values (:approver_id, :is_approved, :comment)`,
		resolution,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *resolutionRepo) ApprovingInfoTx(
	tx *sqlx.Tx,
	approverId int64,
) (res ApprovingInfoEntity, err error) {
	err = tx.Get(
		&res,
		`select 
			 r.id as route_id,
			 sg.id as step_group_id,
			 sg.step_order as step_order,
			 s.id as step_id,
			 s.status as step_status,
			 s.approver_order as approver_order,
			 a.id as approver_id,
			 a.guid as guid,
			 a.status as approver_status
		 from approver a
     inner join step s on s.id = a.step_id
		 inner join step_group sg on sg.id = s.step_group_id
		 inner join route r on r.id = sg.route_id
     where a.id = $1`,
		approverId,
	)
	return res, err
}
