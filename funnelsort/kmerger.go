package funnelsort

import (
	"container/heap"
	"math"
)

type bufEntry struct {
	val int
	buf Buffer
}

// heap of *bufEntry
type entryHeap []*bufEntry

func (h entryHeap) Len() int            { return len(h) }
func (h entryHeap) Less(i, j int) bool  { return h[i].val < h[j].val }
func (h entryHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *entryHeap) Push(x interface{}) { *h = append(*h, x.(*bufEntry)) }
func (h *entryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type kmerger struct {
	in   []Buffer
	heap entryHeap
}

func NewKMerger(arr [][]int) KMerger {
	k := len(arr)
	var in []Buffer

	if k < 512 {
		in = make([]Buffer, k)
		for i := range arr {
			in[i] = NewLeafBuffer(arr[i])
		}
	} else {
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
			in[i] = NewBuffer(m)
		}
	}

	m := &kmerger{in: in}
	m.initHeap()
	return m
}

func (m *kmerger) initHeap() {
	m.heap = make(entryHeap, 0, len(m.in))
	for _, buf := range m.in {
		if v, ok := buf.Peek(); ok {
			m.heap = append(m.heap, &bufEntry{val: v, buf: buf})
		}
	}
	heap.Init(&m.heap)
}

type KMerger interface {
	Next() (int, bool)
	Size() int
}

func (m *kmerger) Size() int { return len(m.in) }

func (m *kmerger) Next() (int, bool) {
	if len(m.heap) == 0 {
		return 0, false
	}
	// pop smallest
	e := heap.Pop(&m.heap).(*bufEntry)
	val := e.val

	// advance buffer; if more, push new entry
	e.buf.Consume()
	if v, ok := e.buf.Peek(); ok {
		e.val = v
		heap.Push(&m.heap, e)
	}
	return val, true
}
