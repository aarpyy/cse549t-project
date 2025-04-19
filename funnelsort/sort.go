package funnelsort

import (
	"math"
	"sort"
	"sync"
)

// Parallel sort for initial subarrays
func parallelSort(arr []int, splits int) [][]int {
	n := len(arr)
	arrs := make([][]int, splits)
	step := n / splits

	var wg sync.WaitGroup
	wg.Add(splits)

	for i := 0; i < splits; i++ {
		start := i * step
		end := start + step
		if i == splits-1 {
			end = n
		}

		go func(i, start, end int) {
			defer wg.Done()
			arrs[i] = make([]int, end-start)
			copy(arrs[i], arr[start:end])
			sort.Ints(arrs[i])
		}(i, start, end)
	}

	wg.Wait()
	return arrs
}

func Sort(arr []int) []int {
	n := len(arr)
	if n <= 1 {
		return arr
	}

	// For small arrays, use simple sort
	if n <= 1024 {
		result := make([]int, n)
		copy(result, arr)
		sort.Ints(result)
		return result
	}

	// Use cube root for splitting as in original implementation
	cbrt := int(math.Cbrt(float64(n)))
	splits := cbrt

	// Limit splits to avoid excessive overhead
	if splits > 32 {
		splits = 32
	} else if splits < 2 {
		splits = 2
	}

	// Parallel sort the subarrays
	sortedArrs := parallelSort(arr, splits)

	// Create KMerger with the sorted arrays
	merger := NewKMerger(sortedArrs)

	// Create output array with pre-allocated capacity
	out := make([]int, 0, n)
	for {
		v, ok := merger.Next()
		if !ok {
			break
		}
		out = append(out, v)
	}

	return out
}
