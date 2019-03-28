package heap

import (
	stdheap "container/heap"
	"math"
	"testing"
)

func TestMinHeap(t *testing.T) {
	inputs := []testElement{3, 0, 9, -4, 3, -5, 8}
	inputMin := inputs[0]
	var input2ndMin testElement = math.MaxInt32
	n := len(inputs)
	h := NewMinHeap(n, false)
	for i := 0; i < n; i++ {
		if inputMin > inputs[i] {
			input2ndMin = inputMin
			inputMin = inputs[i]
		} else if input2ndMin > inputs[i] {
			input2ndMin = inputs[i]
		}
		stdheap.Push(h, inputs[i])
	}
	stdheap.Init(h)
	t.Logf("After init - heap underlying: %v", h.baseHeap)
	c := h.Cap()
	if c != n {
		t.Errorf("cap(%d) != n(%d)", c, n)
	}
	stdheap.Push(h, testElement(4))
	if inputMin > 4 {
		input2ndMin = inputMin
		inputMin = 4
	} else if input2ndMin > 4 {
		input2ndMin = 4
	}
	t.Logf("After push - heap underlying: %v", h.baseHeap)
	h.UpdateTop(testElement(2))
	if inputMin > 2 {
		inputMin = 2
	} else {
		inputMin = input2ndMin
	}
	input2ndMin = math.MaxInt32
	t.Logf("After update max - heap underlying: %v", h.baseHeap)
	if min := h.Top(); min != inputMin {
		t.Errorf("Max item (%d) != %d", min, inputMin)
	}
	t.Log("Start to pop:")
	var last testElement = math.MinInt32
	for h.Len() > 0 {
		x := stdheap.Pop(h).(testElement)
		isLess, err := x.Less(last)
		if err != nil {
			t.Fatal(err)
		}
		if isLess {
			t.Errorf("Pop a value smaller than last one: current = %d, last = %d", x, last)
		} else {
			t.Logf("Pop %d", x)
		}
		last = x
	}
	stdheap.Push(h, testElement(4))
	h.Reset(5)
	if h.Len() != 0 || h.Cap() != 5 {
		t.Fatal("Reset failed.")
	}
	stdheap.Push(h, testElement(4))
	h.Clear()
	if h.Len() != 0 || h.Cap() != 0 || h.a != nil {
		t.Fatal("Clear failed.")
	}
}
