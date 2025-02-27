package model

type StepGroupStatus string
type StepType string

const (
	PARALLEL   StepType = "PARALLEL"
	SEQUENTIAL StepType = "SEQUENTIAL"
)

type StepGroupEntity struct {
	Id       int      `db:"id"`
	RouteId  int      `db:"route_id"`
	Name     string   `db:"name"`
	Number   int      `db:"number"`
	Status   Status   `db:"status"`
	StepType StepType `db:"step_type"`
	Deleted  bool     `db:"deleted"`
}

func NewStepGroupEntity(
	routeId int,
	name string,
	number int,
	stepType StepType,
) *StepGroupEntity {
	return &StepGroupEntity{
		RouteId:  routeId,
		Name:     name,
		Number:   number,
		Status:   NEW,
		StepType: stepType,
		Deleted:  false,
	}
}
