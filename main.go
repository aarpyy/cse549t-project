package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
	"cse549t-project/mergesort"
	"log"
	"runtime"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	log.Println(runtime.GOMAXPROCS(runtime.NumCPU()))
	a := matrix.MagicMatrix(512, 512, 42)
	b := matrix.MagicMatrix(512, 512, 42)

	startTime := time.Now()
	c := matmul.MatrixMultiply(a, b)
	runTime := time.Since(startTime)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication check failed")
	}

	log.Println("Nested For Loop MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")

	startTime = time.Now()
	c = matmul.MatrixMultiplyDandCRec(a, b, len(a))
	runTime = time.Since(startTime)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer check failed")
	}

	log.Println("Divide and Conquer MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")

	startTime = time.Now()
	c = matmul.MatrixMultiplyDandCRecParallel8(a, b, len(a))
	runTime = time.Since(startTime)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 8x Parallel check failed")
	}

	log.Println("Divide and Conquer 8x Parallel MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")

	startTime = time.Now()
	c = matmul.MatrixMultiplyDandCRecParallel4(a, b, len(a))
	runTime = time.Since(startTime)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer 4x Parallel check failed")
	}

	log.Println("Divide and Conquer 4x Parallel MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")

	// -----------------------------------------------------
	// Merge-Sort Algorithms
	// -----------------------------------------------------

	d := matrix.MagicMatrix1D(512, 42)

	startTime = time.Now()
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
