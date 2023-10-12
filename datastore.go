package main

import (
	"bufio"
	csv "encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"time"
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
	fmt.Println("Map loaded")
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

func (L *LSHMap) Query(bucketid string, vector []float64, top int)[]ResType{
	var result []ResType

	b2bsearched := L.GetClosestBucket(bucketid)

	ch := make(chan ResType)
	maxprocs := runtime.GOMAXPROCS(0)
	if len((*L)[b2bsearched]) < maxprocs {
		maxprocs = len((*L)[b2bsearched])
	}

	if top > len((*L)[b2bsearched]) {
		top = len((*L)[b2bsearched])
	}
	fmt.Println(len((*L)[b2bsearched]))

	for i:=0;i<len((*L)[b2bsearched]);i+=maxprocs{
		end := i+maxprocs
		if end > len((*L)[b2bsearched]){
			end = len((*L)[b2bsearched])
		}
		fmt.Println("goroutine spawned")
		go FindCosineSimilarity((*L)[b2bsearched][i:end],vector,ch)
	}

	//read from channel and append to result
		for {
			select {
			case res := <-ch:
				result = append(result, res)
			default:
				//wait for 1 second
				time.Sleep(25 * time.Millisecond)
				//close channel if no more data is coming
				if len(result) == len((*L)[b2bsearched]) {
					close(ch)
					goto end
				}
			}
		}
	end:

	// sort result by score in descending order
	// return top 10 imageids
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	// return top 10 imageids
	var topresults []ResType
	for i := 0; i < top; i++ {
		topresults = append(topresults, result[i])
	}
	return topresults
}

// a method that accepts any number of ResType slices and returns the top k results of all the slices
func MergeResults(k int, slices ...[]ResType) []string {
	var result []ResType
	for _, slice := range slices {
		result = append(result, slice...)
	}
	// sort result by score in descending order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	// return top k results
	var topresults []string
	for i := 0; i < k; i++ {
		topresults = append(topresults, result[i].ImageID)
	}
	return topresults
}
// a struct method to act as a http handler that takes the query vector from the request body and returns the top 3 results


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

func FindCosineSimilarity(recs []ImageRec, vector []float64, ch chan ResType) {
	for _, rec := range recs {
		score := rec.CosineSimilarity(vector)
		ch <- ResType{rec.ImageID, score}
	}
}