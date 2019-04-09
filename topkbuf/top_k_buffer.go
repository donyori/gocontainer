package topkbuf

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
)

type TopKBuffer struct {
	h    *iheap.MinHeap
	k    int
	lock *sync.RWMutex
}

func NewTopKBuffer(k int, isSync bool) *TopKBuffer {
	if k <= 0 {
		panic(fmt.Errorf("gocontainer: k(%d) is non-positive", k))
	}
	tkb := new(TopKBuffer)
	if isSync {
		tkb.lock = new(sync.RWMutex)
	}
	// Not necessary to lock during init.
	tkb.h = iheap.NewMinHeap(k, false)
	tkb.k = k
	heap.Init(tkb.h)
	return tkb
}

func (tkb *TopKBuffer) Len() int {
	if tkb == nil {
		return 0
	}
	if tkb.lock != nil {
		tkb.lock.RLock()
		defer tkb.lock.RUnlock()
	}
	return tkb.h.Len()
}

func (tkb *TopKBuffer) Cap() int {
	if tkb == nil {
		return 0
	}
	if tkb.lock != nil {
		tkb.lock.RLock()
		defer tkb.lock.RUnlock()
	}
	return tkb.h.Cap()
}

func (tkb *TopKBuffer) K() int {
	if tkb == nil {
		return 0
	}
	if tkb.lock != nil {
		tkb.lock.RLock()
		defer tkb.lock.RUnlock()
	}
	return tkb.k
}

func (tkb *TopKBuffer) ResetK(k int) {
	if k <= 0 {
		panic(fmt.Errorf("gocontainer: k(%d) is non-positive", k))
	}
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	// Pop excess items.
	for i := tkb.k - k; i > 0; i-- {
		heap.Pop(tkb.h)
	}
	// Set K.
	tkb.k = k
}

func (tkb *TopKBuffer) Add(x gocontainer.Comparable) {
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	if tkb.h.Len() >= tkb.k {
		var isLess bool
		isLess = (*tkb.h).Top().Less(x)
		if isLess {
			tkb.h.UpdateTop(x)
		}
	} else {
		heap.Push(tkb.h, x)
	}
}

func (tkb *TopKBuffer) Flush() []gocontainer.Comparable {
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	n := tkb.h.Len() // Do NOT call tkb.Len(), which will dead lock!
	if n <= 0 {
		return nil
	}
	xs := make([]gocontainer.Comparable, n)
	// Output in reverse order, in order to let the biggest item at 0 position.
	for i := n - 1; i >= 0; i-- {
		xs[i] = heap.Pop(tkb.h).(gocontainer.Comparable)
	}
	return xs
}

func (tkb *TopKBuffer) Clear() {
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	tkb.h.Reset(tkb.k)
	heap.Init(tkb.h)
}
