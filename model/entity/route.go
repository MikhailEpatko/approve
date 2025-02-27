package entity

type Status string

const (
	TEMPLATE Status = "TEMPLATE"
	NEW      Status = "NEW"
	STARTED  Status = "STARTED"
	FINISHED Status = "FINISHED"
)

type RouteEntity struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Status      Status `db:"status"`
	Deleted     bool   `db:"deleted"`
}

func NewRouteEntity(
	name,
	description string,
) *RouteEntity {
	return &RouteEntity{
		Name:        name,
		Description: description,
		Status:      TEMPLATE,
		Deleted:     false,
	}
}
