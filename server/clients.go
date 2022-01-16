package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ClientStruct struct {
	Fname string `json:"first_name"`
	Lname string `json:"last_name"`
	ID    string `json:"ID,omitempty"`
}

func Client(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var client ClientStruct
	err = json.Unmarshal(body, &client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	client.ID = "x[:]"

	jsonBytes, err := json.Marshal(client)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
