package main

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AutoFields struct {
	next Service
}

func (af *AutoFields) List(ctx context.Context, groupKind string) ([]GenericItem, error) {
	return af.next.List(ctx, groupKind)
}

func (af *AutoFields) Create(ctx context.Context, groupKind string, req GenericItem) error {
	req["uuid"] = uuid.NewString()
	req["createdAt"] = time.Now().Format(time.RFC3339)
	req["updatedAt"] = req["createdAt"]

	group, kind := GetGroupAndKind(groupKind)
	req["group"] = group
	req["kind"] = kind

	return af.next.Create(ctx, groupKind, req)
}

func (af *AutoFields) Read(ctx context.Context, groupKind string, id string) (GenericItem, error) {
	return af.next.Read(ctx, groupKind, id)
}

func (af *AutoFields) Replace(ctx context.Context, groupKind string, id string, req GenericItem) error {
	req["updatedAt"] = time.Now().Format(time.RFC3339)

	group, kind := GetGroupAndKind(groupKind)
	req["group"] = group
	req["kind"] = kind

	return af.next.Replace(ctx, groupKind, id, req)
}

func (af *AutoFields) Delete(ctx context.Context, groupKind string, id string) error {
	return af.next.Delete(ctx, groupKind, id)
}

var _ Service = new(AutoFields)

func NewAutoFields(next Service) *AutoFields {
	return &AutoFields{next: next}
}
