package heap

import (
	"container/heap"

	"github.com/donyori/gocontainer"
)

// Should be used with "container/heap" package.
type MinHeap struct {
	baseHeap
}

func NewMinHeap(length, capacity int, isIndexed bool) *MinHeap {
	return &MinHeap{
		baseHeap: baseHeap{
			a:         make([]gocontainer.Comparable, length, capacity),
			isIndexed: isIndexed,
		},
	}
}

func (h *MinHeap) GetMin() gocontainer.Comparable {
	if h == nil {
		return nil
	}
	return h.Get(0)
}

func (h *MinHeap) UpdateMin(x gocontainer.Comparable) {
	if h == nil {
		panic(ErrNilHeap)
	}
	h.Set(0, x)
	heap.Fix(h, 0)
}

func (h *MinHeap) Less(i, j int) bool {
	if h == nil {
		panic(ErrNilHeap)
	}
	res, err := h.Get(i).Less(h.Get(j))
	if err != nil {
		panic(err)
	}
	return res
}
