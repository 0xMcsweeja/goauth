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

	http.HandleFunc("/", handler.get)
	http.HandleFunc("/dogs", handler.getDog)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (h *Handlers) get(w http.ResponseWriter, r *http.Request) {
	healthcheck := h.store
	jsonBytes, err := json.Marshal(healthcheck)
	if err != nil {
		panic(err)
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
