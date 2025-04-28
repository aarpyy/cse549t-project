package mergesort

import "sort"

func MergeSort(A []int) []int {
	if len(A) <= basecase {
		cpy := make([]int, len(A))
		copy(cpy, A)
		sort.Ints(cpy)
		return cpy
	}

	mid := len(A) / 2
	left := MergeSort(A[:mid])
	right := MergeSort(A[mid:])

	return Merge(left, right)
}
