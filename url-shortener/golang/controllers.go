package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Controller struct {
	Services *Services
}

func (self *Controller) HandleApiV1CreateUrl(writer http.ResponseWriter, req *http.Request) {
	str, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Invalid Json Body")
		return
	}

	var parsedBody struct {
		Id  *string `json:"id"`
		Url *string `json:"url"`
	}
	json.Unmarshal([]byte(str), &parsedBody)
	if parsedBody.Id == nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Id is missing")
		return
	}
	if parsedBody.Url == nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Url is missing")
		return
	}
	url, err := self.Services.DB.UpdateUrl(Url{Id: *parsedBody.Id, Url: *parsedBody.Url})
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Unexpected Server Error")
		return
	}
	outputString, err := json.Marshal(*url)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "Unexpected Server Error")
	}
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, string(outputString))
}

func (self *Controller) HandleURLRedirect(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	urlId := vars["urlId"]
	url, err := self.Services.DB.GetURLById(urlId)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "error")
		return
	}

	if url == nil {
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "")
		return
	}
	fmt.Print(url)
	accept := req.Header["Accept"]
	accept_json := false
	for _, a := range accept {
		if a == "json" {
			accept_json = true
		}
	}
	if accept_json {
		writer.WriteHeader(http.StatusOK)
		result := struct {
			Status int    `json:"status_code"`
			Id     string `json:"url_id"`
			Url    string `json:"url_path"`
		}{http.StatusMovedPermanently, url.Id, url.Url}
		json_data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(writer, string(json_data))
	} else {
		http.Redirect(writer, req, url.Url, http.StatusMovedPermanently)
	}
}
