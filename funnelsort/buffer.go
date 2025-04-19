package funnelsort

import "sort"

type leafbuffer struct {
	arr []int
}

type buffer struct {
	capacity int
	in       KMerger
	arr      []int
}

func NewLeafBuffer(arr []int) Buffer {
	sort.Ints(arr)
	return &leafbuffer{arr: arr}
}

func NewBuffer(capacity int, in KMerger) Buffer {
	return &buffer{
		capacity: capacity,
		in:       in,
		arr:      make([]int, 0, capacity),
	}
}

type Buffer interface {
	Next() (int, bool)
	Peek() (int, bool)
}

func (b *leafbuffer) Next() (int, bool) {
	if len(b.arr) == 0 {
		return 0, false
	}
	v := b.arr[0]
	b.arr = b.arr[1:]
	return v, true
}

func (b *leafbuffer) Peek() (int, bool) {
	if len(b.arr) == 0 {
		return 0, false
	}
	return b.arr[0], true
}

func (b *buffer) fill() {
	// Try to fill from the input merger
	for i := 0; i < b.capacity; i++ {
		v, ok := b.in.Next()
		if !ok {
			break
		}
		b.arr = append(b.arr, v)
	}
}

func (b *buffer) Next() (int, bool) {
	if len(b.arr) == 0 {
		b.fill()
		if len(b.arr) == 0 {
			return 0, false
		}
	}

	v := b.arr[0]
	b.arr = b.arr[1:]
	return v, true
}

func (b *buffer) Peek() (int, bool) {
	if len(b.arr) == 0 {
		b.fill()
		if len(b.arr) == 0 {
			return 0, false
		}
	}

	return b.arr[0], true
}
