package model

type ApproverEntity struct {
	Id      int    `db:"id"`
	StepId  int    `db:"step_id"`
	Guid    string `db:"guid"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Number  int    `db:"number"`
	Deleted bool   `db:"deleted"`
}

func NewApproverEntity(
	stepId int,
	guid string,
	name string,
	email string,
	number int,
) *ApproverEntity {
	return &ApproverEntity{
		StepId:  stepId,
		Guid:    guid,
		Name:    name,
		Email:   email,
		Number:  number,
		Deleted: false,
	}
}
