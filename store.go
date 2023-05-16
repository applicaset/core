package core

import (
	"context"
	"sync"
)

type Store struct {
	db map[string]map[string]GenericItem
	sync.Mutex
}

func (s *Store) List(_ context.Context, groupKind string) ([]GenericItem, error) {
	res := make([]GenericItem, len(s.db))

	table, ok := s.db[groupKind]
	if !ok {
		group, kind := GetGroupAndKind(groupKind)
		return nil, GroupKindNotFoundError{
			Group: group,
			Kind:  kind,
		}
	}

	i := 0
	for k := range table {
		res[i] = table[k]

		i++
	}

	return res, nil
}

func (s *Store) Create(_ context.Context, groupKind string, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.db[groupKind]
	if !ok {
		s.db[groupKind] = make(map[string]GenericItem)
	}

	for k := range s.db[groupKind] {
		if req.GetID() == k {
			return ItemExistsError{ID: k}
		}
	}

	s.db[groupKind][req.GetID()] = req

	return nil
}

func (s *Store) Read(_ context.Context, groupKind string, id string) (GenericItem, error) {
	table, ok := s.db[groupKind]
	if !ok {
		group, kind := GetGroupAndKind(groupKind)
		return nil, GroupKindNotFoundError{
			Group: group,
			Kind:  kind,
		}
	}

	res, ok := table[id]
	if !ok {
		return nil, ItemNotFoundError{ID: id}
	}

	return res, nil
}

func (s *Store) Replace(_ context.Context, groupKind string, id string, req GenericItem) error {
	s.Lock()
	defer s.Unlock()

	table, ok := s.db[groupKind]
	if !ok {
		group, kind := GetGroupAndKind(groupKind)
		return GroupKindNotFoundError{
			Group: group,
			Kind:  kind,
		}
	}

	if _, ok := table[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	s.db[groupKind][id] = req

	return nil
}

func (s *Store) Delete(_ context.Context, groupKind string, id string) error {
	s.Lock()
	defer s.Unlock()

	table, ok := s.db[groupKind]
	if !ok {
		group, kind := GetGroupAndKind(groupKind)
		return GroupKindNotFoundError{
			Kind:  kind,
			Group: group,
		}
	}

	if _, ok := table[id]; !ok {
		return ItemNotFoundError{ID: id}
	}

	delete(s.db[groupKind], id)

	return nil
}

var _ Service = new(Store)

func NewStore() *Store {
	return &Store{
		db: make(map[string]map[string]GenericItem),
	}
}
