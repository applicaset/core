package main

import "context"

type Item interface {
	GetID() string
}

type GenericItem map[string]interface{}

func (item GenericItem) GetID() string {
	return item["id"].(string)
}

type Service interface {
	List(ctx context.Context) (res []GenericItem, err error)
	Create(ctx context.Context, req GenericItem) (err error)
	Read(ctx context.Context, id string) (res GenericItem, err error)
	Replace(ctx context.Context, id string, req GenericItem) (err error)
	Delete(ctx context.Context, id string) (err error)
}
