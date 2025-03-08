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
	DeativateTx(tx *sqlx.Tx, approverId int64) error
	DeactivateApproversByRouteId(tx *sqlx.Tx, routeId int64) error
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
     set active = true
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

func (r *approverRepo) DeativateTx(
	tx *sqlx.Tx,
	approverId int64,
) error {
	_, err := tx.Exec("update approver set active = false where id = $1", approverId)
	return err
}

func (r *approverRepo) DeactivateApproversByRouteId(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec(
		`update approver 
     set active = false
     where step_id in (
       select id from step where step_group_id in (
         select id from step_group where route_id = $1
       )
     )`,
		routeId,
	)
	return err
}
