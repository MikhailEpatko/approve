package repository

import (
	. "approve/internal/approver/model"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository interface {
	FindByStepId(id int64) ([]ApproverEntity, error)
	Save(approver ApproverEntity) (int64, error)
	Update(approver ApproverEntity) error
	SaveAllTx(tx *sqlx.Tx, save []ApproverEntity) error
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
	if err != nil {
		return nil, err
	}
	return approvers, nil
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

func (r *approverRepo) Update(approver ApproverEntity) error {
	_, err := r.db.NamedExec(
		`update approver
     set
       guid = :guid,
       name = :name,
       position = :position,
       email = :email,
       number = :number
     where id = :id`,
		&approver,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *approverRepo) SaveAllTx(
	tx *sqlx.Tx,
	approvers []ApproverEntity,
) error {
	_, err := tx.NamedExec(
		`insert into approver (step_id, guid, name, position, email, number)
     values (:step_id, :guid, :name, :position, :email, :number)`,
		&approvers,
	)
	return err
}
