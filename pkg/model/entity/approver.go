package entity

type ApproverEntity struct {
	Id      int64  `db:"id"`
	StepId  int64  `db:"step_id"`
	Guid    string `db:"guid"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Number  int    `db:"number"`
	Deleted bool   `db:"deleted"`
}

func NewApproverEntity(
	stepId int64,
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
