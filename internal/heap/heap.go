package heap

import (
	"errors"

	"github.com/donyori/gocontainer"
)

type Heap []gocontainer.Comparable

var ErrNilHeap error = errors.New("heap is nil")

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	if h == nil {
		panic(ErrNilHeap)
	}
	res, err := h[i].Less(h[j])
	if err != nil {
		panic(err)
	}
	return res
}

func (h Heap) Swap(i, j int) {
	if h == nil {
		panic(ErrNilHeap)
	}
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	if h == nil {
		panic(ErrNilHeap)
	}
	*h = append(*h, x.(gocontainer.Comparable))
}

func (h *Heap) Pop() interface{} {
	if h == nil {
		panic(ErrNilHeap)
	}
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
