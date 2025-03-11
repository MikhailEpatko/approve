package repository

import (
	cm "approve/internal/common"
	rm "approve/internal/resolution/model"

	"github.com/jmoiron/sqlx"
)

type ResolutionRepository interface {
	FindByApproverId(id int64) ([]rm.ResolutionEntity, error)
	Save(resolution rm.ResolutionEntity) (int64, error)
	Update(resolution rm.ResolutionEntity) error
	SaveTx(tx *sqlx.Tx, resolution rm.ResolutionEntity) (int64, error)
	ApprovingInfoTx(tx *sqlx.Tx, approverId int64) (rm.ApprovingInfoEntity, error)
}

type resolutionRepo struct {
	db *sqlx.DB
}

func NewResolutionRepository(db *sqlx.DB) ResolutionRepository {
	return &resolutionRepo{db}
}

func (r *resolutionRepo) FindByApproverId(id int64) ([]rm.ResolutionEntity, error) {
	var resolutions []rm.ResolutionEntity
	err := r.db.Select(&resolutions, "select * from resolution where approver_id = $1", id)
	return resolutions, err
}

func (r *resolutionRepo) Save(resolution rm.ResolutionEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into resolution (approver_id, is_approved, comment)
     values (:approver_id, :is_approved, :comment)`,
		&resolution,
	)
	return cm.SafeExecuteG(err, func() (int64, error) { return res.LastInsertId() })
}

func (r *resolutionRepo) Update(resolution rm.ResolutionEntity) error {
	_, err := r.db.NamedExec(
		`update resolution
     set
       is_approved = :is_approved,
       comment = :comment
     where id = :id`,
		&resolution,
	)
	return err
}

func (r *resolutionRepo) SaveTx(
	tx *sqlx.Tx,
	resolution rm.ResolutionEntity,
) (int64, error) {
	res, err := tx.NamedExec(
		`insert into resolution (approver_id, is_approved, comment)
     values (:approver_id, :is_approved, :comment)`,
		resolution,
	)
	return cm.SafeExecuteG(err, func() (int64, error) { return res.LastInsertId() })
}

func (r *resolutionRepo) ApprovingInfoTx(
	tx *sqlx.Tx,
	approverId int64,
) (res rm.ApprovingInfoEntity, err error) {
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
