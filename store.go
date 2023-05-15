package main

import (
	"context"
	"sync"
)

type Store struct {
	records map[string]GenericItem
	sync.Mutex
}

func (s *Store) List(ctx context.Context) ([]GenericItem, error) {
	res := make([]GenericItem, len(s.records))

	i := 0
	for k := range s.records {
		res[i] = s.records[k]

		i++
	}

	return res, nil
}

func (s *Store) Create(ctx context.Context, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	for k := range s.records {
		if req.GetID() == k {
			return ItemExistsError{ID: k}
		}
	}

	s.records[req.GetID()] = req

	return nil
}

func (s *Store) Read(ctx context.Context, id string) (GenericItem, error) {
	res, ok := s.records[id]
	if !ok {
		return nil, ItemNotFoundError{ID: id}
	}

	return res, nil
}

func (s *Store) Replace(ctx context.Context, id string, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.records[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	s.records[id] = req

	return nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.records[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	delete(s.records, id)

	return nil
}

var _ Service = new(Store)

func NewStore() *Store {
	return &Store{
		records: make(map[string]GenericItem),
	}
}
