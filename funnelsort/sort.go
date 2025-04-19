package funnelsort

import (
	"math"
)

func Sort(arr []int) []int {
	cbrt := int(math.Cbrt(float64(len(arr))))
	arrs := make([][]int, cbrt)

	step := len(arr) / cbrt
	// Split the array into cbrt parts
	for i := 0; i < cbrt; i++ {
		if i == cbrt-1 {
			arrs[i] = arr[i*step:]
		} else {
			arrs[i] = arr[i*step : (i+1)*step]
		}
	}

	// Create a KMerger with the split arrays
	merger := NewKMerger(arrs)

	// Create output array and fill it
	out := make([]int, 0, len(arr))
	for {
		v, ok := merger.Next()
		if !ok {
			break
		}
		out = append(out, v)
	}

	return out
}
