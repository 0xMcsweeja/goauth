package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Status struct {
	Status string
	Time   time.Time
}

type Handlers struct {
	store map[string]Status
}

func main() {
	handler := newHandler()

	http.HandleFunc("/", handler.healthchecks)
	http.HandleFunc("/dogs", handler.getDog)
	http.HandleFunc("/token", Token)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (h *Handlers) healthchecks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.get(w, r) //need to update for post
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *Handlers) get(w http.ResponseWriter, r *http.Request) {
	healthcheck := h.store
	jsonBytes, err := json.Marshal(healthcheck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	fmt.Fprintln(w, "%s", healthcheck)
	w.Write(jsonBytes)
}

func (h *Handlers) getDog(w http.ResponseWriter, r *http.Request) {
	healthcheck := "dogs here"

	jsonBytes, err := json.Marshal(healthcheck)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func newHandler() *Handlers {
	return &Handlers{
		store: map[string]Status{
			"healthcheck": Status{
				Status: "healthy",
				Time:   time.Now(),
			},
		},
	}
}
