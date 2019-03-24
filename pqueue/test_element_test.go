package pqueue

import "github.com/donyori/gocontainer"

type testElement1 int
type testElement2 float32

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
