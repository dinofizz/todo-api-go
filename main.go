package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	db := &gormdb{}
	db.init()
	defer db.close()
	router := mux.NewRouter()
	app := &Application{db: db, router: router}
	app.initRoutes()

	address := os.Getenv("HOST_ADDRESS")
	log.Printf("Starting web server on %s\n", address)
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
