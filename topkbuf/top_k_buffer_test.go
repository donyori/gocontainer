package topkbuf

import (
	"testing"

	"github.com/donyori/gocontainer"
)

type testElement1 int
type testElement2 float32

func TestTopKBuffer(t *testing.T) {
	tkb, err := NewTopKBuffer(-2)
	if err != nil {
		if err == ErrNonPositiveK {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	k := 5
	tkb, err = NewTopKBuffer(k)
	if err != nil {
		t.Fatal(err)
	}
	inputs := []testElement1{3, 0, 9, -4, 3, -5, 8}
	for i := range inputs {
		err = tkb.Add(&inputs[i])
		if err != nil {
			t.Fatal(err)
		}
		l := tkb.Len()
		k := tkb.K()
		shouldLen := i + 1
		if shouldLen > k {
			shouldLen = k
		}
		if l == shouldLen {
			t.Logf("tkb.Len() = %d", l)
		} else {
			t.Fatalf("tkb.Len(): %d != %d", l, shouldLen)
		}
	}
	var wrongInput testElement2 = 1.2
	err = tkb.Add(&wrongInput)
	if err != nil {
		if err == gocontainer.ErrWrongType {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
	} else {
		t.Fatal("No error but should have one.")
	}
	c := tkb.Cap()
	if c != k {
		t.Errorf("cap(%d) != k(%d)", c, k)
	}
	err = tkb.ResetK(tkb.K() - 1)
	if err != nil {
		t.Fatal(err)
	}
	outputs, err := tkb.Flush()
	if err != nil {
		t.Fatal(err)
	}
	for i, x := range outputs {
		v := x.(*testElement1)
		t.Logf("%d: %v", i, *v)
	}
	for i := range inputs {
		err = tkb.Add(&inputs[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	tkb.Clear()
	if tkb.Len() != 0 {
		t.Fatal("Not empty after Clear().")
	}
	err = tkb.Add(&inputs[0])
	if err != nil {
		t.Fatal(err)
	}
}

func (te *testElement1) Less(another gocontainer.Comparable) (
	res bool, err error) {
	a, ok := another.(*testElement1)
	if !ok {
		return false, gocontainer.ErrWrongType
	}
	if te == nil {
		return a != nil, nil
	} else if a == nil {
		return false, nil
	} else {
		return *te < *a, nil
	}
}

func (te *testElement2) Less(another gocontainer.Comparable) (
	res bool, err error) {
	a, ok := another.(*testElement2)
	if !ok {
		return false, gocontainer.ErrWrongType
	}
	if te == nil {
		return a != nil, nil
	} else if a == nil {
		return false, nil
	} else {
		return *te < *a, nil
	}
}
