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

func NewPriorityQueue(capacity int, isTopMax, isSync bool) (
	pq *PriorityQueue, err error) {
	if capacity < 0 {
		return nil, ErrNegativeCapacity
	}
	pq = new(PriorityQueue)
	if isSync {
		pq.lock = new(sync.RWMutex)
	}
	err = gorecover.Recover(func() {
		// Not necessary to lock during init.
		if isTopMax {
			pq.h = iheap.NewMaxHeap(capacity, false)
		} else {
			pq.h = iheap.NewMinHeap(capacity, false)
		}
	})
	if err != nil {
		return nil, err
	}
	return pq, nil
}

func (pq *PriorityQueue) Top() (x gocontainer.Comparable, err error) {
	if pq == nil {
		return nil, ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	return pq.h.Top(), nil
}

func (pq *PriorityQueue) Enqueue(x gocontainer.Comparable) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	nBefore := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	pErr := gorecover.Recover(func() {
		heap.Push(pq.h, x)
	})
	if pErr == nil {
		return nil
	}
	nAfter := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	if nAfter == nBefore {
		// Push didn't modified the queue.
		// Just return the error.
		return pErr
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
	// Succeed to recover the queue, return the error occurred during enqueue.
	return pErr
}

func (pq *PriorityQueue) Dequeue() (x gocontainer.Comparable, err error) {
	if pq == nil {
		return nil, ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		if pq.h.Len() <= 0 { // Do NOT call pq.Len(), which will dead lock!
			x = nil
			err = ErrEmptyPriorityQueue
			return
		}
		var ok bool
		x, ok = heap.Pop(pq.h).(gocontainer.Comparable)
		if !ok {
			x = nil
			err = gocontainer.ErrWrongType
		}
		err = nil
	})
	if pErr != nil {
		err = pErr
	}
	return
}
