package binarizer_test

import (
	"testing"

	"github.com/algrvvv/owndb/internal/logger"
	"github.com/algrvvv/owndb/internal/storage/binarizer"
)

func TestMarshalNUnmarshal(t *testing.T) {
	log := logger.MustInit("owndb.log", true)
	// log = zerolog.New(io.Discard)

	m := map[string]any{
		"foo":         "bar",
		"bar":         "foo",
		"test":        22,
		"value":       398475,
		"fl":          float64(12),
		"tryfl":       float64(12.21),
		"true_value":  true,
		"false_value": false,
	}

	// marshal data
	marshaller := binarizer.NewBinaryMarshaller(log)
	data, err := marshaller.Marshal(m)
	if err != nil {
		t.Fatalf("FAIL! failed to marshal data: %v", err)
	}

	// unmarshal data for check
	rm, err := marshaller.Unmarshal(data)
	if err != nil {
		t.Fatalf("FAIL! failed to unmarshal data: %v", err)
	}
	t.Logf("unmarshalled map: %v", rm)

	for k, v := range rm {
		srcv, ok := m[k]
		if !ok {
			t.Fatalf("FAIL! key %q not exists in source map", k)
		}

		if srcv != v {
			t.Fatalf("FAIL! value from %q: *%T* '%v' not equal source value: *%T* '%v'", k, v, v, srcv, srcv)
		}

		t.Logf("OK! value from %q: *%T* '%v' equal source value: *%T* '%v'", k, v, v, srcv, srcv)
	}
}
