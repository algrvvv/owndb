package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

func MarshalInt(w *bytes.Buffer, value int64) error {
	const op = "binarizer.MarshalInt"

	// use int type
	err := w.WriteByte(IntType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// get value size
	valueSize := uint16(unsafe.Sizeof(value))
	err = binary.Write(w, binary.LittleEndian, valueSize)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = binary.Write(w, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
