package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

type gormdb struct {
	db     *gorm.DB
}

func (s *gormdb) init() {
	dialect := os.Getenv("GORM_DIALECT")
	connectionString := os.Getenv("GORM_CONNECTION_STRING")
	gormdb, err := gorm.Open(dialect, connectionString)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to %s Database with connection string %s", dialect, connectionString))
	}
	s.db = gormdb
	s.db.AutoMigrate(&GormItem{})
}

func (s *gormdb) createItem(item Item) Item {
	gtd := &GormItem{Description: item.Description, Completed: item.Completed}
	s.db.Create(gtd)
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: gtd.ID}
}

func (s *gormdb) updateItem(id uint, td Item) (Item, error) {
	var gtd GormItem
	if s.db.First(&gtd, id).RecordNotFound() {
		return Item{}, errors.New("item id does not exist")
	}
	s.db.Model(&gtd).Update("Completed", td.Completed).Update("Description", td.Description)
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: gtd.ID}, nil
}

func (s *gormdb) deleteItem(id uint) error {
	var gtd GormItem
	if s.db.First(&gtd, id).RecordNotFound() {
		return errors.New("item id does not exist")
	}
	s.db.Delete(&gtd)
	return nil
}

func (s *gormdb) getItem(id uint) (Item, error) {
	var gtd GormItem
	if s.db.First(&gtd, id).RecordNotFound() {
		return Item{}, errors.New("item id does not exist")
	}
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: gtd.ID}, nil
}

func (s *gormdb) allItems() []Item {
	var gtds []GormItem
	s.db.Find(&gtds)

	tds := make([]Item, len(gtds))
	for i, v := range gtds {
		tds[i] = Item{Description: v.Description, Completed: v.Completed, Id: v.ID}
	}
	return tds
}

func (s *gormdb) close() {
	s.db.Close()
}
