package pqueue

import (
	"container/heap"
	"sync"

	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

// Base type of PriorityQueue and PriorityQueueEx.
type basePriorityQueue struct {
	h    iheap.Heap
	lock *sync.RWMutex
}

func (pq *basePriorityQueue) Len() int {
	if pq == nil {
		return 0
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	return pq.h.Len()
}

func (pq *basePriorityQueue) Cap() int {
	if pq == nil {
		return 0
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	return pq.h.Cap()
}

func (pq *basePriorityQueue) Reset(capacity int) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if capacity < 0 {
		return ErrNegativeCapacity
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	return gorecover.Recover(func() {
		pq.h.Reset(capacity)
		heap.Init(pq.h)
	})
}

func (pq *basePriorityQueue) Clear() error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	return gorecover.Recover(func() {
		pq.h.Clear()
		heap.Init(pq.h)
	})
}
