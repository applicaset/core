package main

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type AutoFields struct {
	next Service
}

func (af *AutoFields) List(ctx context.Context) ([]GenericItem, error) {
	return af.next.List(ctx)
}

func (af *AutoFields) Create(ctx context.Context, req GenericItem) error {
	req["uuid"] = uuid.NewString()
	req["createdAt"] = time.Now().Format(time.RFC3339)
	req["updatedAt"] = req["createdAt"]

	return af.next.Create(ctx, req)
}

func (af *AutoFields) Read(ctx context.Context, id string) (GenericItem, error) {
	return af.next.Read(ctx, id)
}

func (af *AutoFields) Replace(ctx context.Context, id string, req GenericItem) error {
	req["updatedAt"] = time.Now().Format(time.RFC3339)

	return af.next.Replace(ctx, id, req)
}

func (af *AutoFields) Delete(ctx context.Context, id string) error {
	return af.next.Delete(ctx, id)
}

var _ Service = new(AutoFields)

func NewAutoFields(next Service) *AutoFields {
	return &AutoFields{next: next}
}
