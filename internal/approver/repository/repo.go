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
	IsRouteStarted(approverId int64) (bool, error)
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

func (r *approverRepo) IsRouteStarted(approverId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (
       select 1 from approver a
       inner join step s on s.id = a.step_id
       inner join step_group g on g.id = s.step_group_id
       inner join route r on r.id = g.route_id
       where s.id = $1 and r.status in ('STARTED', 'FINISHED'))`,
		approverId,
	)
	return res, err
}
