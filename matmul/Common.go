package matmul

func add(a, b [][]int) [][]int {
	result := make([][]int, len(a))
	for i := range a {
		result[i] = make([]int, len(b[0]))
		for j := range b {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result
}
