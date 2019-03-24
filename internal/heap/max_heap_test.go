package heap

import (
	stdheap "container/heap"
	"math"
	"testing"
)

func TestMaxHeap(t *testing.T) {
	inputs := []testElement{3, 0, 9, -4, 3, -5, 8}
	inputMax := inputs[0]
	var input2ndMax testElement = math.MinInt32
	n := len(inputs)
	h := NewMaxHeap(n, n, false)
	for i := 0; i < n; i++ {
		if inputMax < inputs[i] {
			input2ndMax = inputMax
			inputMax = inputs[i]
		} else if input2ndMax < inputs[i] {
			input2ndMax = inputs[i]
		}
		h.Set(i, inputs[i])
	}
	stdheap.Init(h)
	t.Logf("After init - heap underlying: %v", h.baseHeap)
	c := h.Cap()
	if c != n {
		t.Errorf("cap(%d) != n(%d)", c, n)
	}
	stdheap.Push(h, testElement(4))
	if inputMax < 4 {
		input2ndMax = inputMax
		inputMax = 4
	} else if input2ndMax < 4 {
		input2ndMax = 4
	}
	t.Logf("After push - heap underlying: %v", h.baseHeap)
	h.UpdateMax(testElement(-2))
	if inputMax < -2 {
		inputMax = -2
	} else {
		inputMax = input2ndMax
	}
	input2ndMax = math.MinInt32
	t.Logf("After update max - heap underlying: %v", h.baseHeap)
	if max := h.GetMax(); max != inputMax {
		t.Errorf("Max item (%d) != %d", max, inputMax)
	}
	t.Log("Start to pop:")
	var last testElement = math.MaxInt32
	for h.Len() > 0 {
		x := stdheap.Pop(h).(testElement)
		isLess, err := last.Less(x)
		if err != nil {
			t.Fatal(err)
		}
		if isLess {
			t.Errorf("Pop a value bigger than last one: current = %d, last = %d", x, last)
		} else {
			t.Logf("Pop %d", x)
		}
		last = x
	}
}
