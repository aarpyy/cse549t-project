package matmul

func CheckMatMul(a, b, c [][]int) bool {
	// Check if the dimensions of the matrices are correct
	if len(a[0]) != len(b) {
		panic("Matrix dimensions do not match for multiplication")
	}
	if len(a) != len(c) || len(b[0]) != len(c[0]) {
		panic("Result matrix dimensions do not match")
	}

	// Check if the result of multiplying a and b equals c
	for i := range a {
		for j := range b[0] {
			sum := 0
			for k := range b {
				sum += a[i][k] * b[k][j]
			}
			if sum != c[i][j] {
				return false
			}
		}
	}
	return true
}
