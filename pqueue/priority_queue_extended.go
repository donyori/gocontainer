package pqueue

import (
	"container/heap"
	"fmt"
	"sync"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

type PriorityQueueEx struct {
	basePriorityQueue
}

func NewPriorityQueueEx(capacity int, isTopMax, isSync bool) *PriorityQueueEx {
	if capacity < 0 {
		panic(fmt.Errorf("gocontainer: capacity(%d) is negative", capacity))
	}
	pq := new(PriorityQueueEx)
	if isSync {
		pq.lock = new(sync.RWMutex)
	}
	// Not necessary to lock during init.
	if isTopMax {
		pq.h = iheap.NewMaxHeap(capacity, true)
	} else {
		pq.h = iheap.NewMinHeap(capacity, true)
	}
	return pq
}

func (pq *PriorityQueueEx) Top() *gocontainer.IndexedComparableItem {
	if pq == nil {
		return nil
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	x := pq.h.Top()
	return x.(*gocontainer.IndexedComparableItem)
}

func (pq *PriorityQueueEx) Enqueue(ici *gocontainer.IndexedComparableItem) {
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	nBefore := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	pErr := gorecover.Recover(func() {
		heap.Push(pq.h, ici)
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
	idx := ici.Index()
	if pq.h.Get(idx) != ici {
		// Traverse the whole queue to find wrong item:
		for idx = 0; idx < nAfter; idx++ {
			if pq.h.Get(idx) == ici {
				break
			}
		}
		if idx >= nAfter {
			panic(fmt.Errorf("%v; cannot recover the queue, "+
				"because cannot find the wrong item to remove", pErr))
		}
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

func (pq *PriorityQueueEx) Dequeue() (
	ici *gocontainer.IndexedComparableItem, ok bool) {
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
	ici = heap.Pop(pq.h).(*gocontainer.IndexedComparableItem)
	ok = true
	return
}

func (pq *PriorityQueueEx) Update(ici *gocontainer.IndexedComparableItem,
	newX gocontainer.Comparable) (ok bool) {
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	idx := ici.Index()
	if pq.h.Get(idx) != ici {
		return false
	}
	ici.Set(newX)
	heap.Fix(pq.h, idx)
	return true
}

func (pq *PriorityQueueEx) Remove(ici *gocontainer.IndexedComparableItem) (
	ok bool) {
	if pq == nil {
		return false
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	idx := ici.Index()
	if pq.h.Get(idx) != ici {
		return false
	}
	heap.Remove(pq.h, idx)
	return true
}

func (pq *PriorityQueueEx) Scan(
	f func(ici *gocontainer.IndexedComparableItem) (doesStop bool)) {
	if pq == nil || f == nil {
		return
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	pq.h.Scan(func(x gocontainer.Comparable) bool {
		return f(x.(*gocontainer.IndexedComparableItem))
	})
}
