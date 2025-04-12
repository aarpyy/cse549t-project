package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
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

	startTime2 := time.Now()
	c = matmul.MatrixMultiplyDandCRec(a, b, len(a))
	runTime = time.Since(startTime2)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer check failed")
	}

	log.Println("Divide and Conquer MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")

	startTime3 := time.Now()
	c = matmul.MatrixMultiplyDandCRecParallel(a, b, len(a))
	runTime = time.Since(startTime3)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication Divide and Conquer Parallel check failed")
	}

	log.Println("Divide and Conquer Parallel MM running time: ", runTime.Nanoseconds())
	log.Println("Matrix multiplication check passed.")
}
