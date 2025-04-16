package matmul

import "math"

func ComputeBlockSize(cacheSize string) int {
	// Parse the cache size from a string (e.g., "32KB")
	z := 0
	for i := 0; i < len(cacheSize); i++ {
		if cacheSize[i] >= '0' && cacheSize[i] <= '9' {
			z = z*10 + int(cacheSize[i]-'0')
		} else {
			switch cacheSize[i] {
			case 'K':
				z *= 1024
			case 'M':
				z *= 1024 * 1024
			case 'G':
				z *= 1024 * 1024 * 1024
			case ' ':
				// Ignore spaces
				continue
			}

			// Once we have read the first char of the size unit we are done
			break
		}
	}

	// Block size is (z / 3) ^ (1/2)
	floatZ := math.Sqrt(float64(z) / 3.0)
	return int(math.Floor(floatZ))
}

func blockMM(a, b, c [][]int, i, j, k, blockSize int) {
	for ii := i; ii < min(i+blockSize, len(a)); ii++ {
		for jj := j; jj < min(j+blockSize, len(b[0])); jj++ {
			for kk := k; kk < min(k+blockSize, len(b)); kk++ {
				c[ii][jj] += a[ii][kk] * b[kk][jj]
			}
		}
	}
}

func CacheAwareMM(a, b [][]int, blockSize int) [][]int {
	if len(a[0]) != len(b) {
		panic("Matrix dimensions do not match for multiplication")
	}
	result := make([][]int, len(a))
	for i := range result {
		result[i] = make([]int, len(b[0]))
	}

	for i := 0; i < len(a); i += blockSize {
		for j := 0; j < len(b[0]); j += blockSize {
			for k := 0; k < len(b); k += blockSize {
				blockMM(a, b, result, i, j, k, blockSize)
			}
		}
	}

	return result
}
