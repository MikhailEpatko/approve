package model

import (
	"approve/internal/common"
	"strings"
)

var templatePrefix = "Template "

type CreateRouteTemplateRequest struct {
	Name        string `json:"name"         validate:"required,min=3,max=155"`
	Description string `json:"description"  validate:"max=255"`
}

func (r *CreateRouteTemplateRequest) ToEntity() RouteEntity {
	return RouteEntity{
		Name:        addTemplatePrefix(r.Name),
		Description: r.Description,
		Status:      common.TEMPLATE,
	}
}

func addTemplatePrefix(name string) string {
	if !strings.HasPrefix(name, templatePrefix) {
		name = templatePrefix + name
	}
	return name
}

type UpdateRouteRequest struct {
	Id          int64  `json:"route_id"     validate:"required,min=1"`
	Name        string `json:"name"         validate:"required,min=3,max=155"`
	Description string `json:"description"  validate:"max=255"`
}

func (r *UpdateRouteRequest) ToEntity(status common.Status) RouteEntity {
	return RouteEntity{
		Id:          r.Id,
		Name:        checkTemplatePrefix(r.Name, status),
		Description: r.Description,
	}
}

func checkTemplatePrefix(name string, status common.Status) string {
	if status == common.TEMPLATE && !strings.HasPrefix(name, templatePrefix) {
		name = templatePrefix + name
	}
	if status != common.TEMPLATE && strings.HasPrefix(name, templatePrefix) {
		name, _ = strings.CutPrefix(name, templatePrefix)
	}
	return name
}
