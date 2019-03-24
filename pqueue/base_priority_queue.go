package pqueue

import (
	"container/heap"
	"sync"

	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

// Base type of PriorityQueue and PriorityQueueEx.
type basePriorityQueue struct {
	h    *iheap.MaxHeap
	lock sync.RWMutex
}

func (pq *basePriorityQueue) Len() int {
	if pq == nil {
		return 0
	}
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.h.Len()
}

func (pq *basePriorityQueue) Cap() int {
	if pq == nil {
		return 0
	}
	pq.lock.RLock()
	defer pq.lock.RUnlock()
	return pq.h.Cap()
}

func (pq *basePriorityQueue) Reset(capacity int) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if capacity < 0 {
		return ErrNegativeCapacity
	}
	return gorecover.Recover(func() {
		pq.lock.Lock()
		defer pq.lock.Unlock()
		pq.h.Reset(0, capacity)
		heap.Init(pq.h)
	})
}

func (pq *basePriorityQueue) Clear() error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	return gorecover.Recover(func() {
		pq.lock.Lock()
		defer pq.lock.Unlock()
		pq.h.Clear()
		heap.Init(pq.h)
	})
}
