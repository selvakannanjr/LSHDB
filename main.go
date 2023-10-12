package main

import (
	"fmt"
	"net/http"
)

var (
	seeds []int = []int{69696,420,91101}
	randmatrix [3][][]float64
	L          [3]LSHMap
)

func init() {
	for i := 0; i < 3; i++ {
		randmatrix[i] = GenerateRandomMatrix(8, 2048,int64(seeds[i]) )
		L[i] = make(LSHMap)
	}

	//load the map from the file
	for i := 0; i < 3; i++ {
		L[i].LoadMap("db2.csv", randmatrix[i])
	}
}

func main() {
	// http handler for /query
	//take the query vector from the request body
	// returns the top 3 results
	http.HandleFunc("/query", RateLimiter(HandleRequest,50))
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server started at port 8080")
}