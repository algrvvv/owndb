package binarizer

import "fmt"

type ErrInvalidDataType struct {
	Key string
}

func (e *ErrInvalidDataType) Error() string {
	return fmt.Sprintf("invalid data type for %q key", e.Key)
}
