package repository

import (
	"approve/internal/common"
	sm "approve/internal/step/model"
	gm "approve/internal/stepgroup/model"

	"github.com/jmoiron/sqlx"
)

type StepRepository interface {
	FindByGroupId(id int64) ([]sm.StepEntity, error)
	Save(step sm.StepEntity) (int64, error)
	StartSteps(tx *sqlx.Tx, group gm.StepGroupEntity) ([]sm.StepEntity, error)
	Update(step sm.StepEntity) (int64, error)
	IsRouteStarted(stepId int64) (bool, error)
	FinishStepsByRouteId(tx *sqlx.Tx, routeId int64) error
	FinishStep(tx *sqlx.Tx, stepId int64) error
	CalculateAndSetIsApproved(
		tx *sqlx.Tx,
		stepId int64,
		approverOrder common.OrderType,
		isResolutionApproved bool,
	) (res bool, err error)
	ExistsNotFinishedStepsInGroup(x *sqlx.Tx, stepGroupId int64) (bool, error)
	StartNextStep(tx *sqlx.Tx, stepGroupId int64, stepId int64) (int64, error)
}

type stepRepo struct {
	db *sqlx.DB
}

func NewStepRepository(db *sqlx.DB) StepRepository {
	return &stepRepo{db}
}

func (r *stepRepo) FindByGroupId(id int64) ([]sm.StepEntity, error) {
	var steps []sm.StepEntity
	err := r.db.Select(&steps, "select * from step where step_group_id = $1", id)
	return steps, err
}

func (r *stepRepo) Save(step sm.StepEntity) (id int64, err error) {
	err = r.db.Get(
		&id,
		`insert into step (step_group_id, name, number, status, approver_order, is_approved)
     values ($1, $2, $3, $4, $5, $6) returning id`,
		step.StepGroupId,
		step.Name,
		step.Number,
		step.Status,
		step.ApproverOrder,
		step.IsApproved,
	)
	return id, err
}

func (r *stepRepo) StartSteps(
	tx *sqlx.Tx,
	group gm.StepGroupEntity,
) ([]sm.StepEntity, error) {
	rows, err := tx.NamedQuery(
		`update step 
     set status = 'STARTED'
     where step_group_id = :id 
     and (number = 1 and 'SERIAL' = :step_order or 'SERIAL' <> :step_order)
     returning *`,
		group,
	)
	var saved []sm.StepEntity
	for err == nil && rows.Next() {
		var step sm.StepEntity
		err = rows.StructScan(&step)
		if rows.StructScan(&step) == nil {
			saved = append(saved, step)
		}
	}
	return saved, err
}

func (r *stepRepo) Update(step sm.StepEntity) (stepId int64, err error) {
	_, err = r.db.NamedExec(
		`update step 
     set name = :name,
       number = :number,
       approver_order = :approver_order
     where id = :id`,
		step,
	)
	if err == nil {
		stepId = step.Id
	}
	return stepId, err
}

func (r *stepRepo) IsRouteStarted(stepId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (
       select 1 from step s 
       inner join step_group g on g.id = s.step_group_id
       inner join route r on r.id = g.route_id
       where s.id = $1 and r.status in ('STARTED', 'FINISHED'))`,
		stepId,
	)
	return res, err
}

func (r *stepRepo) FinishStepsByRouteId(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec(
		`update step 
     set status = 'FINISHED'
     where step.step_group_id in (
       select id from step_group where route_id = $1
     )`,
		routeId,
	)
	return err
}

func (r *stepRepo) FinishStep(
	tx *sqlx.Tx,
	stepId int64,
) error {
	_, err := tx.Exec("update step set status = 'FINISHED' where id = $1", stepId)
	return err
}

func (r *stepRepo) CalculateAndSetIsApproved(
	tx *sqlx.Tx,
	stepId int64,
	approverOrder common.OrderType,
	isResolutionApproved bool,
) (res bool, err error) {
	err = tx.Get(
		&res,
		`update step
			set is_approved = (
				case
					when $1 = 'PARALLEL_ANY_OF' and not $2 then exists (
						select 1
						from resolution r
						inner join approver a on r.approver_id = a.id
						inner join step s on a.step_id = s.id
						where s.id = $3
						and r.is_approved = true
					)
				  else $2
				end
			)
			where id = $3
			returning is_approved`,
		approverOrder,
		isResolutionApproved,
		stepId,
	)
	return res, err
}

func (r *stepRepo) ExistsNotFinishedStepsInGroup(
	tx *sqlx.Tx,
	stepGroupId int64,
) (res bool, err error) {
	err = tx.Select(
		&res,
		"select exists (select 1 from step where step_group_id = $1 and status != 'FINISHED')",
		stepGroupId,
	)
	return res, err
}

func (r *stepRepo) StartNextStep(
	tx *sqlx.Tx,
	stepGroupId int64,
	stepId int64,
) (nextStepId int64, err error) {
	err = tx.Get(
		&nextStepId,
		`update step 
     set status = 'STARTED'
     where step_group_id = $1
     and number = (select number + 1 from step where id = $2)
     returning id`,
		stepGroupId,
		stepId,
	)
	return nextStepId, err
}
