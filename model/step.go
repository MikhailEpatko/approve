package model

type ApproveType string

const (
	PARALLEL_ANY_OF    ApproveType = "PARALLEL_ANY_OF"
	PARALLEL_ALL_OF    ApproveType = "PARALLEL_ALL_OF"
	SEQUENTIAL_ALL_OFF ApproveType = "SEQUENTIAL_ALL_OFF"
)

type StepEntity struct {
	Id          int         `db:"id"`
	StepGroupId int         `db:"step_group_id"`
	Name        string      `db:"name"`
	Number      int         `db:"number"`
	Status      Status      `db:"status"`
	ApproveType ApproveType `db:"approve_type"`
	Deleted     bool        `db:"deleted"`
}

func NewStepEntity(
	stepGroupId int,
	name string,
	number int,
	approveType ApproveType,
) *StepEntity {
	return &StepEntity{
		StepGroupId: stepGroupId,
		Name:        name,
		Number:      number,
		Status:      NEW,
		ApproveType: approveType,
		Deleted:     false,
	}
}
