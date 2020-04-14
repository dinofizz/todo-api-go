package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
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
	item, err := a.db.getItem(vars["id"])
	var e *ErrorItemNotFound
	if errors.As(err, &e) {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error occurred")
		return
	}
	respondWithJSON(w, http.StatusOK, item)
}

func (a *Application) getAllToDoItems(w http.ResponseWriter, r *http.Request) {
	items, err := a.db.allItems()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error occurred")
		return
	}
	respondWithJSON(w, http.StatusOK, items)
}

func (a *Application) deleteToDoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := a.db.deleteItem(vars["id"])

	var e *ErrorItemNotFound
	if errors.As(err, &e) {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error occurred")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Application) updateToDoItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var td Item
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&td); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	updatedItem, err := a.db.updateItem(vars["id"], td)
	var e *ErrorItemNotFound
	if errors.As(err, &e) {
		respondWithError(w, http.StatusNotFound, "item not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error occurred")
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

	td, err := a.db.createItem(td)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error occurred")
		return
	}
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
