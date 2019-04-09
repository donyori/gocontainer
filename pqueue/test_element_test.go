package pqueue

type testElement1 int
type testElement2 float32

func (te *testElement1) Less(another interface{}) bool {
	a := another.(*testElement1)
	return a != nil && (te == nil || *te < *a)
}

func (te *testElement2) Less(another interface{}) bool {
	a := another.(*testElement2)
	return a != nil && (te == nil || *te < *a)
}
