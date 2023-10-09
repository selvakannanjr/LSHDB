package main

import (
	"bufio"
	csv "encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
)

type Payload struct {
	Vector []float64 `json:"vector"`
}

type ImageRec struct {
	ImageID string
	Vector  []float64
}

type ResType struct {
	ImageID string  `json:"imageid"`
	Score   float64 `json:"score"`
}

// a struct method to compute and return the cosine similarity with a given vector
func (ir *ImageRec) CosineSimilarity(vector []float64) float64 {
	return cosineSimilarity(ir.Vector, vector)
}

type LSHMap map[string][]ImageRec

// a struct method to insert a new image record into the map
func (L *LSHMap) InsertImageRec(ir ImageRec, matrix [][]float64) {
	bid := ComputeBucketID(ir.Vector,matrix)
	// if bid is not in L, add it to L
	// if bid is in L, append the vector to the bucket
	_, ok := (*L)[bid]
	if !ok {
		(*L)[bid] = []ImageRec{ir}
	} else {
		(*L)[bid] = append((*L)[bid], ir)
	}
}

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
		ir := ImageRec{line[0], vector}
		L.InsertImageRec(ir, matrix)
	}
}

// a struct method to return all the keys in the map
func (L *LSHMap) GetKeys() []string {
	var keys []string
	for key := range *L {
		keys = append(keys, key)
	}
	return keys
}

// a struct method to return the closest bucket to a given bucketid
func (L *LSHMap) GetClosestBucket(bucketid string) string {
	allbuckets := L.GetKeys()
	return FindClosestBucket(bucketid, allbuckets)
}

func FindCosineSimilarity(recs []ImageRec, vector []float64, ch chan ResType) {
	for _, rec := range recs {
		score := rec.CosineSimilarity(vector)
		ch <- ResType{rec.ImageID, score}
	}
}

func (L *LSHMap) Query(bucketid string, vector []float64)[]string{
	var result []string

	b2bsearched := L.GetClosestBucket(bucketid)

	ch := make(chan ResType,4)
	maxprocs := runtime.GOMAXPROCS(0)
	if len((*L)[b2bsearched]) < maxprocs {
		maxprocs = len((*L)[b2bsearched])
	}
	for i:=0;i<len((*L)[b2bsearched]);i+=maxprocs{
		end := i+maxprocs
		if end > len((*L)[b2bsearched]){
			end = len((*L)[b2bsearched])
		}
		go FindCosineSimilarity((*L)[b2bsearched][i:end],vector,ch)
	}



	return result
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
