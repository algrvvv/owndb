package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func UnmarshalFloat(data []byte) (any, error) {
	const op = "binarizer.UnmarshalFloat"

	var v float64
	reader := bytes.NewReader(data)

	err := binary.Read(reader, binary.LittleEndian, &v)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return v, nil
}
