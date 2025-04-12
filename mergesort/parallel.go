package mergesort

import (
	"sync"
)

func PMergeSort(A []int) []int {
	if len(A) <= 1 {
		return A
	}

	mid := len(A) / 2
	var left, right []int

	var wg sync.WaitGroup
	wg.Add(2)

	// Sort left half in a goroutine
	go func() {
		left = MergeSort(A[:mid])
		wg.Done()
	}()

	// Sort right half in a goroutine
	go func() {
		right = MergeSort(A[mid:])
		wg.Done()
	}()

	wg.Wait()
	return Merge(left, right)
}
