package heap

type testElement int

func (te testElement) Less(another interface{}) bool {
	a := another.(testElement)
	return te < a
}
