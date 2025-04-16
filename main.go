package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
	"cse549t-project/mergesort"
	"log"
	"runtime"
	"time"
)

const (
	n    = 10
	dim1 = 512
	dim2 = 512
	dim3 = 512
)

func measurePerformance(f func(a, b [][]int) [][]int) time.Duration {
	total := time.Duration(0)
	for i := 0; i < n; i++ {
		// Generate random matrices
		a := matrix.MagicMatrix(dim1, dim2, 42)
		b := matrix.MagicMatrix(dim2, dim3, 21)

		start := time.Now()
		f(a, b)
		elapsed := time.Since(start)

		total += elapsed
	}

	// Return average time
	return total / time.Duration(n)
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	log.Println("Number of Processors:", runtime.GOMAXPROCS(runtime.NumCPU()))
	a := matrix.MagicMatrix(dim1, dim2, 42)
	b := matrix.MagicMatrix(dim2, dim3, 21)
	log.Println("Matrix A Size: ", len(a), "x", len(a[0]))
	log.Println("Matrix B Size: ", len(b), "x", len(b[0]))

	// Nested For Loop MM
	c := matmul.MatrixMultiply(a, b)
	baseline := measurePerformance(matmul.MatrixMultiply)
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication check failed")
	}
	log.Printf("Nested For Loop MM running time: %dμs", baseline.Microseconds())

	// Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRec(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer check failed")
	}
	runTime := measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRec(a, b, len(a))
	})
	log.Println("Divide and Conquer MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))

	// 8x Parallel Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRecParallel8(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 8x Parallel check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRecParallel8(a, b, len(a))
	})
	log.Println("Divide and Conquer 8x Parallel MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))

	// 4x Parallel Divide and Conquer MM
	c = matmul.MatrixMultiplyDandCRecParallel4(a, b, len(a))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 4x Parallel check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.MatrixMultiplyDandCRecParallel4(a, b, len(a))
	})
	log.Println("Divide and Conquer 4x Parallel MM speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))

	// Cache Aware MM
	c = matmul.CacheAwareMM(a, b, matmul.ComputeBlockSize("1MB"))
	if !matmul.CheckMatMul(a, b, c) {
		panic("Cache aware matrix multiplication check failed")
	}
	runTime = measurePerformance(func(a, b [][]int) [][]int {
		return matmul.CacheAwareMM(a, b, matmul.ComputeBlockSize("1MB"))
	})
	log.Println("Cache aware matrix multiplication speedup: ", float64(baseline.Nanoseconds())/float64(runTime.Nanoseconds()))

	// All checks passed
	log.Println("All Matrix multiplication checks passed.")

	// -----------------------------------------------------
	// Merge-Sort Algorithms
	// -----------------------------------------------------

	d := matrix.MagicMatrix1D(512, 42)

	startTime := time.Now()
	e := mergesort.MergeSort(d)
	runTime = time.Since(startTime)

	if !mergesort.CheckSorted(d, e) {
		panic("Seq. Merge Sort check failed")
	}

	log.Println("Seq. Merge Sort running time: ", runTime.Nanoseconds())
	log.Println("Seq. Merge Sort check passed.")

	startTime = time.Now()
	e = mergesort.PMergeSort(d)
	runTime = time.Since(startTime)

	if !mergesort.CheckSorted(d, e) {
		panic("P Merge Sort check failed")
	}

	log.Println("P Merge Sort running time: ", runTime.Nanoseconds())
	log.Println("P Merge Sort check passed.")
}
