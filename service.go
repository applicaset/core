package core

import (
	"context"
	"strings"
)

type Item interface {
	GetID() string
}

type GenericItem map[string]interface{}

func (item GenericItem) GetID() string {
	return item["id"].(string)
}

type Service interface {
	List(ctx context.Context, groupKind string) (res []GenericItem, err error)
	Create(ctx context.Context, groupKind string, req GenericItem) (err error)
	Read(ctx context.Context, groupKind string, id string) (res GenericItem, err error)
	Replace(ctx context.Context, groupKind string, id string, req GenericItem) (err error)
	Delete(ctx context.Context, groupKind string, id string) (err error)
}

func GetGroupKind(group, kind string) string {
	return strings.Join([]string{group, kind}, "/")
}

func GetGroupAndKind(groupKind string) (group, kind string) {
	parts := strings.Split(groupKind, "/")
	if len(parts) == 1 {
		return "", parts[0]
	}

	return parts[0], parts[1]
}
