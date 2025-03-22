package repository

import (
	cm "approve/internal/common"
	cfg "approve/internal/config"
	gm "approve/internal/stepgroup/model"

	"github.com/jmoiron/sqlx"
)

func FindById(id int64) (group gm.StepGroupEntity, err error) {
	err = cfg.DB.Get(&group, "select * from step_group where id = $1", id)
	return group, err
}

func FindByRouteId(id int64) ([]gm.StepGroupEntity, error) {
	var groups []gm.StepGroupEntity
	err := cfg.DB.Select(&groups, "select * from step_group where route_id = $1", id)
	return groups, err
}

func Save(stepGroup gm.StepGroupEntity) (id int64, err error) {
	err = cfg.DB.Get(
		&id,
		`insert into step_group (route_id, name, number, status, step_order, is_approved)
     values ($1, $2, $3, $4, $5, $6)
     returning id`,
		stepGroup.RouteId,
		stepGroup.Name,
		stepGroup.Number,
		stepGroup.Status,
		stepGroup.StepOrder,
		stepGroup.IsApproved,
	)
	return id, err
}

func StartFirstGroup(
	tx *sqlx.Tx,
	routeId int64,
) (group gm.StepGroupEntity, err error) {
	rows, err := tx.Queryx(
		`update step_group 
     set status = 'STARTED'
     where route_id = $1 and number = 1
     returning *`,
		routeId,
	)
	if err == nil && rows.Next() {
		err = rows.StructScan(&group)
	}
	return group, err
}

func Update(group gm.StepGroupEntity) (groupId int64, err error) {
	_, err = cfg.DB.NamedExec(
		`update step_group 
     set name = :name,
       number = :number,
       step_order = :step_order
     where id = :id`,
		group,
	)
	if err == nil {
		groupId = group.Id
	}
	return groupId, err
}

func IsRouteProcessing(stepGroupId int64) (res bool, err error) {
	err = cfg.DB.Get(
		&res,
		`select exists (
       select 1 from step_group g 
       inner join route r on r.id = g.route_id
       where g.id = $1 and r.status in ('STARTED', 'FINISHED'))`,
		stepGroupId,
	)
	return res, err
}

func FinishGroup(
	tx *sqlx.Tx,
	stepGroupId int64,
) error {
	_, err := tx.Exec("update step_group set status = 'FINISHED' where id = $1", stepGroupId)
	return err
}

func CalculateAndSetIsApproved(
	tx *sqlx.Tx,
	stepGroupId int64,
	stepOrder cm.OrderType,
	isStepApproved bool,
) (res bool, err error) {
	err = tx.Get(
		&res,
		`update step_group
		 set is_approved = (
	     case
			   when $1 = 'PARALLEL_ANY_OF' and not $2 then exists (
					 select 1 
					 from step s
					 inner join step_group g on g.id = s.step_group_id
					 where g.id = $3 and s.is_approved = true
				 )
			   else $2
			 end
		 )
		 where id = $3
		 returning is_approved`,
		stepOrder,
		isStepApproved,
		stepGroupId,
	)
	return res, err
}

func StartNextGroup(
	tx *sqlx.Tx,
	routeId int64,
	stepGroupId int64,
) (nextGroupId int64, err error) {
	err = tx.Get(
		&nextGroupId,
		`update step_group 
     set status = 'STARTED'
     where route_id = $1
     and number = (select number + 1 from step_group where id = $2)
     returning id`,
		routeId,
		stepGroupId,
	)
	return nextGroupId, err
}
