package topkbuf

import (
	"container/heap"
	"errors"
	"sync"

	"github.com/donyori/gocontainer"
	iheap "github.com/donyori/gocontainer/internal/heap"
)

type TopKBuffer struct {
	h    *iheap.Heap
	k    int
	lock sync.RWMutex
}

var (
	ErrNilBuffer    error = errors.New("top-k buffer is nil")
	ErrNonPositiveK error = errors.New("k is non-positive")
	ErrEmptyBuffer  error = errors.New("buffer is empty")
)

func NewTopKBuffer(k int) (tkb *TopKBuffer, err error) {
	tkb = new(TopKBuffer)
	err = tkb.Init(k)
	if err != nil {
		return nil, err
	}
	return tkb, nil
}

func (tkb *TopKBuffer) Init(k int) error {
	if tkb == nil {
		return ErrNilBuffer
	}
	if k <= 0 {
		return ErrNonPositiveK
	}
	tkb.lock.Lock()
	defer tkb.lock.Unlock()
	h := make(iheap.Heap, 0, k)
	tkb.h = &h
	tkb.k = k
	heap.Init(tkb.h)
	return nil
}

func (tkb *TopKBuffer) Len() int {
	if tkb == nil {
		return 0
	}
	tkb.lock.RLock()
	defer tkb.lock.RUnlock()
	return tkb.h.Len()
}

func (tkb *TopKBuffer) K() int {
	if tkb == nil {
		return 0
	}
	tkb.lock.RLock()
	defer tkb.lock.RUnlock()
	return tkb.k
}

func (tkb *TopKBuffer) ResetK(k int) error {
	if tkb == nil {
		return ErrNilBuffer
	}
	if k <= 0 {
		return ErrNonPositiveK
	}
	tkb.lock.Lock()
	defer tkb.lock.Unlock()
	// Pop excess items.
	for i := tkb.k - k; i > 0; i-- {
		heap.Pop(tkb.h)
	}
	tkb.k = k
	return nil
}

func (tkb *TopKBuffer) Add(x gocontainer.Comparable) error {
	if tkb == nil {
		return ErrNilBuffer
	}
	tkb.lock.Lock()
	defer tkb.lock.Unlock()
	if tkb.h.Len() >= tkb.k {
		isLess, err := (*tkb.h)[0].Less(x)
		if err != nil {
			return err
		}
		if isLess {
			(*tkb.h)[0] = x
			heap.Fix(tkb.h, 0)
		}
	} else {
		heap.Push(tkb.h, x)
	}
	return nil
}

func (tkb *TopKBuffer) Flush() (xs []gocontainer.Comparable, err error) {
	if tkb == nil {
		return nil, ErrNilBuffer
	}
	tkb.lock.Lock()
	defer tkb.lock.Unlock()
	n := tkb.h.Len()
	if n <= 0 {
		return nil, nil
	}
	xs = make([]gocontainer.Comparable, n)
	var ok bool
	for i := n - 1; i >= 0; i-- {
		xs[i], ok = heap.Pop(tkb.h).(gocontainer.Comparable)
		if !ok {
			return nil, gocontainer.ErrWrongType
		}
	}
	return xs, nil
}
