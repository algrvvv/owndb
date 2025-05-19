package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func UnmarshalBool(data []byte) (any, error) {
	const op = "binarizer.UnmarshalBool"

	reader := bytes.NewReader(data)

	var boolBytes byte
	err := binary.Read(reader, binary.LittleEndian, &boolBytes)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return boolBytes == 1, nil
}
