package model

type ApproverEntity struct {
	Id       int64  `db:"id"`
	StepId   int64  `db:"step_id"`
	Guid     string `db:"guid"`
	Name     string `db:"name"`
	Position string `db:"position"`
	Email    string `db:"email"`
	Number   int    `db:"number"`
	Deleted  bool   `db:"deleted"`
}
