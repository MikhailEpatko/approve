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
