package model

type RouteRequest struct {
	Name        string `json:"name"        validate:"required,min=3,max=155"`
	Description string `json:"description" validate:"max=255"`
}
