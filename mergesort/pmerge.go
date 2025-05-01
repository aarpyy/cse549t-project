package mergesort

import (
	"sync"
)

func PmergeSort(A []int) []int {

	if len(A) <= basecase {
		return MergeSort(A)
	}

	mid := len(A) / 2
	var left, right []int

	var wg sync.WaitGroup
	wg.Add(2)

	// Sort left half
	go func() {
		left = PmergeSort(A[:mid])
		wg.Done()
	}()

	// Sort right half
	go func() {
		right = PmergeSort(A[mid:])
		wg.Done()
	}()

	wg.Wait()

	C := make([]int, len(A))
	Pmerge(left, right, C, 0, 0)
	return C
}

func Pmerge(A, B, C []int, startC int, depth int) {
	if len(A) == 0 {
		copy(C[startC:], B)
		return
	}
	if len(B) == 0 {
		copy(C[startC:], A)
		return
	}
	if len(A) == 1 && len(B) == 1 {
		if A[0] <= B[0] {
			C[startC] = A[0]
			C[startC+1] = B[0]
		} else {
			C[startC] = B[0]
			C[startC+1] = A[0]
		}
		return
	}

	if len(A)+len(B) <= basecaseMerge {
		MergeC(A, B, C[startC:])
		return
	}

	if len(B) > len(A) {
		A, B = B, A
	}

	midA := len(A) / 2
	val := A[midA]
	posB := binarySearch(B, val)
	posC := midA + posB
	C[startC+posC] = val

	if depth <= maxDepth {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			Pmerge(A[:midA], B[:posB], C, startC, depth+1)
			wg.Done()
		}()
		Pmerge(A[midA+1:], B[posB:], C, startC+posC+1, depth+1)
		wg.Wait()
	} else {
		Pmerge(A[:midA], B[:posB], C, startC, depth+1)
		Pmerge(A[midA+1:], B[posB:], C, startC+posC+1, depth+1)
	}
}

func binarySearch(arr []int, val int) int {
	low, high := 0, len(arr)
	for low < high {
		mid := (low + high) / 2
		if arr[mid] < val {
			low = mid + 1
		} else {
			high = mid
		}
	}
	return low
}
