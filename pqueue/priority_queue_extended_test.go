package pqueue

import (
	"testing"

	"github.com/donyori/gocontainer"
	"github.com/donyori/gorecover"
)

func TestPriorityQueueEx(t *testing.T) {
	err := gorecover.Recover(func() {
		NewPriorityQueueEx(-2, false, false)
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	c := 5
	pq := NewPriorityQueueEx(c, true, true)
	if c2 := pq.Cap(); c != c2 {
		t.Errorf("cap(%d) != %d", c2, c)
	}
	inputsData := []testElement1{3, 0, 9, -4, 3, -5, 8}
	n := len(inputsData)
	inputs := make([]*gocontainer.IndexedComparableItem, n)
	for i := 0; i < n; i++ {
		inputs[i] = gocontainer.NewIndexedComparableItem(&inputsData[i])
	}
	for i := range inputs {
		pq.Enqueue(inputs[i])
	}
	n1 := pq.Len()
	var wrongInput testElement2 = 1.2
	err = gorecover.Recover(func() {
		pq.Enqueue(gocontainer.NewIndexedComparableItem(&wrongInput))
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	n2 := pq.Len()
	if n1 != n2 {
		t.Fatal("Enqueue failed but pushed the item into the queue!")
	}
	ici := pq.Top()
	t.Logf("Top item: %+v(index: %d)", *ici.Get().(*testElement1), ici.Index())
	newX := testElement1(-2)
	ok := pq.Update(ici, &newX)
	if !ok {
		t.Fatal("Update failed!")
	}
	if ici.Index() == 0 {
		t.Fatal("Update failed, item didn't move.")
	}
	ok = pq.Update(gocontainer.NewIndexedComparableItem(&newX), &newX)
	if ok {
		t.Fatal("Update should fail but succeeded!")
	}
	x, ok := pq.Dequeue()
	if !ok {
		t.Fatal("Dequeue failed!")
	}
	t.Logf("Dequeue item: %+v(index: %d)", *x.Get().(*testElement1), x.Index())
	ok = pq.Remove(ici)
	if !ok {
		t.Fatal("Remove failed!")
	}
	ok = pq.Remove(x)
	if ok {
		t.Fatal("Remove should fail but succeeded!")
	}
	err = gorecover.Recover(func() {
		pq.Reset(-2)
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	if pq.Len() <= 0 {
		t.Fatal("Reset failed but clean the queue.")
	}
	c = 4
	pq.Reset(c)
	if c2 := pq.Cap(); c != c2 {
		t.Fatalf("cap(%d) != %d", c2, c)
	}
	pq.Clear()
	if c2 := pq.Cap(); c2 != 0 {
		t.Fatalf("cap(%d) != 0", c2)
	}
	_, ok = pq.Dequeue()
	if ok {
		t.Fatal("No item in the queue but Dequeue succeeded!")
	}
	pq.Enqueue(ici)
	if idx := ici.Index(); idx != 0 {
		t.Fatalf("Index(%d) != 0", idx)
	}
}
