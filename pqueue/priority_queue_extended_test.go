package pqueue

import (
	"testing"

	"github.com/donyori/gocontainer"
)

func TestPriorityQueueEx(t *testing.T) {
	pq, err := NewPriorityQueueEx(-2, false)
	if err != nil {
		if err == ErrNegativeCapacity {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	c := 5
	pq, err = NewPriorityQueueEx(c, true)
	if err != nil {
		t.Fatal(err)
	}
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
		err = pq.Enqueue(inputs[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	n1 := pq.Len()
	var wrongInput testElement2 = 1.2
	err = pq.Enqueue(gocontainer.NewIndexedComparableItem(&wrongInput))
	if err != nil {
		if err == gocontainer.ErrWrongType {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	n2 := pq.Len()
	if n1 != n2 {
		t.Fatal("Enqueue failed but pushed the item into the queue!")
	}
	ici, err := pq.Top()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Top item: %+v", *ici)
	newX := testElement1(-2)
	err = pq.Update(ici, &newX)
	if err != nil {
		t.Fatal(err)
	}
	if ici.Index() == 0 {
		t.Fatal("Update failed, item didn't move.")
	}
	err = pq.Update(gocontainer.NewIndexedComparableItem(&newX), &newX)
	if err != nil {
		if err == ErrItemNotInQueue {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	x, err := pq.Dequeue()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Dequeue item: %+v", *x)
	err = pq.Remove(ici)
	if err != nil {
		t.Fatal(err)
	}
	err = pq.Remove(x)
	if err != nil {
		if err == ErrItemNotInQueue {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	err = pq.Reset(-2)
	if err != nil {
		if err == ErrNegativeCapacity {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	if pq.Len() <= 0 {
		t.Fatal("Reset failed but clean the queue.")
	}
	c = 4
	err = pq.Reset(c)
	if err != nil {
		t.Fatal(err)
	}
	if c2 := pq.Cap(); c != c2 {
		t.Fatalf("cap(%d) != %d", c2, c)
	}
	err = pq.Clear()
	if err != nil {
		t.Fatal(err)
	}
	if c2 := pq.Cap(); c2 != 0 {
		t.Fatalf("cap(%d) != 0", c2)
	}
	_, err = pq.Dequeue()
	if err != nil {
		if err == ErrEmptyPriorityQueue {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	err = pq.Enqueue(ici)
	if err != nil {
		t.Fatal(err)
	}
	if idx := ici.Index(); idx != 0 {
		t.Fatalf("Index(%d) != 0", idx)
	}
}
