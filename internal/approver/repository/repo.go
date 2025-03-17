package repository

import (
	am "approve/internal/approver/model"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository interface {
	FindByStepId(stepId int64) ([]am.ApproverEntity, error)
	Save(approver am.ApproverEntity) (int64, error)
	StartStepApprovers(tx *sqlx.Tx, stepId int64) error
	Update(approver am.ApproverEntity) (int64, error)
	FinishApprover(tx *sqlx.Tx, approverId int64) error
	ExistNotFinishedApproversInStep(tx *sqlx.Tx, stepId int64) (bool, error)
	StartNextApprover(tx *sqlx.Tx, stepId int64, approverId int64) error
	IsRouteStarted(routeId int64) (bool, error)
	FinishStepApprovers(tx *sqlx.Tx, stepId int64) error
}
type approverRepo struct {
	db *sqlx.DB
}

func NewApproverRepository(db *sqlx.DB) ApproverRepository {
	return &approverRepo{db}
}

func (r *approverRepo) FindByStepId(stepId int64) ([]am.ApproverEntity, error) {
	var approvers []am.ApproverEntity
	err := r.db.Select(&approvers, "select * from approver where step_id = $1", stepId)
	return approvers, err
}

func (r *approverRepo) Save(approver am.ApproverEntity) (id int64, err error) {
	err = r.db.Get(
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

func (r *approverRepo) StartStepApprovers(
	tx *sqlx.Tx,
	stepId int64,
) error {
	_, err := tx.Exec(
		`update approver 
     set status = 'STARTED'
     where step_id = $1 
     and (number = 1 and 'SERIAL' = :approver_order or 'SERIAL' <> :approver_order)`,
		stepId,
	)
	return err
}

func (r *approverRepo) Update(approver am.ApproverEntity) (approverId int64, err error) {
	_, err = r.db.NamedExec(
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

func (r *approverRepo) FinishApprover(
	tx *sqlx.Tx,
	approverId int64,
) error {
	_, err := tx.Exec("update approver set status = 'FINISHED' where step_id = $1", approverId)
	return err
}

func (r *approverRepo) FinishStepApprovers(
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

func (r *approverRepo) ExistNotFinishedApproversInStep(
	tx *sqlx.Tx,
	stepId int64,
) (res bool, err error) {
	err = tx.Select(
		&res,
		"select exists (select 1 from approver where step_id = $1 and status != 'FINISHED')",
		stepId,
	)
	return res, err
}

func (r *approverRepo) StartNextApprover(
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

func (a *approverRepo) IsRouteStarted(routeId int64) (res bool, err error) {
	err = a.db.Select(
		&res,
		"select exists (select 1 from route where id = $1 and status = 'STARTED')",
		routeId,
	)
	return res, err
}
