package main

import "math/rand"

//generate random float64 matrix with elements from standard normal distribution with given dimension
// create the same matrix for a given seed
func GenerateRandomMatrix(rows int, cols int, seed int64) [][]float64 {
	rand.NewSource(seed)
	var matrix [][]float64
	for i := 0; i < rows; i++ {
		var row []float64
		for j := 0; j < cols; j++ {
			row = append(row, rand.NormFloat64())
		}
		matrix = append(matrix, row)
	}
	return matrix
}

//compute matrix multiplication between an one-dimensional vector and a two-dimensional matrix
func MatrixMultiplication(vector []float64, matrix [][]float64) []float64 {
	var result []float64
	for i := 0; i < len(matrix); i++ {
		var sum float64
		for j := 0; j < len(vector); j++ {
			sum += vector[j] * matrix[i][j]
		}
		result = append(result, sum)
	}
	return result
}

func ComputeBucketID(vector []float64,matrix [][]float64)string{
	result := MatrixMultiplication(vector,matrix)
	//check if vector element is positive or negative
	//if positive, append 1 to bucketID
	//if negative, append 0 to bucketID
	var bucketID string
	for i := 0; i < len(result); i++ {
		if result[i] >= 0 {
			bucketID += "1"
		} else {
			bucketID += "0"
		}
	}
	return bucketID
}

func FindClosestBucket(querybucket string,allbuckets []string)string{
	//find the hamming distance between querybucket and allbuckets
	//return the bucket with the smallest hamming distance
	var minHammingDistance int
	var closestBucket string
	for i := 0; i < len(allbuckets); i++ {
		var hammingDistance int
		for j := 0; j < len(querybucket); j++ {
			if querybucket[j] != allbuckets[i][j] {
				hammingDistance++
			}
		}
		if i == 0 {
			minHammingDistance = hammingDistance
			closestBucket = allbuckets[i]
		} else if hammingDistance < minHammingDistance {
			minHammingDistance = hammingDistance
			closestBucket = allbuckets[i]
		}
	}
	return closestBucket
}