package binarizer

import (
	"github.com/algrvvv/owndb/internal/storage"
	"github.com/rs/zerolog"
)

const (
	StrType   byte = 0x01
	IntType   byte = 0x02
	FloatType byte = 0x03
	BoolType  byte = 0x04
)

type Binarizer struct {
	logger zerolog.Logger
}

func NewBinaryMarshaller(log zerolog.Logger) storage.Marshaller {
	return &Binarizer{
		logger: log,
	}
}
