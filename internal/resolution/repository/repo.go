package repository

import (
	"approve/internal/database"
	resm "approve/internal/resolution/model"

	"github.com/jmoiron/sqlx"
)

func FindByApproverId(id int64) ([]resm.ResolutionEntity, error) {
	var resolutions []resm.ResolutionEntity
	err := database.DB.Select(&resolutions, "select * from resolution where approver_id = $1", id)
	return resolutions, err
}

func Save(resolution resm.ResolutionEntity) (resolutionId int64, err error) {
	err = database.DB.Get(
		&resolutionId,
		`insert into resolution (approver_id, is_approved, comment)
     values ($1, $2, $3)
     returning id`,
		resolution.ApproverId,
		resolution.IsApproved,
		resolution.Comment,
	)
	return resolutionId, err
}

func SaveTx(
	tx *sqlx.Tx,
	resolution resm.ResolutionEntity,
) (resolutionId int64, err error) {
	err = tx.Get(
		&resolutionId,
		`insert into resolution (approver_id, is_approved, comment)
     values ($1, $2, $3)
     returning id`,
		resolution.ApproverId,
		resolution.IsApproved,
		resolution.Comment,
	)
	return resolutionId, err
}

func ApprovingInfoTx(
	tx *sqlx.Tx,
	approverId int64,
) (res resm.ApprovingInfoEntity, err error) {
	err = tx.Get(
		&res,
		`select 
			 r.id as route_id,
			 sg.id as step_group_id,
			 sg.step_order as step_order,
			 s.id as step_id,
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

func FindByApproverIds(approverIds []int64) (resolutions []resm.ResolutionEntity, err error) {
	query, args, err := sqlx.In(`select * from resolution where approver_id in (?)`, approverIds)
	if err == nil {
		err = database.DB.Select(&resolutions, query, args...)
	}
	return resolutions, err
}
