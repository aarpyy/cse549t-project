package mergesort

func MergeSort(A []int) []int {
	if len(A) <= 1 {
		return A
	}

	mid := len(A) / 2
	left := MergeSort(A[:mid])
	right := MergeSort(A[mid:])

	return Merge(left, right)
}
