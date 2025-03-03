package repository

import (
	"approve/internal/common"
	. "approve/internal/route/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type RouteRepository interface {
	FindByfilter(filter FilterRouteRequest) ([]RouteEntity, int64, error)
	Save(name, description string) (int64, error)
	CopyAsNew(routeTemplateId int64) (int64, error)
	Update(route RouteEntity) (int64, error)
	SaveTx(tx *sqlx.Tx, route RouteEntity) (int64, error)
}

type routeRepo struct {
	db *sqlx.DB
}

type params map[string]interface{}

func NewRouteRepository(db *sqlx.DB) RouteRepository {
	return &routeRepo{db}
}

func (r *routeRepo) FindByfilter(filter FilterRouteRequest) ([]RouteEntity, int64, error) {
	var routes []RouteEntity
	var total int64
	query, params := countByFilterQueryAndParams(filter)
	err := r.db.Select(&total, query, params)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 || int(total) <= filter.PageSize*(filter.PageNumber-1) {
		return routes, 0, nil
	}
	query, params = findByFilterQueryAndParams(filter)
	err = r.db.Select(&routes, query, params)
	if err != nil {
		return nil, 0, err
	}
	return routes, total, nil
}

func findByFilterQueryAndParams(filter FilterRouteRequest) (string, params) {
	params := make(params, 5)
	sb := strings.Builder{}
	sb.WriteString("select r.* from route r")
	if filter.Status != common.TEMPLATE {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id`)
	}
	sb.WriteString(" where r.status = :status and r.deleted = false")
	params["status"] = filter.Status
	if len(filter.Text) >= 3 {
		sb.WriteString(` and (r.name ilike '%' || :text || '%'
                     or r.description ilike '%' || :text || '%')`)
		params["text"] = filter.Text
	}
	if filter.Status != common.TEMPLATE {
		sb.WriteString(" and a.guid = :guid")
		params["guid"] = filter.Guid
	}
	sb.WriteString(" order by r.id desc offset :offset limit :limit")
	params["offset"] = (filter.PageNumber - 1) * filter.PageSize
	params["limit"] = filter.PageSize
	return sb.String(), params
}

func countByFilterQueryAndParams(filter FilterRouteRequest) (string, params) {
	params := make(params, 5)
	sb := strings.Builder{}
	sb.WriteString("select count(r.*) from route r")
	if filter.Status != common.TEMPLATE {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id`)
	}
	sb.WriteString(" where r.status = :status and r.deleted = false")
	params["status"] = filter.Status
	if len(filter.Text) >= 3 {
		sb.WriteString(` and (r.name ilike '%' || :text || '%'
                     or r.description ilike '%' || :text || '%')`)
		params["text"] = filter.Text
	}
	if filter.Status != common.TEMPLATE {
		sb.WriteString(" and a.guid = :guid")
		params["guid"] = filter.Guid
	}
	return sb.String(), params
}

func (r *routeRepo) Save(name, description string) (int64, error) {
	res, err := r.db.NamedExec(
		"insert into route (name, description) values (:name, :description)",
		map[string]interface{}{
			"name":        name,
			"description": description,
		},
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *routeRepo) SaveTx(
	tx *sqlx.Tx,
	route RouteEntity,
) (int64, error) {
	res, err := tx.NamedExec(
		"insert into route (name, description, status) values (:name, :description, :status)",
		route,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *routeRepo) CopyAsNew(routeTemplateId int64) (int64, error) {
	res, err := r.db.Exec(
		`insert into route (name, description, status)
     select name, description, 'NEW' from route where id = $1`,
		routeTemplateId,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *routeRepo) Update(route RouteEntity) (int64, error) {
	res, err := r.db.NamedExec(
		`update route 
     set
			 name = :name, 
			 description = :description
		 where id = :id`,
		route,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
