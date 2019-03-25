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

func NewPriorityQueueEx(capacity int, isSync bool) (
	pq *PriorityQueueEx, err error) {
	if capacity < 0 {
		return nil, ErrNegativeCapacity
	}
	pq = new(PriorityQueueEx)
	if isSync {
		pq.lock = new(sync.RWMutex)
	}
	err = gorecover.Recover(func() {
		// Not necessary to lock during init.
		pq.h = iheap.NewMaxHeap(0, capacity, true)
		heap.Init(pq.h)
	})
	if err != nil {
		return nil, err
	}
	return pq, nil
}

func (pq *PriorityQueueEx) Top() (
	ici *gocontainer.IndexedComparableItem, err error) {
	if pq == nil {
		return nil, ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.RLock()
		defer pq.lock.RUnlock()
	}
	x := pq.h.GetMax()
	ici, ok := x.(*gocontainer.IndexedComparableItem)
	if !ok {
		return nil, gocontainer.ErrWrongType
	}
	return ici, nil
}

func (pq *PriorityQueueEx) Enqueue(
	ici *gocontainer.IndexedComparableItem) error {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	nBefore := pq.h.Len() // Do NOT call pq.Len(), which will dead lock!
	pErr := gorecover.Recover(func() {
		heap.Push(pq.h, ici)
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
	// Succeed to recover the queue, return the error occurred during enqueue.
	return pErr
}

func (pq *PriorityQueueEx) Dequeue() (
	ici *gocontainer.IndexedComparableItem, err error) {
	if pq == nil {
		return nil, ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		if pq.h.Len() <= 0 { // Do NOT call pq.Len(), which will dead lock!
			ici = nil
			err = ErrEmptyPriorityQueue
			return
		}
		var ok bool
		ici, ok = heap.Pop(pq.h).(*gocontainer.IndexedComparableItem)
		if !ok {
			ici = nil
			err = gocontainer.ErrWrongType
		}
		err = nil
	})
	if pErr != nil {
		err = pErr
	}
	return
}

func (pq *PriorityQueueEx) Update(ici *gocontainer.IndexedComparableItem,
	newX gocontainer.Comparable) (err error) {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		idx := ici.Index()
		if pq.h.Get(idx) != ici {
			err = ErrItemNotInQueue
			return
		}
		err = ici.Set(newX)
		if err != nil {
			return
		}
		heap.Fix(pq.h, idx)
	})
	if pErr != nil {
		err = pErr
	}
	return
}

func (pq *PriorityQueueEx) Remove(ici *gocontainer.IndexedComparableItem) (
	err error) {
	if pq == nil {
		return ErrNilPriorityQueue
	}
	if pq.lock != nil {
		pq.lock.Lock()
		defer pq.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		idx := ici.Index()
		if pq.h.Get(idx) != ici {
			err = ErrItemNotInQueue
			return
		}
		heap.Remove(pq.h, idx)
	})
	if pErr != nil {
		err = pErr
	}
	return
}
