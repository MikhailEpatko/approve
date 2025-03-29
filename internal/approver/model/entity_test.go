package model

import (
	cm "approve/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApproverEntity(t *testing.T) {
	a := assert.New(t)

	t.Run("ToNewApprover", func(t *testing.T) {
		approver := ApproverEntity{
			Id:       1,
			StepId:   2,
			Guid:     "guid",
			Name:     "name",
			Position: "position",
			Email:    "email",
			Number:   3,
			Status:   cm.TEMPLATE,
		}

		newApprover := approver.ToNewApprover(4)

		a.Equal(int64(0), newApprover.Id)
		a.Equal(approver.Name, newApprover.Name)
		a.Equal(int64(4), newApprover.StepId)
		a.Equal(approver.Guid, newApprover.Guid)
		a.Equal(approver.Position, newApprover.Position)
		a.Equal(approver.Email, newApprover.Email)
		a.Equal(approver.Number, newApprover.Number)
		a.Equal(cm.NEW, newApprover.Status)
	})
}
