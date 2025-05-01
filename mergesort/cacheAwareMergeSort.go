package mergesort

import (
	"math"
	"math/rand"
)

// cache aware
// k = z/l
func MultiWayMergeSort(A []int) []int {
	// base case
	n := len(A)
	if n <= basecaseCacheAware {
		return MergeSort(A)
	}

	partSize := int(math.Ceil(float64(n) / float64(k)))
	parts := [][]int{}

	for i := 0; i < n; i += partSize {
		end := min(i+partSize, n)
		sortedPart := MultiWayMergeSort(A[i:end])
		parts = append(parts, sortedPart)
	}

	return multiWayMerge(parts)
}

// merge function
func multiWayMerge(parts [][]int) []int {
	merged := []int{}
	indices := make([]int, len(parts))

	for {
		minVal := math.MaxInt64
		minIdx := -1
		for i, part := range parts {
			if indices[i] < len(part) && part[indices[i]] < minVal {
				minVal = part[indices[i]]
				minIdx = i
			}
		}

		if minIdx == -1 {
			break
		}

		merged = append(merged, minVal)
		indices[minIdx]++
	}

	return merged
}

// function to create random array
func randArray(n int) []int {
	A := make([]int, n)
	for i := range A {
		A[i] = rand.Intn(1000000)
	}
	return A
}
