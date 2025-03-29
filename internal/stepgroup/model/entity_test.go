package model

import (
	cm "approve/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStepGroupEntity(t *testing.T) {
	a := assert.New(t)

	t.Run("ToNewStepGroup", func(t *testing.T) {
		stepGroup := StepGroupEntity{
			Id:         1,
			RouteId:    2,
			Name:       "name",
			Number:     3,
			Status:     cm.TEMPLATE,
			StepOrder:  cm.PARALLEL_ALL_OF,
			IsApproved: false,
		}

		newStepGroup := stepGroup.ToNewStepGroup(4)

		a.Equal(int64(0), newStepGroup.Id)
		a.Equal(int64(4), newStepGroup.RouteId)
		a.Equal(stepGroup.Name, newStepGroup.Name)
		a.Equal(stepGroup.Number, newStepGroup.Number)
		a.Equal(stepGroup.StepOrder, newStepGroup.StepOrder)
		a.Equal(stepGroup.IsApproved, newStepGroup.IsApproved)
		a.Equal(cm.NEW, newStepGroup.Status)
	})
}
