package heap

import (
	stdheap "container/heap"

	"github.com/donyori/gocontainer"
)

type Heap interface {
	stdheap.Interface
	Cap() int
	Get(i int) gocontainer.Comparable
	Set(i int, x gocontainer.Comparable)
	Top() gocontainer.Comparable
	UpdateTop(x gocontainer.Comparable)
	Clear()
	Reset(capacity int)
}
