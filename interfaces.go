package gocontainer

type Comparable interface {
	Less(another interface{}) bool
}

type Indexed interface {
	Index() int
	UpdateIndex(idx int)
}
