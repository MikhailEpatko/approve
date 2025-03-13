package common

import "github.com/jmoiron/sqlx"

type DbCleaner interface {
	ClearDb()
}

type cleaner struct {
	db *sqlx.DB
}

func NewCleaner(db *sqlx.DB) *cleaner {
	return &cleaner{db: db}
}
func (c *cleaner) ClearDb() {
	c.db.MustExec("delete from resolution")
	c.db.MustExec("delete from approver")
	c.db.MustExec("delete from step")
	c.db.MustExec("delete from step_group")
	c.db.MustExec("delete from route")
}
