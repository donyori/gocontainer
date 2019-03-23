package gocontainer

type Comparable interface {
	Less(another Comparable) (res bool, err error)
}
