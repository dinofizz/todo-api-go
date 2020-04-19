package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type errorMessage struct {
	Error string `json:"error"`
}

func Test_initRoutes(t *testing.T) {
	router := mux.NewRouter()
	app := &Application{db: nil, router: router}
	app.initRoutes()
}

func TestApplication_getToDoItem(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("getItem", "1").Return(item, nil)

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("GET", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	responseItem := &Item{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "1", item.Id)
	assert.Equal(t, "ABC", item.Description)
	assert.Equal(t, true, item.Completed)
}

func TestApplication_getToDoItem_not_found(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	db.On("getItem", "1").Return(*new(Item), &ErrorItemNotFound{Id: "1"})

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("GET", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "item not found", responseItem.Error)
}

func TestApplication_getToDoItem_db_error(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	db.On("getItem", "1").Return(*new(Item), errors.New("db error"))

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("GET", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Database error occurred", responseItem.Error)
}

func TestApplication_getAllToDoItems(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)

	items := make([]Item, 3)
	items[0] = Item{Description: "A", Completed: true, Id: "1"}
	items[1] = Item{Description: "B", Completed: false, Id: "2"}
	items[2] = Item{Description: "C", Completed: true, Id: "3"}

	db.On("allItems").Return(items, nil)

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("GET", "/todos", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseItems []Item
	err = json.NewDecoder(rr.Body).Decode(&responseItems)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(responseItems))
	assert.Equal(t, "1", responseItems[0].Id)
	assert.Equal(t, "A", responseItems[0].Description)
	assert.True(t, responseItems[0].Completed)
	assert.Equal(t, "2", responseItems[1].Id)
	assert.Equal(t, "B", responseItems[1].Description)
	assert.False(t, responseItems[1].Completed)
	assert.Equal(t, "3", responseItems[2].Id)
	assert.Equal(t, "C", responseItems[2].Description)
	assert.True(t, responseItems[2].Completed)
}

func TestApplication_getAllToDoItems_db_error(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)

	items := make([]Item, 3)
	items[0] = Item{Description: "A", Completed: true, Id: "1"}
	items[1] = Item{Description: "B", Completed: false, Id: "2"}
	items[2] = Item{Description: "C", Completed: true, Id: "3"}

	db.On("allItems").Return(items, errors.New("db error"))

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("GET", "/todos", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Database error occurred", responseItem.Error)
}

func TestApplication_deleteToDoItem(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("deleteItem", "1").Return(nil)

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	responseItem := &Item{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "1", item.Id)
	assert.Equal(t, "ABC", item.Description)
	assert.Equal(t, true, item.Completed)
}

func TestApplication_deleteToDoItem_not_found(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	db.On("deleteItem", "1").Return(&ErrorItemNotFound{})

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "item not found", responseItem.Error)
}

func TestApplication_deleteToDoItem_db_error(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	db.On("deleteItem", "1").Return(errors.New("db error"))

	app := &Application{db: db, router: router}
	app.initRoutes()

	req, err := http.NewRequest("DELETE", "/todo/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Database error occurred", responseItem.Error)
}

func TestApplication_updateToDoItem(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("updateItem", "1", mock.AnythingOfType("Item")).Return(item, nil)

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(item)
	req, err := http.NewRequest("PUT", "/todo/1", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	responseItem := &Item{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "1", item.Id)
	assert.Equal(t, "ABC", item.Description)
	assert.Equal(t, true, item.Completed)
}

func TestApplication_updateToDoItem_not_found(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("updateItem", "1", mock.AnythingOfType("Item")).Return(Item{}, &ErrorItemNotFound{})

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(item)
	req, err := http.NewRequest("PUT", "/todo/1", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "item not found", responseItem.Error)
}

func TestApplication_updateToDoItem_db_error(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("updateItem", "1", mock.AnythingOfType("Item")).Return(Item{}, errors.New("db error"))

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(item)
	req, err := http.NewRequest("PUT", "/todo/1", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Database error occurred", responseItem.Error)
}

func TestApplication_updateToDoItem_invalid_json(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode("foobar")
	req, err := http.NewRequest("PUT", "/todo/1", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", responseItem.Error)
}

func TestApplication_createToDoItem(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("createItem", item).Return(item, nil)

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(item)
	req, err := http.NewRequest("POST", "/todo", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	responseItem := &Item{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "1", item.Id)
	assert.Equal(t, "ABC", item.Description)
	assert.Equal(t, true, item.Completed)
}

func TestApplication_createToDoItem_invalid_json(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode("foobar")
	req, err := http.NewRequest("POST", "/todo", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", responseItem.Error)
}

func TestApplication_createToDoItem_db_error(t *testing.T) {
	router := mux.NewRouter()

	db := new(MockDatabase)
	item := Item{Description: "ABC", Completed: true, Id: "1"}
	db.On("createItem", item).Return(Item{}, errors.New("db error"))

	app := &Application{db: db, router: router}
	app.initRoutes()

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(item)
	req, err := http.NewRequest("POST", "/todo", b)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	responseItem := &errorMessage{}
	err = json.NewDecoder(rr.Body).Decode(responseItem)
	assert.NoError(t, err)
	assert.Equal(t, "Database error occurred", responseItem.Error)
}

func TestApplication_health_live(t *testing.T) {
	items := make([]Item, 0)
	var healthTests = []struct {
		url        string
		items      []Item
		e          error
		statusCode int
	}{
		{"/live", items, nil, http.StatusNoContent},
		{"/live", items, errors.New("Some error"), http.StatusInternalServerError},
		{"/ready", items, nil, http.StatusNoContent},
		{"/ready", items, errors.New("Some error"), http.StatusInternalServerError},
	}

	for _, tt := range healthTests {
		t.Run(tt.url, func(t *testing.T) {
			router := mux.NewRouter()
			db := new(MockDatabase)
			app := &Application{db: db, router: router}
			app.initRoutes()
			db.On("allItems").Return(tt.items, tt.e)
			req, _ := http.NewRequest("GET", tt.url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusCode, rr.Code)
		})
	}
}
