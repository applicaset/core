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
