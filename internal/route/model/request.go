package model

import (
	"approve/internal/common"
)

type CreateRouteTemplateRequest struct {
	Name        string `json:"name"         validate:"required,min=3,max=155"`
	Description string `json:"description"  validate:"max=255"`
}

func (r CreateRouteTemplateRequest) ToEntity() RouteEntity {
	return RouteEntity{
		Name:        r.Name,
		Description: r.Description,
		Status:      common.TEMPLATE,
	}
}
