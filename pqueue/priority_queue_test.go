package pqueue

import (
	"testing"

	"github.com/donyori/gocontainer"
)

func TestPriorityQueue(t *testing.T) {
	pq, err := NewPriorityQueue(-2, false)
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
	pq, err = NewPriorityQueue(c, true)
	if err != nil {
		t.Fatal(err)
	}
	if c2 := pq.Cap(); c != c2 {
		t.Errorf("cap(%d) != %d", c2, c)
	}
	inputs := []testElement1{3, 0, 9, -4, 3, -5, 8}
	for i := range inputs {
		err = pq.Enqueue(&inputs[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	n1 := pq.Len()
	var wrongInput testElement2 = 1.2
	err = pq.Enqueue(&wrongInput)
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
	x, err := pq.Top()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Top item: %+v", *(x.(*testElement1)))
	x, err = pq.Dequeue()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Dequeue item: %+v", *(x.(*testElement1)))
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
	err = pq.Enqueue(x)
	if err != nil {
		t.Fatal(err)
	}
}
