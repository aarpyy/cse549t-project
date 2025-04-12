package matmul

import "sync"

func MatrixMultiplyDandCRecParallel(a, b [][]int, n int) [][]int {
	if len(a[0]) != len(b) {
		panic("Matrix dimensions do not match for multiplication")
	}

	if n == 1 {
		return [][]int{{a[0][0] * b[0][0]}}
	}

	m := n / 2
	a11 := make([][]int, m)
	a12 := make([][]int, m)
	a21 := make([][]int, m)
	a22 := make([][]int, m)
	b11 := make([][]int, m)
	b12 := make([][]int, m)
	b21 := make([][]int, m)
	b22 := make([][]int, m)
	for i := 0; i < m; i++ {
		a11[i] = a[i][:m]
		a12[i] = a[i][m:]
		a21[i] = a[i+m][:m]
		a22[i] = a[i+m][m:]
		b11[i] = b[i][:m]
		b12[i] = b[i][m:]
		b21[i] = b[i+m][:m]
		b22[i] = b[i+m][m:]
	}
	var wg sync.WaitGroup
	wg.Add(8)
	var x1, x2, x3, x4, x5, x6, x7, x8 [][]int
	go func() {
		defer wg.Done()
		x1 = MatrixMultiplyDandCRecParallel(a11, b11, m)
	}()
	go func() {
		defer wg.Done()
		x2 = MatrixMultiplyDandCRecParallel(a12, b21, m)
	}()
	go func() {
		defer wg.Done()
		x3 = MatrixMultiplyDandCRecParallel(a11, b12, m)
	}()
	go func() {
		defer wg.Done()
		x4 = MatrixMultiplyDandCRecParallel(a12, b22, m)
	}()
	go func() {
		defer wg.Done()
		x5 = MatrixMultiplyDandCRecParallel(a21, b11, m)
	}()
	go func() {
		defer wg.Done()
		x6 = MatrixMultiplyDandCRecParallel(a22, b21, m)
	}()
	go func() {
		defer wg.Done()
		x7 = MatrixMultiplyDandCRecParallel(a21, b12, m)
	}()
	go func() {
		defer wg.Done()
		x8 = MatrixMultiplyDandCRecParallel(a22, b22, m)
	}()
	wg.Wait()

	c11 := add(x1, x2)
	c12 := add(x3, x4)
	c21 := add(x5, x6)
	c22 := add(x7, x8)

	c := make([][]int, n)
	for i := range a {
		c[i] = make([]int, n)
		if i < m {
			c[i] = append(c11[i], c12[i]...)
		} else {
			c[i] = append(c21[i-m], c22[i-m]...)
		}
	}
	return c
}
