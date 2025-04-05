package matrix

import "math/rand"

func MagicMatrix(m, n, seed int) [][]int {
	// Create random generator using seed
	rnd := rand.New(rand.NewSource(int64(seed)))

	// Create random mxn matrix
	matrix := make([][]int, m)
	for i := range matrix {
		matrix[i] = make([]int, n)
		for j := range matrix[i] {
			matrix[i][j] = rnd.Intn(100) // Random number between 0 and 99
		}
	}

	return matrix
}
