package binarizer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func UnmarshalInt(data []byte) (any, error) {
	const op = "binarizer.UnmarshalInt"

	var value int64
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &value)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if value < math.MinInt32 || value > math.MaxInt32 {
		return value, nil
	}

	return int(value), nil
}
