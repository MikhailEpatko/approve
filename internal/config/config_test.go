package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildServer(t *testing.T) {
	a := assert.New(t)
	err := BuildServer()
	a.Nil(err)
}
