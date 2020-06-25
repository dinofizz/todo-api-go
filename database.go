package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type GormItem struct {
	gorm.Model
	Description string
	Completed   bool
}

type Item struct {
	Id          string
	Description string
	Completed   bool
}

type Database interface {
	init()
	ping() error
	createItem(item Item) (Item, error)
	deleteItem(id string) error
	updateItem(id string, td Item) (Item, error)
	getItem(id string) (Item, error)
	allItems() ([]Item, error)
	close()
}

type ErrorItemNotFound struct {
	Id string
}

func (e *ErrorItemNotFound) Error() string {
	return fmt.Sprintf("Unable to find item with id %s", e.Id)
}

