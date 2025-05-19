package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

func MarshalBool(w *bytes.Buffer, value bool) error {
	const op = "binarizer.MarshalBool"

	err := w.WriteByte(BoolType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	valueSize := uint16(unsafe.Sizeof(value))
	err = binary.Write(w, binary.LittleEndian, valueSize)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var bytesValue byte
	if value {
		bytesValue = 1
	} else {
		bytesValue = 0
	}

	err = binary.Write(w, binary.LittleEndian, bytesValue)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
