package funnelsort

import (
	"math"
)

type kmerger struct {
	in     []Buffer
	tree   []int // indices of losers
	winner int   // index of current winner
}

func NewKMerger(arr [][]int) KMerger {
	k := len(arr)
	var in []Buffer
	if k < 1000 {
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
			in[i] = NewBuffer(NewKMerger(a))
		}
	}
	m := &kmerger{
		in:   in,
		tree: make([]int, k),
	}
	m.initTree()
	return m
}

type KMerger interface {
	Next() (int, bool)
	Size() int
}

func (m *kmerger) initTree() {
	k := len(m.in)
	// initialize all tree nodes to -1
	for i := range m.tree {
		m.tree[i] = -1
	}
	// build tournament
	for i := 0; i < k; i++ {
		m.play(i)
	}
}

// play takes in the index of a recently consumed buffer and updates the tree so that the new top of this buffer
// is inserted properly according to the sort order
func (m *kmerger) play(idx int) {
	n := len(m.in)

	// Start at the leaf node
	node := (idx + n) / 2
	challenger := idx

	// While we aren't at the root
	for node > 0 {
		loser := m.tree[node-1]
		// get values
		lv, lok := sentinelVal(m.in, loser)
		cv, cok := sentinelVal(m.in, challenger)

		// If we don't have an index of a loser yet or the valid challenger is less than the valid loser, then swap
		if loser < 0 || (cv < lv && lok && cok) {
			m.tree[node-1], challenger = challenger, loser
		}

		// Use implicit indexing to get the parent node
		node /= 2
	}
	m.winner = challenger
}

// sentinelVal returns the value of the i-th buffer
func sentinelVal(in []Buffer, i int) (int, bool) {
	if i < 0 {
		return 0, false
	}
	return in[i].Peek()
}

func (m *kmerger) Size() int { return len(m.in) }

func (m *kmerger) Next() (int, bool) {
	// If we don't have a winner index, we have no next value
	if m.winner < 0 {
		return 0, false
	}
	val, ok := sentinelVal(m.in, m.winner)
	if !ok {
		return 0, false
	}
	m.in[m.winner].Consume()
	m.play(m.winner)
	return val, true
}
