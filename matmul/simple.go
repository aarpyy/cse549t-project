package matmul

func MatrixMultiply(a, b [][]int) [][]int {
	if len(a[0]) != len(b) {
		panic("Matrix dimensions do not match for multiplication")
	}
	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}
	for i := range a {
		for j := range b[0] {
			for k := range b {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}
