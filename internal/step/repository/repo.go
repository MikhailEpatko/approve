package repository

import (
	"approve/internal/common"
	cfg "approve/internal/config"
	sm "approve/internal/step/model"
	gm "approve/internal/stepgroup/model"

	"github.com/jmoiron/sqlx"
)

func FindById(stepId int64) (step sm.StepEntity, err error) {
	err = cfg.DB.Get(
		&step,
		`select * from step where id = $1`,
		stepId,
	)
	return step, err
}

func FindByGroupId(id int64) ([]sm.StepEntity, error) {
	var steps []sm.StepEntity
	err := cfg.DB.Select(&steps, "select * from step where step_group_id = $1", id)
	return steps, err
}

func Save(step sm.StepEntity) (id int64, err error) {
	err = cfg.DB.Get(
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

func StartSteps(
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

func Update(step sm.StepEntity) (stepId int64, err error) {
	_, err = cfg.DB.NamedExec(
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

func IsRouteStarted(stepId int64) (res bool, err error) {
	err = cfg.DB.Get(
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

func FinishStep(
	tx *sqlx.Tx,
	stepId int64,
) error {
	_, err := tx.Exec("update step set status = 'FINISHED' where id = $1", stepId)
	return err
}

func CalculateAndSetIsApproved(
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

func ExistsNotFinishedStepsInGroup(
	tx *sqlx.Tx,
	stepGroupId int64,
) (res bool, err error) {
	err = tx.Get(
		&res,
		"select exists (select 1 from step where step_group_id = $1 and status != 'FINISHED')",
		stepGroupId,
	)
	return res, err
}

func StartNextStep(
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
