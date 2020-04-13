package main

import (
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"time"
)

func main() {
	db := &gormdb{}
	db.init()
	defer db.close()
	router := mux.NewRouter()
	app := &Application{db: db, router: router}
	app.initRoutes()

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
