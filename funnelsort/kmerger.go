package funnelsort

import (
	"math"
	"sync"
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
			in[i] = NewBuffer(sqrtK, m)
		}
	}

	return &kmerger{
		k:  k,
		in: in,
	}
}

type KMerger interface {
	Next() (int, bool)
}

func (m *kmerger) Next() (int, bool) {
	var minVal *int
	minIndex := -1

	var wg sync.WaitGroup
	peeked := make([]*int, len(m.in))
	for i, b := range m.in {
		wg.Add(1)
		go func(i int, b Buffer) {
			defer wg.Done()
			if v, ok := b.Peek(); ok {
				peeked[i] = &v
			} else {
				peeked[i] = nil
			}
		}(i, b)
	}

	wg.Wait()

	// Iterate over input buffers
	for i, v := range peeked {
		if v == nil {
			continue
		}
		if minVal == nil || *v < *minVal {
			minVal = v
			minIndex = i
		}
	}

	if minVal == nil {
		return 0, false
	}

	// Consume the top value
	m.in[minIndex].Consume()
	return *minVal, true
}
