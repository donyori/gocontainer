package heap

import (
	"errors"

	"github.com/donyori/gocontainer"
)

// Base type of MinHeap and MaxHeap.
type baseHeap struct {
	a         []gocontainer.Comparable
	isIndexed bool
}

func (h *baseHeap) Len() int {
	if h == nil {
		return 0
	}
	return len(h.a)
}

func (h *baseHeap) Swap(i, j int) {
	if h == nil || h.a == nil {
		panic(ErrNilHeap)
	}
	h.a[i], h.a[j] = h.a[j], h.a[i]
	if h.isIndexed {
		x := h.a[i].(gocontainer.Indexed)
		y := h.a[j].(gocontainer.Indexed)
		err := x.UpdateIndex(i)
		if err != nil {
			panic(err)
		}
		err = y.UpdateIndex(j)
		if err != nil {
			panic(err)
		}
	}
}

// This method should be called by "container/heap" package.
// Do NOT call it directly.
func (h *baseHeap) Push(x interface{}) {
	if h == nil {
		panic(ErrNilHeap)
	}
	cx, ok := x.(gocontainer.Comparable)
	if !ok {
		panic(gocontainer.ErrWrongType)
	}
	if h.isIndexed {
		ix, ok := x.(gocontainer.Indexed)
		if !ok {
			panic(ErrNotIndexed)
		}
		err := ix.UpdateIndex(len(h.a))
		if err != nil {
			panic(err)
		}
	}
	h.a = append(h.a, cx)
}

// This method should be called by "container/heap" package.
// Do NOT call it directly.
func (h *baseHeap) Pop() interface{} {
	if h == nil || h.a == nil {
		panic(ErrNilHeap)
	}
	old := h.a
	last := len(old) - 1
	x := old[last]
	old[last] = nil // To avoid potential memory leak.
	if h.isIndexed {
		ix, ok := x.(gocontainer.Indexed)
		if !ok {
			panic(ErrNotIndexed)
		}
		err := ix.UpdateIndex(-1) // for safety
		if err != nil {
			panic(err)
		}
	}
	h.a = old[:last]
	return x
}

func (h *baseHeap) Cap() int {
	if h == nil {
		return 0
	}
	return cap(h.a)
}

func (h *baseHeap) Get(i int) gocontainer.Comparable {
	if h == nil || i < 0 || i >= len(h.a) {
		return nil
	}
	return h.a[i]
}

// For convenience to implement Set() in Heap interface,
//   and set an item directly without fix the heap in test.
// Set() should fix the heap by container/heap.Fix() after calling this method.
func (h *baseHeap) set(i int, x gocontainer.Comparable) {
	if h == nil || h.a == nil {
		panic(ErrNilHeap)
	}
	if i < 0 || i >= len(h.a) {
		panic(errors.New("index out of range"))
	}
	if h.isIndexed {
		ix, ok := x.(gocontainer.Indexed)
		if !ok {
			panic(ErrNotIndexed)
		}
		err := ix.UpdateIndex(i)
		if err != nil {
			panic(err)
		}
	}
	h.a[i] = x
}

func (h *baseHeap) Top() gocontainer.Comparable {
	return h.Get(0)
}

func (h *baseHeap) Clear() {
	if h == nil {
		panic(ErrNilHeap)
	}
	h.a = nil
}

func (h *baseHeap) Reset(capacity int) {
	if h == nil {
		panic(ErrNilHeap)
	}
	h.a = make([]gocontainer.Comparable, 0, capacity)
}
