package gocontainer

import "sync"

type IndexedComparableItem struct {
	x    Comparable
	idx  int
	lock sync.RWMutex
}

func NewIndexedComparableItem(x Comparable) *IndexedComparableItem {
	return &IndexedComparableItem{x: x}
}

func (ici *IndexedComparableItem) Get() Comparable {
	if ici == nil {
		return nil
	}
	ici.lock.RLock()
	defer ici.lock.RUnlock()
	return ici.x
}

// This method should be called by other container.
// Do NOT call it directly.
func (ici *IndexedComparableItem) Set(x Comparable) error {
	if ici == nil {
		return ErrNilItem
	}
	ici.lock.Lock()
	defer ici.lock.Unlock()
	ici.x = x
	return nil
}

func (ici *IndexedComparableItem) Index() int {
	if ici == nil {
		return 0
	}
	ici.lock.RLock()
	defer ici.lock.RUnlock()
	return ici.idx
}

// This method should be called by other container.
// Do NOT call it directly.
func (ici *IndexedComparableItem) UpdateIndex(idx int) error {
	if ici == nil {
		return ErrNilItem
	}
	ici.lock.Lock()
	defer ici.lock.Unlock()
	ici.idx = idx
	return nil
}

func (ici *IndexedComparableItem) Less(another Comparable) (
	res bool, err error) {
	a, ok := another.(*IndexedComparableItem)
	if !ok {
		return false, ErrWrongType
	}
	if ici == nil {
		return a != nil, nil
	} else if a == nil {
		return false, nil
	} else {
		return ici.x.Less(a.x)
	}
}
