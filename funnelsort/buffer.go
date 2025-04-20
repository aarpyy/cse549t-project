package funnelsort

import (
	"math"
	"sort"
)

type leafbuffer struct {
	arr []int
}

type buffer struct {
	capacity      int
	fillThreshold int
	in            KMerger
	arr           []int
}

func NewLeafBuffer(arr []int) Buffer {
	sort.Ints(arr)
	return &leafbuffer{arr: arr}
}

func NewBuffer(in KMerger) Buffer {
	k := in.Size()
	// Calculate capacity as 2k^(3/2)
	fillThreshold := int(math.Pow(float64(k), 1.5)) // Fill when less than k^(3/2) elements
	capacity := 2 * fillThreshold

	return &buffer{
		capacity:      capacity,
		fillThreshold: fillThreshold,
		in:            in,
		arr:           make([]int, 0, capacity),
	}
}

type Buffer interface {
	Consume()
	Peek() (int, bool)
	Fill()
}

func (b *leafbuffer) Consume() {
	if len(b.arr) == 0 {
		return
	}
	b.arr = b.arr[1:]
}

func (b *leafbuffer) Peek() (int, bool) {
	if len(b.arr) == 0 {
		return 0, false
	}
	return b.arr[0], true
}

func (b *leafbuffer) Fill() {
	// Leaf buffers don't need filling
}

func (b *buffer) Fill() {
	// Only fill if below threshold
	if len(b.arr) >= b.fillThreshold {
		return
	}

	// Fill up to capacity
	for len(b.arr) < b.capacity {
		v, ok := b.in.Next()
		if !ok {
			break
		}
		b.arr = append(b.arr, v)
	}
}

func (b *buffer) Consume() {
	if len(b.arr) == 0 {
		return
	}
	b.arr = b.arr[1:]

	// Check if we need to fill
	if len(b.arr) < b.fillThreshold {
		b.Fill()
	}
}

func (b *buffer) Peek() (int, bool) {
	if len(b.arr) == 0 {
		return 0, false
	}

	return b.arr[0], true
}
