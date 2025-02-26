package model

import "time"

type Status string

const (
	TEMPLATE Status = "TEMPLATE"
	NEW      Status = "NEW"
	STARTED  Status = "STARTED"
	FINISHED Status = "FINISHED"
)

type RouteEntity struct {
	Id          int
	Name        string
	Description string
	Status      Status
	CreatedAt   time.Time
	ModifiedAt  time.Time
	Deleted     bool
}

func NewRouteEntity(
	name,
	description string,
) *RouteEntity {
	return &RouteEntity{
		Name:        name,
		Description: description,
		Status:      TEMPLATE,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Deleted:     false,
	}
}
