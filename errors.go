package gocontainer

import "errors"

var (
	ErrWrongType error = errors.New("gocontainer: type is wrong")
	ErrNilItem   error = errors.New("gocontainer: item is nil")
)
