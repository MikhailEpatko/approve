package model

import "time"

type Status string

const (
	TEMPLATE Status = "TEMPLATE"
	NEW      Status = "NEW"
	STARTED  Status = "STARTED"
	FINISHED Status = "FINISHED"
)

type Route struct {
	Id          int
	Name        string
	Description string
	Status      Status
	CreatedAt   time.Time
	ModifiedAt  time.Time
	Deleted     bool
}

func NewRoute(
	name,
	description string,
) *Route {
	return &Route{
		Name:        name,
		Description: description,
		Status:      TEMPLATE,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Deleted:     false,
	}
}
