package heap

import (
	stdheap "container/heap"

	"github.com/donyori/gocontainer"
)

// Should be used with "container/heap" package.
type MinHeap struct {
	baseHeap
}

func NewMinHeap(capacity int, isIndexed bool) *MinHeap {
	var a []gocontainer.Comparable
	if capacity != 0 {
		a = make([]gocontainer.Comparable, 0, capacity)
	}
	h := &MinHeap{
		baseHeap: baseHeap{
			a:         a,
			isIndexed: isIndexed,
		},
	}
	stdheap.Init(h)
	return h
}

func (h *MinHeap) Less(i, j int) bool {
	if h == nil {
		panic(ErrNilHeap)
	}
	isLess, err := h.Get(i).Less(h.Get(j))
	if err != nil {
		panic(err)
	}
	return isLess
}

func (h *MinHeap) Set(i int, x gocontainer.Comparable) {
	h.set(i, x)
	stdheap.Fix(h, 0)
}

func (h *MinHeap) UpdateTop(x gocontainer.Comparable) {
	h.Set(0, x)
}
