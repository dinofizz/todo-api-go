package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Application struct {
	router *mux.Router
	db     Database
}

func (a *Application) initRoutes() {
	a.router.HandleFunc("/todo", a.createTodoItem).Methods("POST")
	a.router.HandleFunc("/todos", a.getAllToDoItems).Methods("GET")
	a.router.HandleFunc("/todo/{id}", a.getToDoItem).Methods("GET")
	a.router.HandleFunc("/todo/{id}", a.updateToDoItem).Methods("PUT")
	a.router.HandleFunc("/todo/{id}", a.deleteToDoItem).Methods("DELETE")
}

func (a *Application) getToDoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ToDo item ID")
		return
	}

	item, err := a.db.getItem(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

func (a *Application) getAllToDoItems(w http.ResponseWriter, r *http.Request) {
	items := a.db.allItems()
	respondWithJSON(w, http.StatusOK, items)
}

func (a *Application) deleteToDoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ToDo item ID")
		return
	}

	err = a.db.deleteItem(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Application) updateToDoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ToDo item ID")
		return
	}
	var td Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&td); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	updatedItem, err := a.db.updateItem(uint(id), td)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	}
	respondWithJSON(w, http.StatusOK, updatedItem)
}

func (a *Application) createTodoItem(w http.ResponseWriter, r *http.Request) {
	var td Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&td); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	td = a.db.createItem(td)
	respondWithJSON(w, http.StatusCreated, td)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
