package heap

import "github.com/donyori/gocontainer"

type testElement int

func (te testElement) Less(another gocontainer.Comparable) (
	res bool, err error) {
	a, ok := another.(testElement)
	if !ok {
		return false, gocontainer.ErrWrongType
	}
	return te < a, nil
}
