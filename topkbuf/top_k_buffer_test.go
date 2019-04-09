package topkbuf

import (
	"testing"

	"github.com/donyori/gorecover"
)

type testElement1 int
type testElement2 float32

func TestTopKBuffer(t *testing.T) {
	err := gorecover.Recover(func() {
		NewTopKBuffer(-2, false)
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	k := 5
	tkb := NewTopKBuffer(k, true)
	inputs := []testElement1{3, 0, 9, -4, 3, -5, 8}
	for i := range inputs {
		tkb.Add(&inputs[i])
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
	err = gorecover.Recover(func() {
		tkb.Add(&wrongInput)
	})
	if err != nil {
		t.Log(err)
	} else {
		t.Fatal("No error but should have one.")
	}
	c := tkb.Cap()
	if c != k {
		t.Errorf("cap(%d) != k(%d)", c, k)
	}
	tkb.ResetK(tkb.K() - 1)
	outputs := tkb.Flush()
	for i, x := range outputs {
		v := x.(*testElement1)
		t.Logf("%d: %v", i, *v)
	}
	for i := range inputs {
		tkb.Add(&inputs[i])
	}
	tkb.Clear()
	if tkb.Len() != 0 {
		t.Fatal("Not empty after Clear().")
	}
	tkb.Add(&inputs[0])
}

func (te *testElement1) Less(another interface{}) bool {
	a := another.(*testElement1)
	return a != nil && (te == nil || *te < *a)
}

func (te *testElement2) Less(another interface{}) bool {
	a := another.(*testElement2)
	return a != nil && (te == nil || *te < *a)
}
