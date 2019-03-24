package heap

import "errors"

var (
	ErrNilHeap    error = errors.New("gocontainer: heap is nil")
	ErrNotIndexed error = errors.New("gocontainer: item is NOT indexed")
)
