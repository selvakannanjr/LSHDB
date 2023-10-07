package main

import (
	"bufio"
	csv "encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Payload struct {
	Vector []float64 `json:"vector"`
}

type ImageRec struct {
	ImageID string
	Vector  []float64
}

// a struct method to compute and return the cosine similarity with a given vector
func (ir *ImageRec) CosineSimilarity(vector []float64) float64 {
	return cosineSimilarity(ir.Vector, vector)
}

type LSHMap map[string][]ImageRec

// a struct method to load the map from a csv file
func (L *LSHMap) LoadMap(filename string, matrix [][]float64) {
	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		var vector []float64
		for i := 1; i < len(line); i++ {
			element, _ := strconv.ParseFloat(line[i], 64)
			vector = append(vector, element)
		}
		bid := ComputeBucketID(vector,matrix)
		// if bid is not in L, add it to L
		// if bid is in L, append the vector to the bucket
		_, ok := (*L)[bid]
		if !ok {
			(*L)[bid] = []ImageRec{{ImageID: line[0], Vector: vector}}
		} else {
			(*L)[bid] = append((*L)[bid], ImageRec{ImageID: line[0], Vector: vector})
		}
	}
}



// compute and return the dot product of two float64 vectors
func dotProduct(vector1 []float64, vector2 []float64) float64 {
	var sum float64 = 0
	for i := 0; i < len(vector1); i++ {
		sum += vector1[i] * vector2[i]
	}
	return sum
}

// compute and return cosine similarity of two float64 vectors
func cosineSimilarity(vector1 []float64, vector2 []float64) float64 {
	return dotProduct(vector1, vector2) / (math.Sqrt(dotProduct(vector1, vector1)) * math.Sqrt(dotProduct(vector2, vector2)))
}
