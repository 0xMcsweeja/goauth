package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//add more properties for a client eg a key or password
// implement a middleware/func that filters the properties based on something eg auth token
type ClientStruct struct {
	Fullname      string      `json:"Fullname,omitempty"`
	Fname         string      `json:"First_Name,omitempty"`
	Lname         string      `json:"Last_Name,omitempty"`
	ID            string      `json:"ID,omitempty"`
	Credentials   Credentials `json:"Credentials, omitempty"`
	Authenticated bool        `json:"Authenticated,omitempty"`
}

type Credentials struct {
	Credential_type string `json:"credential_type,omitempty"`
	Secret          string `json:"secret,omitempty"`
}

func VisibleClientDetails(w http.ResponseWriter, r *http.Request, client ClientStruct) {
	// if not not authenticated only return partial struct
	//if client.Authenticated {
	//	return &client
	//}
	if client.Authenticated {
		jsonBytes, err := json.Marshal(client)
		if err != nil {
			log.Println("Couldn't parse client information")
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	} else {
		redacted := ClientStruct{
			Fullname:      client.Fullname,
			ID:            client.ID,
			Authenticated: false,
		}

		jsonBytes, err := json.Marshal(redacted)
		if err != nil {
			log.Println("Couldn't parse client information")
		}
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}

}
