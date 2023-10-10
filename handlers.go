package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Received: %v\n", payload.Vector)
	fmt.Println("Received: ", payload.Vector)
}

func RateLimiter(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	sema := make(chan struct{}, maxClients)

	return func(w http.ResponseWriter, req *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		f(w, req)
	}
}