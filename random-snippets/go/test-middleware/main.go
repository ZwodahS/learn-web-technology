package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type HandleFunWrap struct {
	f func(http.ResponseWriter, *http.Request)
}

func (wrapper HandleFunWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wrapper.f(w, r)
}

func h(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return HandleFunWrap{f: f}
}

func handleTest(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("Test")
}

func main() {
	port := 9999

	var err error

	router := mux.NewRouter()
	router.Handle("/test", handlers.ContentTypeHandler(h(handleTest),
		"application/json", "application/xml")).Methods("PUT")

	http.Handle("/", router)

	fmt.Printf("Server Starting\n")
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	fmt.Printf("Server Stopped\n")

	if err != nil {
		panic(err)
	}
}
