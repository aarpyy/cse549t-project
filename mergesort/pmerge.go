package mergesort

import (
	"sort"
	"sync"
)

// Entry point
func PmergeSort(arr []int) []int {
	if len(arr) <= 1 {
		cpy := make([]int, len(arr))
		copy(cpy, arr)
		sort.Ints(cpy)
		return cpy
	}

	output := make([]int, len(arr))
	var wg sync.WaitGroup
	wg.Add(1)
	go pmergeSort(arr, output, &wg)
	wg.Wait()
	return output
}

// Internal recursive parallel mergesort
func pmergeSort(arr []int, output []int, wg *sync.WaitGroup) {
	defer wg.Done()

	n := len(arr)
	if n <= 1 {
		copy(output, arr)
		return
	}

	mid := n / 2
	left := make([]int, mid)
	right := make([]int, n-mid)

	var wgInner sync.WaitGroup
	wgInner.Add(2)
	go pmergeSort(arr[:mid], left, &wgInner)
	go pmergeSort(arr[mid:], right, &wgInner)
	wgInner.Wait()

	wgInner.Add(1)
	go Pmerge(left, right, output, 0, &wgInner)
	wgInner.Wait()
}

func Pmerge(A, B, C []int, startC int, wg *sync.WaitGroup) {
	defer wg.Done()

	if len(A) == 0 {
		copy(C[startC:], B)
		return
	}
	if len(B) == 0 {
		copy(C[startC:], A)
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

	var wgInner sync.WaitGroup
	wgInner.Add(2)
	go Pmerge(A[:midA], B[:posB], C, startC, &wgInner)
	go Pmerge(A[midA+1:], B[posB:], C, startC+posC+1, &wgInner)
	wgInner.Wait()
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
