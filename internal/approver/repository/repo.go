package repository

import (
	. "approve/internal/approver/model"
	sm "approve/internal/step/model"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository interface {
	FindByStepId(id int64) ([]ApproverEntity, error)
	Save(approver ApproverEntity) (int64, error)
	StartApproversTx(tx *sqlx.Tx, step sm.StepEntity) error
	Update(approver ApproverEntity) (int64, error)
	FinishApproverTx(tx *sqlx.Tx, approverId int64) error
	FinishApproversByRouteId(tx *sqlx.Tx, routeId int64) error
	AreActiveInStepExist(tx *sqlx.Tx, stepId int64) (bool, error)
	StartActiveApprovers(tx *sqlx.Tx, stepId int64, approverId int64) error
}

type approverRepo struct {
	db *sqlx.DB
}

func NewApproverRepository(db *sqlx.DB) ApproverRepository {
	return &approverRepo{db}
}

func (r *approverRepo) FindByStepId(id int64) ([]ApproverEntity, error) {
	var approvers []ApproverEntity
	err := r.db.Select(&approvers, "select * from approver where step_id = $1", id)
	return approvers, err
}

func (r *approverRepo) Save(approver ApproverEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`insert into approver (step_id, guid, name, position, email, number)
     values (:step_id, :guid, :name, :position, :email, :number)`,
		&approver,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *approverRepo) StartApproversTx(
	tx *sqlx.Tx,
	step sm.StepEntity,
) error {
	_, err := tx.NamedQuery(
		`update approver 
     set status = 'STARTED'
     where step_id = :id 
     and (number = 1 and 'SEQUENTIAL_ALL_OFF' = :approver_order or 'SEQUENTIAL_ALL_OFF' <> :approver_order)
     returning *`,
		step,
	)
	return err
}

func (r *approverRepo) Update(approver ApproverEntity) (approverId int64, err error) {
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
	if err == nil {
		approverId = approver.Id
	}
	return approverId, err
}

func (r *approverRepo) FinishApproverTx(
	tx *sqlx.Tx,
	approverId int64,
) error {
	_, err := tx.Exec("update approver set status = 'FINISHED' where id = $1", approverId)
	return err
}

func (r *approverRepo) FinishApproversByRouteId(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec(
		`update approver 
     set status = 'FINISHED'
     where step_id in (
       select id from step where step_group_id in (
         select id from step_group where route_id = $1
       )
     )`,
		routeId,
	)
	return err
}

func (r *approverRepo) AreActiveInStepExist(
	tx *sqlx.Tx,
	stepId int64,
) (bool, error) {
	var res bool
	err := tx.Select(
		&res,
		`select exists (
  		 select 1 
  		 from approver
		   where step_id = $1 
		   and status != 'FINISHED'
		 )`,
		stepId,
	)
	return res, err
}

func (r *approverRepo) StartActiveApprovers(
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
