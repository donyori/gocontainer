package pqueue

import (
	"testing"

	"github.com/donyori/gorecover"
)

func TestPriorityQueue(t *testing.T) {
	err := gorecover.Recover(func() {
		NewPriorityQueue(-2, false, false)
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	c := 5
	pq := NewPriorityQueue(c, true, true)
	if c2 := pq.Cap(); c != c2 {
		t.Errorf("cap(%d) != %d", c2, c)
	}
	inputs := []testElement1{3, 0, 9, -4, 3, -5, 8}
	for i := range inputs {
		pq.Enqueue(&inputs[i])
	}
	n1 := pq.Len()
	var wrongInput testElement2 = 1.2
	err = gorecover.Recover(func() {
		pq.Enqueue(&wrongInput)
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
	x := pq.Top()
	t.Logf("Top item: %+v", *(x.(*testElement1)))
	x, ok := pq.Dequeue()
	if !ok {
		t.Fatal("Dequeue failed!")
	}
	t.Logf("Dequeue item: %+v", *(x.(*testElement1)))
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
	pq.Enqueue(x)
}
