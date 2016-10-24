package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func handleTest(writer http.ResponseWriter, req *http.Request) {
	time.Sleep(3 * time.Second)
	fmt.Fprintf(writer, "Value: %v\n", mux.Vars(req)["value"])
}

func runWeb(port int) {
	flag.Parse()

	var err error

	router := mux.NewRouter()
	router.HandleFunc("/test/{value}", handleTest).Methods("GET")

	http.Handle("/", router)

	fmt.Printf("Server Starting at %v\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	fmt.Printf("Server Stopped\n")

	if err != nil {
		panic(err)
	}
}

type Response struct {
	Url      string
	Response *http.Response
	Err      error
	Body     string
}

func AsyncGet(host string) chan Response {
	c := make(chan Response)
	go func() {
		response, err := http.Get(host)
		body, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		c <- Response{Url: host, Response: response, Err: err, Body: string(body)}
	}()
	return c
}

func runAsync(host string) {
	fmt.Println("Sending async request")
	c1 := AsyncGet(host)
	c2 := AsyncGet(host)
	fmt.Println("Sent async request")

	fmt.Println("Waiting for response 2")
	response2 := <-c2
	fmt.Println(response2.Body)
	fmt.Println("Waiting for response 1")
	response1 := <-c1
	fmt.Println(response1.Body)
}

func main() {
	run := flag.String("app", "web", "what to run")
	port := flag.Int("port", 9999, "port")
	value := flag.String("value", "temp", "")
	flag.Parse()

	run_host := fmt.Sprintf("http://localhost:%v/test/%v", *port, *value)

	if *run == "web" {
		runWeb(*port)
	} else {
		runAsync(run_host)
	}
}
