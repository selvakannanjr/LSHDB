package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func RateLimiter(f http.HandlerFunc, maxClients int) http.HandlerFunc {
	semaphore := make(chan struct{}, maxClients)

	return func(w http.ResponseWriter, req *http.Request) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()
		f(w, req)
	}
}

// a http handler that takes the query vector from the request body
//queries all the maps and returns the top 3 results
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

	//compute bucketid of the query vector for each map
	var bucketids [3]string;
	for i := 0; i < 3; i++ {
		bucketids[i] = ComputeBucketID(payload.Vector, randmatrix[i])
	}
	//query all the maps and merge the results
	topresults := MergeResults(3, L[0].Query(bucketids[0], payload.Vector,3),
		L[1].Query(bucketids[1], payload.Vector, 3),
		L[2].Query(bucketids[2], payload.Vector, 3))
	
	//return the top 3 results in an anonymous struct
	json.NewEncoder(w).Encode(struct {
		Results []string `json:"results"`
	}{topresults})
}