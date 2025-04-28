package mergesort

import (
	"sync"
)

func ParMergeSort(A []int) []int {

	if len(A) <= basecase {
		return MergeSort(A)
	}

	mid := len(A) / 2
	var left, right []int

	var wg sync.WaitGroup
	wg.Add(2)

	// Sort left half in a goroutine
	go func() {
		left = ParMergeSort(A[:mid])
		wg.Done()
	}()

	// Sort right half in a goroutine
	go func() {
		right = ParMergeSort(A[mid:])
		wg.Done()
	}()

	wg.Wait()
	return Merge(left, right)
}
