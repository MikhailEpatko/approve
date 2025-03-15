package repository

import (
	"approve/internal/common"
	rm "approve/internal/route/model"
	"github.com/jmoiron/sqlx"
	"strings"
)

type FindByFilterRouteRepository interface {
	FindByfilter(filter rm.FilterRouteRequest) ([]rm.RouteEntity, int64, error)
}

func NewFindByFilterRouteRepository(db *sqlx.DB) FindByFilterRouteRepository {
	return &routeRepo{db}
}

type params map[string]interface{}

func (r *routeRepo) FindByfilter(filter rm.FilterRouteRequest) ([]rm.RouteEntity, int64, error) {
	var routes []rm.RouteEntity
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

func findByFilterQueryAndParams(filter rm.FilterRouteRequest) (string, params) {
	params := make(params, 5)
	sb := strings.Builder{}
	sb.WriteString("select r.* from route r")
	if filter.Status != common.TEMPLATE {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id
                     where r.status = :status`)
	} else {
		sb.WriteString(" where r.status = :status")
	}
	params["status"] = filter.Status
	if filter.Guid != "" {
		sb.WriteString(" and a.guid = :guid")
		params["guid"] = filter.Guid
	}
	if len(filter.Text) >= 3 {
		sb.WriteString(` and (r.name ilike '%' || :text || '%'
                     or r.description ilike '%' || :text || '%')`)
		params["text"] = filter.Text
	}
	sb.WriteString(" order by r.id desc offset :offset limit :limit")
	params["offset"] = (filter.PageNumber - 1) * filter.PageSize
	params["limit"] = filter.PageSize
	return sb.String(), params
}

func countByFilterQueryAndParams(filter rm.FilterRouteRequest) (string, params) {
	params := make(params, 5)
	sb := strings.Builder{}
	sb.WriteString("select count(r.*) from route r")
	if filter.Status != common.TEMPLATE {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id`)
		sb.WriteString(" where r.status = :status and a.guid = :guid")
		params["guid"] = filter.Guid
	} else {
		sb.WriteString(" where r.status = :status")
	}
	params["status"] = filter.Status
	if len(filter.Text) >= 3 {
		sb.WriteString(` and (r.name ilike '%' || :text || '%'
                     or r.description ilike '%' || :text || '%')`)
		params["text"] = filter.Text
	}
	return sb.String(), params
}
