package main

import (
	"fmt"
)

func main() {
	//generate a random matrix with 3 rows and 4 columns
	//make a LSHMap
	//and load the csv file into the LSHMap
	matrix := GenerateRandomMatrix(3, 4)
	L := make(LSHMap)
	L.LoadMap("sample.csv", matrix)
	vector := []float64{1, 2, 3, 4}
	bid := ComputeBucketID(vector, matrix)
	//Query the LSHMap with the query bucket id and vector
	results := L.Query(bid, vector,3)
	fmt.Println("query results:", results)

}