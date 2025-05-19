package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func MarshalString(w *bytes.Buffer, value string) error {
	const op = "binarizer.MarshalString"

	err := w.WriteByte(StrType)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// преобразуем в uint16
	// так как минимальный размер + не отрицательное
	strlen := uint16(len(value))
	err = binary.Write(w, binary.LittleEndian, strlen)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = w.WriteString(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
