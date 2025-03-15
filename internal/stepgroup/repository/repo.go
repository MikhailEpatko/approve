package repository

import (
	cm "approve/internal/common"
	gm "approve/internal/stepgroup/model"

	"github.com/jmoiron/sqlx"
)

type StepGroupRepository interface {
	FindByRouteId(id int64) ([]gm.StepGroupEntity, error)
	Save(stepGroup gm.StepGroupEntity) (int64, error)
	StartGroups(tx *sqlx.Tx, routeId int64) (gm.StepGroupEntity, error)
	Update(group gm.StepGroupEntity) (int64, error)
	IsRouteStarted(stepGroupId int64) (bool, error)
	FinishGroupsByRouteId(tx *sqlx.Tx, routeId int64) error
	FinishGroup(tx *sqlx.Tx, stepGroupId int64) error
	CalculateAndSetIsApproved(
		tx *sqlx.Tx,
		stepGroupId int64,
		stepOrder cm.OrderType,
		isStepApproved bool,
	) (bool, error)
	ExistsNotFinishedStepGroupsInRoute(tx *sqlx.Tx, routeId int64) (bool, error)
	StartNextGroup(tx *sqlx.Tx, routeId int64, stepGroupId int64) (int64, error)
}

type stepGroupRepo struct {
	db *sqlx.DB
}

func NewStepGroupRepository(db *sqlx.DB) StepGroupRepository {
	return &stepGroupRepo{db}
}

func (r *stepGroupRepo) FindByRouteId(id int64) ([]gm.StepGroupEntity, error) {
	var groups []gm.StepGroupEntity
	err := r.db.Select(&groups, "select * from step_group where route_id = $1", id)
	return groups, err
}

func (r *stepGroupRepo) Save(stepGroup gm.StepGroupEntity) (id int64, err error) {
	err = r.db.Get(
		&id,
		`insert into step_group (route_id, name, number, status, step_order)
     values ($1, $2, $3, $4, $5)
     returning id`,
		stepGroup.RouteId,
		stepGroup.Name,
		stepGroup.Number,
		stepGroup.Status,
		stepGroup.StepOrder,
	)
	return id, err
}

func (r *stepGroupRepo) StartGroups(
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

func (r *stepGroupRepo) Update(group gm.StepGroupEntity) (groupId int64, err error) {
	_, err = r.db.NamedExec(
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

func (r *stepGroupRepo) IsRouteStarted(stepGroupId int64) (res bool, err error) {
	err = r.db.Get(
		&res,
		`select exists (
       select 1 from step_group g 
       inner join route r on r.id = g.route_id
       where g.id = $1 and r.status in ('STARTED', 'FINISHED'))`,
		stepGroupId,
	)
	return res, err
}

func (r *stepGroupRepo) FinishGroupsByRouteId(
	tx *sqlx.Tx,
	routeId int64,
) error {
	_, err := tx.Exec(
		`update step_group 
     set status = 'FINISHED'
     where route_id = $1`,
		routeId,
	)
	return err
}

func (r *stepGroupRepo) FinishGroup(
	tx *sqlx.Tx,
	stepGroupId int64,
) error {
	_, err := tx.Exec("update step_group set status = 'FINISHED' where id = $1", stepGroupId)
	return err
}

func (r *stepGroupRepo) CalculateAndSetIsApproved(
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

func (r *stepGroupRepo) ExistsNotFinishedStepGroupsInRoute(
	tx *sqlx.Tx,
	routeId int64,
) (res bool, err error) {
	err = tx.Select(
		&res,
		"select exists (select 1 from step_group where route_id = $1 and status != 'FINISHED')",
		routeId,
	)
	return res, err
}

func (r *stepGroupRepo) StartNextGroup(
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
