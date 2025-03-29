package model

import (
	cm "approve/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteEntity(t *testing.T) {
	a := assert.New(t)

	t.Run("ToNewRoute", func(t *testing.T) {
		route := RouteEntity{
			Id:          1,
			Name:        "Template name",
			Description: "description",
			Status:      cm.TEMPLATE,
			IsApproved:  false,
		}

		newRoute := route.ToNewRoute()

		a.Equal(int64(0), newRoute.Id)
		a.Equal("name", newRoute.Name)
		a.Equal(route.Description, newRoute.Description)
		a.Equal(cm.NEW, newRoute.Status)
		a.Equal(route.IsApproved, newRoute.IsApproved)
	})

	t.Run("removeTemplatePrefix", func(t *testing.T) {
		s := "Template name loren template ipsum Template"
		want := "name loren template ipsum Template"

		got := removeTemplatePrefix(s)

		a.Equal(want, got)
	})

	t.Run("removeTemplatePrefix", func(t *testing.T) {
		want := "name loren template ipsum Template"

		got := removeTemplatePrefix(want)

		a.Equal(want, got)
	})
}
