package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"        // beautiful color :D
	"github.com/gorilla/mux"        // router
	_ "github.com/mattn/go-sqlite3" // sqlite3
	"github.com/namsral/flag"       // allow environment variable in flags
)

type Services struct {
	DB DatabaseWrapper
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	var dbFile = flag.String("database", "url.db", "database file")
	var port = flag.Int("port", 9999, "port to run the server on")
	flag.Parse()

	// set up the various services
	db, err := sql.Open("sqlite3", *dbFile)
	checkError(err)
	defer db.Close()

	wrapper := DatabaseWrapper{DB: db}
	services := Services{DB: wrapper}

	controller := Controller{Services: &services}

	router := mux.NewRouter()
	router.HandleFunc("/v1/url", controller.HandleApiV1CreateUrl).Methods("POST")
	router.HandleFunc("/r/{urlId}", controller.HandleURLRedirect).Methods("GET")

	http.Handle("/", router)

	red_color := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("Server Started at %v\n", red_color(*port))
	err = http.ListenAndServe(fmt.Sprintf(":%v", *port), nil)
	fmt.Printf("Server Stopped\n")

	if err != nil {
		log.Fatal("Error: ", err)
	}
}
