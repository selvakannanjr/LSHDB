package main

import (
	"fmt"
	"runtime"
	"sync"
)
func main() {
	binary := []string{"0000", "0001", "0010", "0011", "0100", "0101","0110", "0111", "1000", "1001", "1010", "1011", "1100", "1101","1110", "1111"}
	maxprocs := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup;
	if len(binary) < maxprocs {
		maxprocs = len(binary)
	}
	// iterate over the binary slice slicing it into maxprocs chunks
	for i := 0; i < len(binary); i += maxprocs {
		end := i + maxprocs
		if end > len(binary) {
			end = len(binary)
		}
		// launch a goroutine for each chunk
		wg.Add(1)
		fmt.Println("Launching goroutine")
		go func(binary []string) {
			defer wg.Done()
			for _, b := range binary {
				fmt.Println(b)
			}
		}(binary[i:end])
	}
	wg.Wait()
	
}
