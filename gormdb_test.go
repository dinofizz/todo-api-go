package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initOsEnv() {
	os.Setenv("GORM_DIALECT", "sqlite3")
	os.Setenv("GORM_CONNECTION_STRING", ":memory:")
}

func initDB() *gormdb {
	initOsEnv()
	db := &gormdb{}
	db.init()
	return db
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	return func(t *testing.T) {
		t.Log("teardown test case")
	}
}

func Test_init(t *testing.T) {
	db := initDB()
	defer db.close()
}

func Test_init_error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	os.Setenv("GORM_DIALECT", "some_other_db")
	os.Setenv("GORM_CONNECTION_STRING", ":unknown:")
	db := &gormdb{}
	db.init()
	defer db.close()
}

func Test_createItem(t *testing.T) {
	db := initDB()
	defer db.close()

	item := db.createItem(Item{Description: "test_description", Completed: false})

	assert.Equal(t, uint(1), item.Id)
	assert.Equal(t, "test_description", item.Description)
	assert.Equal(t, false, item.Completed)
}

func Test_updateItem(t *testing.T) {
	db := initDB()
	defer db.close()

	createdItem := db.createItem(Item{Description: "test_description", Completed: false})
	update := Item{Description: "updated description", Completed: true}
	item, err := db.updateItem(createdItem.Id, update)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), item.Id)
	assert.Equal(t, "updated description", item.Description)
	assert.Equal(t, true, item.Completed)
}

func Test_updateItem_not_exists(t *testing.T) {
	db := initDB()
	defer db.close()

	update := Item{Description: "updated description", Completed: true}
	item, err := db.updateItem(1234, update)

	assert.EqualError(t, err, "item id does not exist")
	assert.Equal(t, uint(0), item.Id)
	assert.Equal(t, "", item.Description)
	assert.Equal(t, false, item.Completed)
}

func Test_deleteItem(t *testing.T) {
	db := initDB()
	defer db.close()

	createdItem := db.createItem(Item{Description: "test_description", Completed: false})
	err := db.deleteItem(createdItem.Id)
	assert.NoError(t, err)
	item, err := db.getItem(createdItem.Id)
	assert.EqualError(t, err, "item id does not exist")
	assert.Equal(t, uint(0), item.Id)
	assert.Equal(t, "", item.Description)
	assert.Equal(t, false, item.Completed)
}

func Test_deleteItem_not_exists(t *testing.T) {
	db := initDB()
	defer db.close()

	err := db.deleteItem(1327)
	assert.EqualError(t, err, "item id does not exist")
}

func Test_getItem(t *testing.T) {
	db := initDB()
	defer db.close()

	createdItem := db.createItem(Item{Description: "test_description", Completed: false})
	item, err := db.getItem(createdItem.Id)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), item.Id)
	assert.Equal(t, "test_description", item.Description)
	assert.Equal(t, false, item.Completed)
}

func Test_getItem_not_exists(t *testing.T) {
	db := initDB()
	defer db.close()

	item, err := db.getItem(1327)
	assert.EqualError(t, err, "item id does not exist")
	assert.Equal(t, uint(0), item.Id)
	assert.Equal(t, "", item.Description)
	assert.Equal(t, false, item.Completed)
}

func Test_allItems(t *testing.T) {
	db := initDB()
	defer db.close()

	db.createItem(Item{Description: "A", Completed: false})
	db.createItem(Item{Description: "B", Completed: true})
	db.createItem(Item{Description: "C", Completed: false})
	items  := db.allItems()
	assert.Equal(t, 3, len(items))
	assert.Equal(t, uint(1), items[0].Id)
	assert.Equal(t, "A", items[0].Description)
	assert.Equal(t, false, items[0].Completed)
	assert.Equal(t, uint(2), items[1].Id)
	assert.Equal(t, "B", items[1].Description)
	assert.Equal(t, true, items[1].Completed)
	assert.Equal(t, uint(3), items[2].Id)
	assert.Equal(t, "C", items[2].Description)
	assert.Equal(t, false, items[2].Completed)
}


