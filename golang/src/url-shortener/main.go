package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"  // beautiful color :D
	"github.com/gorilla/mux"  // router
	_ "github.com/lib/pq"     // pg
	"github.com/namsral/flag" // allow environment variable in flags
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

	var dbPort = flag.Int("dbport", 5432, "database port")
	var dbHost = flag.String("dbstring", "localhost", "database host")
	var dbUser = flag.String("dbuser", "user", "database user")
	var dbPassword = flag.String("dbpassword", "password", "database password")
	var dbName = flag.String("dbname", "url", "database name")

	var port = flag.Int("port", 9999, "port to run the server on")
	flag.Parse()

	// set up the various services
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%v dbname=%v password=%v port=%v host=%v sslmode=disable",
			*dbUser, *dbName, *dbPassword, *dbPort, *dbHost))
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
