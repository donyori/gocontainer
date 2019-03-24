package heap

import (
	"container/heap"

	"github.com/donyori/gocontainer"
)

// Should be used with "container/heap" package.
type MaxHeap struct {
	baseHeap
}

func NewMaxHeap(length, capacity int, isIndexed bool) *MaxHeap {
	return &MaxHeap{
		baseHeap: baseHeap{
			a:         make([]gocontainer.Comparable, length, capacity),
			isIndexed: isIndexed,
		},
	}
}

func (h *MaxHeap) GetMax() gocontainer.Comparable {
	if h == nil {
		return nil
	}
	return h.Get(0)
}

func (h *MaxHeap) UpdateMax(x gocontainer.Comparable) {
	if h == nil {
		panic(ErrNilHeap)
	}
	h.Set(0, x)
	heap.Fix(h, 0)
}

func (h *MaxHeap) Less(i, j int) bool {
	if h == nil {
		panic(ErrNilHeap)
	}
	isLess, err := h.Get(i).Less(h.Get(j))
	if err != nil {
		panic(err)
	}
	return !isLess
}
