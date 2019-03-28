package topkbuf

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
	"github.com/donyori/gorecover"
)

type TopKBuffer struct {
	h    *iheap.MinHeap
	k    int
	lock *sync.RWMutex
}

var (
	ErrNilBuffer    error = errors.New("gocontainer: top-k buffer is nil")
	ErrNonPositiveK error = errors.New("gocontainer: k is non-positive")
)

func NewTopKBuffer(k int, isSync bool) (tkb *TopKBuffer, err error) {
	if k <= 0 {
		return nil, ErrNonPositiveK
	}
	tkb = new(TopKBuffer)
	if isSync {
		tkb.lock = new(sync.RWMutex)
	}
	err = gorecover.Recover(func() {
		// Not necessary to lock during init.
		tkb.h = iheap.NewMinHeap(k, false)
		tkb.k = k
		heap.Init(tkb.h)
	})
	if err != nil {
		return nil, err
	}
	return tkb, nil
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

func (tkb *TopKBuffer) ResetK(k int) error {
	if tkb == nil {
		return ErrNilBuffer
	}
	if k <= 0 {
		return ErrNonPositiveK
	}
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	return gorecover.Recover(func() {
		// Pop excess items.
		for i := tkb.k - k; i > 0; i-- {
			heap.Pop(tkb.h)
		}
		// Set K.
		tkb.k = k
	})
}

func (tkb *TopKBuffer) Add(x gocontainer.Comparable) (err error) {
	if tkb == nil {
		return ErrNilBuffer
	}
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		if tkb.h.Len() >= tkb.k {
			var isLess bool
			isLess, err = (*tkb.h).Top().Less(x)
			if err != nil {
				return
			}
			if isLess {
				tkb.h.UpdateTop(x)
			}
		} else {
			heap.Push(tkb.h, x)
		}
	})
	if pErr != nil {
		err = pErr
	}
	return
}

func (tkb *TopKBuffer) Flush() (xs []gocontainer.Comparable, err error) {
	if tkb == nil {
		return nil, ErrNilBuffer
	}
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	pErr := gorecover.Recover(func() {
		n := tkb.h.Len() // Do NOT call tkb.Len(), which will dead lock!
		if n <= 0 {
			xs = nil
			err = nil
			return
		}
		xs = make([]gocontainer.Comparable, n)
		var ok bool
		for i := n - 1; i >= 0; i-- {
			xs[i], ok = heap.Pop(tkb.h).(gocontainer.Comparable)
			if !ok {
				xs = nil
				err = gocontainer.ErrWrongType
				return
			}
		}
	})
	if pErr != nil {
		err = pErr
	}
	return
}

func (tkb *TopKBuffer) Clear() error {
	if tkb == nil {
		return ErrNilBuffer
	}
	if tkb.lock != nil {
		tkb.lock.Lock()
		defer tkb.lock.Unlock()
	}
	return gorecover.Recover(func() {
		tkb.h.Reset(tkb.k)
		heap.Init(tkb.h)
	})
}
