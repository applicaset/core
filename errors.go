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

type KindNotFoundError struct {
	Name string
}

func (err KindNotFoundError) Error() string {
	return fmt.Sprintf("kind with name '%s' not found", err.Name)
}
