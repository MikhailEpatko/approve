package model

import (
	cm "approve/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouteRequest(t *testing.T) {
	a := assert.New(t)

	t.Run("addTemplatePrefix should add prefix", func(t *testing.T) {
		s := "name loren template ipsum Template"
		want := "Template name loren template ipsum Template"
		got := addTemplatePrefix(s)
		a.Equal(want, got)
	})

	t.Run("addTemplatePrefix should not add prefix", func(t *testing.T) {
		want := "Template name loren template ipsum Template"
		got := addTemplatePrefix(want)
		a.Equal(want, got)
	})

	t.Run("checkTemplatePrefix - if status is TEMPLATE then should add prefix", func(t *testing.T) {
		s := "name loren template ipsum Template"
		want := "Template name loren template ipsum Template"
		got := checkTemplatePrefix(s, cm.TEMPLATE)
		a.Equal(want, got)
	})

	t.Run("checkTemplatePrefix - if status is TEMPLATE then should not add prefix", func(t *testing.T) {
		want := "Template name loren template ipsum Template"
		got := checkTemplatePrefix(want, cm.TEMPLATE)
		a.Equal(want, got)
	})

	t.Run("checkTemplatePrefix - if status is NEW then should not add prefix", func(t *testing.T) {
		want := "name"
		got := checkTemplatePrefix(want, cm.NEW)
		a.Equal(want, got)
	})

	t.Run("checkTemplatePrefix - if status is NEW then should remove prefix", func(t *testing.T) {
		s := "Template name"
		want := "name"
		got := checkTemplatePrefix(s, cm.NEW)
		a.Equal(want, got)
	})
}
