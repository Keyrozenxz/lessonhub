package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type StudentResponse struct {
	Error string `json:"error"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	hr := HealthResponse{Status: "OK"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(hr)
}

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	sr := StudentResponse{Error: "student not found"}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad req: ", id)
		return
	}
	for i := 0; i < 10; i++ {
		id++
		id--
		//эмуляция работы
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(sr)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /student/{id}", StudentHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
