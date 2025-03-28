package repository

import (
	am "approve/internal/approver/model"
	"approve/internal/database"
	sm "approve/internal/step/model"
	"github.com/jmoiron/sqlx"
)

func FindById(approverId int64) (approver am.ApproverEntity, err error) {
	err = database.DB.Get(&approver, "select * from approver where id = $1", approverId)
	return approver, err
}

func FindByIdTx(
	tx *sqlx.Tx,
	approverId int64,
) (approver am.ApproverEntity, err error) {
	err = tx.Get(&approver, "select * from approver where id = $1", approverId)
	return approver, err
}

func FindByStepId(stepId int64) (approvers []am.ApproverEntity, err error) {
	err = database.DB.Select(&approvers, "select * from approver where step_id = $1", stepId)
	return approvers, err
}

func FindByStepIdTx(
	tx *sqlx.Tx,
	stepId int64,
) (approvers []am.ApproverEntity, err error) {
	err = tx.Select(&approvers, "select * from approver where step_id = $1", stepId)
	return approvers, err
}

func Save(approver am.ApproverEntity) (id int64, err error) {
	err = database.DB.Get(
		&id,
		`insert into approver (step_id, guid, name, position, email, number, status)
     values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		approver.StepId,
		approver.Guid,
		approver.Name,
		approver.Position,
		approver.Email,
		approver.Number,
		approver.Status,
	)
	return id, err
}

func StartStepApprovers(
	tx *sqlx.Tx,
	step sm.StepEntity,
) error {
	_, err := tx.NamedExec(
		`update approver 
     set status = 'STARTED'
     where step_id = :id 
     and (number = 1 and 'SERIAL' = :approver_order or 'SERIAL' <> :approver_order)`,
		step,
	)
	return err
}

func Update(
	tx *sqlx.Tx,
	approver am.ApproverEntity,
) (approverId int64, err error) {
	_, err = tx.NamedExec(
		`update approver 
     set name = :name,
       guid = :guid,
       position = :position,
       email = :email,
       number = :number
     where id = :id`,
		approver,
	)
	return approver.Id, err
}

func FinishApprover(
	tx *sqlx.Tx,
	approverId int64,
) error {
	_, err := tx.Exec("update approver set status = 'FINISHED' where id = $1", approverId)
	return err
}

func FinishStepApprovers(
	tx *sqlx.Tx,
	stepId int64,
) error {
	_, err := tx.Exec(
		`update approver 
     set status = 'FINISHED'
     where step_id = $1`,
		stepId,
	)
	return err
}

func ExistNotFinishedApproversInStep(
	tx *sqlx.Tx,
	stepId int64,
) (res bool, err error) {
	err = tx.Get(
		&res,
		"select exists (select 1 from approver where step_id = $1 and status != 'FINISHED')",
		stepId,
	)
	return res, err
}

func StartNextApprover(
	tx *sqlx.Tx,
	stepId int64,
	approverId int64,
) error {
	_, err := tx.Exec(
		`update approver 
     set status = 'STARTED'
     where step_id = $1
     and number = (select number + 1 from approver where id = $2)`,
		stepId,
		approverId,
	)
	return err
}

func IsRouteProcessing(
	tx *sqlx.Tx,
	approverId int64,
) (res bool, err error) {
	err = tx.Get(
		&res,
		`select exists (select 1 
                    from route r
                    inner join step_group sg on sg.route_id = r.id
                    inner join step s on s.step_group_id = sg.id
                    inner join approver a on s.id = a.step_id
                    where a.id = $1 
                    and r.status in ('STARTED', 'FINISHED'))`,
		approverId,
	)
	return res, err
}

func FindByStepIds(stepIds []int64) (approvers []am.ApproverEntity, err error) {
	query, args, err := sqlx.In(`select * from approver where step_id in (?)`, stepIds)
	if err == nil {
		err = database.DB.Select(&approvers, query, args...)
	}
	return approvers, err
}

func DeleteById(approver int64) error {
	_, err := database.DB.Exec("delete from approver where id = $1", approver)
	return err
}

func SaveAll(
	tx *sqlx.Tx,
	toSave []am.ApproverEntity,
) (err error) {
	_, err = tx.NamedExec(`insert into approver (step_id, guid, name, position, email, number, status)
                         values (:step_id, :guid, :name, :position, :email, :number, :status)`,
		toSave)
	return err
}
