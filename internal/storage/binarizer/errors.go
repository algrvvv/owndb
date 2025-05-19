package binarizer

import "fmt"

type ErrInvalidDataType struct {
	Key string
}

func (e *ErrInvalidDataType) Error() string {
	return fmt.Sprintf("invalid data type for %q key", e.Key)
}

type ErrTypeMismatch struct {
	Key  string
	Got  any
	Want string
}

func (e *ErrTypeMismatch) Error() string {
	return fmt.Sprintf("type mismatch for %q: got: %T; want: %v", e.Key, e.Got, e.Want)
}
