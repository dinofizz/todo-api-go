package main

import (
	"github.com/jinzhu/gorm"
)

type GormItem struct {
	gorm.Model
	Description string
	Completed bool
}

type Item struct {
	Id uint
	Description string
	Completed   bool
}

type Database interface {
	init()
	createItem(item Item) Item
	deleteItem(id uint) error
	updateItem(id uint, td Item) (Item, error)
	getItem(id uint) (Item, error)
	allItems() []Item
	close()
}

