package repository

import (
	cfg "approve/internal/database"
	rm "approve/internal/route/model"
	"strings"
)

type parameters map[string]any

func FindByfilter(filter rm.FilterRouteRequest) ([]rm.RouteEntity, int64, error) {
	var routes []rm.RouteEntity
	var total int64
	query, params := countByFilterQueryAndParams(filter)
	rows, err := cfg.DB.NamedQuery(query, params)
	if err != nil {
		return nil, 0, err
	}
	if rows.Next() {
		err = rows.Scan(&total)
	}
	if err != nil {
		return nil, 0, err
	}
	if total == 0 || int(total) <= filter.PageSize*(filter.PageNumber-1) {
		return routes, 0, nil
	}
	query, params = findByFilterQueryAndParams(filter)
	rows, err = cfg.DB.NamedQuery(query, params)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var route rm.RouteEntity
		err = rows.StructScan(&route)
		if err != nil {
			return nil, 0, err
		}
		routes = append(routes, route)
	}
	return routes, total, nil
}

func findByFilterQueryAndParams(filter rm.FilterRouteRequest) (string, parameters) {
	params := make(parameters, 5)
	sb := strings.Builder{}
	sb.WriteString("select distinct r.* from route r")
	if filter.Guid != "" {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id`)
	}
	sb.WriteString(" where 1=1")
	if filter.Status != "" {
		sb.WriteString(" and r.status = :status")
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

func countByFilterQueryAndParams(filter rm.FilterRouteRequest) (string, parameters) {
	params := make(parameters, 5)
	sb := strings.Builder{}
	sb.WriteString("select count(distinct r.id) from route r")
	if filter.Guid != "" {
		sb.WriteString(` inner join step_group sg on r.id = sg.route_id
                     inner join step s on sg.id = s.step_group_id 
                     inner join approver a on s.id = a.step_id`)
	}
	sb.WriteString(" where 1=1")
	if filter.Status != "" {
		sb.WriteString(" and r.status = :status")
		params["status"] = filter.Status
	}
	if filter.Guid != "" {
		sb.WriteString(" and a.guid = :guid")
		params["guid"] = filter.Guid
	}
	if len(filter.Text) >= 3 {
		sb.WriteString(` and (r.name ilike '%' || :text || '%'
                     or r.description ilike '%' || :text || '%')`)
		params["text"] = filter.Text
	}
	return sb.String(), params
}
