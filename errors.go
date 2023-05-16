package main

import "fmt"

type ItemExistsError struct {
	ID string
}

func (err ItemExistsError) Error() string {
	return fmt.Sprintf("item with id '%s' already exists", err.ID)
}

type ItemNotFoundError struct {
	ID string
}

func (err ItemNotFoundError) Error() string {
	return fmt.Sprintf("item with id '%s' not found", err.ID)
}

type GroupKindNotFoundError struct {
	Group string
	Kind  string
}

func (err GroupKindNotFoundError) Error() string {
	return fmt.Sprintf("kind with group '%s' and name '%s' not found", err.Group, err.Kind)
}
