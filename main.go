package main

import (
	"fmt"
	"net/http"
)

var (
	randmatrix [][]float64
	L          LSHMap
)

func init() {
	randmatrix = GenerateRandomMatrix(3, 4)
	L = make(LSHMap)
	L.LoadMap("sample.csv", randmatrix)
}

func main() {
	// http handler for /query
	//take the query vector from the request body
	// returns the top 3 results
	http.HandleFunc("/query", L.HandleRequest)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server started at port 8080")
}