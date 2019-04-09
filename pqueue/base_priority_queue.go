package pqueue

import (
	"container/heap"
	"fmt"
	"sync"

	iheap "github.com/donyori/gocontainer/internal/heap"
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

func (pq *basePriorityQueue) Reset(capacity int) {
	if capacity < 0 {
		panic(fmt.Errorf("gocontainer: capacity(%d) is negative", capacity))
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pq.h.Reset(capacity)
	heap.Init(pq.h)
}

func (pq *basePriorityQueue) Clear() {
	if pq == nil {
		return
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pq.h.Clear()
	heap.Init(pq.h)
}
