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
	Id          uint
	Description string
	Completed   bool
}

type Database interface {
	init()
	createItem(item Item) (Item, error)
	deleteItem(id uint) error
	updateItem(id uint, td Item) (Item, error)
	getItem(id uint) (Item, error)
	allItems() ([]Item, error)
	close()
}

type ErrorItemNotFound struct {
	Id uint
}

func (e *ErrorItemNotFound) Error() string {
	return fmt.Sprintf("Unable to find item with id %d", e.Id)
}

