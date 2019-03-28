package heap

import (
	stdheap "container/heap"
	"testing"

	"github.com/donyori/gocontainer"
	"github.com/donyori/gorecover"
)

func TestIndexed(t *testing.T) {
	inputs := []testElement{3, 0, 9, -4, 3, -5, 8}
	n := len(inputs)
	h := NewMinHeap(n, true)
	err := gorecover.Recover(func() {
		h.Set(0, inputs[0])
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	for i := 0; i < n; i++ {
		stdheap.Push(h, gocontainer.NewIndexedComparableItem(inputs[i]))
	}
	checkIndex(t, h)
	stdheap.Init(h)
	checkIndex(t, h)
	oldA := make([]gocontainer.Comparable, n)
	m := copy(oldA, h.a)
	if m != n {
		t.Fatalf("Copy failed, copied %d of %d", m, n)
	}
	err = gorecover.Recover(func() {
		stdheap.Push(h, testElement(4))
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	for i := 0; i < n; i++ {
		if h.a[i] != oldA[i] {
			t.Fatal("Push failed but modified heap!")
		}
	}
	stdheap.Push(h, gocontainer.NewIndexedComparableItem(testElement(4)))
	n++
	checkIndex(t, h)
	for h.Len() > 0 {
		x := stdheap.Pop(h).(*gocontainer.IndexedComparableItem)
		t.Logf("Pop %+v", *x)
		checkIndex(t, h)
	}
}

func checkIndex(t *testing.T, h *MinHeap) {
	n := h.Len()
	for i := 0; i < n; i++ {
		idx := h.Get(i).(*gocontainer.IndexedComparableItem).Index()
		if i != idx {
			t.Errorf("Index error: index %d at pos %d", idx, i)
		}
	}
}
