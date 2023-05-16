package core_test

import (
	"github.com/applicaset/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGroupKind(t *testing.T) {
	group := "group1"
	kind := "kind1"

	groupKind := core.GetGroupKind(group, kind)
	group2, kind2 := core.GetGroupAndKind(groupKind)

	assert.Equal(t, group, group2)
	assert.Equal(t, kind, kind2)
}
