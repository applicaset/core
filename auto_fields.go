package main

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AutoFields struct {
	next Service
}

func (af *AutoFields) List(ctx context.Context, kind string) ([]GenericItem, error) {
	return af.next.List(ctx, kind)
}

func (af *AutoFields) Create(ctx context.Context, kind string, req GenericItem) error {
	req["uuid"] = uuid.NewString()
	req["createdAt"] = time.Now().Format(time.RFC3339)
	req["updatedAt"] = req["createdAt"]
	req["kind"] = kind

	return af.next.Create(ctx, kind, req)
}

func (af *AutoFields) Read(ctx context.Context, kind string, id string) (GenericItem, error) {
	return af.next.Read(ctx, kind, id)
}

func (af *AutoFields) Replace(ctx context.Context, kind string, id string, req GenericItem) error {
	req["updatedAt"] = time.Now().Format(time.RFC3339)
	req["kind"] = kind

	return af.next.Replace(ctx, kind, id, req)
}

func (af *AutoFields) Delete(ctx context.Context, kind string, id string) error {
	return af.next.Delete(ctx, kind, id)
}

var _ Service = new(AutoFields)

func NewAutoFields(next Service) *AutoFields {
	return &AutoFields{next: next}
}
