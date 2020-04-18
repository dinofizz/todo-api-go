package main

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"strconv"
)

type gormdb struct {
	db *gorm.DB
}

func (s *gormdb) init() {
	dialect := os.Getenv("GORM_DIALECT")
	if dialect == "" {
		panic(fmt.Sprint("Missing value for GORM_DIALECT environment variable."))
	}
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		panic(fmt.Sprint("Missing value for CONNECTION_STRING environment variable."))
	}

	gormdb, err := gorm.Open(dialect, connectionString)
	if err != nil {
		fmt.Println(err)
		panic(fmt.Sprintf("failed to connect to %s Database with connection string %s", dialect, connectionString))
	}
	s.db = gormdb
	s.db.AutoMigrate(&GormItem{})
}

func (s *gormdb) createItem(item Item) (Item, error) {
	gtd := &GormItem{Description: item.Description, Completed: item.Completed}
	if err := s.db.Create(gtd).Error; err != nil {
		return Item{}, err
	}
	id := strconv.FormatUint(uint64(gtd.ID), 10)
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: id}, nil
}

func (s *gormdb) updateItem(id string, td Item) (Item, error) {
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return Item{}, errors.New("Invalid ID type.")
	}
	var gtd GormItem
	if err := s.db.First(&gtd, uintId).Error; gorm.IsRecordNotFoundError(err) {
		return Item{}, &ErrorItemNotFound{Id: id}
	} else if err != nil {
		return Item{}, err
	}
	s.db.Model(&gtd).Update("Completed", td.Completed).Update("Description", td.Description)
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: id}, nil
}

func (s *gormdb) deleteItem(id string) error {
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return errors.New("Invalid ID type.")
	}
	var gtd GormItem
	if err := s.db.First(&gtd, uintId).Error; gorm.IsRecordNotFoundError(err) {
		return &ErrorItemNotFound{Id: id}
	} else if err != nil {
		return err
	}
	s.db.Delete(&gtd)
	return nil
}

func (s *gormdb) getItem(id string) (Item, error) {
	uintId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return Item{}, errors.New("Invalid ID type.")
	}
	var gtd GormItem
	if err := s.db.First(&gtd, uintId).Error; gorm.IsRecordNotFoundError(err) {
		return Item{}, &ErrorItemNotFound{Id: id}
	} else if err != nil {
		return Item{}, err
	}
	return Item{Description: gtd.Description, Completed: gtd.Completed, Id: id}, nil
}

func (s *gormdb) allItems() ([]Item, error) {
	var gtds []GormItem
	if err := s.db.Find(&gtds).Error; err != nil {
		return make([]Item, 0), err
	}

	tds := make([]Item, len(gtds))
	for i, v := range gtds {
		id := strconv.FormatUint(uint64(v.ID), 10)
		tds[i] = Item{Description: v.Description, Completed: v.Completed, Id: id}
	}
	return tds, nil
}

func (s *gormdb) close() {
	s.db.Close()
}
