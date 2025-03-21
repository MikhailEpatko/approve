package model

import cm "approve/internal/common"

type RouteEntity struct {
	Id          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      cm.Status `db:"status"`
	IsApproved  bool      `db:"is_approved"`
}

func (e RouteEntity) ToResponse() RouteResponse {
	return RouteResponse{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		Status:      e.Status,
	}
}
