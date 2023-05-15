package main

import (
	"context"
	"sync"
)

type Store struct {
	db map[string]map[string]GenericItem
	sync.Mutex
}

func (s *Store) List(ctx context.Context, kind string) ([]GenericItem, error) {
	res := make([]GenericItem, len(s.db))

	table, ok := s.db[kind]
	if !ok {
		return nil, KindNotFoundError{Name: kind}
	}

	i := 0
	for k := range table {
		res[i] = table[k]

		i++
	}

	return res, nil
}

func (s *Store) Create(ctx context.Context, kind string, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.db[kind]
	if !ok {
		s.db[kind] = make(map[string]GenericItem)
	}

	for k := range s.db[kind] {
		if req.GetID() == k {
			return ItemExistsError{ID: k}
		}
	}

	s.db[kind][req.GetID()] = req

	return nil
}

func (s *Store) Read(ctx context.Context, kind string, id string) (GenericItem, error) {
	table, ok := s.db[kind]
	if !ok {
		return nil, KindNotFoundError{Name: kind}
	}

	res, ok := table[id]
	if !ok {
		return nil, ItemNotFoundError{ID: id}
	}

	return res, nil
}

func (s *Store) Replace(ctx context.Context, kind string, id string, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	table, ok := s.db[kind]
	if !ok {
		return KindNotFoundError{Name: kind}
	}

	if _, ok := table[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	s.db[kind][id] = req

	return nil
}

func (s *Store) Delete(ctx context.Context, kind string, id string) error {
	s.Lock()
	defer s.Unlock()

	table, ok := s.db[kind]
	if !ok {
		return KindNotFoundError{Name: kind}
	}

	if _, ok := table[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	delete(s.db[kind], id)

	return nil
}

var _ Service = new(Store)

func NewStore() *Store {
	return &Store{
		db: make(map[string]map[string]GenericItem),
	}
}
