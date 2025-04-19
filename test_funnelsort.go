package main

import (
	"cse549t-project/funnelsort"
	"log"
	"math/rand"
	"sort"
	"time"
)

func randomArray(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Intn(n)
	}
	return arr
}

func testCorrectness() {
	arr := randomArray(10000)
	sorted := funnelsort.Sort(arr)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] < sorted[i-1] {
			panic("Array is not sorted")
		}
	}
	log.Println("Array is sorted successfully")
}

func testSpeedup() {
	// Compare to sort.Ints from average of 100 runs
	var base time.Duration
	for i := 0; i < 1; i++ {
		arr := randomArray(500000)
		start := time.Now()
		sort.Ints(arr)
		base += time.Since(start)
	}

	var fs time.Duration
	for i := 0; i < 1; i++ {
		arr := randomArray(500000)
		start := time.Now()
		funnelsort.Sort(arr)
		fs += time.Since(start)
	}

	log.Printf("Average speedup: %.4f", float64(base)/float64(fs))
}

func main() {
	testCorrectness()
	testSpeedup()
	log.Println("All tests passed successfully")
}
