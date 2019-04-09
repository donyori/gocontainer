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
func (ici *IndexedComparableItem) Set(x Comparable) {
	ici.lock.Lock()
	defer ici.lock.Unlock()
	ici.x = x
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
func (ici *IndexedComparableItem) UpdateIndex(idx int) {
	ici.lock.Lock()
	defer ici.lock.Unlock()
	ici.idx = idx
}

func (ici *IndexedComparableItem) Less(another interface{}) bool {
	a := another.(*IndexedComparableItem)
	return a != nil && (ici == nil || ici.x.Less(a.x))
	// Equivalent to:
	// if ici == nil {
	//     return a != nil
	// } else if a == nil {
	//     return false
	// } else {
	//     return ici.x.Less(a.x)
	// }
}
