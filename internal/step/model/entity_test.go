package model

import (
	cm "approve/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepEntity(t *testing.T) {
	a := assert.New(t)

	t.Run("ToNewStep", func(t *testing.T) {
		step := StepEntity{
			Id:            1,
			StepGroupId:   2,
			Name:          "name",
			Number:        3,
			Status:        cm.TEMPLATE,
			ApproverOrder: cm.SERIAL,
			IsApproved:    false,
		}

		newStep := step.ToNewStep(4)

		a.Equal(int64(0), newStep.Id)
		a.Equal(int64(4), newStep.StepGroupId)
		a.Equal(step.Name, newStep.Name)
		a.Equal(step.Number, newStep.Number)
		a.Equal(step.ApproverOrder, newStep.ApproverOrder)
		a.Equal(step.IsApproved, newStep.IsApproved)
		a.Equal(cm.NEW, newStep.Status)
	})
}
