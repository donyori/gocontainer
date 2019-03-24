package gocontainer

type Comparable interface {
	Less(another Comparable) (res bool, err error)
}

type Indexed interface {
	Index() int
	UpdateIndex(idx int) error
}
