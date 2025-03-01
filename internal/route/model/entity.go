package model

import "approve/internal/common"

type RouteEntity struct {
	Id          int64         `db:"id"`
	Name        string        `db:"name"`
	Description string        `db:"description"`
	Status      common.Status `db:"status"`
	Deleted     bool          `db:"deleted"`
}

func NewRouteEntity(
	name,
	description string,
) *RouteEntity {
	return &RouteEntity{
		Name:        name,
		Description: description,
		Status:      common.TEMPLATE,
		Deleted:     false,
	}
}

func (e RouteEntity) ToResponse() RouteResponse {
	return RouteResponse{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		Status:      e.Status,
	}
}
