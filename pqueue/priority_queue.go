package pqueue

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

type PriorityQueue struct {
	basePriorityQueue
}

func NewPriorityQueue(capacity int, isTopMax, isSync bool) *PriorityQueue {
	if capacity < 0 {
		panic(fmt.Errorf("gocontainer: capacity(%d) is negative", capacity))
	}
	pq := new(PriorityQueue)
	if isSync {
		pq.lock = new(sync.RWMutex)
	}
	// Not necessary to lock during init.
	if isTopMax {
		pq.h = iheap.NewMaxHeap(capacity, false)
	} else {
		pq.h = iheap.NewMinHeap(capacity, false)
	}
	return pq
}

func (pq *PriorityQueue) Top() gocontainer.Comparable {
	if pq == nil {
		return nil
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	return pq.h.Top()
}

func (pq *PriorityQueue) Enqueue(x gocontainer.Comparable) {
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	nBefore := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	pErr := gorecover.Recover(func() {
		heap.Push(pq.h, x)
	})
	if pErr == nil {
		return
	}
	// Try to recover the queue and then panic.
	defer panic(pErr)
	nAfter := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	if nAfter == nBefore {
		// Push didn't modified the queue.
		// Just return.
		return
	}
	// Try to recover the queue:
	var idx int
	// Traverse the whole queue to find wrong item:
	for idx = 0; idx < nAfter; idx++ {
		if pq.h.Get(idx) == x {
			break
		}
	}
	if idx >= nAfter {
		panic(fmt.Errorf("%v; cannot recover the queue, "+
			"because cannot find the wrong item to remove", pErr))
	}
	pErr2 := gorecover.Recover(func() {
		heap.Remove(pq.h, idx)
	})
	if pErr2 != nil {
		panic(fmt.Errorf("%v; error occurs when recover the queue: %v",
			pErr, pErr2))
	}
	// Succeed to recover the queue.
}

func (pq *PriorityQueue) Dequeue() (x gocontainer.Comparable, ok bool) {
	if pq == nil {
		return // nil, false
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	if pq.h.Len() <= 0 { // Do NOT call pq.Len(), which will dead lock!
		return // nil, false
	}
	x = heap.Pop(pq.h).(gocontainer.Comparable)
	ok = true
	return
}

func (pq *PriorityQueue) Scan(
	f func(x gocontainer.Comparable) (doesStop bool)) {
	if pq == nil || f == nil {
		return
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	pq.h.Scan(f)
}
