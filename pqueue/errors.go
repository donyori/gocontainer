package pqueue

import "errors"

var (
	ErrNilPriorityQueue   error = errors.New("gocontainer: priority queue is nil")
	ErrNegativeCapacity   error = errors.New("gocontainer: capacity is negative")
	ErrEmptyPriorityQueue error = errors.New("gocontainer: queue is empty")
	ErrItemNotInQueue     error = errors.New("gocontainer: item is NOT in the queue")
)
