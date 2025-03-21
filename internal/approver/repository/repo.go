package repository

import (
	am "approve/internal/approver/model"
	sm "approve/internal/step/model"
	"github.com/jmoiron/sqlx"
)

type ApproverRepository struct {
	db *sqlx.DB
}

func NewApproverRepository(db *sqlx.DB) *ApproverRepository {
	return &ApproverRepository{db}
}

func (r *ApproverRepository) FindById(approverId int64) (approver am.ApproverEntity, err error) {
	err = r.db.Get(&approver, "select * from approver where id = $1", approverId)
	return approver, err
}

func (r *ApproverRepository) FindByStepId(stepId int64) ([]am.ApproverEntity, error) {
	var approvers []am.ApproverEntity
	err := r.db.Select(&approvers, "select * from approver where step_id = $1", stepId)
	return approvers, err
}

func (r *ApproverRepository) Save(approver am.ApproverEntity) (id int64, err error) {
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

func (r *ApproverRepository) StartStepApprovers(
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

func (r *ApproverRepository) Update(approver am.ApproverEntity) (approverId int64, err error) {
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

func (r *ApproverRepository) FinishApprover(
	tx *sqlx.Tx,
	approverId int64,
) error {
	_, err := tx.Exec("update approver set status = 'FINISHED' where id = $1", approverId)
	return err
}

func (r *ApproverRepository) FinishStepApprovers(
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

func (r *ApproverRepository) ExistNotFinishedApproversInStep(
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

func (r *ApproverRepository) StartNextApprover(
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

func (a *ApproverRepository) IsRouteStarted(routeId int64) (res bool, err error) {
	err = a.db.Get(
		&res,
		`select exists (select 1 
                    from route r
                    inner join step_group sg on sg.route_id = r.id
                    inner join step s on s.step_group_id = sg.id
                    inner join approver a on s.id = a.step_id
                    where a.id = $1 
                    and r.status = 'STARTED')`,
		routeId,
	)
	return res, err
}
