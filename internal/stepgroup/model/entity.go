package model

import (
	"approve/internal/common"
)

type StepGroupEntity struct {
	Id        int64            `db:"id"`
	RouteId   int64            `db:"route_id"`
	Name      string           `db:"name"`
	Number    int              `db:"number"`
	Status    common.Status    `db:"status"`
	StepOrder common.OrderType `db:"step_order"`
	Deleted   bool             `db:"deleted"`
}
