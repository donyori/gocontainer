package pqueue

import (
	"container/heap"
	"fmt"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

type PriorityQueue struct {
	basePriorityQueue
}

func NewPriorityQueue(capacity int) (pq *PriorityQueue, err error) {
	pq = new(PriorityQueue)
	err = pq.Init(capacity)
	if err != nil {
		return nil, err
	}
	return pq, nil
}

func (pq *PriorityQueue) Init(capacity int) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if capacity < 0 {
		return ErrNegativeCapacity
	}
	return gorecover.Recover(func() {
		pq.lock.Lock()
		defer pq.lock.Unlock()
		pq.h = iheap.NewMaxHeap(0, capacity, false)
		heap.Init(pq.h)
	})
}

func (pq *PriorityQueue) Top() (x gocontainer.Comparable, err error) {
	if pq == nil {
		return nil, ErrNilPriorityQueue
	}
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.h.GetMax(), nil
}

func (pq *PriorityQueue) Enqueue(x gocontainer.Comparable) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	pq.lock.Lock()
	defer pq.lock.Unlock()
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
	pErr := gorecover.Recover(func() {
		pq.lock.Lock()
		defer pq.lock.Unlock()
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
