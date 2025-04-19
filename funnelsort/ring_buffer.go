package funnelsort

type RingBuffer struct {
	size   int
	buffer []int
	index  int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		size:   size,
		buffer: make([]int, size),
		index:  0,
	}
}

func (rb *RingBuffer) Push(value int) {
	rb.buffer[rb.index] = value
	rb.index = (rb.index + 1) % rb.size
	rb.size++
}

func (rb *RingBuffer) Pop() int {
	value := rb.buffer[rb.index]
	rb.index = (rb.index + 1) % rb.size
	rb.size--
	return value
}

func (rb *RingBuffer) Take(count int) []int {
	if rb.index+count < rb.size {
		result := rb.buffer[rb.index : rb.index+count]
		rb.index += count
		rb.size -= count
		return result
	}

	result := make([]int, count)
	copy(result, rb.buffer[rb.index:])
	copy(result[len(rb.buffer)-rb.index:], rb.buffer[:count-len(rb.buffer)+rb.index])
	rb.index = (rb.index + count) % rb.size
	rb.size -= count
	return result
}
