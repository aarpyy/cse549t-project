package main

import (
	"cse549t-project/matmul"
	"cse549t-project/matrix"
	"log"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	a := matrix.MagicMatrix(8, 9, 42)
	b := matrix.MagicMatrix(9, 8, 42)

	c := matmul.MatrixMultiply(a, b)

	if !matmul.CheckMatMul(a, b, c) {
		panic("Matrix multiplication check failed")
	}

	log.Println("Matrix multiplication check passed")
}
