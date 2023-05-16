package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGroupKind(t *testing.T) {
	group := "group1"
	kind := "kind1"

	groupKind := GetGroupKind(group, kind)
	group2, kind2 := GetGroupAndKind(groupKind)

	assert.Equal(t, group, group2)
	assert.Equal(t, kind, kind2)
}
