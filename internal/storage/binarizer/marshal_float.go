package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

func MarshalFloat(w *bytes.Buffer, value float64) error {
	const op = "binarizer.MarshalFloat"

	if err := w.WriteByte(FloatType); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	vlen := uint16(unsafe.Sizeof(value))
	err := binary.Write(w, binary.LittleEndian, vlen)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = binary.Write(w, binary.LittleEndian, value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
