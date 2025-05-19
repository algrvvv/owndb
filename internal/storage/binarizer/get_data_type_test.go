package binarizer_test

import (
	"testing"

	"github.com/algrvvv/owndb/internal/storage/binarizer"
)

func TestGetDataType(t *testing.T) {
	tests := []struct {
		key     string
		data    any
		expType byte
		expErr  bool
	}{
		{
			key:     "bar",
			data:    "foo",
			expType: binarizer.StrType,
			expErr:  false,
		},
		{
			key:     "int type",
			data:    12,
			expType: binarizer.IntType,
			expErr:  false,
		},
		{
			key:     "another string",
			data:    "bar",
			expType: binarizer.StrType,
			expErr:  false,
		},
		{
			key:     "float 64",
			data:    12.1234,
			expType: binarizer.FloatType,
			expErr:  false,
		},
		{
			key:     "map value",
			data:    map[string]string{},
			expType: 0x00,
			expErr:  true,
		},
	}

	for _, test := range tests {
		dt, err := binarizer.GetDataType(test.key, test.data)

		if err == nil && test.expErr {
			t.Fatalf("failed: exp err for %q", test.key)
		}

		if err != nil && !test.expErr {
			t.Fatalf("failed: exp err: %v; got error: %s", test.expErr, err)
		}

		if dt != test.expType {
			t.Fatalf("failed: exp type: %v; got: %v", test.expType, dt)
		}

		t.Logf("%q passed", test.key)
	}
}
