package funnelsort

import (
	"math"
)

type kmerger struct {
	k  int
	in []Buffer
}

func NewKMerger(arr [][]int) KMerger {
	k := len(arr)
	var in []Buffer

	if k < 512 {
		in = make([]Buffer, k)
		for i := 0; i < k; i++ {
			in[i] = NewLeafBuffer(arr[i])
		}
	} else {
		// Create sqrt(k) buffers/mergers
		sqrtK := int(math.Sqrt(float64(k)))
		in = make([]Buffer, sqrtK)
		for i := 0; i < sqrtK; i++ {
			var a [][]int
			if i == sqrtK-1 {
				a = arr[i*sqrtK:]
			} else {
				a = arr[i*sqrtK : (i+1)*sqrtK]
			}
			m := NewKMerger(a)
			in[i] = NewBuffer(m) // Pass size k for buffer calculation
		}
	}

	return &kmerger{
		k:  k,
		in: in,
	}
}

type KMerger interface {
	Next() (int, bool)
	Size() int
}

func (m *kmerger) Size() int {
	return m.k
}

func (m *kmerger) Next() (int, bool) {
	var minVal int
	minIndex := -1
	hasValue := false

	// Sequential access is more efficient than creating goroutines for each peek
	for i, b := range m.in {
		if v, ok := b.Peek(); ok {
			if !hasValue || v < minVal {
				minVal = v
				minIndex = i
				hasValue = true
			}
		}
	}

	if !hasValue {
		return 0, false
	}

	// Consume the top value
	m.in[minIndex].Consume()
	return minVal, true
}
